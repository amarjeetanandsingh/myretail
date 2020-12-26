package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/amarjeetanandsingh/myRetail/arango"
	"github.com/go-resty/resty/v2"
)

const (
	ProductCollection = "Product"
)

var (
	ErrNotFound = errors.New("product not found")
)

type Repository interface {
	FindProductNameByID(productID int) (string, error)
	FindProductPriceByID(productID int) (CurrentPrice, error)
	UpdateProductPrice(key int, fields interface{}) error
}

type repo struct {
	Database  arango.DbServer
	NameStore *resty.Client
}

func (r repo) FindProductNameByID(productID int) (string, error) {
	req := r.NameStore.R().
		SetQueryParam("fields", "descriptions").
		SetQueryParam("id_type", "TCIN")

	resp, err := req.Get("/" + strconv.Itoa(productID))
	if err != nil {
		return "", fmt.Errorf("error getting name for product id: %v :: %w", productID, err)
	}
	fmt.Printf("api.target request: %+v\n", resp.Request)

	prodName := struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(resp.Body(), &prodName); err != nil {
		return "", fmt.Errorf("error unmarshling name for product id: %v :: %w", productID, err)
	}
	return prodName.Name, nil
}

func (r *repo) FindProductPriceByID(productID int) (CurrentPrice, error) {
	prodCurrPrice := CurrentPrice{}
	productIDStr := strconv.Itoa(productID)

	if err := r.Database.FindByID(ProductCollection, productIDStr, &prodCurrPrice); err != nil {
		return prodCurrPrice, fmt.Errorf("error fetching price for product id: %v: %w", productID, err)
	}
	if prodCurrPrice.CurrencyCode == "" {
		return prodCurrPrice, fmt.Errorf("curr price not found for product id: %v: %w", productID, ErrNotFound)
	}
	return prodCurrPrice, nil
}

func (r *repo) UpdateProductPrice(key int, fields interface{}) error {
	productID := strconv.Itoa(key)
	if err := r.Database.UpdateDoc(ProductCollection, productID, fields); err != nil {
		return fmt.Errorf("error updating product id: %v, fields: %v: %w", key, fields, err)
	}
	return nil
}

// TODO fetch from config
//https://api.target.com/products/v3/13860428?fields=descriptions&id_type=TCIN
func NewRepo(dbServer arango.DbServer, nameStore *resty.Client) Repository {
	return &repo{
		Database:  dbServer,
		NameStore: nameStore,
	}
}
