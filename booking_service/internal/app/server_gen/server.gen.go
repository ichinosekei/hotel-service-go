// Package server_gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package server_gen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Booking defines model for Booking.
type Booking struct {
	// BookingId Unique ID of the booking.
	BookingId *int `json:"bookingId,omitempty"`

	// CheckInDate Check-in date of the booking.
	CheckInDate *openapi_types.Date `json:"checkInDate,omitempty"`

	// ClientFullName Full name of the client.
	ClientFullName *string `json:"clientFullName,omitempty"`

	// ClientPhoneNumber Phone number of the client.
	ClientPhoneNumber *string `json:"clientPhoneNumber,omitempty"`

	// Duration Duration of the booking in nights.
	Duration *int `json:"duration,omitempty"`

	// HotelId id of the hotel.
	HotelId *int `json:"hotelId,omitempty"`

	// RoomNumber Room number in the hotel.
	RoomNumber *string `json:"roomNumber,omitempty"`

	// TotalPrice Total price of the booking.
	TotalPrice *float32 `json:"totalPrice,omitempty"`
}

// BookingRequest defines model for BookingRequest.
type BookingRequest struct {
	// CheckInDate Check-in date of the booking.
	CheckInDate *openapi_types.Date `json:"checkInDate,omitempty"`

	// ClientFullName Full name of the client.
	ClientFullName *string `json:"clientFullName,omitempty"`

	// ClientPhoneNumber Phone number of the client.
	ClientPhoneNumber *string `json:"clientPhoneNumber,omitempty"`

	// Duration Duration of the booking in nights.
	Duration *int `json:"duration,omitempty"`

	// HotelId Id of the hotel.
	HotelId *int `json:"hotelId,omitempty"`

	// RoomNumber Room number in the hotel.
	RoomNumber *int `json:"roomNumber,omitempty"`
}

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

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+yWzXIbNwzHX4XD9riW5CY9VDc5blId6tG4H5dOpkMtIQnJLkmDWCWajN69A1K7lrRr",
	"100aTw+9SQT5BwH8AO4nXfo6eAeOo55+0rHcQG3Szyvv36Nby89APgAxQjIss2Fu5Y+FWBIGRu/0VP/m",
	"8K4BNb9WfqV4A+qwd6QLzbsAeqrRMayB9L7Q5QbK93N3bRj6Uq/EeIFOWcMwILfyVBvWUy32e/nIJJcW",
	"9QrB8eumqm5MPeBALMqZuhPPB0YPay023sFNUy+B+nLJqFyyPkHRNmTyyXOh64PlLGaFTjlcbzgOJ3Pj",
	"GaqhmqBtldKW4dPkff1QZLfe121g6AaV7uNiz6ZaEJYDGf9VbCqI8bGCripv+F47e9b7fbfil++gZPF2",
	"YPQW7hqI3Ef1f8Keg7D5MxDWKfUxkCV0K9/XnS3mauVJ1caZtUSY06WMs9lBG3qOGLkS3QNU6hegraA6",
	"W8x1obdAMatejiajiYTkAzgTUE/1i9Fk9EIXOhjeJPDGJuB4ezlu9ROaPiN6RiGB4GeUgw9dJT4gb1Ia",
	"AvktWrDKAhus0j0F8FQ/qYVe+MizgL9fXrWuCk25H6683aUu8I7BJd8mhArLdHr8LmY48syXX98SrPRU",
	"fzO+fxTGhxdhfNZqKeviBwmsnjI1kBZi8C7m5vtucvlve89uTxPYlis2ZQkxrpqq2qkyZdWOpEwvJ5N+",
	"2q+MVYdY1IX6GWMUDS/8bU2FVqELDctoMAeRHwZq592qwlIU2luYisDYnYKPGDkm+qSOa9yCU9IDCb40",
	"cci4NSTx74duOHcM5EyVOARSPxJ5Uhdq5lTj4GOAksEqSKu+LBuiFO++0LGpa0O7Y7iWXf6KHpvj3BVy",
	"gzUMEHoLTAhbkakwsvR6ezTFZ1QMUOIKy7a/liaCVT51MpIKR4OrD/AbOOX3Vb6NNBOZGhgo6ukfvbdE",
	"muPRgYiy7a4B2ulCuzSgdTiar+f4FkcodqM/HehP1/3bHuuTf8Q6MtTxydB3/g2R2Q01wawrTVuCbrB9",
	"dgecli3JvOzL3HhWr33jrLpQN/6YDFk75b+n+NXAfwN8nolh9tMz8AXoH56RFnjkqObXfw/5T8ntExhv",
	"zr+ou4dxiO/2iX6M7d6D+t9l+STXEvxh0IBt0/DZcOe6Sam+DOyhG31Vqh8avRnklM6YPAwRdf5pI9wq",
	"+WjRhW6o0lPddsafMW/R+7f7vwIAAP//TqiXEyUOAAA=",
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
