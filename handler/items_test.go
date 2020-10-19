package handler

import (
	"encoding/json"
	"errors"
	"github.com/MihaPecnik/restaurant/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type dbMock struct {
	mock.Mock
}

func (d dbMock) CreateProduct(request models.ProductPrice) error {
	return d.Called(request).Error(0)
}

func (d dbMock) ListProducts() ([]models.ProductPrice, error) {
	panic("implement me")
}

func (d dbMock) CreateOrder(order models.Order) error {
	panic("implement me")
}

func (d dbMock) AddItemToOrder(orderId, itemId int64) error {
	panic("implement me")
}

func (d dbMock) PayTheOrder(orderId int64, payment float64) (models.Receipt, error) {
	panic("implement me")
}

func (d dbMock) UpdateThePrice(id int64, cost float64) error {
	panic("implement me")
}

func TestNewHandler(t *testing.T) {
	dbMock := &dbMock{}
	handler := NewHandler(dbMock)
	assert.NotNil(t, handler)
}

func TestCreateItem(t *testing.T) {
	dbMock := &dbMock{}
	item := models.ProductPrice{Name: "burger", Cost: 10.01}
	dbMock.On("CreateProduct", item).Return(nil)

	h := NewHandler(dbMock)
	bytes, _ := json.Marshal(&item)
	req, err := http.NewRequest("POST", "/item", strings.NewReader(string(bytes)))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	h.CreateItem(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusCreated)

	expected := "\"Product created\""
	assert.Equal(t, expected, rr.Body.String())
}

func TestCreateItemInternalFail(t *testing.T) {
	dbMock := &dbMock{}
	item := models.ProductPrice{Name: "burger", Cost: 10.01}
	dbMock.On("CreateProduct", item).Return(errors.New("InternalError"))

	h := NewHandler(dbMock)
	bytes, _ := json.Marshal(&item)
	req, err := http.NewRequest("POST", "/item", strings.NewReader(string(bytes)))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	h.CreateItem(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusInternalServerError)

	expected := "{\"error\":\"InternalError\"}"
	assert.Equal(t, expected, rr.Body.String())
}

func TestCreateItemFail(t *testing.T) {
	dbMock := &dbMock{}
	note := `{s}`
	dbMock.On("CreateProduct", mock.Anything).Return(nil)

	h := NewHandler(dbMock)
	req, err := http.NewRequest("POST", "/item", strings.NewReader(note))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.CreateItem(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusBadRequest)

	expected := `{"error":"invalid character 's' looking for beginning of object key string"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestCreateItemNoName(t *testing.T) {
	dbMock := &dbMock{}
	item := models.ProductPrice{Cost: 10.01}
	dbMock.On("CreateProduct", item).Return(errors.New("InternalError"))

	h := NewHandler(dbMock)
	bytes, _ := json.Marshal(&item)
	req, err := http.NewRequest("POST", "/item", strings.NewReader(string(bytes)))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	h.CreateItem(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusBadRequest)

	expected := "{\"error\":\"cost should be higher then 0.00 and name should exist\"}"
	assert.Equal(t, expected, rr.Body.String())
}

func TestCreateItemNegativePrice(t *testing.T) {
	dbMock := &dbMock{}
	item := models.ProductPrice{Name: "Burger", Cost: -1.10}
	dbMock.On("CreateProduct", item).Return(errors.New("InternalError"))

	h := NewHandler(dbMock)
	bytes, _ := json.Marshal(&item)
	req, err := http.NewRequest("POST", "/item", strings.NewReader(string(bytes)))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	h.CreateItem(rr, req)

	assert.Equal(t,
		rr.Code, http.StatusBadRequest)

	expected := "{\"error\":\"cost should be higher then 0.00 and name should exist\"}"
	assert.Equal(t, expected, rr.Body.String())
}
