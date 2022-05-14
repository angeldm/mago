package attribute

import (
	"fmt"
	"github.com/angeldm/mago/api"
	"strings"

	"github.com/angeldm/mago/internal/utils"
)

type MAttribute struct {
	Route     string
	Attribute *Attribute
	APIClient *api.Client
}

func CreateAttribute(a *Attribute, apiClient *api.Client) (*MAttribute, error) {
	mAttribute := &MAttribute{
		Attribute: &Attribute{},
		APIClient: apiClient,
	}
	endpoint := productsAttribute
	httpClient := apiClient.HTTPClient

	payLoad := createAttributePayload{
		Attribute: *a,
	}

	resp, err := httpClient.R().SetBody(payLoad).SetResult(mAttribute.Attribute).Post(endpoint)
	mAttribute.Route = productsAttribute + "/" + mAttribute.Attribute.AttributeCode

	return mAttribute, utils.MayReturnErrorForHTTPResponse(err, resp, "create attribute")
}

func GetAttributeByAttributeCode(attributeCode string, apiClient *api.Client) (*MAttribute, error) {
	mAttributeSet := &MAttribute{
		Route:     fmt.Sprintf("%s/%s", productsAttribute, attributeCode),
		Attribute: &Attribute{},
		APIClient: apiClient,
	}

	err := mAttributeSet.UpdateAttributeFromRemote()

	return mAttributeSet, err
}

func (mas *MAttribute) UpdateAttributeOnRemote() error {
	resp, err := mas.APIClient.HTTPClient.R().SetResult(mas.Attribute).SetBody(mas.Attribute).Put(mas.Route)
	return utils.MayReturnErrorForHTTPResponse(err, resp, "update remote attribute from local")
}

func (mas *MAttribute) UpdateAttributeFromRemote() error {
	resp, err := mas.APIClient.HTTPClient.R().SetResult(mas.Attribute).Get(mas.Route)
	return utils.MayReturnErrorForHTTPResponse(err, resp, "update local attribute from remote")
}

func (mas *MAttribute) AddOption(option Option) (string, error) {
	endpoint := mas.Route + "/" + productsAttributeOptions
	httpClient := mas.APIClient.HTTPClient

	payLoad := addOptionPayload{
		Option: option,
	}

	resp, err := httpClient.R().SetBody(payLoad).Post(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "assign option to attribute")
	if err != nil {
		return "", err
	}

	optionValue := utils.MayTrimSurroundingQuotes(resp.String())
	optionValue = strings.TrimPrefix(optionValue, "id_")

	return optionValue, mas.UpdateAttributeFromRemote()
}

func (mas *MAttribute) DeleteOption(option Option) error {
	endpoint := mas.Route + "/" + productsAttributeOptions + "/" + option.Value
	httpClient := mas.APIClient.HTTPClient

	payLoad := addOptionPayload{
		Option: option,
	}

	resp, err := httpClient.R().SetBody(payLoad).Delete(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "assign option to attribute")
	if err != nil {
		return err
	}

	return err
}
