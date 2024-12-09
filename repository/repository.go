package repository

import (
	managementvoucher "voucher_system/repository/management_voucher"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	User    UserRepository
	Manage  managementvoucher.ManagementVoucherInterface
	Voucher VoucherRepository
	Redeem  RedeemRepository
	History HistoryRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) Repository {
	return Repository{
		User:    NewUserRepository(db, log),
		Manage:  managementvoucher.NewManagementVoucherRepo(db, log),
		Voucher: NewVoucherRepository(db, log),
		Redeem:  NewRedeemRepository(db, log),
		History: NewHistoryRepository(db, log),
	}
}
