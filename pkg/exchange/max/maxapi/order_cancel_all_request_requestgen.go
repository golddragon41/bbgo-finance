// Code generated by "requestgen -method POST -url v2/orders/clear -type OrderCancelAllRequest -responseType []Order"; DO NOT EDIT.

package max

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
)

func (o *OrderCancelAllRequest) Side(side string) *OrderCancelAllRequest {
	o.side = &side
	return o
}

func (o *OrderCancelAllRequest) Market(market string) *OrderCancelAllRequest {
	o.market = &market
	return o
}

func (o *OrderCancelAllRequest) GroupID(groupID uint32) *OrderCancelAllRequest {
	o.groupID = &groupID
	return o
}

// GetQueryParameters builds and checks the query parameters and returns url.Values
func (o *OrderCancelAllRequest) GetQueryParameters() (url.Values, error) {
	var params = map[string]interface{}{}

	query := url.Values{}
	for k, v := range params {
		query.Add(k, fmt.Sprintf("%v", v))
	}

	return query, nil
}

// GetParameters builds and checks the parameters and return the result in a map object
func (o *OrderCancelAllRequest) GetParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}
	// check side field -> json key side
	if o.side != nil {
		side := *o.side

		// assign parameter of side
		params["side"] = side
	} else {
	}
	// check market field -> json key market
	if o.market != nil {
		market := *o.market

		// assign parameter of market
		params["market"] = market
	} else {
	}
	// check groupID field -> json key groupID
	if o.groupID != nil {
		groupID := *o.groupID

		// assign parameter of groupID
		params["groupID"] = groupID
	} else {
	}

	return params, nil
}

// GetParametersQuery converts the parameters from GetParameters into the url.Values format
func (o *OrderCancelAllRequest) GetParametersQuery() (url.Values, error) {
	query := url.Values{}

	params, err := o.GetParameters()
	if err != nil {
		return query, err
	}

	for k, v := range params {
		if o.isVarSlice(v) {
			o.iterateSlice(v, func(it interface{}) {
				query.Add(k+"[]", fmt.Sprintf("%v", it))
			})
		} else {
			query.Add(k, fmt.Sprintf("%v", v))
		}
	}

	return query, nil
}

// GetParametersJSON converts the parameters from GetParameters into the JSON format
func (o *OrderCancelAllRequest) GetParametersJSON() ([]byte, error) {
	params, err := o.GetParameters()
	if err != nil {
		return nil, err
	}

	return json.Marshal(params)
}

// GetSlugParameters builds and checks the slug parameters and return the result in a map object
func (o *OrderCancelAllRequest) GetSlugParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}

	return params, nil
}

func (o *OrderCancelAllRequest) applySlugsToUrl(url string, slugs map[string]string) string {
	for k, v := range slugs {
		needleRE := regexp.MustCompile(":" + k + "\\b")
		url = needleRE.ReplaceAllString(url, v)
	}

	return url
}

func (o *OrderCancelAllRequest) iterateSlice(slice interface{}, f func(it interface{})) {
	sliceValue := reflect.ValueOf(slice)
	for i := 0; i < sliceValue.Len(); i++ {
		it := sliceValue.Index(i).Interface()
		f(it)
	}
}

func (o *OrderCancelAllRequest) isVarSlice(v interface{}) bool {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	}
	return false
}

func (o *OrderCancelAllRequest) GetSlugsMap() (map[string]string, error) {
	slugs := map[string]string{}
	params, err := o.GetSlugParameters()
	if err != nil {
		return slugs, nil
	}

	for k, v := range params {
		slugs[k] = fmt.Sprintf("%v", v)
	}

	return slugs, nil
}

func (o *OrderCancelAllRequest) Do(ctx context.Context) ([]Order, error) {

	params, err := o.GetParameters()
	if err != nil {
		return nil, err
	}
	query := url.Values{}

	apiURL := "v2/orders/clear"

	req, err := o.client.NewAuthenticatedRequest(ctx, "POST", apiURL, query, params)
	if err != nil {
		return nil, err
	}

	response, err := o.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse []Order
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}
	return apiResponse, nil
}
