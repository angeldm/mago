package main

import (
	"fmt"
	"github.com/angeldm/mago/api"
	configurableproducts2 "github.com/angeldm/mago/configurableproducts"
	attribute2 "github.com/angeldm/mago/products/attribute"
	"log"
	"strconv"
)

// TODO: FINISH EXAMPLE

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

	// define your attribute
	attr := &attribute2.Attribute{
		AttributeCode:        "myselectattribute",
		FrontendInput:        "select",
		DefaultFrontendLabel: "aw",
		IsRequired:           false,
	}

	// create attribute on remote
	mAttribute, err := attribute2.CreateAttribute(attr, apiClient)
	if err != nil {
		panic(err)
	}

	optionValue, err := mAttribute.AddOption(attribute2.Option{
		Label: "spaget",
		Value: "spaget",
	})
	if err != nil {
		panic(err)
	}

	optionValueInt, err := strconv.Atoi(optionValue)
	if err != nil {
		panic(err)
	}

	option := &configurableproducts2.Option{
		AttributeID:  fmt.Sprintf("%d", mAttribute.Attribute.AttributeID),
		Label:        mAttribute.Attribute.DefaultFrontendLabel,
		Position:     0,
		IsUseDefault: false,
		Values: []configurableproducts2.Value{
			{
				ValueIndex: optionValueInt,
			},
		},
	}

	mOption, err := configurableproducts2.SetOptionForExistingConfigurableProduct("configurableSpaget", option, apiClient)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created configurable Product: '%+v'", mOption)
	fmt.Printf("Created configurable Product options: '%+v'", mOption.Options)
}
