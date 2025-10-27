package delivery

import (
	"ecom_promotion_v2/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AppHandlers interface {
	// MiddleWare
	RequireTokenLocal(*fiber.Ctx) error
	RequireXKeyTelegramPortal(*fiber.Ctx) error
	// RequireTokenApp(ctx *fiber.Ctx) error
	PromotionsHandlers
	CategoryHandlers
	ProgramPromotionHandler
	ProgramCategoryHandlers
	ProgramProductHandlers
	ProgramDirectHandlers
	TelegramHandlers
	ProgramRedirectChannelHandlers
	LoginHandlers
}

type appHandlers struct {
	svc *services.AppServices
	PromotionsHandlers
	CategoryHandlers
	ProgramPromotionHandler
	ProgramCategoryHandlers
	ProgramProductHandlers
	ProgramDirectHandlers
	TelegramHandlers
	LoginHandlers
	ProgramRedirectChannelHandlers
}

func NewAppHandlers(appService *services.AppServices) AppHandlers {
	return &appHandlers{
		appService,
		NewPromotionsHandlers(appService),
		NewCategoryHandlers(appService),
		NewProgramPromotionHandlers(appService),
		NewProgramCategoryHandlers(appService),
		NewProgramProductHandlers(appService),
		NewProgramDirectHandlers(appService),
		NewTelegramHandlers(),
		NewLoginHandlers(appService),
		NewProgramRedirectHandlers(appService),
	}
}
