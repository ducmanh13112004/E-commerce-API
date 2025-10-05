package delivery

import (
	"ecom_promotion_v2/internal/services"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
)

type AppHandlers interface {
	// MiddleWare
	RateLimit(*fiber.Ctx) error
	CustomRateLimit(*fiber.Ctx) error
	RequireTokenWeb(*fiber.Ctx) error
	RequireTokenLocal(*fiber.Ctx) error
	RequireTokenClient(*fiber.Ctx) error
	RequireTokenPortal(*fiber.Ctx) error
	RequireXKeyTelegramPortal(*fiber.Ctx) error
	RequireTokenApp(ctx *fiber.Ctx) error
	Health(*fiber.Ctx) error
	FiberRateLimit(maxRequests int, duration time.Duration) fiber.Handler
	PromotionsHandlers
	CategoryHandlers
	ProgramPromotionHandler
	ProgramCategoryHandlers
	ProgramProductHandlers
	ProgramDirectHandlers
	TelegramHandlers
	ProgramRedirectChannelHandlers
}

type appHandlers struct {
	svc *services.AppServices
	rdb *redis.Client
	PromotionsHandlers
	CategoryHandlers
	ProgramPromotionHandler
	ProgramCategoryHandlers
	ProgramProductHandlers
	ProgramDirectHandlers
	TelegramHandlers

	ProgramRedirectChannelHandlers
}

func NewAppHandlers(appService *services.AppServices, rdb *redis.Client) AppHandlers {
	return &appHandlers{
		appService,
		rdb,
		NewPromotionsHandlers(appService),
		NewCategoryHandlers(appService),
		NewProgramPromotionHandlers(appService),
		NewProgramCategoryHandlers(appService),
		NewProgramProductHandlers(appService),
		NewProgramDirectHandlers(appService),
		NewTelegramHandlers(),
		NewProgramRedirectHandlers(appService),
	}
}

func (bk *appHandlers) Health(ctx *fiber.Ctx) error {
	statusCode, result := bk.svc.HealthCheckService.HealthCheck(ctx)
	return ctx.Status(statusCode).JSON(result)
}
