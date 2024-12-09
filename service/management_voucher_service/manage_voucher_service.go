package managementvoucherservice

import (
	"voucher_system/models"
	"voucher_system/repository"
	managementvoucher "voucher_system/repository/management_voucher"

	"go.uber.org/zap"
)

type ManageVoucherService interface {
	CreateVoucher(voucher *models.Voucher) error
	SoftDeleteVoucher(voucherID int) error
	UpdateVoucher(voucher *models.Voucher, voucherID int) error
	ShowRedeemPoints() (*[]managementvoucher.RedeemPoint, error)
	GetVouchersByQueryParams(status, area, voucher_type string) (*[]models.Voucher, error)
	CreateRedeemVoucher(redeem *models.Redeem, points int) error
}

type ManagementVoucherservice struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewManagementVoucherService(repo repository.Repository, log *zap.Logger) ManageVoucherService {
	return &ManagementVoucherservice{repo: repo, log: log}
}

func (ms *ManagementVoucherservice) CreateVoucher(voucher *models.Voucher) error {

	if err := ms.repo.Manage.CreateVoucher(voucher); err != nil {
		ms.log.Error("Error from service creating voucher: " + err.Error())
		return err
	}
	return nil
}

func (ms *ManagementVoucherservice) SoftDeleteVoucher(voucherID int) error {

	if err := ms.repo.Manage.SoftDeleteVoucher(voucherID); err != nil {
		ms.log.Error("Error from service soft-deletes: " + err.Error())
		return err
	}

	return nil
}

func (ms *ManagementVoucherservice) UpdateVoucher(voucher *models.Voucher, voucherID int) error {

	if err := ms.repo.Manage.UpdateVoucher(voucher, voucherID); err != nil {
		ms.log.Error("Error from service Update Voucher: " + err.Error())
		return err
	}

	return nil
}

func (ms *ManagementVoucherservice) ShowRedeemPoints() (*[]managementvoucher.RedeemPoint, error) {

	vouchers, err := ms.repo.Manage.ShowRedeemPoints()
	if err != nil {
		ms.log.Error("Error from service Show redeem points: " + err.Error())
		return nil, err
	}

	return vouchers, nil
}

func (ms *ManagementVoucherservice) GetVouchersByQueryParams(status, area, voucher_type string) (*[]models.Voucher, error) {

	vouchers, err := ms.repo.Manage.GetVouchersByQueryParams(status, area, voucher_type)
	if err != nil {
		ms.log.Error("Error from service Show redeem points: " + err.Error())
		return nil, err
	}

	return vouchers, nil
}

func (ms *ManagementVoucherservice) CreateRedeemVoucher(redeem *models.Redeem, points int) error {

	err := ms.repo.Manage.CreateRedeemVoucher(redeem, points)
	if err != nil {
		ms.log.Error("Error from service create redeem voucher: " + err.Error())
		return err
	}

	return nil
}
