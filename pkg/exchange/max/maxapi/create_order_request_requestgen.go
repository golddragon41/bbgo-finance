// Code generated by "requestgen -method POST -url /api/v2/orders -type CreateOrderRequest -responseType .Order"; DO NOT EDIT.

package max

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
)

func (c *CreateOrderRequest) Market(market string) *CreateOrderRequest {
	c.market = market
	return c
}

func (c *CreateOrderRequest) Side(side string) *CreateOrderRequest {
	c.side = side
	return c
}

func (c *CreateOrderRequest) Volume(volume string) *CreateOrderRequest {
	c.volume = volume
	return c
}

func (c *CreateOrderRequest) OrderType(orderType OrderType) *CreateOrderRequest {
	c.orderType = orderType
	return c
}

func (c *CreateOrderRequest) Price(price string) *CreateOrderRequest {
	c.price = &price
	return c
}

func (c *CreateOrderRequest) StopPrice(stopPrice string) *CreateOrderRequest {
	c.stopPrice = &stopPrice
	return c
}

func (c *CreateOrderRequest) ClientOrderID(clientOrderID string) *CreateOrderRequest {
	c.clientOrderID = &clientOrderID
	return c
}

func (c *CreateOrderRequest) GroupID(groupID string) *CreateOrderRequest {
	c.groupID = &groupID
	return c
}

// GetQueryParameters builds and checks the query parameters and returns url.Values
func (c *CreateOrderRequest) GetQueryParameters() (url.Values, error) {
	var params = map[string]interface{}{}

	query := url.Values{}
	for _k, _v := range params {
		query.Add(_k, fmt.Sprintf("%v", _v))
	}

	return query, nil
}

// GetParameters builds and checks the parameters and return the result in a map object
func (c *CreateOrderRequest) GetParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}
	// check market field -> json key market
	market := c.market

	// TEMPLATE check-required
	if len(market) == 0 {
		return nil, fmt.Errorf("market is required, empty string given")
	}
	// END TEMPLATE check-required

	// assign parameter of market
	params["market"] = market
	// check side field -> json key side
	side := c.side

	// TEMPLATE check-required
	if len(side) == 0 {
		return nil, fmt.Errorf("side is required, empty string given")
	}
	// END TEMPLATE check-required

	// assign parameter of side
	params["side"] = side
	// check volume field -> json key volume
	volume := c.volume

	// TEMPLATE check-required
	if len(volume) == 0 {
		return nil, fmt.Errorf("volume is required, empty string given")
	}
	// END TEMPLATE check-required

	// assign parameter of volume
	params["volume"] = volume
	// check orderType field -> json key ord_type
	orderType := c.orderType

	// assign parameter of orderType
	params["ord_type"] = orderType
	// check price field -> json key price
	if c.price != nil {
		price := *c.price

		// assign parameter of price
		params["price"] = price
	} else {
	}
	// check stopPrice field -> json key stop_price
	if c.stopPrice != nil {
		stopPrice := *c.stopPrice

		// assign parameter of stopPrice
		params["stop_price"] = stopPrice
	} else {
	}
	// check clientOrderID field -> json key client_oid
	if c.clientOrderID != nil {
		clientOrderID := *c.clientOrderID

		// assign parameter of clientOrderID
		params["client_oid"] = clientOrderID
	} else {
	}
	// check groupID field -> json key group_id
	if c.groupID != nil {
		groupID := *c.groupID

		// assign parameter of groupID
		params["group_id"] = groupID
	} else {
	}

	return params, nil
}

// GetParametersQuery converts the parameters from GetParameters into the url.Values format
func (c *CreateOrderRequest) GetParametersQuery() (url.Values, error) {
	query := url.Values{}

	params, err := c.GetParameters()
	if err != nil {
		return query, err
	}

	for _k, _v := range params {
		if c.isVarSlice(_v) {
			c.iterateSlice(_v, func(it interface{}) {
				query.Add(_k+"[]", fmt.Sprintf("%v", it))
			})
		} else {
			query.Add(_k, fmt.Sprintf("%v", _v))
		}
	}

	return query, nil
}

// GetParametersJSON converts the parameters from GetParameters into the JSON format
func (c *CreateOrderRequest) GetParametersJSON() ([]byte, error) {
	params, err := c.GetParameters()
	if err != nil {
		return nil, err
	}

	return json.Marshal(params)
}

// GetSlugParameters builds and checks the slug parameters and return the result in a map object
func (c *CreateOrderRequest) GetSlugParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}

	return params, nil
}

func (c *CreateOrderRequest) applySlugsToUrl(url string, slugs map[string]string) string {
	for _k, _v := range slugs {
		needleRE := regexp.MustCompile(":" + _k + "\\b")
		url = needleRE.ReplaceAllString(url, _v)
	}

	return url
}

func (c *CreateOrderRequest) iterateSlice(slice interface{}, _f func(it interface{})) {
	sliceValue := reflect.ValueOf(slice)
	for _i := 0; _i < sliceValue.Len(); _i++ {
		it := sliceValue.Index(_i).Interface()
		_f(it)
	}
}

func (c *CreateOrderRequest) isVarSlice(_v interface{}) bool {
	rt := reflect.TypeOf(_v)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	}
	return false
}

func (c *CreateOrderRequest) GetSlugsMap() (map[string]string, error) {
	slugs := map[string]string{}
	params, err := c.GetSlugParameters()
	if err != nil {
		return slugs, nil
	}

	for _k, _v := range params {
		slugs[_k] = fmt.Sprintf("%v", _v)
	}

	return slugs, nil
}

func (c *CreateOrderRequest) Do(ctx context.Context) (*Order, error) {

	params, err := c.GetParameters()
	if err != nil {
		return nil, err
	}
	query := url.Values{}

	apiURL := "/api/v2/orders"

	req, err := c.client.NewAuthenticatedRequest(ctx, "POST", apiURL, query, params)
	if err != nil {
		return nil, err
	}

	response, err := c.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse Order
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
