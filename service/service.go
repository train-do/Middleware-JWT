package service

import (
	"voucher_system/repository"
	managementvoucherservice "voucher_system/service/management_voucher_service"

	"go.uber.org/zap"
)

type Service struct {
	User    UserService
	Manage  managementvoucherservice.ManageVoucherService
	Voucher VoucherService
	History HistoryService
}

func NewService(repo repository.Repository, log *zap.Logger) Service {
	return Service{
		User:    NewUserService(repo, log),
		Manage:  managementvoucherservice.NewManagementVoucherService(repo, log),
		Voucher: NewVoucherService(repo, log),
		History: NewHistoryService(repo, log),
	}
}
