package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Định nghĩa interface và struct cho handler
type ProgramRedirectChannelHandlers interface {
	// CreateProgramDirect(*fiber.Ctx) error
	CreateProgramRedirectHandler(ctx *fiber.Ctx) error
	GetProgramRedirectListAllHandler(ctx *fiber.Ctx) error
	GetProgramRedirectIDHandler(ctx *fiber.Ctx) error
	UpdateProgramRedirectHandler(ctx *fiber.Ctx) error
	DeleteProgramRedirectHandler(ctx *fiber.Ctx) error
	GetProgramRedirectListNameAllHandler(ctx *fiber.Ctx) error
	GetProgramListNavigationHandler(ctx *fiber.Ctx) error
	CreateProgramNavigationHandler(ctx *fiber.Ctx) error
	UpdateProgramNavigationHandler(ctx *fiber.Ctx) error
}

type programredirectHandlers struct {
	svc *services.AppServices
}

func NewProgramRedirectHandlers(
	appService *services.AppServices,
) ProgramRedirectChannelHandlers {
	return &programredirectHandlers{
		svc: appService,
	}
}

func (h *programredirectHandlers) GetProgramRedirectIDHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// 1. Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := new(models.InputGetPromotion)
			return body, ctx.BodyParser(body)
		},

		// 2. Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*models.InputGetPromotion)
			if err := models.Validate.Struct(bodyData); err != nil {
				internal.Log.Error(funcName, zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		},

		// 3. Gọi service
		func(customer models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*models.InputGetPromotion)
			result, errService := h.svc.ProgramRedirectService.GetProgramRedirectIDService(funcName, bodyData.Data.RedirectId)
			if errService != nil {
				return nil, errService
			}
			return map[string]interface{}{
				"direct_program_ID": result,
			}, nil
		},
	)
}

func (h *programredirectHandlers) GetProgramRedirectListAllHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// 1. Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := &struct {
				Data models.FilterProgramRedirectTb `json:"data"`
			}{}
			return body, ctx.BodyParser(body)
		},

		// 2. Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*struct {
				Data models.FilterProgramRedirectTb `json:"data"`
			})
			if err := models.Validate.Struct(bodyData); err != nil {
				internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		},

		// 3. Gọi service
		func(customer models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*struct {
				Data models.FilterProgramRedirectTb `json:"data"`
			})
			result, errService := h.svc.ProgramRedirectService.GetProgramRedirectListAllService(funcName, &bodyData.Data)
			if errService != nil {
				return nil, errService
			}
			return map[string]interface{}{
				"list_direct_links": result,
			}, nil
		},
	)
}

func (h *programredirectHandlers) DeleteProgramRedirectHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// 1. Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := &struct {
				Data models.DeleteProgramRedirectRequest `json:"data"`
			}{}
			return body, ctx.BodyParser(body)
		},

		// 2. Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*struct {
				Data models.DeleteProgramRedirectRequest `json:"data"`
			})
			if err := models.Validate.Struct(bodyData.Data); err != nil {
				internal.Log.Error("Validation error", zap.Any("input", bodyData), zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		},

		// 3. Gọi service
		func(user models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*struct {
				Data models.DeleteProgramRedirectRequest `json:"data"`
			})

			updateBy := user.Email

			errService := h.svc.ProgramRedirectService.DeleteProgramRedirect(funcName, &bodyData.Data, updateBy)
			if errService != nil {
				return nil, errService
			}
			return map[string]interface{}{
				"message": "OK",
			}, nil
		},
	)
}

func (h *programredirectHandlers) UpdateProgramRedirectHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// 1. Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := new(models.InputProgramRedirect)
			return body, ctx.BodyParser(body)
		},

		// 2. Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*models.InputProgramRedirect)
			if err := models.Validate.Struct(bodyData); err != nil {
				internal.Log.Error("Validation error", zap.Any("input", bodyData), zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		},

		// 3. Gọi service
		func(user models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*models.InputProgramRedirect)

			bodyData.Email = user.Email

			errService := h.svc.ProgramRedirectService.UpdateProgramRedirect(funcName, bodyData, bodyData.Email)
			if errService != nil {
				return nil, errService
			}

			return map[string]interface{}{
				"message": "OK",
			}, nil
		},
	)
}

