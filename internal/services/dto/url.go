package dto

type URL struct {
	URL string `json:"url" validate:"required,url"`
}
