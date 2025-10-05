package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"
	"ecom_promotion_v2/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PromotionsHandlers interface {
	//local
	LocalReceiveVoucherFrProgramId(*fiber.Ctx) error

	//public
	AfiliateGetProgramListFrCategoryId(*fiber.Ctx) error

	//portal
	SavePromotionHandler(ctx *fiber.Ctx) error
	GenerateListPromotionCode(ctx *fiber.Ctx) error
	CreatePromotion(ctx *fiber.Ctx) error
	GetPromotionByProgramId(ctx *fiber.Ctx) error
}
type promotionsHandlers struct {
	svc *services.AppServices
	// repo *repositories.Repositories
}

func NewPromotionsHandlers(
	appService *services.AppServices,
) PromotionsHandlers {
	return &promotionsHandlers{
		svc: appService,
	}
}

func (tk *promotionsHandlers) LocalReceiveVoucherFrProgramId(ctx *fiber.Ctx) error {
	defer utils.ExecTime(utils.GetTimeUTC7(), "LocalReceiveVoucherFrProgramId", nil)
	funcName := "LocalReceiveVoucherFrProgramId"
	resultLocal := models.RespLocal{
		StatusCode: internal.SysStatus.WrongParams.Status,
		Message:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { LocalSendKibana(ctx, funcName, uri, tokenAuth, body, resultLocal, startTime) }()
	// Handler input , validate input
	bodyData := &models.InputReceiveVoucher{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Cannot parse", zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(internal.SysStatus.WrongParams)
	}
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError", zap.Any("input", bodyData), zap.Error(validateError))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}
	resp, errService := tk.svc.PromotionsService.ReceiveVoucherFrProgramId(bodyData)
	if errService != nil {
		resultLocal = models.RespLocal{
			StatusCode: errService.Status,
			Message:    errService.Msg,
			Data:       errService.Detail,
		}
	} else {
		resultLocal = models.RespLocal{
			StatusCode: resp.StatusCode,
			Message:    resp.Message,
			Data:       resp.Data,
		}
	}
	return ctx.Status(http.StatusOK).JSON(resultLocal)
}

func (tk *promotionsHandlers) AfiliateGetProgramListFrCategoryId(ctx *fiber.Ctx) error {
	defer utils.ExecTime(utils.GetTimeUTC7(), "AfiliateGetProgramListFrCategoryId", nil)
	funcName := "AfiliateGetProgramListFrCategoryId"
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)
	var customer models.UserInfo
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth), zap.Any("body", body))
	defer func() { WebkitSendKibana(ctx, funcName, uri, tokenAuth, body, resultFe, startTime, customer) }()
	//get info user
	customer = models.UserInfo{
		CustomerId:  ctx.Locals("customer_id").(string),
		PhoneNb:     ctx.Locals("customer_phone").(string),
		AppVersion:  ctx.Locals("app_version").(string),
		TokenWebkit: ctx.Locals("token_webkit").(string),
		CustomerIp:  ctx.Locals("customer_ip").(string),
	}
	// Handler input , validate input
	bodyData := &models.InputAfiliateGetProgramListFrCategoryId{}
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
	resp, errService := tk.svc.PromotionsService.GetProgramListFrCategoryId(&customer, bodyData)
	if errService != nil {
		resultFe = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
	} else {
		resultFe = models.RespWeb{
			Status: resp.Status,
			Msg:    resp.Msg,
			Detail: resp.Detail,
		}
	}
	return ctx.Status(http.StatusOK).JSON(resultFe)
}

func (p *promotionsHandlers) SavePromotionHandler(ctx *fiber.Ctx) error {
	funcName := "CreatePromotionHandler"
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
		Data struct {
			Promotion []models.CreatePromotion `json:"promotion"`
		} `json:"data"`
	}
	bodyData := &InputGetPromotion{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	errService := p.svc.SaveListPromotion(bodyData.Data.Promotion)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}
	return ctx.Status(fiber.StatusOK).JSON(resultFe)

}

