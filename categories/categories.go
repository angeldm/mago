package categories

import (
	"errors"
	"fmt"
	"github.com/angeldm/mago"
	"github.com/angeldm/mago/api"

	"github.com/angeldm/mago/internal/utils"
)

type MCategory struct {
	Route     string
	Category  *Category
	Products  *[]ProductLink
	APIClient *api.Client
}

func CreateCategory(c *Category, apiClient *api.Client) (*MCategory, error) {
	mC := &MCategory{
		Category:  &Category{},
		Products:  &[]ProductLink{},
		APIClient: apiClient,
	}
	endpoint := categories
	httpClient := apiClient.HTTPClient

	payLoad := createCategoryPayload{
		Category: *c,
	}

	resp, err := httpClient.R().SetBody(payLoad).SetResult(mC.Category).Post(endpoint)
	mC.Route = fmt.Sprintf("%s/%d", categories, mC.Category.ID)

	err = utils.MayReturnErrorForHTTPResponse(err, resp, "create category")
	return mC, err
}

func GetCategoryByName(name string, apiClient *api.Client) (*MCategory, error) {
	mC := &MCategory{
		Category:  &Category{},
		Products:  &[]ProductLink{},
		APIClient: apiClient,
	}
	searchQuery := utils.BuildSearchQuery("name", name, "in")
	endpoint := categoriesList + "?" + searchQuery
	httpClient := apiClient.HTTPClient

	response := &categorySearchQueryResponse{}

	resp, err := httpClient.R().SetResult(response).Get(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "get category by name from remote")
	if err != nil {
		return nil, err
	}

	if len(response.Categories) == 0 {
		return nil, mago.ErrNotFound
	}

	mC.Category = &response.Categories[0]
	mC.Route = fmt.Sprintf("%s/%d", categories, mC.Category.ID)

	err = utils.MayReturnErrorForHTTPResponse(mC.UpdateCategoryFromRemote(), resp, "get detailed category by name from remote")

	return mC, err
}

func GetCategories(rootCategoy string, apiClient *api.Client) ([]Category, error) {
	endpoint := categories + "?" + "rootCategoryId=" + rootCategoy
	httpClient := apiClient.HTTPClient

	response := &categorySearchResponse{}

	resp, err := httpClient.R().SetResult(response).Get(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "get category by name from remote")
	if err != nil {
		return nil, err
	}

	// if len(response.Categories) == 0 {
	// 	return nil, errors.New("not found")
	// }

	return response.Categories, err
}

func RemoveCategoriesByParent(parent string, apiClient *api.Client) error {
	searchQuery := utils.BuildSearchQuery("parent_id", parent, "in")
	endpoint := categoriesList + "?" + searchQuery
	httpClient := apiClient.HTTPClient

	response := &categorySearchQueryResponse{}

	resp, err := httpClient.R().SetResult(response).Get(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "get category by name from remote")
	if err != nil {
		return err
	}

	if len(response.Categories) == 0 {
		return errors.New("not found")
	}

	endpoint = categories
	for _, v := range response.Categories {
		imC := &MCategory{
			Category:  &v,
			Products:  &[]ProductLink{},
			APIClient: apiClient,
		}
		imC.Route = fmt.Sprintf("%s/%d", categories, imC.Category.ID)

		err := imC.RemoveCategoryFromRemote()
		if err != nil {
			return err
		}
	}

	return nil
}

func (mC *MCategory) UpdateCategoryFromRemote() error {
	resp, err := mC.APIClient.HTTPClient.R().SetResult(mC.Category).Get(mC.Route)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "get category from remote")
	if err != nil {
		return err
	}

	return mC.UpdateCategoryProductsFromRemote()
}

func (mC *MCategory) UpdateCategoryProductsFromRemote() error {
	resp, err := mC.APIClient.HTTPClient.R().SetResult(mC.Products).Get(fmt.Sprintf("%s/%s", mC.Route, categoriesProductsRelative))
	return utils.MayReturnErrorForHTTPResponse(err, resp, "get category products from remote")
}

func (mC *MCategory) AssignProductByProductLink(pl *ProductLink) error {
	if pl.CategoryID == "" {
		pl.CategoryID = fmt.Sprintf("%d", mC.Category.ID)
	}

	httpClient := mC.APIClient.HTTPClient
	endpoint := fmt.Sprintf("%s/%s", mC.Route, categoriesProductsRelative)

	payLoad := assignProductPayload{ProductLink: *pl}

	resp, err := httpClient.R().SetBody(payLoad).Put(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "assign product to category")

	if err == nil {
		*mC.Products = append(*mC.Products, *pl)
	}

	return err
}
func (mC *MCategory) RemoveCategoryFromRemote() error {
	resp, err := mC.APIClient.HTTPClient.R().Delete(mC.Route)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "get category from remote")
	if err != nil {
		return err
	}
	return nil
}
