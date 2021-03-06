package main

import (
	"github.com/angeldm/mago/api"
	"github.com/angeldm/mago/products"
	"log"
)

func main() {
	// initiate storeconfig
	storeConfig := &api.StoreConfig{
		Scheme:    "https",
		HostName:  "mago.hermsi.localhost",
		StoreCode: "default",
	}
	// initiate bearer payload
	bearerToken := "yd1o9zs1hb1qxnn8ek68eu8nwqjg5hrv"

	// create a new apiClient
	apiClient, err := api.NewAPIClientFromIntegration(storeConfig, bearerToken)
	if err != nil {
		panic(err)
	}
	log.Printf("Obtained client: '%v'", apiClient)

	// get product by SKU
	mProduct, err := products.GetProductBySKU("spaget1234", apiClient)
	if err != nil {
		panic(err)
	}

	// here you go
	log.Printf("Got product: %+v", mProduct)
}
