package router

import (
	"voucher_system/infra"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoutes(ctx infra.ServiceContext) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// authMiddleware := ctx.Middleware.Authentication()
	jwtMiddleware := ctx.Middleware.JWTMiddleware()
	rateLimter := ctx.Middleware.RateLimiter()

	// allowedIPs := []string{"127.0.0.1", "192.168.1.100"}
	// r.Use(ctx.Middleware.IPWhitelistMiddleware(allowedIPs))
	
	r.POST("/login", rateLimter, ctx.Ctl.User.Login)
	r.POST("/register", ctx.Ctl.User.Register)
	
	router := r.Group("/vouchers", jwtMiddleware)
	{
		router.POST("/create", ctx.Ctl.Manage.CreateVoucher)
		router.DELETE("/:id", ctx.Ctl.Manage.SoftDeleteVoucher)
		router.PUT("/:id", ctx.Ctl.Manage.UpdateVoucher)
		router.GET("/redeem-points", ctx.Ctl.Manage.ShowRedeemPoints)
		router.GET("/", ctx.Ctl.Manage.GetVouchersByQueryParams)
		router.POST("/redeem", ctx.Ctl.Manage.CreateRedeemVoucher)
		router.GET("/:user_id", ctx.Ctl.Voucher.FindVouchers)
		router.GET("/:user_id/validate", ctx.Ctl.Voucher.ValidateVoucher)
		router.POST("/", ctx.Ctl.Voucher.UseVoucher)
		router.GET("/redeem-history/:user_id", ctx.Ctl.Voucher.GetRedeemHistoryByUser)
		router.GET("/usage-history/:user_id", ctx.Ctl.Voucher.GetUsageHistoryByUser)
		router.GET("/users-by-voucher/:voucher_code", ctx.Ctl.Voucher.GetUsersByVoucherCode)

	}

	return r
}
