package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"
	"ecom_promotion_v2/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoginHandlers interface {
	LoginHandler(*fiber.Ctx) error
}

type loginHandlers struct {
	svc *services.AppServices
}

func NewLoginHandlers(svc *services.AppServices) LoginHandlers {
	return &loginHandlers{svc: svc}
}

func (h *loginHandlers) LoginHandler(ctx *fiber.Ctx) error {
	funcName := "LoginHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	// Mặc định response lỗi
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Parse body
	var body struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		internal.Log.Error("Cannot parse body", zap.Error(err))
		return ctx.Status(http.StatusOK).JSON(resultFe)
	}

	// Validate input
	if err := models.Validate.Struct(body); err != nil {
		internal.Log.Error("Validate error", zap.Any("input", body), zap.Error(err))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}

	// Gọi service Login
	token, err := h.svc.LoginService.Login(body.Username, funcName, body.Password)
	if err != nil {
		resultFe = models.RespWeb{
			Status: internal.SysStatus.WrongParams.Status,
			Msg:    internal.SysStatus.WrongParams.Msg,
			Detail: err.Error(),
		}
		return ctx.Status(http.StatusOK).JSON(resultFe)
	}

	// Trả về token
	resultFe = models.RespWeb{
		Status: 1,
		Msg:    "Login thành công",
		Detail: map[string]interface{}{
			"token": token,
		},
	}

	return ctx.Status(http.StatusOK).JSON(resultFe)
}
