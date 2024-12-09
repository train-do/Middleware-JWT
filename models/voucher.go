package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Voucher struct {
	ID              int             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" swaggerignore:"true"`
	VoucherName     string          `gorm:"type:varchar(255);not null" json:"voucher_name,omitempty" binding:"required" example:"PROMO GAJIAN"`
	VoucherCode     string          `gorm:"type:varchar(50);unique;not null" json:"voucher_code,omitempty" binding:"required" example:"DESCERIA100"`
	VoucherType     string          `gorm:"type:varchar(20);not null;check:voucher_type in ('e-commerce', 'redeem points')" json:"voucher_type,omitempty" binding:"required" example:"redeem points"`
	PointsRequired  int             `gorm:"default:0" json:"points_required,omitempty" example:"220"`
	Description     string          `gorm:"type:text;not null" json:"description,omitempty" example:"10% off for purchases above 200.000"`
	VoucherCategory string          `gorm:"type:varchar(20);not null;check:voucher_category in ('Free Shipping', 'Discount')" json:"voucher_category,omitempty" binding:"required" example:"Free Shipping"`
	DiscountValue   float64         `gorm:"type:numeric(10,2);not null" json:"discount_value,omitempty" binding:"required" example:"10.0"`
	MinimumPurchase float64         `gorm:"type:numeric(10,2);default:0" json:"minimum_purchase,omitempty" binding:"required" example:"200000.0"`
	PaymentMethods  []string        `gorm:"type:jsonb" json:"payment_methods,omitempty" binding:"required" swaggertype:"array,string" example:"Credit Card"`
	StartDate       time.Time       `gorm:"type:timestamp with time zone;not null" json:"start_date,omitempty" binding:"required" example:"2024-12-01T00:00:00Z"`
	EndDate         time.Time       `gorm:"type:timestamp with time zone;not null" json:"end_date,omitempty" binding:"required" example:"2024-12-07T00:00:00Z"`
	ApplicableAreas []string        `gorm:"type:jsonb" json:"applicable_areas,omitempty" binding:"required" swaggertype:"array,string" example:"Jawa"`
	Quota           int             `gorm:"default:0" json:"quota,omitempty" binding:"required" example:"50"`
	Status          bool            `gorm:"type:boolean" json:"status,omitempty" example:"true"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (v *Voucher) BeforeSave(tx *gorm.DB) (err error) {
	currentDate := time.Now()
	v.Status = currentDate.After(v.StartDate) && currentDate.Before(v.EndDate)
	// Marshal PaymentMethods and ApplicableAreas to JSON before saving
	if len(v.PaymentMethods) > 0 {
		// Marshal to JSON to ensure proper formatting
		paymentMethodsJSON, err := json.Marshal(v.PaymentMethods)
		if err != nil {
			return err
		}
		v.PaymentMethods = nil // Clear the original array to use the marshaled value
		v.PaymentMethods = []string{string(paymentMethodsJSON)}
	}

	if len(v.ApplicableAreas) > 0 {
		// Marshal to JSON to ensure proper formatting
		applicableAreasJSON, err := json.Marshal(v.ApplicableAreas)
		if err != nil {
			return err
		}
		v.ApplicableAreas = nil // Clear the original array to use the marshaled value
		v.ApplicableAreas = []string{string(applicableAreasJSON)}
	}

	return nil
}

type Redeem struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"not null" json:"user_id"`
	VoucherID  int       `gorm:"not null" json:"voucher_id"`
	RedeemDate time.Time `gorm:"default:current_date" json:"redeem_date"`
	User       User      `gorm:"foreignKey:UserID;references:ID" swaggerignore:"true" json:"-"`
	Voucher    Voucher   `gorm:"foreignKey:VoucherID;references:ID" swaggerignore:"true" json:"-"`
}

type History struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int       `gorm:"not null" json:"user_id"`
	VoucherID         int       `gorm:"not null" json:"voucher_id"`
	UsageDate         time.Time `gorm:"default:current_date" json:"usage_date"`
	TransactionAmount float64   `gorm:"type:numeric(10,2);not null" json:"transaction_amount"`
	BenefitValue      float64   `gorm:"type:numeric(10,2);not null" json:"benefit_value"`
	User              User      `gorm:"foreignKey:UserID;references:ID" swaggerignore:"true" json:"-"`
	Voucher           Voucher   `gorm:"foreignKey:VoucherID;references:ID" swaggerignore:"true" json:"-"`
}

