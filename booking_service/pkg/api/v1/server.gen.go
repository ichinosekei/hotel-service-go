// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Defines values for BookingPaymentStatus.
const (
	NotPaid BookingPaymentStatus = "not paid"
	Paid    BookingPaymentStatus = "paid"
)

// Defines values for PaymentWebhookRequestStatus.
const (
	Failed PaymentWebhookRequestStatus = "failed"
	Ok     PaymentWebhookRequestStatus = "ok"
)

// Booking defines model for Booking.
type Booking struct {
	// BookingId Unique ID of the booking.
	BookingId *string `json:"bookingId,omitempty"`

	// CheckInDate Check-in date of the booking.
	CheckInDate *string `json:"checkInDate,omitempty"`

	// CheckOutDate Check-in date of the booking.
	CheckOutDate *string `json:"checkOutDate,omitempty"`

	// ClientFullName Full name of the client.
	ClientFullName *string `json:"clientFullName,omitempty"`

	// ClientPhoneNumber Phone number of the client.
	ClientPhoneNumber *string `json:"clientPhoneNumber,omitempty"`

	// HotelId id of the hotel.
	HotelId *int `json:"hotelId,omitempty"`

	// HotelierPhoneNumber Phone number of the hotelier.
	HotelierPhoneNumber *string `json:"hotelierPhoneNumber,omitempty"`

	// PaymentStatus Status of the payment
	PaymentStatus *BookingPaymentStatus `json:"paymentStatus,omitempty"`

	// RoomNumber Room number in the hotel.
	RoomNumber *int `json:"roomNumber,omitempty"`

	// TotalPrice Total price of the booking.
	TotalPrice *float64 `json:"totalPrice,omitempty"`
}

// BookingPaymentStatus Status of the payment
type BookingPaymentStatus string

// BookingRequest defines model for BookingRequest.
type BookingRequest struct {
	// CheckInDate Check-in date of the booking.
	CheckInDate string `json:"checkInDate"`

	// CheckOutDate Check-in date of the booking.
	CheckOutDate string `json:"checkOutDate"`

	// ClientFullName Full name of the client.
	ClientFullName string `json:"clientFullName"`

	// ClientPhoneNumber Phone number of the client.
	ClientPhoneNumber string `json:"clientPhoneNumber"`

	// HotelId Id of the hotel.
	HotelId int `json:"hotelId"`

	// HotelierPhoneNumber Phone number of the hotelier.
	HotelierPhoneNumber string `json:"hotelierPhoneNumber"`

	// RoomNumber Room number in the hotel.
	RoomNumber int `json:"roomNumber"`
}

// PaymentWebhookRequest defines model for PaymentWebhookRequest.
type PaymentWebhookRequest struct {
	// BookingId Identifier for the booking associated with this payment.
	BookingId string `json:"bookingId"`

	// PaymentId Unique identifier for the payment.
	PaymentId string `json:"paymentId"`

	// Status Status of the payment (success or failure).
	Status PaymentWebhookRequestStatus `json:"status"`
}

// PaymentWebhookRequestStatus Status of the payment (success or failure).
type PaymentWebhookRequestStatus string

// GetApiV1BookingsClientParams defines parameters for GetApiV1BookingsClient.
type GetApiV1BookingsClientParams struct {
	// PhoneNumber The phone number of the client.
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber"`
}

// GetApiV1BookingsHotelParams defines parameters for GetApiV1BookingsHotel.
type GetApiV1BookingsHotelParams struct {
	// HotelId The unique ID of the hotel.
	HotelId int `form:"hotelId" json:"hotelId"`
}

// PostApiV1BookingsJSONRequestBody defines body for PostApiV1Bookings for application/json ContentType.
type PostApiV1BookingsJSONRequestBody = BookingRequest

