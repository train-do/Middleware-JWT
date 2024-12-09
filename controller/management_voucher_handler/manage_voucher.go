package managementvoucherhandler

import (
	"net/http"
	"strconv"
	"voucher_system/helper"
	"voucher_system/models"
	"voucher_system/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ManageVoucherHandler interface {
	CreateVoucher(c *gin.Context)
	SoftDeleteVoucher(c *gin.Context)
	UpdateVoucher(c *gin.Context)
	ShowRedeemPoints(c *gin.Context)
	GetVouchersByQueryParams(c *gin.Context)
	CreateRedeemVoucher(c *gin.Context)
}

type ManagementVoucherHandler struct {
	service service.Service
	log     *zap.Logger
}

func NewManagementVoucherHanlder(service service.Service, log *zap.Logger) ManageVoucherHandler {
	return &ManagementVoucherHandler{service: service, log: log}
}

// CreateVoucher godoc
// @Summary Create a new voucher
// @Description Create a new voucher with provided details
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param voucher body models.Voucher true "Voucher details"
// @Success 200 {object} utils.ResponseOK{data=models.Voucher} "Created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid payload"
// @Failure 500 {object} utils.ErrorResponse "Failed to create voucher"
// @Security Authentication
// @Security UserID
// @Router /vouchers/create [post]
func (mh *ManagementVoucherHandler) CreateVoucher(c *gin.Context) {

	var voucher models.Voucher

	err := c.ShouldBindJSON(&voucher)
	if err != nil {
		mh.log.Error("Invalid payload", zap.Error(err))
		helper.ResponseError(c, "INVALID", "Invalid Payload"+err.Error(), http.StatusInternalServerError)
		return
	}

	err = mh.service.Manage.CreateVoucher(&voucher)
	if err != nil {
		mh.log.Error("Failed to create", zap.Error(err))
		helper.ResponseError(c, "FAILED", "Failed to create Voucher", http.StatusBadRequest)
		return
	}

	mh.log.Info("Create Voucher successfully")
	helper.ResponseOK(c, voucher, "Created successfully", http.StatusOK)
}

// SoftDeleteVoucher godoc
// @Summary Soft delete a voucher
// @Description Soft delete a voucher by ID
// @Tags Vouchers
// @Param id path int true "Voucher ID"
// @Success 200 {object} utils.ResponseOK "Deleted successfully"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete voucher"
// @Security Authentication
// @Security UserID
// @Router /vouchers/{id} [delete]
func (mh *ManagementVoucherHandler) SoftDeleteVoucher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := mh.service.Manage.SoftDeleteVoucher(id)
	if err != nil {
		mh.log.Error("Failed to Deleted", zap.Error(err))
		helper.ResponseError(c, "FAILED", "Failed to deleted Voucher", http.StatusInternalServerError)
		return
	}

	mh.log.Info("Deleted Voucher successfully")
	helper.ResponseOK(c, id, "Deleted succesfully", http.StatusOK)
}

// UpdateVoucher godoc
// @Summary Update a voucher
// @Description Update a voucher by ID
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param id path int true "Voucher ID"
// @Param voucher body models.Voucher true "Updated voucher details"
// @Success 200 {object} utils.ResponseOK{data=models.Voucher} "Updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid payload"
// @Failure 500 {object} utils.ErrorResponse "Failed to update voucher"
// @Security Authentication
// @Security UserID
// @Router /vouchers/{id} [put]
func (mh *ManagementVoucherHandler) UpdateVoucher(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	voucher := models.Voucher{}

	if err := c.ShouldBindJSON(&voucher); err != nil {
		mh.log.Error("Invalid payload", zap.Error(err))
		helper.ResponseError(c, "INVALID", "Invalid Payload"+err.Error(), http.StatusBadRequest)
		return
	}

	err := mh.service.Manage.UpdateVoucher(&voucher, id)
	if err != nil {
		mh.log.Error("Failed to Updated Voucher", zap.Error(err))
		helper.ResponseError(c, "FAILED", "Failed to Updated Voucher", http.StatusInternalServerError)
		return
	}

	mh.log.Info("Updated Voucher successfully")
	helper.ResponseOK(c, id, "updated succesfully", http.StatusOK)
}

