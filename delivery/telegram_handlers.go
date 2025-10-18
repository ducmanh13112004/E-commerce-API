package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TelegramHandlers interface {
	BotTelegramHandlers(ctx *fiber.Ctx) error
	BotTelegramLibraryHandlers(ctx *fiber.Ctx) error
}
type telegramHandlers struct {
}

func NewTelegramHandlers() TelegramHandlers {
	return &telegramHandlers{}
}

func (t *telegramHandlers) BotTelegramHandlers(ctx *fiber.Ctx) error {
	requestId := string(ctx.Request().Header.Peek("X-Request-Id")) // uuid
	funcName := "BotTelegramHandlers-" + requestId
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	var body interface{}
	ctx.BodyParser(&body)

	bodyData := &models.ForwardMessage{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Cannot parse",
			zap.String("func", funcName),
			zap.Any("body", body),
			zap.Error(err),
		)
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError",
			zap.String("func", funcName),
			zap.Any("input", bodyData),
			zap.Error(validateError),
		)
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}

	// err = utils_call.CallSendMessageTelegram(bodyData.Message)
	if err != nil {
		internal.Log.Error("Send telegram message failed",
			zap.String("func", funcName),
			zap.String("service", bodyData.ServiceName),
			zap.Error(err),
		)
		resultFe.Msg = "Gửi tin nhắn thất bại"
		resultFe.Status = internal.SysStatus.SystemError.Status
		return ctx.Status(http.StatusOK).JSON(resultFe)
	}

	internal.Log.Info("Send telegram message success",
		zap.String("func", funcName),
		zap.String("service", bodyData.ServiceName),
	)

	resultFe.Msg = "Thành công"
	resultFe.Status = 1
	return ctx.Status(http.StatusOK).JSON(resultFe)
}

func (t *telegramHandlers) BotTelegramLibraryHandlers(ctx *fiber.Ctx) error {
	funcName := "BotTelegramLibraryHandlers"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	// funcName := "GetPromotionListHandler"
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)
	bodyData := &models.ForwardMessage{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Cannot parse", zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError", zap.Any("input", bodyData), zap.Error(validateError))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}
	_, err = time.Parse("2006-01-02 15:04:05", bodyData.Time)
	if err != nil {
		resultFe.Msg = "Định dạng ngày tháng năm không hợp lệ"
		resultFe.Status = internal.SysStatus.SystemError.Status
		return ctx.Status(http.StatusOK).JSON(resultFe)
	}
	// err = utils.ForwardMessage(bodyData)
	if err != nil {
		resultFe.Msg = "Gửi tin nhắn thất bại"
		resultFe.Status = internal.SysStatus.SystemError.Status
		return ctx.Status(http.StatusOK).JSON(resultFe)
	}
	resultFe.Msg = "Thành công"
	resultFe.Status = 1
	return ctx.Status(http.StatusOK).JSON(resultFe)
}
