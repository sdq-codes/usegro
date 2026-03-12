package authentication

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/apps/base/repositories"
	"github.com/sdq-codes/usegro-api/internal/apps/base/services/verification"
	notificationModels "github.com/sdq-codes/usegro-api/internal/apps/notifications/models"
	notification "github.com/sdq-codes/usegro-api/internal/apps/notifications/services"
	"github.com/sdq-codes/usegro-api/internal/helper/auth"
	"github.com/sdq-codes/usegro-api/internal/helper/encryption"
	"github.com/sdq-codes/usegro-api/internal/helper/random"
	"github.com/sdq-codes/usegro-api/internal/interface/resources/templates/emails"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"github.com/sdq-codes/usegro-api/pkg/amplitude"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

const (
	verifyCodeTTL       = 20 * time.Minute
	pwdResetTTL         = 20 * time.Minute
	pwdResetTokenPrefix = "auth:pwd_reset:token:" // hashedToken → userID
	pwdResetUserPrefix  = "auth:pwd_reset:user:"  // userID → hashedToken (for invalidation)
	pwdResetRatePrefix  = "auth:pwd_reset:rate:"  // email → request count
	pwdResetRateMax     = 3                       // max reset requests per window
	pwdResetRateWindow  = time.Hour
	loginCodeTTL        = 10 * time.Minute
	loginCodePrefix     = "auth:login_code:"
)

type Service struct {
	db                     *gorm.DB
	rdb                    *redis.Client
	userRepository         *repositories.UserRepository
	verificationRepository repositories.VerificationRepositoryInterface
	verificationService    *verification.Service
	jobRepository          repositories.JobRepositoryInterface
}

func NewAuthenticationService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		db:                     db,
		rdb:                    rdb,
		userRepository:         repositories.NewUserRepository(),
		verificationRepository: repositories.NewVerificationRepository(db),
		verificationService:    verification.NewVerificationService(db, rdb),
		jobRepository:          repositories.NewJobRepository(db),
	}
}

// RegisterUser creates a new account and returns the user plus a token pair.
func (s *Service) RegisterUser(ctx context.Context, userDTI dto.RegisterUserDTI) (*models.User, string, string, error) {
	tx := s.db.Begin()

	isUserEmailExist, err := s.userRepository.IsUserEmailExist(ctx, tx, userDTI.Email)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}
	if isUserEmailExist {
		tx.Rollback()
		return nil, "", "", exception.UserEmailAlreadyTakenError
	}

	hashedPassword, err := encryption.HashPassword(userDTI.Password)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	userModel := &models.User{Email: userDTI.Email, Password: hashedPassword}
	if err = s.userRepository.CreateUser(ctx, tx, userModel); err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	s.verificationService = verification.NewVerificationService(tx, s.rdb)
	if err = s.verificationService.CreateVerification(ctx, &models.Verification{
		UserID: userModel.ID,
		Type:   "EMAIL",
		Status: "PENDING",
	}); err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	// Generate code and store in Redis (replaces VerificationToken table).
	code := random.GenerateRandomCode(6)
	if err := s.verificationService.StoreEmailVerifyCode(ctx, userModel.ID.String(), code); err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	if err := notification.QueueEmail(ctx, tx, s.rdb, notificationModels.EmailNotification{
		FromEmail: emails.NO_REPLY_EMAIL,
		ToEmails:  []string{userModel.Email},
		Template:  emails.EMAIL_VERIFICATION_TEMPLATE,
		Data:      map[string]string{"token": code},
		Subject:   "Welcome to useGro",
	}, s.jobRepository); err != nil {
		logger.Log.Error(fmt.Sprintf("Verification email failed for %s: %v", userModel.Email, err))
	}

	minimalUser := models.User{ID: userModel.ID, Email: userModel.Email}

	accessToken, err := auth.CreateToken(minimalUser)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	refreshToken, err := auth.CreateRefreshToken(ctx, s.rdb, minimalUser)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	tx.Commit()
	amplitude.Track(userModel.ID.String(), amplitude.EventUserSignedUp, map[string]interface{}{
		amplitude.PropEmail:      userModel.Email,
		amplitude.PropAuthMethod: "email_password",
	})
	return userModel, accessToken, refreshToken, nil
}