// ShowRedeemPoints godoc
// @Summary Show redeem points
// @Description Retrieve the list of redeem points
// @Tags Vouchers
// @Produce json
// @Success 200 {object} utils.ResponseOK{data=[]models.Redeem} "Redeem points retrieved successfully"
// @Failure 404 {object} utils.ErrorResponse "Redeem points not found"
// @Security Authentication
// @Security UserID
// @Router /vouchers/redeem-points [get]
func (mh *ManagementVoucherHandler) ShowRedeemPoints(c *gin.Context) {

	voucher, err := mh.service.Manage.ShowRedeemPoints()
	if err != nil {
		mh.log.Error("Failed to Get Reedem Points List", zap.Error(err))
		helper.ResponseError(c, "NOT FOUND", "Reedem Points List Not Found", http.StatusNotFound)
		return
	}

	mh.log.Info("Redeem points retrieved successfully")
	helper.ResponseOK(c, voucher, "Redeem points retrieved successfully", http.StatusOK)

}

// GetVouchersByQueryParams godoc
// @Summary Get vouchers by query parameters
// @Description Retrieve vouchers based on status, area, and voucher type
// @Tags Vouchers
// @Produce json
// @Param status query string false "Voucher status"
// @Param area query string false "Voucher area"
// @Param voucher_type query string false "Voucher type"
// @Success 200 {object} utils.ResponseOK{data=[]models.Voucher} "Vouchers retrieved successfully"
// @Failure 404 {object} utils.ErrorResponse "Vouchers not found"
// @Security Authentication
// @Security UserID
// @Router /vouchers [get]
func (mh *ManagementVoucherHandler) GetVouchersByQueryParams(c *gin.Context) {

	status := c.Query("status")
	area := c.Query("area")
	voucher_type := c.Query("voucher_type")

	voucher, err := mh.service.Manage.GetVouchersByQueryParams(status, area, voucher_type)
	if err != nil {
		mh.log.Error("Failed to Get Voucher List", zap.Error(err))
		helper.ResponseError(c, "NOT FOUND", "Voucher Not Found", http.StatusNotFound)
		return
	}

	mh.log.Info("Voucher retrieved successfully")
	helper.ResponseOK(c, voucher, "Voucher retrieved successfully", http.StatusOK)

}

type RedeemRequest struct {
	VoucherID int `json:"voucher_id" binding:"required"`
	UserID    int `json:"user_id" binding:"required"`
	Points    int `json:"points" binding:"required"`
}

// CreateRedeemVoucher godoc
// @Summary Create a redeem voucher
// @Description Redeem a voucher using points
// @Tags Vouchers
// @Accept json
// @Produce json
// @Param redeemRequest body RedeemRequest true "Redeem request payload"
// @Success 200 {object} utils.ResponseOK{data=models.Redeem} "Redeem created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid payload"
// @Failure 500 {object} utils.ErrorResponse "Failed to create redeem voucher"
// @Security Authentication
// @Security UserID
// @Router /vouchers/redeem [post]
func (mh *ManagementVoucherHandler) CreateRedeemVoucher(c *gin.Context) {
	var RedeemRequest RedeemRequest
	err := c.ShouldBindJSON(&RedeemRequest)
	if err != nil {
		mh.log.Error("Invalid payload", zap.Error(err))
		helper.ResponseError(c, "INVALID", "Invalid Payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	redeem := models.Redeem{
		VoucherID: RedeemRequest.VoucherID,
		UserID:    RedeemRequest.UserID,
	}

	err = mh.service.Manage.CreateRedeemVoucher(&redeem, RedeemRequest.Points)
	if err != nil {
		mh.log.Error("Failed to create redeem voucher", zap.Error(err))
		helper.ResponseError(c, "FAILED", "Failed to create redeem voucher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mh.log.Info("Create Redeem Voucher successfully")
	helper.ResponseOK(c, redeem, "Created successfully", http.StatusOK)
}