// PostApiV1WebhookPaymentJSONRequestBody defines body for PostApiV1WebhookPayment for application/json ContentType.
type PostApiV1WebhookPaymentJSONRequestBody = PaymentWebhookRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostApiV1BookingsWithBody request with any body
	PostApiV1BookingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostApiV1Bookings(ctx context.Context, body PostApiV1BookingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetApiV1BookingsClient request
	GetApiV1BookingsClient(ctx context.Context, params *GetApiV1BookingsClientParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetApiV1BookingsHotel request
	GetApiV1BookingsHotel(ctx context.Context, params *GetApiV1BookingsHotelParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostApiV1WebhookPaymentWithBody request with any body
	PostApiV1WebhookPaymentWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostApiV1WebhookPayment(ctx context.Context, body PostApiV1WebhookPaymentJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostApiV1BookingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1BookingsRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1Bookings(ctx context.Context, body PostApiV1BookingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1BookingsRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetApiV1BookingsClient(ctx context.Context, params *GetApiV1BookingsClientParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetApiV1BookingsClientRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetApiV1BookingsHotel(ctx context.Context, params *GetApiV1BookingsHotelParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetApiV1BookingsHotelRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1WebhookPaymentWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1WebhookPaymentRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1WebhookPayment(ctx context.Context, body PostApiV1WebhookPaymentJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1WebhookPaymentRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostApiV1BookingsRequest calls the generic PostApiV1Bookings builder with application/json body
func NewPostApiV1BookingsRequest(server string, body PostApiV1BookingsJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostApiV1BookingsRequestWithBody(server, "application/json", bodyReader)
}

// NewPostApiV1BookingsRequestWithBody generates requests for PostApiV1Bookings with any type of body
func NewPostApiV1BookingsRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/bookings")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetApiV1BookingsClientRequest generates requests for GetApiV1BookingsClient
func NewGetApiV1BookingsClientRequest(server string, params *GetApiV1BookingsClientParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/bookings/client")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "phoneNumber", runtime.ParamLocationQuery, params.PhoneNumber); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetApiV1BookingsHotelRequest generates requests for GetApiV1BookingsHotel
func NewGetApiV1BookingsHotelRequest(server string, params *GetApiV1BookingsHotelParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/bookings/hotel")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "hotelId", runtime.ParamLocationQuery, params.HotelId); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostApiV1WebhookPaymentRequest calls the generic PostApiV1WebhookPayment builder with application/json body
func NewPostApiV1WebhookPaymentRequest(server string, body PostApiV1WebhookPaymentJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostApiV1WebhookPaymentRequestWithBody(server, "application/json", bodyReader)
}

// NewPostApiV1WebhookPaymentRequestWithBody generates requests for PostApiV1WebhookPayment with any type of body
func NewPostApiV1WebhookPaymentRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/webhook/payment")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostApiV1BookingsWithBodyWithResponse request with any body
	PostApiV1BookingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1BookingsResponse, error)

	PostApiV1BookingsWithResponse(ctx context.Context, body PostApiV1BookingsJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1BookingsResponse, error)

	// GetApiV1BookingsClientWithResponse request
	GetApiV1BookingsClientWithResponse(ctx context.Context, params *GetApiV1BookingsClientParams, reqEditors ...RequestEditorFn) (*GetApiV1BookingsClientResponse, error)

	// GetApiV1BookingsHotelWithResponse request
	GetApiV1BookingsHotelWithResponse(ctx context.Context, params *GetApiV1BookingsHotelParams, reqEditors ...RequestEditorFn) (*GetApiV1BookingsHotelResponse, error)

	// PostApiV1WebhookPaymentWithBodyWithResponse request with any body
	PostApiV1WebhookPaymentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1WebhookPaymentResponse, error)

	PostApiV1WebhookPaymentWithResponse(ctx context.Context, body PostApiV1WebhookPaymentJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1WebhookPaymentResponse, error)
}

type PostApiV1BookingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Booking
}

// Status returns HTTPResponse.Status
func (r PostApiV1BookingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostApiV1BookingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetApiV1BookingsClientResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Booking
}

// Status returns HTTPResponse.Status
func (r GetApiV1BookingsClientResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetApiV1BookingsClientResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetApiV1BookingsHotelResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Booking
}

// Status returns HTTPResponse.Status
func (r GetApiV1BookingsHotelResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetApiV1BookingsHotelResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostApiV1WebhookPaymentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r PostApiV1WebhookPaymentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostApiV1WebhookPaymentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostApiV1BookingsWithBodyWithResponse request with arbitrary body returning *PostApiV1BookingsResponse
func (c *ClientWithResponses) PostApiV1BookingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1BookingsResponse, error) {
	rsp, err := c.PostApiV1BookingsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1BookingsResponse(rsp)
}

func (c *ClientWithResponses) PostApiV1BookingsWithResponse(ctx context.Context, body PostApiV1BookingsJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1BookingsResponse, error) {
	rsp, err := c.PostApiV1Bookings(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1BookingsResponse(rsp)
}

// GetApiV1BookingsClientWithResponse request returning *GetApiV1BookingsClientResponse
func (c *ClientWithResponses) GetApiV1BookingsClientWithResponse(ctx context.Context, params *GetApiV1BookingsClientParams, reqEditors ...RequestEditorFn) (*GetApiV1BookingsClientResponse, error) {
	rsp, err := c.GetApiV1BookingsClient(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetApiV1BookingsClientResponse(rsp)
}

// GetApiV1BookingsHotelWithResponse request returning *GetApiV1BookingsHotelResponse
func (c *ClientWithResponses) GetApiV1BookingsHotelWithResponse(ctx context.Context, params *GetApiV1BookingsHotelParams, reqEditors ...RequestEditorFn) (*GetApiV1BookingsHotelResponse, error) {
	rsp, err := c.GetApiV1BookingsHotel(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetApiV1BookingsHotelResponse(rsp)
}

// PostApiV1WebhookPaymentWithBodyWithResponse request with arbitrary body returning *PostApiV1WebhookPaymentResponse
func (c *ClientWithResponses) PostApiV1WebhookPaymentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1WebhookPaymentResponse, error) {
	rsp, err := c.PostApiV1WebhookPaymentWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1WebhookPaymentResponse(rsp)
}

func (c *ClientWithResponses) PostApiV1WebhookPaymentWithResponse(ctx context.Context, body PostApiV1WebhookPaymentJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1WebhookPaymentResponse, error) {
	rsp, err := c.PostApiV1WebhookPayment(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1WebhookPaymentResponse(rsp)
}

// ParsePostApiV1BookingsResponse parses an HTTP response from a PostApiV1BookingsWithResponse call
func ParsePostApiV1BookingsResponse(rsp *http.Response) (*PostApiV1BookingsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostApiV1BookingsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Booking
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseGetApiV1BookingsClientResponse parses an HTTP response from a GetApiV1BookingsClientWithResponse call
func ParseGetApiV1BookingsClientResponse(rsp *http.Response) (*GetApiV1BookingsClientResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetApiV1BookingsClientResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Booking
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetApiV1BookingsHotelResponse parses an HTTP response from a GetApiV1BookingsHotelWithResponse call
func ParseGetApiV1BookingsHotelResponse(rsp *http.Response) (*GetApiV1BookingsHotelResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetApiV1BookingsHotelResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Booking
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostApiV1WebhookPaymentResponse parses an HTTP response from a PostApiV1WebhookPaymentWithResponse call
func ParsePostApiV1WebhookPaymentResponse(rsp *http.Response) (*PostApiV1WebhookPaymentResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostApiV1WebhookPaymentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a booking
	// (POST /api/v1/bookings)
	PostApiV1Bookings(ctx echo.Context) error
	// Get client bookings
	// (GET /api/v1/bookings/client)
	GetApiV1BookingsClient(ctx echo.Context, params GetApiV1BookingsClientParams) error
	// Get bookings for a specific hotel
	// (GET /api/v1/bookings/hotel)
	GetApiV1BookingsHotel(ctx echo.Context, params GetApiV1BookingsHotelParams) error
	// Webhook for processing payment status
	// (POST /api/v1/webhook/payment)
	PostApiV1WebhookPayment(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostApiV1Bookings converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1Bookings(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1Bookings(ctx)
	return err
}

// GetApiV1BookingsClient converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1BookingsClient(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetApiV1BookingsClientParams
	// ------------- Required query parameter "phoneNumber" -------------

	err = runtime.BindQueryParameter("form", true, true, "phoneNumber", ctx.QueryParams(), &params.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter phoneNumber: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1BookingsClient(ctx, params)
	return err
}

// GetApiV1BookingsHotel converts echo context to params.
func (w *ServerInterfaceWrapper) GetApiV1BookingsHotel(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetApiV1BookingsHotelParams
	// ------------- Required query parameter "hotelId" -------------

	err = runtime.BindQueryParameter("form", true, true, "hotelId", ctx.QueryParams(), &params.HotelId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter hotelId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetApiV1BookingsHotel(ctx, params)
	return err
}

// PostApiV1WebhookPayment converts echo context to params.
func (w *ServerInterfaceWrapper) PostApiV1WebhookPayment(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostApiV1WebhookPayment(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/api/v1/bookings", wrapper.PostApiV1Bookings)
	router.GET(baseURL+"/api/v1/bookings/client", wrapper.GetApiV1BookingsClient)
	router.GET(baseURL+"/api/v1/bookings/hotel", wrapper.GetApiV1BookingsHotel)
	router.POST(baseURL+"/api/v1/webhook/payment", wrapper.PostApiV1WebhookPayment)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+yXwXLbNhCGXwWD9tDOyJLcpIfqZidNqkNdTZK2h4wngciltDEJwMBSjiajd+8sQNKk",
	"SCuyY0976I0iiB/A7reLX19kYgprNGjycvZF+mQNhQqP58ZcoV7xo3XGgiOEMLCMA/OUf6TgE4eW0Gg5",
	"k39qvC5BzF8Kkwlag6i+HcuRpK0FOZOeHKvuRjJZQ3I11y8VQV/pBQ+eoBapIjha7Y+SHksuR9D0qszz",
	"C1UMCPKI0KpoxOKEA1qLtdFwURZLcH25MCh0GD1CcW0I8qEEYFrPDp+0JqMmWIFrZiO4e++onji4J6u2",
	"BWh6S4pK39eL72ul6mM5kqDLQs7eS6swlSOpDYnweDmwhDOmuGu/b4wp6u2i/loIyJDKFw6TgdS+4zFh",
	"eXCAlMy4QpGcydSUyxxu1ePacrdr3pjlJ0iIl6tq6Q1cl+CpX1L/l8Ljl8L83yuFx+KUpeC6RAcpl0ib",
	"kr0s97I0FOrh496GsLPxywGMF7Fq/4bl2pirO2k+cEHMU9CEGYITmXFt9ITy3iSoCFJxg7QWtEZf94lD",
	"DefAPYT91Q4J+vu0LvGDL5MEvBfGiUxhXjr4cdxqaOaKu4XCHIaa2V5qb88yaoWv2VI/GSyAOjP93Z4t",
	"5uG0hdJqxYGNIAil00haHXEfgoCUs27VocRbcBtufGeLuRzJDTgfVU/H0/GUg2QsaGVRzuSz8XT8THIe",
	"aB2iNlEWJ5vTSa0fyDCRkL0e5ICbjxIabhoAqqyDsM5sMIVUpEAK87BP5kvxbE63XBhPZxb/Oj2vl4rh",
	"BE/nJt2Glmo08Q0z+yKVtTkmYfbkk+cN1EaHn753kMmZ/G5y64QmlQ2a7PXtXTdt5EoIL7w12kf2f5qe",
	"PvbqcdluAOt0VQxmZZ5vRRKimo45Tc+n037Yz1UqqrOIE/E7es8ahhvRRuWYCtS2JL4YVCXyy0DujM5y",
	"TFih3oXKHah0K+AzevJNra1wA1pwTwnwhfvGKb2CIP7z0A7nmsBplQcOwYlfnTNOnIgzLUoNny0k3B8g",
	"vDVJUjoXzsvFWxaFcts2XMsmfqMem5NYFbyDFQwQ+gbIIWxYJkdPXPz11HA+JbyFBDNM6vpaKg+pMKGl",
	"oxO2dWn0AX4NXX5fxN1wMTlVAIHzcva+50y4OA5ej8ifXZfgtuylwnUtbafZd/EdtVBsnE2Y0G+Qu8se",
	"69N7sY4EhT8a+mZ95ZzaDhXBWZOaOgVNY3twBXTTFmSe92UuDIlXptSpOBEXpk0Gv+vy31N8MvBfA+1H",
	"Ypj9cA18A/rVNVIDj+TF/OXXIf8tLHsE4+X+38jGIQ3x3bIvd7Ldc1b/XZY7sebDV40G0joMD4Y75o1T",
	"9W1gD+3oSam+q/VGkNuM30RvOqn/YN5pQd6xvQSdWoOahIMEcAO+4/Ci+RKZM0X7/Qe/9QTxTist32q+",
	"Y2SraSpJjEtRr/LtAQNTeelF83/4KWzMsG8/ys0MJLWS6XoP6wz/ON59zCsm2WyIePUch+X5vl28pfFj",
	"Y50/Cm0oUvsgOrtIips15lCfkRfndatM7fNah4dRbc3oQhWD78O6Q11w345zrxVstOVIli6XM1l38w8+",
	"fiJ3l7t/AgAA//9AjikbzhMAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
