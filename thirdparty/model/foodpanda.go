package model

type ProductVariation struct {
	Price float64 `json:"price"`
}

type Image struct {
	ImageURL string `json:"image_url"`
}

type Product struct {
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Images            []Image            `json:"images"`
	ProductVariations []ProductVariation `json:"product_variations"`
}

type MenuCategory struct {
	Name     string    `json:"name"`
	Products []Product `json:"products"`
}

type Menu struct {
	MenuCategories []MenuCategory `json:"menu_categories"`
}

type Data struct {
	Menus []Menu `json:"menus"`
}

type Source struct {
	Data Data `json:"data"`
}
