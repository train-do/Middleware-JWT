package managementvoucherhandler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	managementvoucherhandler "voucher_system/controller/management_voucher_handler"
	"voucher_system/models"
	managementvoucher "voucher_system/repository/management_voucher"
	"voucher_system/service"
	managementvoucherservice "voucher_system/service/management_voucher_service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSoftDeleteVoucher(t *testing.T) {
	// Create a new No-op logger and mock service for each test
	log := *zap.NewNop()

	t.Run("Successfully delete voucher", func(t *testing.T) {
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		r := gin.Default() // Always create a new router for each test
		r.DELETE("/voucher/:id", handler.SoftDeleteVoucher)

		voucherID := 123

		// Mock the service call to SoftDeleteVoucher to return no error (successful deletion)
		mockService.On("SoftDeleteVoucher", voucherID).Return(nil)

		// Create a request to delete the voucher with ID 123
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/voucher/%d", voucherID), nil)
		w := httptest.NewRecorder()

		// Call the handler
		r.ServeHTTP(w, req)

		// Assert the response code is 200 (OK)
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert that the service method was called with the correct voucher ID
		mockService.AssertCalled(t, "SoftDeleteVoucher", voucherID)

		// Check the response body
		expectedResponse := `{"status":true,"data":123,"message":"Deleted succesfully"}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("Fail to delete voucher due to service error", func(t *testing.T) {
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		r := gin.Default() // Always create a new router for each test
		r.DELETE("/voucher/:id", handler.SoftDeleteVoucher)

		voucherID := 123

		// Mock the service call to SoftDeleteVoucher to return an error (deletion failed)
		mockService.On("SoftDeleteVoucher", voucherID).Return(fmt.Errorf("failed to delete voucher"))

		// Create a request to delete the voucher with ID 123
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/voucher/%d", voucherID), nil)
		w := httptest.NewRecorder()

		// Call the handler
		r.ServeHTTP(w, req)

		// Assert the response code is 500 (Internal Server Error)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Assert that the service method was called with the correct voucher ID
		mockService.AssertCalled(t, "SoftDeleteVoucher", voucherID)

		// Check the response body
		expectedResponse := `{"error_msg":"FAILED", "message":"Failed to deleted Voucher", "status":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestCreateVoucher(t *testing.T) {
	// Setup logger
	log := *zap.NewNop()

	t.Run("Successfully create voucher", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.POST("/vouchers", handler.CreateVoucher)

		// Mock Data
		mockVoucher := models.Voucher{
			VoucherName:     "Test Voucher",
			VoucherCode:     "TEST123",
			VoucherType:     "e-commerce",
			VoucherCategory: "Discount",
			DiscountValue:   10.0,
			MinimumPurchase: 100.0,
			PaymentMethods:  []string{"Credit Card"},
			StartDate:       time.Now().Round(0),
			EndDate:         time.Now().AddDate(0, 1, 0).Round(0),
			ApplicableAreas: []string{"Jawa"},
			Quota:           50,
		}

		// Mock Response
		mockService.On("CreateVoucher", &mockVoucher).Return(nil)

		// Create Request Body
		body, _ := json.Marshal(mockVoucher)

		// Perform Request
		req := httptest.NewRequest(http.MethodPost, "/vouchers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "CreateVoucher", &mockVoucher)

		// Assert the JSON Response Body
		expectedResponse := map[string]interface{}{
			"status":  true,
			"message": "Created succesfully",
			"data": map[string]interface{}{
				"voucher_name":     "Test Voucher",
				"voucher_code":     "TEST123",
				"voucher_type":     "e-commerce",
				"voucher_category": "Discount",
				"discount_value":   10.0,
				"minimum_purchase": 100.0,
				"payment_methods":  []interface{}{"Credit Card"},
				"start_date":       mockVoucher.StartDate.Format(time.RFC3339Nano),
				"end_date":         mockVoucher.EndDate.Format(time.RFC3339Nano),
				"applicable_areas": []interface{}{"Jawa"},
				"quota":            float64(mockVoucher.Quota),
				"created_at":       "0001-01-01T00:00:00Z",
				"updated_at":       "0001-01-01T00:00:00Z",
			},
		}

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("Fail to create voucher due to service error", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.POST("/vouchers", handler.CreateVoucher)

		// Mock Data
		mockVoucher := models.Voucher{
			VoucherName:     "Test Voucher",
			VoucherCode:     "TEST123",
			VoucherType:     "e-commerce",
			VoucherCategory: "Discount",
			DiscountValue:   10.0,
			MinimumPurchase: 100.0,
			PaymentMethods:  []string{"Credit Card"},
			StartDate:       time.Now().Round(0),
			EndDate:         time.Now().AddDate(0, 1, 0).Round(0),
			ApplicableAreas: []string{"Jawa"},
			Quota:           50,
		}

		// Mock Response
		mockService.On("CreateVoucher", &mockVoucher).Return(fmt.Errorf("failed to create voucher"))

		// Create Request Body
		body, _ := json.Marshal(mockVoucher)

		// Perform Request
		req := httptest.NewRequest(http.MethodPost, "/vouchers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertCalled(t, "CreateVoucher", &mockVoucher)

		// Assert the JSON Response Body
		expectedResponse := `{"error_msg":"FAILED", "message":"Failed to create Voucher", "status":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestUpdateVoucher(t *testing.T) {
	// Setup logger
	log := *zap.NewNop()

	t.Run("Successfully update voucher", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.PUT("/vouchers/:id", handler.UpdateVoucher)

		// Mock Data
		mockVoucher := models.Voucher{
			VoucherName:     "Updated Voucher",
			VoucherCode:     "TEST123",
			VoucherType:     "e-commerce",
			VoucherCategory: "Discount",
			DiscountValue:   15.0,
			MinimumPurchase: 150.0,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now().Round(0),
			EndDate:         time.Now().AddDate(0, 1, 0).Round(0),
			ApplicableAreas: []string{"Bali"},
			Quota:           75,
		}

		// Mock Response
		mockService.On("UpdateVoucher", &mockVoucher, 1).Return(nil)

		// Create Request Body
		body, _ := json.Marshal(mockVoucher)

		// Perform Request
		req := httptest.NewRequest(http.MethodPut, "/vouchers/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "UpdateVoucher", &mockVoucher, 1)

		// Assert the JSON Response Body
		expectedResponse := map[string]interface{}{
			"status":  true,
			"message": "updated succesfully",
			"data":    float64(1), // ID of the updated voucher
		}

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("Fail to update voucher due to service error", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.PUT("/vouchers/:id", handler.UpdateVoucher)

		// Mock Data
		mockVoucher := models.Voucher{
			VoucherName:     "Updated Voucher",
			VoucherCode:     "TEST123",
			VoucherType:     "e-commerce",
			VoucherCategory: "Discount",
			DiscountValue:   15.0,
			MinimumPurchase: 150.0,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now().Round(0),
			EndDate:         time.Now().AddDate(0, 1, 0).Round(0),
			ApplicableAreas: []string{"Bali"},
			Quota:           75,
		}

		// Mock Response
		mockService.On("UpdateVoucher", &mockVoucher, 1).Return(fmt.Errorf("failed to update voucher"))

		// Create Request Body
		body, _ := json.Marshal(mockVoucher)

		// Perform Request
		req := httptest.NewRequest(http.MethodPut, "/vouchers/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "UpdateVoucher", &mockVoucher, 1)

		// Assert the JSON Response Body
		expectedResponse := `{"error_msg":"FAILED", "message":"Failed to Updated Voucher", "status":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestShowRedeemPoints(t *testing.T) {
	// Setup logger
	log := *zap.NewNop()

	t.Run("Successfully retrieve redeem points", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.GET("/vouchers/redeem-points", handler.ShowRedeemPoints)

		// Mock Data
		mockRedeemPoints := []managementvoucher.RedeemPoint{
			{
				VoucherName:    "Discount 10%",
				PointsRequired: 100,
				DiscountValue:  10.0,
			},
			{
				VoucherName:    "Discount 20%",
				PointsRequired: 200,
				DiscountValue:  20.0,
			},
		}

		// Mock Response
		mockService.On("ShowRedeemPoints").Return(&mockRedeemPoints, nil)

		// Perform Request
		req := httptest.NewRequest(http.MethodGet, "/vouchers/redeem-points", nil)
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "ShowRedeemPoints")

		// Assert the JSON Response Body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Redeem points retrieved successfully", actualResponse["message"])
		assert.True(t, actualResponse["status"].(bool))
	})

	t.Run("Fail to retrieve redeem points due to service error", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.GET("/vouchers/redeem-points", handler.ShowRedeemPoints)

		// Mock Response
		mockService.On("ShowRedeemPoints").Return(nil, fmt.Errorf("failed to retrieve redeem points"))

		// Perform Request
		req := httptest.NewRequest(http.MethodGet, "/vouchers/redeem-points", nil)
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertCalled(t, "ShowRedeemPoints")

		// Assert the JSON Response Body
		expectedResponse := `{"error_msg":"NOT FOUND", "message":"Reedem Points List Not Found", "status":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestGetVouchersByQueryParams(t *testing.T) {
	// Setup logger
	log := *zap.NewNop()

	t.Run("Successfully retrieve vouchers with query params", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.GET("/vouchers", handler.GetVouchersByQueryParams)

		// Mock Data
		mockVouchers := []models.Voucher{
			{
				VoucherName: "Active Voucher 1",
				VoucherType: "e-commerce",
			},
			{
				VoucherName: "Active Voucher 2",
				VoucherType: "e-commerce",
			},
		}

		// Mock Response
		mockService.On("GetVouchersByQueryParams", "active", "Jawa", "e-commerce").Return(&mockVouchers, nil)

		// Perform Request
		req := httptest.NewRequest(http.MethodGet, "/vouchers?status=active&area=Jawa&voucher_type=e-commerce", nil)
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "GetVouchersByQueryParams", "active", "Jawa", "e-commerce")

		// Assert the JSON Response Body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Voucher retrieved successfully", actualResponse["message"])
		assert.True(t, actualResponse["status"].(bool))
	})

	t.Run("Fail to retrieve vouchers due to service error", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.GET("/vouchers", handler.GetVouchersByQueryParams)

		// Mock Response
		mockService.On("GetVouchersByQueryParams", "expired", "", "").Return(nil, fmt.Errorf("failed to retrieve vouchers"))

		// Perform Request
		req := httptest.NewRequest(http.MethodGet, "/vouchers?status=expired", nil)
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertCalled(t, "GetVouchersByQueryParams", "expired", "", "")

		// Assert the JSON Response Body
		expectedResponse := `{"error_msg":"NOT FOUND", "message":"Voucher Not Found", "status":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestCreateRedeemVoucher(t *testing.T) {
	log := *zap.NewNop() // Setup logger

	t.Run("Successfully create redeem voucher", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.POST("/redeem", handler.CreateRedeemVoucher)

		// Input Payload
		payload := `{"voucher_id": 1, "user_id": 2, "points": 100}`

		// Mock Response
		mockRedeem := models.Redeem{
			VoucherID: 1,
			UserID:    2,
		}
		mockService.On("CreateRedeemVoucher", &mockRedeem, 100).Return(nil)

		// Perform Request
		req := httptest.NewRequest(http.MethodPost, "/redeem", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "CreateRedeemVoucher", &mockRedeem, 100)

		// Assert the JSON Response Body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Created successfully", actualResponse["message"])
		assert.True(t, actualResponse["status"].(bool))
	})

	t.Run("Fail to create redeem voucher due to invalid payload", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.POST("/redeem", handler.CreateRedeemVoucher)

		// Invalid Payload
		payload := `{"voucher_id": 1, "user_id": "invalid", "points": 100}`

		// Perform Request
		req := httptest.NewRequest(http.MethodPost, "/redeem", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateRedeemVoucher")

		// Assert the JSON Response Body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid Payload: json: cannot unmarshal string into Go struct field .user_id of type int", actualResponse["message"])
		assert.False(t, actualResponse["status"].(bool))
	})

	t.Run("Fail to create redeem voucher due to service error", func(t *testing.T) {
		// Mock Service
		mockService := &managementvoucherservice.ManagementVoucherServiceMock{}
		service := service.Service{
			Manage: mockService,
		}
		handler := managementvoucherhandler.NewManagementVoucherHanlder(service, &log)

		// Router and Endpoint
		r := gin.Default()
		r.POST("/redeem", handler.CreateRedeemVoucher)

		// Input Payload
		payload := `{"voucher_id": 1, "user_id": 2, "points": 100}`

		// Mock Response
		mockRedeem := models.Redeem{
			VoucherID: 1,
			UserID:    2,
		}
		mockService.On("CreateRedeemVoucher", &mockRedeem, 100).Return(fmt.Errorf("service error"))

		// Perform Request
		req := httptest.NewRequest(http.MethodPost, "/redeem", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call the Handler
		r.ServeHTTP(w, req)

		// Assert the Response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "CreateRedeemVoucher", &mockRedeem, 100)

		// Assert the JSON Response Body
		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create redeem voucher: service error", actualResponse["message"])
		assert.False(t, actualResponse["status"].(bool))
	})
}
