package main

import (
	"log"
	"time"

	"github.com/amarjeetanandsingh/myRetail/arango"
	"github.com/amarjeetanandsingh/myRetail/config"
	"github.com/amarjeetanandsingh/myRetail/handler"
	"github.com/amarjeetanandsingh/myRetail/product"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// create arango client
	arangoClient := NewArangoClient()
	nameStoreClient := NewNameStoreClient()
	productRepo := product.NewRepo(arangoClient, nameStoreClient)
	productService := product.NewService(productRepo)
	h := handler.New(productService)
	h.InstallRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func NewNameStoreClient() *resty.Client {
	nameStoreClient := resty.New()
	nameStoreClient.SetTimeout(time.Second * 5)
	nameStoreClient.SetHostURL("https://api.target.com/products/v3")
	nameStoreClient.SetHeader("key", "43cJWpLjH8Z8oR18KdrZDBKAgLLQKJjz")
	return nameStoreClient
}

func NewArangoClient() arango.DbServer {
	arangoHost := config.GetEnv("ARANGO_HOST", "localhost")
	arangoPort := config.GetEnv("ARANGO_PORT", "8529")
	arangoUsername := config.GetEnv("ARANGO_USERNAME", "root")
	arangoPassword := config.GetEnv("ARANGO_PASSWORD", "root")
	dbServer, err := arango.New(arangoHost, arangoPort, arangoUsername, arangoPassword)
	if err != nil {
		log.Fatal("Arango init failed: ", err)
	}
	return dbServer
}