func (p *promotionsHandlers) GenerateListPromotionCode(ctx *fiber.Ctx) error {
	funcName := "CreateListPromotionCode"
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
		Data struct {
			Quantity int    `json:"quantity" validate:"required"`
			Prefix   string `json:"prefix" validate:"required"`
		} `json:"data"`
	}
	bodyData := &InputGetPromotion{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	resp, errService := p.svc.GenerateListPromotionCode(bodyData.Data.Quantity, bodyData.Data.Prefix)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
		resultFe.Detail = map[string]interface{}{
			"list_code": resp,
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *promotionsHandlers) CreatePromotion(ctx *fiber.Ctx) error {
	funcName := "CreatePromotion"
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

	input := &models.InputCreatePromotion{}
	err := ctx.BodyParser(&input)
	if err != nil {
		internal.Log.Error("Body Parser Error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	if input.Type == 0 {
		if input.Data.Prefix == nil || input.Data.Quantity == nil || input.Data.PromotionEndUsable == nil {
			internal.Log.Error("Validate input error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
	} else if input.Type == 1 {
		if input.Data.PromotionEndUsable == nil || input.Data.TotalRelease == nil || input.Data.PromotionCode == nil {
			internal.Log.Error("Validate input error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}

	} else if input.Type == 2 {
		const excelHeaderOffset = 2
		lenListPromotion := len(input.Data.ListPromotion)
		if lenListPromotion > 10000 {
			internal.Log.Error("Validate input error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
			resultFe.Msg = "Dữ liệu không được quá 10000 dòng"
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		}
		if lenListPromotion == 0 {
			internal.Log.Error("Validate input error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
			resultFe.Msg = "Dữ liệu trong excel không được rỗng"
			return ctx.Status(fiber.StatusOK).JSON(resultFe)
		} else {
			var msg strings.Builder
			var isValid bool = true
			for index, promotion := range input.Data.ListPromotion {
				_, err := time.Parse("2006-01-02", promotion.PromotionEndUsable)
				if err != nil && len(promotion.PromotionCode) == 0 {
					msg.WriteString("Dòng thứ ")
					msg.WriteString(strconv.Itoa(index + excelHeaderOffset))
					msg.WriteString(" không có mã khuyến mãi và định dạng ngày không hợp lệ")
					isValid = false
					break
				} else if err != nil {
					msg.WriteString("Dòng thứ ")
					msg.WriteString(strconv.Itoa(index + excelHeaderOffset))
					msg.WriteString(" định dạng ngày không hợp lệ")
					isValid = false
					break
				} else if len(promotion.PromotionCode) == 0 {
					msg.WriteString("Dòng thứ ")
					msg.WriteString(strconv.Itoa(index + excelHeaderOffset))
					msg.WriteString(" không có mã khuyến mãi")
					isValid = false
					break
				}
			}
			if !isValid {
				resultFe.Msg = msg.String()
				internal.Log.Error("Validate input error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
				return ctx.Status(fiber.StatusOK).JSON(resultFe)
			}
		}
	} else {
		internal.Log.Error("Validate input error", zap.Any("funcName", funcName), zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(internal.SysStatus.WrongParams)
	}

	errService := p.svc.CreatePromotion(funcName, input)
	if errService != nil {
		internal.Log.Error("Response error", zap.Any("funcName", funcName), zap.Any("input", input))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}
	internal.Log.Info("Response success", zap.Any("funcName", funcName), zap.Any("input", input))
	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *promotionsHandlers) GetPromotionByProgramId(ctx *fiber.Ctx) error {
	funcName := "GetPromotionByProgramId"
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
	type Input struct {
		Data struct {
			ProgramId int64 `json:"program_id" validate:"required"`
		} `json:"data" validate:"required"`
	}
	bodyData := &Input{}
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
	resp, errService := p.svc.PromotionsService.GetListPromotionByProgramId(bodyData.Data.ProgramId, funcName)

	if errService != nil {
		return ctx.Status(http.StatusOK).JSON(errService)
	}
	resultFe.Msg = "Thành công"
	resultFe.Status = 1
	resultFe.Detail = resp
	return ctx.Status(http.StatusOK).JSON(resultFe)
}
