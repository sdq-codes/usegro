package dto

type CreateCategoryDTO struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateCategoryDTO struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}
