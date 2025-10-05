package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"
	"ecom_promotion_v2/internal/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProgramPromotionHandler interface {
	GetProgramPromotionListHandler(ctx *fiber.Ctx) error
	GetProgramPromotionByIDHandler(ctx *fiber.Ctx) error
	CreateProgramPromotionHandler(ctx *fiber.Ctx) error
	UpdateProgramPromotionHandler(ctx *fiber.Ctx) error
	OnOffProgramPromotionHandler(ctx *fiber.Ctx) error
	DeleteProgramPromotionHandler(ctx *fiber.Ctx) error
	GetReportProgramPromotionHandler(ctx *fiber.Ctx) error
	GetReportProgramPromotionLocalHandler(ctx *fiber.Ctx) error
}

type programPromotionHandlers struct {
	svc *services.AppServices
}

func NewProgramPromotionHandlers(appService *services.AppServices) ProgramPromotionHandler {
	return &programPromotionHandlers{
		svc: appService,
	}
}

func (p *programPromotionHandlers) GetProgramPromotionListHandler(ctx *fiber.Ctx) error {
	funcName := "GetProgramPromotionList"
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
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	if requestId == "" {
		requestId = uuid.New().String()
	}
	funcName += "-" + requestId
	startTime := utils.GetTimeUTC7()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, employee) }()
	bodyData := &models.FilterProgramPromotion{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Cannot parse", zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	if bodyData.ProgramId == nil {
		_, err = time.Parse("2006-01-02", bodyData.FromDate)
		if err != nil {
			internal.Log.Error("Incorrect format fromDate", zap.Any("funcName", funcName), zap.Any("body", body))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
		_, err = time.Parse("2006-01-02", bodyData.ToDate)
		if err != nil {
			internal.Log.Error("Incorrect format toDate", zap.Any("funcName", funcName), zap.Any("body", body))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
	}
	resp, errService := p.svc.ProgramPromotionService.GetProgramPromotionList(bodyData, funcName)
	if errService != nil {
		resultFe = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultFe = models.RespWeb{
			Status: 1,
			Msg:    "Thành công",
			Detail: map[string]interface{}{
				"list_program": resp,
			},
		}
	}
	return ctx.Status(http.StatusOK).JSON(resultFe)

}

func (p *programPromotionHandlers) GetProgramPromotionByIDHandler(ctx *fiber.Ctx) error {
	funcName := "GetProgramPromotionByID"
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
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	if requestId == "" {
		requestId = uuid.New().String()
	}
	funcName += "-" + requestId
	startTime := utils.GetTimeUTC7()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, employee) }()

	type InputGetPromotion struct {
		Data struct {
			ProgramID int64 `json:"program_id" validate:"required"`
		} `json:"data" validate:"required"`
	}
	bodyData := &InputGetPromotion{}
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
	resp, errService := p.svc.ProgramPromotionService.GetProgramPromotionByID(bodyData.Data.ProgramID, funcName)
	if errService != nil {
		resultFe = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultFe = models.RespWeb{
			Status: 1,
			Msg:    "Thành công",
			Detail: map[string]interface{}{
				"program_promotion": resp,
			},
		}
	}
	return ctx.Status(http.StatusOK).JSON(resultFe)
}

