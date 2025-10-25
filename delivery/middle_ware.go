package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string, limit int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, limit)
		visitors[ip] = limiter
	}

	return limiter
}
func (s *appHandlers) RateLimit(ctx *fiber.Ctx) error {
	// customer_id := ctx.Locals("customer_id").(string)
	// limiter := getVisitor(customer_id, 10)
	// if !limiter.Allow() {
	// 	internal.Log.Info("Rate limit", zap.Any("Customer_id", customer_id))
	// 	return ctx.Status(fiber.StatusOK).JSON(models.RespWeb{
	// 		Status: internal.SysStatus.RateLimit.Status,
	// 		Msg:    internal.SysStatus.RateLimit.Msg,
	// 	})
	// }
	return ctx.Next()
}

func (s *appHandlers) CustomRateLimit(ctx *fiber.Ctx) error {
	customer_id := ctx.Locals("customer_id").(string)
	limiter := getVisitor(customer_id, 30)
	if !limiter.Allow() {
		internal.Log.Info("Rate limit", zap.Any("Customer_id", customer_id))
		return ctx.Status(fiber.StatusOK).JSON(models.RespWeb{
			Status: internal.SysStatus.RateLimit.Status,
			Msg:    internal.SysStatus.RateLimit.Msg,
		})
	}
	return ctx.Next()
}
func (s *appHandlers) FiberRateLimit(maxRequests int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: duration,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(models.RespWeb{
				Status: internal.SysStatus.RateLimit.Status,
				Msg:    internal.SysStatus.RateLimit.Msg,
			})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			customerID := c.Locals("customer_id").(string)

			if customerID == "" {
				// Fallback to IP or reject
				customerID = c.IP()
			}
			return "rate_limit_customer_" + customerID
		},
	})
}
func (s *appHandlers) RequireTokenLocal(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("TOKEN"))
	resultFe := models.RespLocal{}
	defer func() {
		internal.Log.Info("RequireTokenLocal", zap.Any("url", ctx.Context().URI()), zap.Any("authen", token), zap.Any("result", resultFe))
	}()
	if utils.IsEmpty(token) {
		resultFe = models.RespLocal{
			StatusCode: internal.SysStatus.TokenRequired.Status,
			Message:    internal.SysStatus.TokenRequired.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	getToken := utils.CreateEcomToken(internal.Keys.LOCAL_ECOM_CLIENT_KEY, internal.Keys.LOCAL_ECOM_SECRET_KEY)

	if token != getToken {
		detail := map[string]interface{}{}
		if !internal.Envs.IsProduction {
			detail["token"] = getToken
		}
		resultFe = models.RespLocal{
			StatusCode: internal.SysStatus.InvalidToken.Status,
			Message:    internal.SysStatus.InvalidToken.Msg,
			Data:       detail,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	return ctx.Next()
}

func (s *appHandlers) RequireXKeyTelegramPortal(ctx *fiber.Ctx) error {
	funcName := "RequireXKeyTelegramPortal"
	xApiKey := string(ctx.Request().Header.Peek("X-API-KEY"))

	internal.Log.Info(funcName, zap.Any("ip", ctx.IP()), zap.Any("url", ctx.Context().URI()), zap.Any("authen", xApiKey))

	// Kiểm tra nếu API Key bị thiếu
	if utils.IsEmpty(xApiKey) {
		internal.Log.Error(funcName, zap.String("error", "Missing API Key"), zap.String("client_ip", ctx.IP()))
		resultFe := models.RespWeb{
			Status: internal.CODE_TOKEN_REQUIRED,
			Msg:    "API Key is required",
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}

	// Kiểm tra nếu API Key không hợp lệ
	// if xApiKey != internal.Keys.X_API_KEY_TELEGRAM {
	// 	internal.Log.Error(funcName, zap.String("error", "Invalid API Key"), zap.String("provided_key", xApiKey), zap.String("client_ip", ctx.IP()))
	// 	resultFe := models.RespWeb{
	// 		Status: internal.CODE_INVALID_TOKEN,
	// 		Msg:    "Invalid API Key",
	// 	}
	// 	return ctx.Status(fiber.StatusUnauthorized).JSON(resultFe)
	// }

	// internal.Log.Info(funcName, zap.String("message", "API Key validated successfully"), zap.String("client_ip", ctx.IP()))

	return ctx.Next()
}