// LoginUser authenticates a user and returns a token pair.
func (s *Service) LoginUser(ctx context.Context, loginDTO dto.RegisterUserDTI) (*models.User, string, string, error) {
	tx := s.db.Begin()

	userExist, err := s.userRepository.IsUserEmailExist(ctx, tx, loginDTO.Email)
	if err != nil || !userExist {
		tx.Rollback()
		return nil, "", "", exception.IncorrectUserNameAndPasswordRequestBodyError
	}

	user, err := s.userRepository.GetUser(ctx, tx, loginDTO.Email)
	if err != nil {
		tx.Rollback()
		return nil, "", "", exception.IncorrectUserNameAndPasswordRequestBodyError
	}

	if err := encryption.ComparePassword(user.Password, loginDTO.Password); err != nil {
		tx.Rollback()
		return nil, "", "", exception.IncorrectUserNameAndPasswordRequestBodyError
	}

	minimalUser := models.User{ID: user.ID, Email: user.Email}

	accessToken, err := auth.CreateToken(minimalUser)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	refreshToken, err := auth.CreateRefreshToken(ctx, s.rdb, minimalUser)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	tx.Commit()
	amplitude.Track(user.ID.String(), amplitude.EventUserLoggedIn, map[string]interface{}{
		amplitude.PropEmail:      user.Email,
		amplitude.PropAuthMethod: "email_password",
	})
	return user, accessToken, refreshToken, nil
}

// RefreshTokens validates and rotates the refresh token, returning a new token pair.
func (s *Service) RefreshTokens(ctx context.Context, rawRefreshToken string) (string, string, error) {
	_, accessToken, refreshToken, err := auth.ValidateAndRotateRefreshToken(ctx, s.rdb, rawRefreshToken)
	return accessToken, refreshToken, err
}

// Logout revokes the refresh token so it cannot be used again.
func (s *Service) Logout(ctx context.Context, rawRefreshToken string) error {
	userID := s.userIDFromRefreshToken(ctx, rawRefreshToken)
	err := auth.RevokeRefreshToken(ctx, s.rdb, rawRefreshToken)
	if err == nil && userID != "" {
		amplitude.Track(userID, amplitude.EventUserLoggedOut, nil)
	}
	return err
}

// userIDFromRefreshToken resolves the user ID stored in Redis for the given raw
// refresh token without consuming it — used solely for analytics.
func (s *Service) userIDFromRefreshToken(ctx context.Context, rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	key := "refresh:" + hex.EncodeToString(h[:])
	b, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return ""
	}
	var u models.User
	if json.Unmarshal(b, &u) != nil {
		return ""
	}
	return u.ID.String()
}

// SendPasswordResetEmail generates an opaque token, stores it in Redis, and emails a reset link.
func (s *Service) SendPasswordResetEmail(ctx context.Context, email string) error {
	// Rate limiting: max 3 requests per email per hour.
	rateKey := pwdResetRatePrefix + email
	count, err := s.rdb.Incr(ctx, rateKey).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		// First request in the window — set the expiry.
		s.rdb.Expire(ctx, rateKey, pwdResetRateWindow)
	}
	if count > pwdResetRateMax {
		return exception.TooManyRequestsError
	}

	// Read-only queries — no transaction needed.
	exists, err := s.userRepository.IsUserEmailExist(ctx, s.db, email)
	if err != nil {
		return err
	}
	if !exists {
		// Silently succeed — never reveal whether an email is registered.
		return nil
	}

	user, err := s.userRepository.GetUser(ctx, s.db, email)
	if err != nil {
		return err
	}

	// Generate a secure opaque token. Store sha256(rawToken) → userID in Redis.
	rawToken, hashedToken, err := random.GenerateSecureToken(32)
	if err != nil {
		return err
	}

	// Invalidate any existing reset token for this user before issuing a new one.
	if oldHash, err := s.rdb.GetDel(ctx, pwdResetUserPrefix+user.ID.String()).Result(); err == nil && oldHash != "" {
		s.rdb.Del(ctx, pwdResetTokenPrefix+oldHash)
	}

	// Store dual keys: token → userID (for lookup) and userID → token (for invalidation).
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, pwdResetTokenPrefix+hashedToken, user.ID.String(), pwdResetTTL)
	pipe.Set(ctx, pwdResetUserPrefix+user.ID.String(), hashedToken, pwdResetTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s/authentication/reset-password/%s", config.GetConfig().FrontEnd.Url, rawToken)

	if err := notification.QueueEmail(ctx, s.db, s.rdb, notificationModels.EmailNotification{
		FromEmail: emails.NO_REPLY_EMAIL,
		ToEmails:  []string{email},
		Template:  emails.PASSWORD_RESET_TEMPLATE,
		Data:      map[string]string{"resetLink": resetLink},
		Subject:   "Password Reset Request",
	}, s.jobRepository); err != nil {
		logger.Log.Error(fmt.Sprintf("Password reset email failed for %s: %v", email, err))
	}

	return nil
}

