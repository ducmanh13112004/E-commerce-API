package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/utils"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func localHandler(
	ctx *fiber.Ctx,
	parseBody func(*fiber.Ctx) (interface{}, error),
	serviceCall func(interface{}, string) (interface{}, *internal.SystemStatus),
) error {
	handlerName := utils.GetFunctionName()
	// Khởi tạo response mặc định
	resultFe := models.RespLocal{
		StatusCode: internal.SysStatus.WrongParams.Status,
		Message:    internal.SysStatus.WrongParams.Msg,
	}
	// Parse body để logging
	var rawBody interface{}
	ctx.BodyParser(&rawBody)

	// Lấy thông tin request
	uri := string(ctx.Request().URI().RequestURI())
	requestID := string(ctx.Request().Header.Peek("X-Request-Id"))
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := utils.GetTimeUTC7()
	// Xây dựng funcName
	funcName := handlerName + "-" + requestID
	internal.Log.Info("Start Handler", zap.Any("funcName", funcName), zap.Any("uri", uri), zap.Any("tokenAuth", tokenAuth), zap.Any("body", rawBody))
	defer func() {
		LocalSendKibana(ctx, funcName, uri, tokenAuth, rawBody, resultFe, startTime)
	}()

	// Parse body input đặc thù
	bodyData, err := parseBody(ctx)
	if err != nil {
		internal.Log.Error("Parse body failed", zap.Any("funcName", funcName), zap.Any("body", rawBody), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	// Gọi service và xử lý response
	result, errService := serviceCall(bodyData, funcName)
	if errService != nil {
		resultFe = models.RespLocal{
			StatusCode: errService.Status,
			Message:    errService.Msg,
			Data:       errService.Detail,
		}
	} else {
		resultFe = models.RespLocal{
			StatusCode: 0,
			Message:    "OK",
			Data:       result,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}
func webkitHandler(
	ctx *fiber.Ctx,
	parseBody func(*fiber.Ctx, string) (interface{}, error),
	validateInput func(interface{}, string) *internal.SystemStatus,
	serviceCall func(models.UserInfo, interface{}, string) (interface{}, *internal.SystemStatus),
) error {
	handlerName := utils.GetFunctionName()
	// Khởi tạo response mặc định
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	// Parse body để logging
	var rawBody interface{}
	ctx.BodyParser(&rawBody)

	// Lấy thông tin request
	var customer models.UserInfo
	uri := string(ctx.Request().URI().RequestURI())
	requestID := string(ctx.Request().Header.Peek("X-Request-Id"))
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	// Xây dựng funcName
	funcName := handlerName + "-" + requestID
	internal.Log.Info("Start Handler", zap.Any("funcName", funcName), zap.Any("uri", uri), zap.Any("tokenAuth", tokenAuth), zap.Any("body", rawBody))
	defer func() {
		WebkitSendKibana(ctx, funcName, uri, tokenAuth, rawBody, resultFe, startTime, customer)
	}()
	// Kiểm tra thông tin user
	customer, err := getUserInfo(ctx)
	if err != nil {
		internal.Log.Error("Customer invalid", zap.Any("funcName", funcName), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	funcName += "-" + customer.PhoneNb

	// Parse body input đặc thù
	bodyData, err := parseBody(ctx, funcName)
	if err != nil {
		internal.Log.Error("Parse body failed", zap.Any("funcName", funcName), zap.Any("body", rawBody), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	// Validate input
	errValidate := validateInput(bodyData, funcName)
	if errValidate != nil {
		resultFe.Status = errValidate.Status
		resultFe.Msg = errValidate.Msg
		resultFe.Detail = errValidate.Detail
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	// Gọi service và xử lý response
	result, errService := serviceCall(customer, bodyData, funcName)
	if errService != nil {
		resultFe = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultFe = models.RespWeb{
			Status: 1,
			Msg:    "OK",
			Detail: result,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

// Hàm helper lấy thông tin user
func getUserInfo(ctx *fiber.Ctx) (models.UserInfo, error) {
	customerID, ok1 := ctx.Locals("customer_id").(string)
	phone, ok2 := ctx.Locals("customer_phone").(string)
	appVer, ok3 := ctx.Locals("app_version").(string)
	token, ok4 := ctx.Locals("access_token").(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return models.UserInfo{}, errors.New("invalid user info")
	}

	return models.UserInfo{
		CustomerId:  customerID,
		PhoneNb:     phone,
		AppVersion:  appVer,
		AccessToken: token,
	}, nil
	// userInfo, ok := ctx.Locals("user_info").(models.UserInfo)
	// if !ok {
	// 	return models.UserInfo{}, errors.New("invalid user info")
	// }
	// return userInfo, nil
}
func portalHandler(
	ctx *fiber.Ctx,
	parseBody func(*fiber.Ctx, string) (interface{}, error),
	validateInput func(interface{}, string) *internal.SystemStatus,
	serviceCall func(models.PortalPayload, interface{}, string) (interface{}, *internal.SystemStatus),
) error {
	handlerName := utils.GetFunctionName()
	// Khởi tạo response mặc định
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	// Parse body để logging
	var rawBody interface{}
	ctx.BodyParser(&rawBody)

	// Kiểm tra thông tin user

	uri := string(ctx.Request().URI().RequestURI())
	requestID := string(ctx.Request().Header.Peek("X-Request-Id"))
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	// Xây dựng funcName
	funcName := handlerName + "-" + requestID
	userInfo, ok := ctx.Locals("user_info").(models.PortalPayload)
	if !ok {
		internal.Log.Error("Locals user_info type assertion failed", zap.Any("funcName", funcName))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	internal.Log.Info("Start Handler", zap.Any("funcName", funcName), zap.Any("uri", uri), zap.Any("tokenAuth", tokenAuth), zap.Any("body", rawBody))
	defer func() {
		PortalSendKibana(ctx, funcName, uri, tokenAuth, rawBody, resultFe, startTime, userInfo)
	}()

	funcName += "-" + userInfo.Email

	// Parse body input đặc thù
	bodyData, err := parseBody(ctx, funcName)
	if err != nil {
		internal.Log.Error("Parse body failed", zap.Any("funcName", funcName), zap.Any("body", rawBody), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	// Validate input
	errValidate := validateInput(bodyData, funcName)
	if errValidate != nil {
		resultFe.Status = errValidate.Status
		resultFe.Msg = errValidate.Msg
		resultFe.Detail = errValidate.Detail
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	// Gọi service và xử lý response
	result, errService := serviceCall(userInfo, bodyData, funcName)
	if errService != nil {
		resultFe = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultFe = models.RespWeb{
			Status: 1,
			Msg:    "OK",
			Detail: result,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}
