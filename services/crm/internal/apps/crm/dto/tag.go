package dto

type TagCreateDTO struct {
	Tag string `json:"tag" validate:"required,min=2,max=50"`
}

type TagUpdateDTO struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}
