package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"

	// "go/token"
	"ecom_promotion_v2/internal/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Định nghĩa interface và struct cho handler
type CategoryHandlers interface {
	CreateCategory(*fiber.Ctx) error
	UpdateCategory(*fiber.Ctx) error
	DeleteCategory(*fiber.Ctx) error
	GetCategoryListHandler(*fiber.Ctx) error
	GetCategoryListNameHandler(*fiber.Ctx) error
}

type categoryHandlers struct {
	svc *services.AppServices
}

func NewCategoryHandlers(
	appService *services.AppServices,
) CategoryHandlers {
	return &categoryHandlers{
		svc: appService,
	}
}
func (tk *categoryHandlers) CreateCategory(ctx *fiber.Ctx) error {
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	funcName := "CreateCategory-" + requestId
	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	var userPayload models.PortalPayload
	defer func() {
		PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultWeb, startTime, userPayload)
	}()
	type InputCategory struct {
		Data models.CategoryTb `json:"data"`
	}
	bodyData := &InputCategory{}
	if err := ctx.BodyParser(bodyData); err != nil {
		internal.Log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}
	var ok bool
	userPayload, ok = ctx.Locals("payload").(models.PortalPayload)
	if !ok {
		internal.Log.Error("Failed to get user payload from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}
	bodyData.Data.UpdateBy = userPayload.Email
	if validateError := models.Validate.Struct(bodyData.Data); validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}
	errService := tk.svc.CategoryService.CreateCategory(&bodyData.Data)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Any("error", errService))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	}
	resultWeb.Status = 1
	resultWeb.Msg = "Thành công"

	return ctx.Status(fiber.StatusOK).JSON(resultWeb)
}

func (tk *categoryHandlers) DeleteCategory(ctx *fiber.Ctx) error {
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	funcName := "DeleteCategory-" + requestId
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))
	var rawBody interface{}
	ctx.BodyParser(&rawBody)

	type InputDeleteCategory struct {
		Data models.CategoryTb `json:"data"`
	}
	bodyData := &InputDeleteCategory{}
	if err := ctx.BodyParser(&bodyData); err != nil {
		internal.Log.Error(funcName, zap.Any("body", rawBody), zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(resultWeb)
	}
	userPayload, ok := ctx.Locals("payload").(models.PortalPayload)
	if !ok {
		internal.Log.Error("Failed to get user payload from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}
	defer func() {
		PortalSendKibana(ctx, funcName, uri, tokenAuth, bodyData, resultWeb, startTime, userPayload)
	}()
	bodyData.Data.UpdateBy = userPayload.Email
	errService := tk.svc.CategoryService.DeleteCategory(&bodyData.Data, userPayload.Email, funcName)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Any("error", errService))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	}
	resultWeb.Status = 1
	resultWeb.Msg = "Thành công"
	return ctx.Status(fiber.StatusOK).JSON(resultWeb)
}

func (tk *categoryHandlers) UpdateCategory(ctx *fiber.Ctx) error {
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	funcName := "UpdateCategory-" + requestId
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))
	type InputUpdateCategory struct {
		Data models.CategoryTb `json:"data"`
	}
	bodyData := &InputUpdateCategory{}
	if err := ctx.BodyParser(&bodyData); err != nil {
		internal.Log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(resultWeb)
	}
	userPayload, ok := ctx.Locals("payload").(models.PortalPayload)
	if !ok {
		internal.Log.Error("Failed to get user payload from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}
	defer func() {
		PortalSendKibana(ctx, funcName, uri, tokenAuth, bodyData, resultWeb, startTime, userPayload)
	}()
	bodyData.Data.UpdateBy = userPayload.Email
	if validateError := models.Validate.Struct(bodyData); validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}
	errService := tk.svc.CategoryService.UpdateCategory(&bodyData.Data, userPayload.Email)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Any("error", errService))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	}
	resultWeb.Status = 1
	resultWeb.Msg = "Thành công"
	return ctx.Status(fiber.StatusOK).JSON(resultWeb)
}

func (p *categoryHandlers) GetCategoryListHandler(ctx *fiber.Ctx) error {
	funcName := "GetCategoryListHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	// Parse raw body để log khi lỗi
	var rawBody interface{}
	ctx.BodyParser(&rawBody)

	type Filter struct {
		Data models.FilterCategory `json:"data"`
	}
	bodyData := &Filter{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Cannot parse body", zap.Any("body", rawBody), zap.Error(err))
		// Gửi về Kibana trước khi trả lỗi
		PortalSendKibana(ctx, funcName, uri, tokenAuth, bodyData, resultWeb, startTime, models.PortalPayload{})
		return ctx.Status(fiber.StatusOK).JSON(resultWeb)
	}

	// Lấy user payload
	userPayload, ok := ctx.Locals("payload").(models.PortalPayload)
	if !ok {
		internal.Log.Error("Failed to get payload from context")
		PortalSendKibana(ctx, funcName, uri, tokenAuth, bodyData, resultWeb, startTime, models.PortalPayload{})
		return ctx.Status(fiber.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}

	// Validate
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("Validation error", zap.Any("input", bodyData), zap.Error(validateError))
		PortalSendKibana(ctx, funcName, uri, tokenAuth, bodyData, resultWeb, startTime, userPayload)
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}

	// Gọi service
	resp, errService := p.svc.CategoryService.GetCategoryListService(&bodyData.Data)
	if errService != nil {
		resultWeb = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultWeb = models.RespWeb{
			Status: 1,
			Msg:    "Thành công",
			Detail: map[string]interface{}{
				"list_categories": resp,
			},
		}
	}

	// Gửi về Kibana
	PortalSendKibana(ctx, funcName, uri, tokenAuth, bodyData, resultWeb, startTime, userPayload)

	return ctx.Status(http.StatusOK).JSON(resultWeb)
}
func (p *categoryHandlers) GetCategoryListNameHandler(ctx *fiber.Ctx) error {
	funcName := "GetCategoryListNameHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	// Lấy payload từ context
	userPayload, ok := ctx.Locals("payload").(models.PortalPayload)
	if !ok {
		internal.Log.Error("Failed to get payload from context")
		resultWeb = models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		}
		PortalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultWeb, startTime, models.PortalPayload{})
		return ctx.Status(http.StatusUnauthorized).JSON(resultWeb)
	}

	// Gọi service
	resp, errService := p.svc.CategoryService.GetCategoryListNameService()
	if errService != nil {
		resultWeb = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultWeb = models.RespWeb{
			Status: 1,
			Msg:    "Thành công",
			Detail: map[string]interface{}{
				"list_name_categories": resp,
			},
		}
	}

	// Gửi về Kibana
	PortalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultWeb, startTime, userPayload)

	return ctx.Status(http.StatusOK).JSON(resultWeb)
}
