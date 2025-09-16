package requests

type ProductRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Brand       string `json:"brand" binding:"required,min=2,max=50"`
	Model2      string `json:"model" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"required,min=5,max=255"`
	Stock       string `json:"stock" binding:"required"`
	Price       string `json:"price" binding:"required"`
	StatusID    string `json:"status_id" binding:"required"`
	CategoryID  string `json:"category_id" binding:"required"`
}
