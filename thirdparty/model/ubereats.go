package model

type MenuItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Offers      struct {
		Price string `json:"price"`
	} `json:"offers"`
}

type Category struct {
	Name        string     `json:"name"`
	HasMenuItem []MenuItem `json:"hasMenuItem"`
}

type JsonData struct {
	HasMenu struct {
		HasMenuSection []Category `json:"hasMenuSection"`
	} `json:"hasMenu"`
}
