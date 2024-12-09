package controller

import (
	managementvoucherhandler "voucher_system/controller/management_voucher_handler"
	"voucher_system/database"
	"voucher_system/service"

	"go.uber.org/zap"
)

type Controller struct {
	User    AuthController
	Manage  managementvoucherhandler.ManageVoucherHandler
	Voucher VoucherController
}

func NewController(service service.Service, logger *zap.Logger, cacher database.Cacher) *Controller {
	return &Controller{
		User:    NewAuthController(service, logger, cacher),
		Manage:  managementvoucherhandler.NewManagementVoucherHanlder(service, logger),
		Voucher: *NewVoucherController(service, logger),
	}
}
