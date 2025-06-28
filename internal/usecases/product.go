package usecases

type Product struct {
	ProductId    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	ProductImage string  `json:"product_image"`
}
