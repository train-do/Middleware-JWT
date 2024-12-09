package controller

import (
	"net/http"
	"strconv"
	"time"
	"voucher_system/helper"
	"voucher_system/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VoucherController struct {
	service service.Service
	log     *zap.Logger
}

func NewVoucherController(service service.Service, log *zap.Logger) *VoucherController {
	return &VoucherController{service: service, log: log}
}

func (c *VoucherController) FindVouchers(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		helper.ResponseError(ctx, err.Error(), "Invalid user ID", http.StatusBadRequest)
		return
	}
	voucherType := ctx.Query("type")

	voucher, err := c.service.Voucher.FindVouchers(userID, voucherType)
	if err != nil {
		if err.Error() == "no vouchers available" {
			helper.ResponseError(ctx, err.Error(), "", http.StatusBadRequest)
			return
		}
		c.log.Error("Error fetching vouchers", zap.Error(err))
		c.log.Debug("Error fetching vouchers", zap.Error(err))
		helper.ResponseError(ctx, err.Error(), "", http.StatusInternalServerError)
		return
	}

	result := gin.H{
		"voucher": voucher,
	}

	helper.ResponseOK(ctx, result, "", http.StatusOK)
}

func (c *VoucherController) ValidateVoucher(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		helper.ResponseError(ctx, err.Error(), "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request struct {
		VoucherCode       string  `json:"voucher_code" binding:"required"`
		TransactionAmount float64 `json:"transaction_amount" binding:"required"`
		ShippingAmount    float64 `json:"shipping_amount" binding:"required"`
		Area              string  `json:"area" binding:"required"`
		PaymentMethod     string  `json:"payment_method" binding:"required"`
		TransactionDate   string  `json:"transaction_date" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		helper.ResponseError(ctx, "Invalid input", err.Error(), http.StatusBadRequest)
		return
	}

	transactionDate, err := time.Parse("2006-01-02", request.TransactionDate)
	if err != nil {
		c.log.Error("Handler: Invalid date format", zap.Error(err))
		helper.ResponseError(ctx, "Invalid date format", "Transaction date must be in format YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	voucher, benefitValue, err := c.service.Voucher.ValidateVoucher(userID, request.VoucherCode, request.TransactionAmount, request.ShippingAmount, request.Area, request.PaymentMethod, transactionDate)
	if err != nil {
		c.log.Error("Error fetching voucher", zap.Error(err))
		c.log.Debug("Error fetching voucher", zap.Error(err))
		helper.ResponseError(ctx, "Voucher validation failed", err.Error(), http.StatusBadRequest)
		return
	}

	var msg string
	if voucher.Status {
		msg = "valid"
	} else {
		msg = "invalid"
	}

	response := gin.H{
		"benefit_value": benefitValue,
		"status":        msg,
	}
	helper.ResponseOK(ctx, response, "Voucher is valid", http.StatusOK)
}

func (c *VoucherController) UseVoucher(ctx *gin.Context) {
	var request struct {
		UserID            int     `json:"user_id"`
		VoucherCode       string  `json:"voucher_code"`
		TransactionAmount float64 `json:"transaction_amount"`
		PaymentMethod     string  `json:"payment_method"`
		Area              string  `json:"area"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		c.log.Error("Error invalid input", zap.Error(err))
		c.log.Debug("Error invalid input", zap.Error(err))
		helper.ResponseError(ctx, "Invalid input", err.Error(), http.StatusBadRequest)
		return
	}

	err := c.service.Voucher.UseVoucher(request.UserID, request.VoucherCode, request.TransactionAmount, request.PaymentMethod, request.Area)
	if err != nil {
		c.log.Error("Error failed used voucher", zap.Error(err))
		c.log.Debug("Error failed used voucher", zap.Error(err))
		helper.ResponseError(ctx, "Failed used voucher", err.Error(), http.StatusBadRequest)
		return
	}
	helper.ResponseOK(ctx, nil, "voucher used successfully", http.StatusCreated)

}