func (h *programredirectHandlers) CreateProgramRedirectHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) { // Parse body
			body := new(models.InputProgramRedirect)
			return body, ctx.BodyParser(body)
		}, func(input interface{}, funcName string) *internal.SystemStatus { // Validate input
			bodyData := input.(*models.InputProgramRedirect)
			if err := models.Validate.Struct(bodyData); err != nil {
				internal.Log.Error(funcName, zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		}, func(user models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) { // Gọi service
			errService := h.svc.ProgramRedirectService.CreateProgramRedirect(funcName, input.(*models.InputProgramRedirect), user.Email)
			return nil, errService
		})
}
func (h *programredirectHandlers) GetProgramRedirectListNameAllHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			return nil, nil
		},
		func(input interface{}, funcName string) *internal.SystemStatus {
			return nil
		},
		func(customer models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			resp, errService := h.svc.ProgramRedirectService.GetProgramRedirectListNameAllService(funcName)
			if errService != nil {
				return nil, errService
			}
			return map[string]interface{}{
				"list_direct_links": resp,
			}, nil
		})
}

func (h *programredirectHandlers) GetProgramListNavigationHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// 1. Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := &struct {
				Code string `json:"code" validate:"required"`
			}{}
			return body, ctx.BodyParser(body)
		},

		// 2. Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*struct {
				Code string `json:"code" validate:"required"`
			})
			if err := models.Validate.Struct(bodyData); err != nil {
				internal.Log.Error("Validation error", zap.Any("input", bodyData), zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		},

		// 3. Call service
		func(customer models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*struct {
				Code string `json:"code" validate:"required"`
			})

			result, errService := h.svc.ProgramRedirectService.GetProgramListNavigationService(funcName, bodyData.Code)
			if errService != nil {
				return nil, errService
			}

			return map[string]interface{}{
				"navigation": result,
			}, nil
		},
	)
}

func (h *programredirectHandlers) CreateProgramNavigationHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// 1. Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := new(models.MasterDataInput)
			return body, ctx.BodyParser(body)
		},

		// 2. Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*models.MasterDataInput)
			if bodyData.Code == "" || bodyData.DataKeyValue == nil {
				internal.Log.Error("Validation error: missing required fields", zap.Any("input", bodyData))
				return internal.SysStatus.WrongParams
			}
			if err := models.Validate.Struct(bodyData); err != nil {
				internal.Log.Error("Validation error", zap.Any("input", bodyData), zap.Error(err))
				return internal.SysStatus.WrongParams
			}
			return nil
		},

		// 3. Call service
		func(user models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*models.MasterDataInput)

			bodyData.UpdateBy = user.Email

			result, errService := h.svc.ProgramRedirectService.CreateProgramNavigationService(funcName, bodyData)
			if errService != nil {
				return nil, errService
			}

			return map[string]interface{}{
				"message": "OK",
				"detail":  result,
			}, nil
		},
	)
}

func (h *programredirectHandlers) UpdateProgramNavigationHandler(ctx *fiber.Ctx) error {
	return portalHandler(ctx,
		// Parse body
		func(ctx *fiber.Ctx, funcName string) (interface{}, error) {
			body := new(models.MasterDataInput)
			return body, ctx.BodyParser(body)
		},
		// Validate input
		func(input interface{}, funcName string) *internal.SystemStatus {
			bodyData := input.(*models.MasterDataInput)
			if bodyData.Code == "" || bodyData.DataKeyValue == nil {
				internal.Log.Error(funcName, zap.String("error", "Missing required fields"))
				return internal.SysStatus.WrongParams
			}
			return nil
		},
		// Gọi service
		func(user models.PortalPayload, input interface{}, funcName string) (interface{}, *internal.SystemStatus) {
			bodyData := input.(*models.MasterDataInput)
			bodyData.UpdateBy = user.Email
			resp, errService := h.svc.ProgramRedirectService.UpdateProgramNavigationService(funcName, bodyData)
			return resp, errService
		},
	)
}
