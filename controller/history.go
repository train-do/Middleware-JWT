package controller

import (
	"net/http"
	"strconv"
	"voucher_system/helper"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (c *VoucherController) GetRedeemHistoryByUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		helper.ResponseError(ctx, "Invalid user ID", err.Error(), http.StatusBadRequest)
		return
	}
	redeems, err := c.service.History.GetRedeemHistoryByUser(userID)
	if err != nil {
		c.log.Error("Error fetching redeem history", zap.Error(err))
		helper.ResponseError(ctx, "Failed to fetch history", err.Error(), http.StatusInternalServerError)
		return
	}
	helper.ResponseOK(ctx, gin.H{"redeem_history": redeems}, "Redeem history fetched successfully", http.StatusOK)
}

func (c *VoucherController) GetUsageHistoryByUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		helper.ResponseError(ctx, "Invalid user ID", err.Error(), http.StatusBadRequest)
		return
	}
	histories, err := c.service.History.GetUsageHistoryByUser(userID)
	if err != nil {
		c.log.Error("Error fetching usage history", zap.Error(err))
		helper.ResponseError(ctx, "Failed to fetch usage history", err.Error(), http.StatusInternalServerError)
		return
	}
	helper.ResponseOK(ctx, gin.H{"usage_history": histories}, "Usage history fetched successfully", http.StatusOK)
}

func (c *VoucherController) GetUsersByVoucherCode(ctx *gin.Context) {
	voucherCode := ctx.Param("voucher_code")
	redeems, err := c.service.History.GetUsersByVoucherCode(voucherCode)
	if err != nil {
		c.log.Error("Error fetching users by voucher", zap.Error(err))
		helper.ResponseError(ctx, "Failed to fetch users", err.Error(), http.StatusInternalServerError)
		return
	}
	helper.ResponseOK(ctx, gin.H{"users": redeems}, "Users fetched successfully", http.StatusOK)
}
