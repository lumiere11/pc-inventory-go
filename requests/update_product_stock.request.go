package requests

type UpdateProductRequest struct {
	Stock int `json:"stock" binding:"required,min=1"`
}
