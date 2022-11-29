package registrantalert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// RegistrantAlert is an interface for Registrant Alert API.
type RegistrantAlert interface {
	// BasicPreview returns only the number of domains for the basic search. No credits deducted.
	BasicPreview(ctx context.Context, basicSearchTerms *BasicSearchTerms, option ...Option) (int, *Response, error)

	// BasicPurchase returns parsed Registrant Alert API response for the basic search.
	BasicPurchase(ctx context.Context, basicSearchTerms *BasicSearchTerms, option ...Option) (*RegistrantAlertResponse, *Response, error)

	// BasicRawData returns raw Registrant Alert API response for the basic search.
	BasicRawData(ctx context.Context, basicSearchTerms *BasicSearchTerms, option ...Option) (*Response, error)

	// AdvancedPreview returns only the number of domains for the advanced search. No credits deducted.
	AdvancedPreview(ctx context.Context, advancedSearchTerms []AdvancedSearchTerm, option ...Option) (int, *Response, error)

	// AdvancedPurchase returns parsed Registrant Alert API response for the advanced search.
	AdvancedPurchase(ctx context.Context, advancedSearchTerms []AdvancedSearchTerm, option ...Option) (*RegistrantAlertResponse, *Response, error)

	// AdvancedRawData returns raw Registrant Alert API response for the advanced search.
	AdvancedRawData(ctx context.Context, advancedSearchTerms []AdvancedSearchTerm, option ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body.
	Body []byte
}

// registrantAlertServiceOp is the type implementing the RegistrantAlert interface.
type registrantAlertServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ RegistrantAlert = &registrantAlertServiceOp{}

// newRequest creates the API request with default parameters and specified body.
func (service registrantAlertServiceOp) newRequest(body []byte) (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodPost, service.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return req, nil
}

// apiResponse is used for parsing Registrant Alert API response as a model instance.
type apiResponse struct {
	RegistrantAlertResponse
	ErrorMessage
}

// validateBasicSearchTerms validates the terms of the basic search.
func validateBasicSearchTerms(basicSearchTerms *BasicSearchTerms) error {
	const limitOfSearchTerms = 4

	if basicSearchTerms == nil {
		return &ArgError{"basicSearchTerms.include", "is required."}
	}

	if basicSearchTerms.Include == nil || len(basicSearchTerms.Include) == 0 ||
		len(basicSearchTerms.Include) > limitOfSearchTerms {
		return &ArgError{"basicSearchTerms.include", "must have between 1 and 4 items."}
	}

	if basicSearchTerms.Exclude != nil && len(basicSearchTerms.Exclude) > limitOfSearchTerms {
		return &ArgError{"basicSearchTerms.exclude", "must have between 0 and 4 items."}
	}

	return nil
}

// validateAdvancedSearchTerms validates the terms of the advanced search.
func validateAdvancedSearchTerms(advancedSearchTerms []AdvancedSearchTerm) error {
	if advancedSearchTerms == nil {
		return &ArgError{"advancedSearchTerms", "is required."}
	}
	if len(advancedSearchTerms) == 0 || len(advancedSearchTerms) > 4 {
		return &ArgError{"advancedSearchTerms", "must have between 1 and 4 items."}
	}

	for i, searchTerm := range advancedSearchTerms {
		if searchTerm.Field == "" {
			return &ArgError{"advancedSearchTerms." + strconv.Itoa(i) + ".Field", "is required."}
		}
		if searchTerm.Term == "" {
			return &ArgError{"advancedSearchTerms." + strconv.Itoa(i) + ".Term", "is required."}
		}
	}

	return nil
}

// validateOptions validates options.
func validateOptions(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			return &ArgError{"Option", "can not be nil"}
		}
	}
	return nil
}

// request returns intermediate API response for further actions.
func (service registrantAlertServiceOp) request(
	ctx context.Context,
	basicSearchTerms *BasicSearchTerms,
	advancedSearchTerms []AdvancedSearchTerm,
	purchase bool,
	opts ...Option) (*Response, error) {
	var request = &registrantAlertRequest{
		service.client.apiKey,
		basicSearchTerms,
		advancedSearchTerms,
		"",
		"preview",
		true,
		"json",
		"",
		"",
		"",
		"",
		"",
		"",
	}

	if purchase {
		request.Mode = "purchase"
	}

	if err := validateOptions(opts...); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(request)
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := service.newRequest(requestBody)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Registrant Alert API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// BasicPurchase returns parsed Registrant Alert API response.
func (service registrantAlertServiceOp) BasicPurchase(
	ctx context.Context,
	basicSearchTerms *BasicSearchTerms,
	opts ...Option,
) (registrantAlertResponse *RegistrantAlertResponse, resp *Response, err error) {
	err = validateBasicSearchTerms(basicSearchTerms)
	if err != nil {
		return nil, nil, err
	}
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionResponseFormat("json"))

	resp, err = service.request(ctx, basicSearchTerms, nil, true, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	registrantAlertResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if registrantAlertResp.Message != nil || registrantAlertResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    registrantAlertResp.Code,
			Message: registrantAlertResp.Message,
		}
	}

	return &registrantAlertResp.RegistrantAlertResponse, resp, nil
}

