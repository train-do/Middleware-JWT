package middleware

import (
	"net/http"
	"time"
	"voucher_system/database"
	"voucher_system/helper"
	"voucher_system/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

type Middleware struct {
	log    *zap.Logger
	Cacher database.Cacher
}

func NewMiddleware(log *zap.Logger, cacher database.Cacher) Middleware {
	return Middleware{
		log:    log,
		Cacher: cacher,
	}
}

func (m *Middleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userID := c.GetHeader("User-ID")

		if token == "" || userID == "" {
			m.log.Warn("Authentication failed", zap.String("userID", userID), zap.String("token", token))
			helper.ResponseError(c, "Token and User-ID are required", "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		m.log.Info("Authenticating user", zap.String("userID", userID), zap.String("token", token))

		storedToken, err := m.Cacher.Get(userID)
		if err != nil {
			m.log.Error("Failed to retrieve token from cache", zap.Error(err))
			helper.ResponseError(c, "Failed to retrieve token", "Server error", http.StatusInternalServerError)
			c.Abort()
			return
		}

		if storedToken == "" || storedToken != token {
			m.log.Warn("Invalid token", zap.String("userID", userID), zap.String("storedToken", storedToken))
			helper.ResponseError(c, "Invalid token", "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		m.log.Info("Authentication successful", zap.String("userID", userID))
		c.Next()
	}
}

func (m *Middleware) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			m.log.Warn("Missing JWT token")
			helper.ResponseError(c, "Missing token", "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtKey, nil
		})

		if err != nil {
			m.log.Warn("Error parsing token", zap.Error(err), zap.String("token", tokenString))
			helper.ResponseError(c, "Invalid or expired token", "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !token.Valid {
			m.log.Warn("Token is invalid", zap.String("token", tokenString))
			helper.ResponseError(c, "Invalid or expired token", "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}
		if claims.ExpiresAt.Time.Before(time.Now()) {
			m.log.Warn("Token has expired", zap.String("token", tokenString))
			helper.ResponseError(c, "Token has expired", "Unauthorized", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// if err != nil || !token.Valid {
		// 	m.log.Warn("Invalid or expired JWT token", zap.Error(err))
		// 	m.log.Error("Invalid or expired JWT token", zap.String("token: ", tokenString), zap.Any("token parse: ", token))
		// 	helper.ResponseError(c, "Invalid or expired token", "Unauthorized", http.StatusUnauthorized)
		// 	c.Abort()
		// 	return
		// }

		m.log.Info("JWT token valid", zap.String("userID", claims.Subject))
		c.Set("userID", claims.Subject)
		c.Next()
	}
}

// RateLimiter middleware with logging and helper response
func (m *Middleware) RateLimiter() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  3, // Maksimal 3 kegagalan dalam 1 menit
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	blockedIPs := make(map[string]time.Time) // Untuk menyimpan IP yang diblokir

	return func(c *gin.Context) {
		ip := c.ClientIP()
		m.log.Info("Checking login rate limit", zap.String("clientIP", ip))

		// Cek apakah IP diblokir
		if unblockTime, exists := blockedIPs[ip]; exists {
			if time.Now().Before(unblockTime) {
				m.log.Warn("Blocked IP tried to access", zap.String("clientIP", ip))
				helper.ResponseError(c, "Your IP is temporarily blocked due to multiple failed login attempts. Try again later.", "Forbidden", http.StatusForbidden)
				c.Abort()
				return
			}
			// Hapus dari blok jika waktu sudah habis
			delete(blockedIPs, ip)
		}

		// Dapatkan konteks pembatasan
		context, err := instance.Get(c, ip)
		if err != nil {
			m.log.Error("Rate limiter error", zap.Error(err))
			helper.ResponseError(c, "Rate limiter error", "Server error", http.StatusInternalServerError)
			c.Abort()
			return
		}

		if context.Reached {
			m.log.Warn("Login rate limit reached", zap.String("clientIP", ip))
			// Blokir IP selama 5 menit
			blockedIPs[ip] = time.Now().Add(5 * time.Minute)
			helper.ResponseError(c, "Too many failed login attempts. Your IP is now blocked for 5 minutes.", "Too Many Requests", http.StatusTooManyRequests)
			c.Abort()
			return
		}

		c.Next()

		// Cek apakah login gagal
		if status, ok := c.Get("login_failed"); ok && status.(bool) {
			m.log.Warn("Login failed", zap.String("clientIP", ip))
			// Tambahkan hitungan ke limiter
			_, _ = instance.Peek(c, ip) // Update limiter state secara manual
		}
	}
}

func (m *Middleware) IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	allowed := make(map[string]bool)
	for _, ip := range allowedIPs {
		allowed[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		m.log.Info("Checking IP whitelist", zap.String("clientIP", clientIP))

		if !allowed[clientIP] {
			m.log.Warn("Access denied for IP", zap.String("clientIP", clientIP))
			helper.ResponseError(c, "Access denied", "Forbidden", http.StatusForbidden)
			c.Abort()
			return
		}

		m.log.Info("IP allowed", zap.String("clientIP", clientIP))
		c.Next()
	}
}