func VoucherSeed() []Voucher {
	return []Voucher{
		{
			VoucherName:     "10% Discount",
			VoucherCode:     "DISCOUNT10",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "10% off for purchases above $100",
			VoucherCategory: "Discount",
			DiscountValue:   10.0,
			MinimumPurchase: 100.0,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now().AddDate(0, 0, -5), // StartDate 5 days ago
			EndDate:         time.Now().AddDate(0, 0, -1), // EndDate 1 day ago
			ApplicableAreas: []string{"US", "Canada"},
			Quota:           100,
		},
		{
			VoucherName:     "Free Shipping",
			VoucherCode:     "FREESHIP50",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "Free shipping for orders above $50",
			VoucherCategory: "Free Shipping",
			DiscountValue:   0.0,
			MinimumPurchase: 50.0,
			PaymentMethods:  []string{"All"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 2, 0), // 2 months valid
			ApplicableAreas: []string{"Worldwide"},
			Quota:           200,
		},
		{
			VoucherName:     "Redeem 500 Points",
			VoucherCode:     "POINTS500",
			VoucherType:     "redeem points",
			PointsRequired:  500,
			Description:     "Redeem 500 points for a $20 discount",
			VoucherCategory: "Discount",
			DiscountValue:   20.0,
			MinimumPurchase: 0.0,
			PaymentMethods:  []string{"Credit Card"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 3, 0), // 3 months valid
			ApplicableAreas: []string{"US"},
			Quota:           150,
		},
		{
			VoucherName:     "5% Discount",
			VoucherCode:     "DISCOUNT5",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "5% discount on all purchases",
			VoucherCategory: "Discount",
			DiscountValue:   5.0,
			MinimumPurchase: 0.0,
			PaymentMethods:  []string{"PayPal"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 1, 0), // 1 month valid
			ApplicableAreas: []string{"Europe"},
			Quota:           500,
		},
		{
			VoucherName:     "Black Friday Sale",
			VoucherCode:     "BLACKFRIDAY",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "20% off for Black Friday",
			VoucherCategory: "Discount",
			DiscountValue:   20.0,
			MinimumPurchase: 200.0,
			PaymentMethods:  []string{"Credit Card", "Bank Transfer"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 0, 7), // 1 week valid
			ApplicableAreas: []string{"Worldwide"},
			Quota:           300,
		},
		{
			VoucherName:     "Holiday Free Shipping",
			VoucherCode:     "HOLIDAYSHIP",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "Free shipping during the holiday season",
			VoucherCategory: "Free Shipping",
			DiscountValue:   0.0,
			MinimumPurchase: 75.0,
			PaymentMethods:  []string{"All"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 1, 0), // 1 month valid
			ApplicableAreas: []string{"US", "Canada"},
			Quota:           400,
		},
		{
			VoucherName:     "Cyber Monday Special",
			VoucherCode:     "CYBERMON",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "15% off for Cyber Monday",
			VoucherCategory: "Discount",
			DiscountValue:   15.0,
			MinimumPurchase: 150.0,
			PaymentMethods:  []string{"Credit Card"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 0, 5), // 5 days valid
			ApplicableAreas: []string{"Worldwide"},
			Quota:           100,
		},
		{
			VoucherName:     "Student Discount",
			VoucherCode:     "STUDENT15",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "15% discount for students",
			VoucherCategory: "Discount",
			DiscountValue:   15.0,
			MinimumPurchase: 0.0,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 2, 0), // 2 months valid
			ApplicableAreas: []string{"Europe"},
			Quota:           200,
		},
		{
			VoucherName:     "New Year Sale",
			VoucherCode:     "NEWYEAR50",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "Flat $50 off for the New Year sale",
			VoucherCategory: "Discount",
			DiscountValue:   50.0,
			MinimumPurchase: 300.0,
			PaymentMethods:  []string{"All"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 1, 0), // 1 month valid
			ApplicableAreas: []string{"US"},
			Quota:           150,
		},
		{
			VoucherName:     "Valentine's Free Shipping",
			VoucherCode:     "VALSHIP",
			VoucherType:     "e-commerce",
			PointsRequired:  0,
			Description:     "Free shipping for Valentine's Day",
			VoucherCategory: "Free Shipping",
			DiscountValue:   0.0,
			MinimumPurchase: 100.0,
			PaymentMethods:  []string{"Credit Card", "PayPal"},
			StartDate:       time.Now(),
			EndDate:         time.Now().AddDate(0, 1, 14), // 1 month 14 days valid
			ApplicableAreas: []string{"Worldwide"},
			Quota:           300,
		},
	}
}