// BasicPreview returns only the number of domains. No credits deducted.
func (service registrantAlertServiceOp) BasicPreview(
	ctx context.Context,
	basicSearchTerms *BasicSearchTerms,
	opts ...Option,
) (domainsCount int, resp *Response, err error) {
	err = validateBasicSearchTerms(basicSearchTerms)
	if err != nil {
		return 0, nil, err
	}

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionResponseFormat("json"))

	resp, err = service.request(ctx, basicSearchTerms, nil, false, optsJSON...)
	if err != nil {
		return 0, resp, err
	}

	registrantAlertResp, err := parse(resp.Body)
	if err != nil {
		return 0, resp, err
	}

	if registrantAlertResp.Message != nil || registrantAlertResp.Code != 0 {
		return 0, nil, &ErrorMessage{
			Code:    registrantAlertResp.Code,
			Message: registrantAlertResp.Message,
		}
	}

	return registrantAlertResp.DomainsCount, resp, nil
}

// BasicRawData returns raw Registrant Alert API response as the Response struct with Body saved as a byte slice.
func (service registrantAlertServiceOp) BasicRawData(
	ctx context.Context,
	basicSearchTerms *BasicSearchTerms,
	opts ...Option,
) (resp *Response, err error) {
	err = validateBasicSearchTerms(basicSearchTerms)
	if err != nil {
		return nil, err
	}

	resp, err = service.request(ctx, basicSearchTerms, nil, true, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// AdvancedPreview returns only the number of domains. No credits deducted.
func (service registrantAlertServiceOp) AdvancedPreview(
	ctx context.Context,
	advancedSearchTerms []AdvancedSearchTerm,
	opts ...Option,
) (domainsCount int, resp *Response, err error) {
	err = validateAdvancedSearchTerms(advancedSearchTerms)
	if err != nil {
		return 0, nil, err
	}

	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionResponseFormat("json"))

	resp, err = service.request(ctx, nil, advancedSearchTerms, false, optsJSON...)
	if err != nil {
		return 0, resp, err
	}

	registrantAlertResp, err := parse(resp.Body)
	if err != nil {
		return 0, resp, err
	}

	if registrantAlertResp.Message != nil || registrantAlertResp.Code != 0 {
		return 0, nil, &ErrorMessage{
			Code:    registrantAlertResp.Code,
			Message: registrantAlertResp.Message,
		}
	}

	return registrantAlertResp.DomainsCount, resp, nil
}

// AdvancedPurchase returns parsed Registrant Alert API response.
func (service registrantAlertServiceOp) AdvancedPurchase(
	ctx context.Context,
	advancedSearchTerms []AdvancedSearchTerm,
	opts ...Option,
) (registrantAlertResponse *RegistrantAlertResponse, resp *Response, err error) {
	err = validateAdvancedSearchTerms(advancedSearchTerms)
	if err != nil {
		return nil, nil, err
	}
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionResponseFormat("json"))

	resp, err = service.request(ctx, nil, advancedSearchTerms, true, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	registrantAlertResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if registrantAlertResp.Message != nil || registrantAlertResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    registrantAlertResp.Code,
			Message: registrantAlertResp.Message,
		}
	}

	return &registrantAlertResp.RegistrantAlertResponse, resp, nil
}

// AdvancedRawData returns raw Registrant Alert API response as the Response struct with Body saved as a byte slice.
func (service registrantAlertServiceOp) AdvancedRawData(
	ctx context.Context,
	advancedSearchTerms []AdvancedSearchTerm,
	opts ...Option,
) (resp *Response, err error) {
	err = validateAdvancedSearchTerms(advancedSearchTerms)
	if err != nil {
		return nil, err
	}

	resp, err = service.request(ctx, nil, advancedSearchTerms, true, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
