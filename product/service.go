package product

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type ProductServer interface {
	FindProductByID(ID int) (*Product, error)
	UpdateProductPrice(ID int, prodPrice CurrentPrice) error
}

type service struct {
	repo Repository
}

func (service *service) FindProductByID(productID int) (*Product, error) {
	prod := &Product{}
	errGrp, _ := errgroup.WithContext(context.Background())

	// 1) find curr price from db
	errGrp.Go(func() error {
		currentPrice, err := service.repo.FindProductPriceByID(productID)
		if err != nil {
			return err
		}
		prod.CurrentPrice = currentPrice
		return nil
	})

	// 2) find product name from http client
	errGrp.Go(func() error {
		prodName, err := service.repo.FindProductNameByID(productID)
		if err != nil {
			return err
		}
		prod.Name = prodName
		return nil
	})

	if err := errGrp.Wait(); err != nil {
		return nil, err
	}

	return prod, nil
}

func (service *service) UpdateProductPrice(productID int, prodPrice CurrentPrice) error {
	if productID == 0 {
		return fmt.Errorf("invalid productID")
	}
	if prodPrice.CurrencyCode == "" {
		return fmt.Errorf("currencyCode must not be empty")
	}

	if err := service.repo.UpdateProductPrice(productID, prodPrice); err != nil {
		return err
	}
	return nil
}

func NewService(repo Repository) ProductServer {
	return &service{repo: repo}
}
