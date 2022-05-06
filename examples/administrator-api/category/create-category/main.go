package main

import (
	"fmt"
	"github.com/angeldm/mago/api"
	categories2 "github.com/angeldm/mago/categories"
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

	c := &categories2.Category{
		Name:     "spagetegory",
		Level:    2,
		IsActive: true,
	}

	mC, err := categories2.CreateCategory(c, apiClient)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n MagentoCategory struct: '%+v'", mC)
	fmt.Printf("\n MagentoCategory remote: %+v", mC.Category)
}
