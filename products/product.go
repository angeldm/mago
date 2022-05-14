package products

import (
	"github.com/angeldm/mago/api"
	"github.com/angeldm/mago/internal/utils"
)

const (
	products = "/products"
)

type MProduct struct {
	Route     string
	Product   *Product
	APIClient *api.Client
}

func CreateOrReplaceProduct(product *Product, saveOptions bool, apiClient *api.Client) (*MProduct, error) {
	mp := &MProduct{
		Product:   product,
		APIClient: apiClient,
	}

	err := mp.createOrReplaceProduct(saveOptions)

	return mp, err
}

func GetProductBySKU(sku string, apiClient *api.Client) (*MProduct, error) {
	mProduct := &MProduct{
		Route:     products + "/" + sku,
		Product:   &Product{},
		APIClient: apiClient,
	}

	err := mProduct.UpdateProductFromRemote()

	return mProduct, err
}

func GetAllProducts(apiClient *api.Client) ([]*MProduct, error) { // searchQuery := BuildSearchQuery("currentPage", "0", "in")
	endpoint := products + "?" + "searchCriteria[currentPage]=0"
	httpClient := apiClient.HTTPClient

	response := &ProductSearchResponse{}

	resp, err := httpClient.R().SetResult(response).Get(endpoint)
	err = utils.MayReturnErrorForHTTPResponse(err, resp, "get product by name from remote")
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, nil
	}

	items := make([]*MProduct, len(response.Items))
	for k, tr := range response.Items {
		mProduct := &MProduct{
			Product:   &Product{},
			APIClient: apiClient,
			Route:     products + "/" + tr.Sku,
		}
		err := utils.MayReturnErrorForHTTPResponse(mProduct.UpdateProductFromRemote(), resp, "get detailed product by name from remote")
		if err != nil {
			return nil, err
		}
		items[k] = mProduct
	}

	return items, err
}

func DeleteProductBySKU(sku string, apiClient *api.Client) (*MProduct, error) {
	mProduct := &MProduct{
		Route:     products + "/" + sku,
		Product:   &Product{},
		APIClient: apiClient,
	}

	err := mProduct.DeleteFromRemote()

	return mProduct, err
}

func (mProduct *MProduct) DeleteFromRemote() error {
	httpClient := mProduct.APIClient.HTTPClient

	resp, err := httpClient.R().Delete(mProduct.Route)
	return utils.MayReturnErrorForHTTPResponse(err, resp, "get detailed product from remote")
}

func (mProduct *MProduct) createOrReplaceProduct(saveOptions bool) error {
	endpoint := products
	httpClient := mProduct.APIClient.HTTPClient

	payLoad := AddProductPayload{
		Product:     *mProduct.Product,
		SaveOptions: saveOptions,
	}

	resp, err := httpClient.R().SetBody(payLoad).SetResult(mProduct.Product).Post(endpoint)
	productSKU := utils.MayTrimSurroundingQuotes(mProduct.Product.Sku)
	mProduct.Route = products + "/" + productSKU

	return utils.MayReturnErrorForHTTPResponse(err, resp, "create new product on remote")
}

func (mProduct *MProduct) UpdateProductFromRemote() error {
	httpClient := mProduct.APIClient.HTTPClient

	resp, err := httpClient.R().SetResult(mProduct.Product).Get(mProduct.Route)
	return utils.MayReturnErrorForHTTPResponse(err, resp, "get detailed product from remote")
}

func (mProduct *MProduct) UpdateQuantityForStockItem(stockItem string, quantity int, isInStock bool) error {
	httpClient := mProduct.APIClient.HTTPClient

	updateStockPayload := updateStockPayload{StockItem: StockItem{Qty: quantity, IsInStock: isInStock}}

	resp, err := httpClient.R().SetBody(updateStockPayload).Put(mProduct.Route + "/" + stockItemsRelative + "/" + stockItem)
	return utils.MayReturnErrorForHTTPResponse(err, resp, "update stock for product")
}
