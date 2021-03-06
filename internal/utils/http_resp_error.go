package utils

import (
	"fmt"
	"github.com/angeldm/mago"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var ErrBadRequest = fmt.Errorf("%s", "bad request")

func wrapError(err error, triedTo string, response ...map[string]interface{}) error {
	if len(response) == 0 {
		return fmt.Errorf("error while trying to %w - %s", err, triedTo)
	}
	return fmt.Errorf("error while trying to %w - %s. %+v", err, triedTo, response)
}

func MayReturnErrorForHTTPResponse(err error, resp *resty.Response, triedTo string) error {
	if err != nil {
		err = wrapError(err, triedTo)
	} else if resp.StatusCode() == http.StatusNotFound {
		err = mago.ErrNotFound
	} else if resp.StatusCode() >= http.StatusBadRequest {
		additional := map[string]interface{}{
			"statusCode": resp.StatusCode(),
			"response":   string(resp.Body()),
		}
		err = wrapError(ErrBadRequest, triedTo, additional)
	}

	return errors.WithStack(err)
}