// ResetUserPassword validates the reset token from Redis and updates the password.
func (s *Service) ResetUserPassword(ctx context.Context, req dto.ResetPasswordDTI) error {
	// Hash the incoming raw token to look it up in Redis.
	h := sha256.Sum256([]byte(req.Token))
	hashedToken := hex.EncodeToString(h[:])

	// Validate and consume the token atomically (single-use).
	userID, err := s.rdb.GetDel(ctx, pwdResetTokenPrefix+hashedToken).Result()
	if err != nil || userID == "" {
		return exception.ResetPasswordTokenExpiredError
	}

	// Clean up the user-index key so no other token can be issued against a stale index.
	s.rdb.Del(ctx, pwdResetUserPrefix+userID)

	tx := s.db.Begin()

	user, err := s.userRepository.GetUserById(ctx, tx, userID)
	if err != nil {
		tx.Rollback()
		return exception.UserEmailNotFoundError
	}

	hashedPassword, err := encryption.HashPassword(req.Password)
	if err != nil {
		tx.Rollback()
		return err
	}
	user.Password = hashedPassword

	if err := s.userRepository.UpdateUser(ctx, tx, user); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	amplitude.Track(userID, amplitude.EventPasswordResetCompleted, nil)
	return nil
}

func (s *Service) IsUserEmailExist(ctx context.Context, userDTI dto.UserExistDTI) (bool, error) {
	exists, err := s.userRepository.IsUserEmailExist(ctx, s.db, userDTI.Email)
	if err != nil {
		return false, err
	}
	if exists {
		return true, exception.UserEmailAlreadyTakenError
	}
	return false, nil
}

// ── Google OAuth ──────────────────────────────────────────────────────────────

