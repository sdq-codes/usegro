package catalogControllers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/usegro/services/catalog/config"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	catalogModels "github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/interface/response"
	"github.com/usegro/services/shared/pkg/storage"
	"gorm.io/gorm"
)

type MediaController struct {
	db *gorm.DB
}

func NewMediaController(db *gorm.DB) *MediaController {
	return &MediaController{db: db}
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

var allowedPresignTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

const maxFileSizeBytes = 20 * 1024 * 1024 // 20 MB

func (mc *MediaController) UploadMedia(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	itemIDStr := c.FormValue("item_id")
	if itemIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "item_id is required",
		})
	}

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "invalid item_id",
		})
	}

	// Verify item belongs to this CRM
	var item catalogModels.CatalogItem
	if err := mc.db.Where("id = ? AND crm_id = ?", itemID, crmID).First(&item).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusNotFound,
			ResponseMessage: "item not found",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "invalid multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "no files provided",
		})
	}

	r2cfg := config.GetConfig().R2
	r2 := storage.NewR2Client(storage.R2Config{
		AccountID:       r2cfg.AccountID,
		AccessKeyID:     r2cfg.AccessKeyID,
		SecretAccessKey: r2cfg.SecretAccessKey,
		BucketName:      r2cfg.BucketName,
		PublicURL:       r2cfg.PublicURL,
	})

	// Get current max position for this item
	var maxPos int
	mc.db.Model(&catalogModels.CatalogItemMedia{}).
		Where("item_id = ?", itemID).
		Select("COALESCE(MAX(position), -1)").
		Scan(&maxPos)

	var results []catalogModels.CatalogItemMedia

	for i, fileHeader := range files {
		if fileHeader.Size > maxFileSizeBytes {
			return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusBadRequest,
				ResponseMessage: fmt.Sprintf("file %s exceeds 20MB limit", fileHeader.Filename),
			})
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if !allowedMimeTypes[contentType] {
			return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusBadRequest,
				ResponseMessage: fmt.Sprintf("unsupported file type: %s", contentType),
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusInternalServerError,
				ResponseMessage: "failed to open file",
			})
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		key := fmt.Sprintf("catalog/%s/%s-%d%s", itemID, time.Now().Format("20060102150405"), i, ext)

		url, err := r2.Upload(c.Context(), key, file, contentType, fileHeader.Size)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusInternalServerError,
				ResponseMessage: "upload failed",
			})
		}

		media := catalogModels.CatalogItemMedia{
			ItemID:   itemID,
			URL:      url,
			Key:      key,
			MimeType: contentType,
			Size:     fileHeader.Size,
			Position: maxPos + 1 + i,
		}

		if err := mc.db.Create(&media).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusInternalServerError,
				ResponseMessage: "failed to save media record",
			})
		}

		results = append(results, media)
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    fiber.StatusCreated,
		ResponseMessage: "uploaded",
		Data:            results,
	})
}

func (mc *MediaController) DeleteMedia(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	mediaID := c.Params("mediaId")

	var media catalogModels.CatalogItemMedia
	if err := mc.db.First(&media, "id = ?", mediaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusNotFound,
			ResponseMessage: "media not found",
		})
	}

	// Verify item belongs to this CRM
	var item catalogModels.CatalogItem
	if err := mc.db.Where("id = ? AND crm_id = ?", media.ItemID, crmID).First(&item).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusForbidden,
			ResponseMessage: "forbidden",
		})
	}

	r2cfg := config.GetConfig().R2
	r2 := storage.NewR2Client(storage.R2Config{
		AccountID:       r2cfg.AccountID,
		AccessKeyID:     r2cfg.AccessKeyID,
		SecretAccessKey: r2cfg.SecretAccessKey,
		BucketName:      r2cfg.BucketName,
		PublicURL:       r2cfg.PublicURL,
	})
	_ = r2.Delete(c.Context(), media.Key)

	mc.db.Delete(&media)

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    fiber.StatusOK,
		ResponseMessage: "deleted",
	})
}

func (mc *MediaController) PresignUpload(c *fiber.Ctx) error {
	var req dto.PresignRequest
	if err := c.BodyParser(&req); err != nil || len(req.Files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "files array is required",
		})
	}
	if len(req.Files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "maximum 10 files per request",
		})
	}

	r2cfg := config.GetConfig().R2
	r2 := storage.NewR2Client(storage.R2Config{
		AccountID:       r2cfg.AccountID,
		AccessKeyID:     r2cfg.AccessKeyID,
		SecretAccessKey: r2cfg.SecretAccessKey,
		BucketName:      r2cfg.BucketName,
		PublicURL:       r2cfg.PublicURL,
	})

	var results []dto.PresignResponse
	for i, f := range req.Files {
		if !allowedPresignTypes[f.ContentType] {
			return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusBadRequest,
				ResponseMessage: fmt.Sprintf("unsupported type: %s", f.ContentType),
			})
		}

		ext := strings.ToLower(filepath.Ext(f.Name))
		key := fmt.Sprintf("catalog/%s/%d-%d%s", uuid.New().String(), time.Now().UnixMilli(), i, ext)

		uploadURL, err := r2.PresignPutURL(c.Context(), key, f.ContentType, 15*time.Minute)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusInternalServerError,
				ResponseMessage: "failed to generate upload URL",
			})
		}

		results = append(results, dto.PresignResponse{
			Key:       key,
			UploadURL: uploadURL,
			PublicURL: fmt.Sprintf("%s/%s", r2cfg.PublicURL, key),
		})
	}

	return c.JSON(response.CommonResponse{
		ResponseCode:    fiber.StatusOK,
		ResponseMessage: "ok",
		Data:            results,
	})
}

// UploadTemp accepts multipart files, uploads them to R2 server-side,
// and returns keys + public URLs without requiring an item_id.
// The caller includes the returned keys in the product creation payload.
func (mc *MediaController) UploadTemp(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "invalid multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "no files provided",
		})
	}
	if len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusBadRequest,
			ResponseMessage: "maximum 10 files per request",
		})
	}

	r2cfg := config.GetConfig().R2
	r2 := storage.NewR2Client(storage.R2Config{
		AccountID:       r2cfg.AccountID,
		AccessKeyID:     r2cfg.AccessKeyID,
		SecretAccessKey: r2cfg.SecretAccessKey,
		BucketName:      r2cfg.BucketName,
		PublicURL:       r2cfg.PublicURL,
	})

	var results []dto.PresignResponse
	for i, fileHeader := range files {
		if fileHeader.Size > maxFileSizeBytes {
			return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusBadRequest,
				ResponseMessage: fmt.Sprintf("file %s exceeds 20MB limit", fileHeader.Filename),
			})
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if !allowedMimeTypes[contentType] {
			return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusBadRequest,
				ResponseMessage: fmt.Sprintf("unsupported file type: %s", contentType),
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusInternalServerError,
				ResponseMessage: "failed to open file",
			})
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		key := fmt.Sprintf("catalog/%s/%d-%d%s", uuid.New().String(), time.Now().UnixMilli(), i, ext)

		publicURL, err := r2.Upload(c.Context(), key, file, contentType, fileHeader.Size)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusInternalServerError,
				ResponseMessage: "upload failed",
			})
		}

		results = append(results, dto.PresignResponse{
			Key:       key,
			UploadURL: "",
			PublicURL: publicURL,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    fiber.StatusCreated,
		ResponseMessage: "uploaded",
		Data:            results,
	})
}
