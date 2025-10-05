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

type ProgramProductHandlers interface {
	CreateListSkuHandler(ctx *fiber.Ctx) error
}

type programProductHandlers struct {
	svc *services.AppServices
}

func NewProgramProductHandlers(appService *services.AppServices) ProgramProductHandlers {
	return &programProductHandlers{
		svc: appService,
	}
}

func (p *programProductHandlers) CreateListSkuHandler(ctx *fiber.Ctx) error {
	funcName := "CreateListSkuHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)
	employee := ctx.Locals("payload").(models.PortalPayload)
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := utils.GetTimeUTC7()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, employee) }()
	type InputProductProgram struct {
		Data models.InsertListSku `json:"data"`
	}
	bodyData := &InputProductProgram{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError", zap.Any("input", bodyData), zap.Error(validateError))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}

	errService := p.svc.CreateListSku(&bodyData.Data)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}
	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}
