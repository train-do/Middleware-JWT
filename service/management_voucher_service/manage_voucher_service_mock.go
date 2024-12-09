package managementvoucherservice

import (
	"voucher_system/models"
	managementvoucher "voucher_system/repository/management_voucher"

	"github.com/stretchr/testify/mock"
)

type ManagementVoucherServiceMock struct {
	mock.Mock
}

func (m *ManagementVoucherServiceMock) CreateVoucher(voucher *models.Voucher) error {
	args := m.Called(voucher)
	return args.Error(0)
}

func (m *ManagementVoucherServiceMock) SoftDeleteVoucher(voucherID int) error {
	args := m.Called(voucherID)
	return args.Error(0)
}

func (m *ManagementVoucherServiceMock) UpdateVoucher(voucher *models.Voucher, voucherID int) error {
	args := m.Called(voucher, voucherID)
	return args.Error(0)
}

func (m *ManagementVoucherServiceMock) ShowRedeemPoints() (*[]managementvoucher.RedeemPoint, error) {
	args := m.Called()
	if points, ok := args.Get(0).(*[]managementvoucher.RedeemPoint); ok {
		return points, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ManagementVoucherServiceMock) GetVouchersByQueryParams(status, area, voucher_type string) (*[]models.Voucher, error) {
	args := m.Called(status, area, voucher_type)
	if vouchers := args.Get(0); vouchers != nil {
		return vouchers.(*[]models.Voucher), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ManagementVoucherServiceMock) CreateRedeemVoucher(redeem *models.Redeem, points int) error {
	args := m.Called(redeem, points)
	return args.Error(0)
}