func googleOAuthConfig() *oauth2.Config {
	cfg := config.GetConfig().Google
	return &oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin exchanges an OAuth code for a token pair, creating the user on first login.
func (s *Service) GoogleLogin(ctx context.Context, code string) (*models.User, string, string, error) {
	oauthToken, err := googleOAuthConfig().Exchange(ctx, code)
	if err != nil {
		return nil, "", "", exception.UnauthorizedError
	}

	client := googleOAuthConfig().Client(ctx, oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", "", err
	}
	defer resp.Body.Close()

	var googleUser dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, "", "", err
	}
	if !googleUser.VerifiedEmail {
		return nil, "", "", exception.UnauthorizedError
	}

	tx := s.db.Begin()

	userExists, err := s.userRepository.IsUserEmailExist(ctx, tx, googleUser.Email)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	var user *models.User

	if userExists {
		user, err = s.userRepository.GetUser(ctx, tx, googleUser.Email)
		if err != nil {
			tx.Rollback()
			return nil, "", "", err
		}
	} else {
		b := make([]byte, 32)
		if _, err := cryptoRand.Read(b); err != nil {
			tx.Rollback()
			return nil, "", "", err
		}
		hashedPassword, err := encryption.HashPassword(base64.StdEncoding.EncodeToString(b))
		if err != nil {
			tx.Rollback()
			return nil, "", "", err
		}

		user = &models.User{Email: googleUser.Email, Password: hashedPassword}
		if err = s.userRepository.CreateUser(ctx, tx, user); err != nil {
			tx.Rollback()
			return nil, "", "", err
		}

		now := time.Now()
		s.verificationService = verification.NewVerificationService(tx, s.rdb)
		if err := s.verificationService.CreateVerification(ctx, &models.Verification{
			UserID:     user.ID,
			Type:       "EMAIL",
			Status:     "VERIFIED",
			VerifiedAt: &now,
		}); err != nil {
			tx.Rollback()
			return nil, "", "", err
		}
	}

	minimalUser := models.User{ID: user.ID, Email: user.Email}

	accessToken, err := auth.CreateToken(minimalUser)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	refreshToken, err := auth.CreateRefreshToken(ctx, s.rdb, minimalUser)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}

	tx.Commit()
	if userExists {
		amplitude.Track(user.ID.String(), amplitude.EventUserLoggedIn, map[string]interface{}{
			amplitude.PropEmail:      googleUser.Email,
			amplitude.PropAuthMethod: "google",
		})
	} else {
		amplitude.Track(user.ID.String(), amplitude.EventUserSignedUp, map[string]interface{}{
			amplitude.PropEmail:      googleUser.Email,
			amplitude.PropAuthMethod: "google",
		})
	}
	return user, accessToken, refreshToken, nil
}

// GoogleAuthURL returns the Google OAuth consent screen URL for the given state token.
func (s *Service) GoogleAuthURL(state string) string {
	return googleOAuthConfig().AuthCodeURL(state)
}

// ── Email code login ──────────────────────────────────────────────────────────

// RequestLoginCode generates a 6-character code, stores it in Redis, and emails it.
// Always returns nil to avoid leaking whether an email is registered.
func (s *Service) RequestLoginCode(ctx context.Context, req dto.RequestEmailCodeDTI) error {
	tx := s.db.Begin()
	exists, err := s.userRepository.IsUserEmailExist(ctx, tx, req.Email)
	tx.Commit()
	if err != nil || !exists {
		return nil
	}

	code := random.GenerateRandomCode(6)
	if err := s.rdb.Set(ctx, loginCodePrefix+req.Email, code, loginCodeTTL).Err(); err != nil {
		return err
	}

	tx = s.db.Begin()
	if err := notification.QueueEmail(ctx, tx, s.rdb, notificationModels.EmailNotification{
		FromEmail: emails.NO_REPLY_EMAIL,
		ToEmails:  []string{req.Email},
		Template:  emails.LOGIN_CODE_TEMPLATE,
		Data:      map[string]string{"code": code},
		Subject:   "Your useGro login code",
	}, s.jobRepository); err != nil {
		logger.Log.Error(fmt.Sprintf("Login code email failed for %s: %v", req.Email, err))
	}
	tx.Commit()

	return nil
}

// VerifyLoginCode validates the code and returns a token pair on success.
func (s *Service) VerifyLoginCode(ctx context.Context, req dto.VerifyEmailCodeDTI) (*models.User, string, string, error) {
	stored, err := s.rdb.GetDel(ctx, loginCodePrefix+req.Email).Result()
	if err != nil || stored == "" || stored != req.Code {
		return nil, "", "", exception.InvalidLoginCodeError
	}

	tx := s.db.Begin()
	user, err := s.userRepository.GetUser(ctx, tx, req.Email)
	tx.Commit()
	if err != nil {
		return nil, "", "", exception.UserEmailNotFoundError
	}

	minimalUser := models.User{ID: user.ID, Email: user.Email}

	accessToken, err := auth.CreateToken(minimalUser)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := auth.CreateRefreshToken(ctx, s.rdb, minimalUser)
	if err != nil {
		return nil, "", "", err
	}

	amplitude.Track(user.ID.String(), amplitude.EventUserLoggedIn, map[string]interface{}{
		amplitude.PropEmail:      req.Email,
		amplitude.PropAuthMethod: "email_code",
	})
	return user, accessToken, refreshToken, nil
}