func (p *programPromotionHandlers) CreateProgramPromotionHandler(ctx *fiber.Ctx) error {
	funcName := "CreateProgramPromotion"
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
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	if requestId == "" {
		requestId = uuid.New().String()
	}
	funcName += "-" + requestId
	startTime := utils.GetTimeUTC7()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, employee) }()

	type InputGetPromotion struct {
		Data models.CreateProgramPromotion `json:"data" validate:"required"`
	}
	bodyData := &InputGetPromotion{}
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

	update_by := ctx.Locals("email").(string)
	resp, errService := p.svc.CreateProgramPromotion(&bodyData.Data, update_by, funcName)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
		resultFe.Detail = map[string]interface{}{
			"program_id": resp,
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programPromotionHandlers) UpdateProgramPromotionHandler(ctx *fiber.Ctx) error {
	requestId := string(ctx.Request().Header.Peek("X-Request-Id")) // uuid
	funcName := "UpdateProgramPromotion-" + requestId
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

	type InputGetPromotion struct {
		Data models.UpdateProgramPromotion `json:"data" validate:"required"`
	}
	bodyData := &InputGetPromotion{}
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

	update_by := ctx.Locals("email").(string)
	errService := p.svc.UpdateProgramPromotion(&bodyData.Data, update_by, funcName)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programPromotionHandlers) OnOffProgramPromotionHandler(ctx *fiber.Ctx) error {
	funcName := "OnOffProgramPromotion"
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
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	if requestId == "" {
		requestId = uuid.New().String()
	}
	funcName += "-" + requestId
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, employee) }()

	type InputGetPromotion struct {
		Data models.OnOffProgramPromotion `json:"data"`
	}
	bodyData := &InputGetPromotion{}
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

	update_by := ctx.Locals("email").(string)
	errService := p.svc.OnOffProgramPromotion(&bodyData.Data, update_by, funcName)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programPromotionHandlers) DeleteProgramPromotionHandler(ctx *fiber.Ctx) error {
	funcName := "DeleteProgramPromotion"
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
	requestId := string(ctx.Request().Header.Peek("X-Request-Id"))
	if requestId == "" {
		requestId = uuid.New().String()
	}
	funcName += "-" + requestId
	startTime := utils.GetTimeUTC7()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { PortalSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, employee) }()

	type InputDelete struct {
		Data struct {
			ProgramID int64 `json:"program_id" validate:"required"`
		} `json:"data" validate:"required"`
	}
	bodyData := &InputDelete{}
	err := ctx.BodyParser(&bodyData)
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError", zap.Any("input", bodyData), zap.Error(validateError))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	update_by := ctx.Locals("email").(string)
	errService := p.svc.DeleteProgramPromotion(bodyData.Data.ProgramID, update_by, funcName)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programPromotionHandlers) GetReportProgramPromotionHandler(ctx *fiber.Ctx) error {
	funcName := "GetReportProgramPromotionHandler"

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

	bodyData := &models.FilterReportProgramPromotion{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Can not parse body", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	if len(bodyData.ListProgramId) == 0 && len(bodyData.FromDate) == 0 && len(bodyData.ToDate) == 0 {
		internal.Log.Error("Body is empty", zap.Any("funcName", funcName), zap.Any("body", body))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	if len(bodyData.FromDate) != 0 && len(bodyData.ToDate) != 0 && len(bodyData.ListProgramId) == 0 {
		_, err := time.Parse("2006-01-02", bodyData.FromDate)
		if err != nil {
			internal.Log.Error("Incorrect format fromDate", zap.Any("funcName", funcName), zap.Any("body", body))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}

		_, err = time.Parse("2006-01-02", bodyData.ToDate)
		if err != nil {
			internal.Log.Error("Incorrect format toDate", zap.Any("funcName", funcName), zap.Any("body", body))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
		if bodyData.FromDate > bodyData.ToDate {
			internal.Log.Error("Incorrect FromDate > ToDate", zap.Any("funcName", funcName), zap.Any("body", body))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
	}
	resp, errService := p.svc.GetReportProgramPromotion(bodyData)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
		resultFe.Detail = map[string]interface{}{
			"program_promotion": resp,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programPromotionHandlers) GetReportProgramPromotionLocalHandler(ctx *fiber.Ctx) error {
	funcName := "GetReportProgramPromotionLocal"

	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	resultFe := models.RespLocal{
		StatusCode: internal.SysStatus.WrongParams.Status,
		Message:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)

	bodyData := &models.ReportProgramPromotionLocal{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Can not parse body", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError", zap.Any("input", bodyData), zap.Error(validateError))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}
	lengthListID := len(bodyData.ListProgramId)
	if lengthListID == 0 || lengthListID > 10 {
		internal.Log.Error("Body is empty", zap.Any("funcName", funcName), zap.Any("body", body))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	resp, errService := p.svc.GetReportProgramPromotionLocal(bodyData.ListProgramId)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.StatusCode = 0
		resultFe.Message = "Ok"
		resultFe.Data = map[string]interface{}{
			"program_promotion": resp,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}
