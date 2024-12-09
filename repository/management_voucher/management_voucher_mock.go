package managementvoucher

import (
	"voucher_system/models"

	"github.com/stretchr/testify/mock"
)

type ManagementVoucherRepoMock struct {
	mock.Mock
}

func (m *ManagementVoucherRepoMock) CreateVoucher(voucher *models.Voucher) error {

	args := m.Called(voucher)
	return args.Error(0)
}

func (m *ManagementVoucherRepoMock) SoftDeleteVoucher(voucherID int) error {
	args := m.Called(voucherID)
	return args.Error(0)
}

func (m *ManagementVoucherRepoMock) UpdateVoucher(voucher *models.Voucher, voucherID int) error {
	args := m.Called(voucher, voucherID)
	return args.Error(0)
}

func (m *ManagementVoucherRepoMock) ShowRedeemPoints() (*[]RedeemPoint, error) {
	args := m.Called()
	if points := args.Get(0); points != nil {
		return points.(*[]RedeemPoint), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ManagementVoucherRepoMock) GetVouchersByQueryParams(status, area, voucher_type string) (*[]models.Voucher, error) {
	args := m.Called(status, area, voucher_type)
	if vouchers := args.Get(0); vouchers != nil {
		return vouchers.(*[]models.Voucher), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ManagementVoucherRepoMock) CreateRedeemVoucher(redeem *models.Redeem, points int) error {
	args := m.Called(redeem, points)
	return args.Error(0)
}
