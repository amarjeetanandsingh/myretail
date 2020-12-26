package product

type Product struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	CurrentPrice CurrentPrice `json:"current_price"`
}

type CurrentPrice struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currency_code"`
}
