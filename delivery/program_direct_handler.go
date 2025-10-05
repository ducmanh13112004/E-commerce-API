package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"
	"ecom_promotion_v2/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Định nghĩa interface và struct cho handler
type ProgramDirectHandlers interface {
	CreateProgramDirect(*fiber.Ctx) error
	GetProgramDirectListAllHandler(*fiber.Ctx) error
	GetProgramDirectIDHandler(*fiber.Ctx) error
	UpdateProgramDirect(*fiber.Ctx) error
	DeleteProgramDirect(*fiber.Ctx) error
}

type programdirectHandlers struct {
	svc *services.AppServices
}

func NewProgramDirectHandlers(
	appService *services.AppServices,
) ProgramDirectHandlers {
	return &programdirectHandlers{
		svc: appService,
	}
}

func (p *programdirectHandlers) GetProgramDirectIDHandler(ctx *fiber.Ctx) error {
	funcName := "GetProgramDirectIDHandler"
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	// Biến lưu kết quả trả về
	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Đảm bảo Kibana logging luôn được gọi
	defer func() {
		resultLocalConverted := models.RespLocal{
			StatusCode: resultWeb.Status,
			Message:    resultWeb.Msg,
		}
		LocalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultLocalConverted, startTime)
	}()

	// Lấy redirect_id từ request body
	var body interface{}
	ctx.BodyParser(&body)
	type InputGetPromotion struct {
		Data models.ProgramDirectScreenMobileTb `json:"data"`
	}
	bodyData := &InputGetPromotion{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultWeb)
	}

	// Gọi service với redirect_id
	redirectID := bodyData.Data.RedirectId
	ListdirectID, errService := p.svc.ProgramDirectService.GetProgramDirectIDService(redirectID)

	if errService != nil {
		resultWeb = models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
		return ctx.Status(http.StatusInternalServerError).JSON(resultWeb)
	}

	// Trả về kết quả thành công
	resultWeb = models.RespWeb{
		Status: 1,
		Msg:    "Thành công",
		Detail: map[string]interface{}{
			"direct_program_ID": ListdirectID,
		},
	}
	return ctx.Status(http.StatusOK).JSON(resultWeb)
}

func (p *programdirectHandlers) GetProgramDirectListAllHandler(ctx *fiber.Ctx) error {
	funcName := "GetProgramDirectListAllHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	// Khởi tạo response mặc định
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Parse body từ request
	var body interface{}
	ctx.BodyParser(&body)
	type Filter struct {
		Data models.FilterProgramDirectScreenMobileTb `json:"data"`
	}
	bodyData := &Filter{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error("Cannot parse", zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	// Kiểm tra tính hợp lệ của dữ liệu đầu vào
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("validateError", zap.Any("input", bodyData), zap.Error(validateError))
		return ctx.Status(http.StatusOK).JSON(internal.SysStatus.WrongParams)
	}

	// Gọi service xử lý logic
	resp, errService := p.svc.ProgramDirectService.GetProgramDirectListAllService(&bodyData.Data)
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
				"list_direct_links": resp,
			},
		}
	}

	// Trả về JSON response
	return ctx.Status(http.StatusOK).JSON(resultFe)
}

func (tk *programdirectHandlers) CreateProgramDirect(ctx *fiber.Ctx) error {
	funcName := "CreateProgramDirect"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Lấy thông tin request
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	// Định nghĩa struct cho dữ liệu đầu vào
	type InputProgramDirect struct {
		Data models.ProgramDirectScreenMobileTb `json:"data"`
	}

	// Parse body
	bodyData := &InputProgramDirect{}
	if err := ctx.BodyParser(bodyData); err != nil {
		internal.Log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}

	// Lấy thông tin người cập nhật
	updateBy, ok := ctx.Locals("email").(string)
	if !ok {
		internal.Log.Error("Failed to get update_by from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}

	// Gán thông tin người cập nhật
	bodyData.Data.UpdateBy = updateBy

	// Validate dữ liệu đầu vào
	if validateError := models.Validate.Struct(bodyData.Data); validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}

	// Gọi service để tạo program direct
	errService := tk.svc.ProgramDirectService.CreateProgramDirect(funcName, &bodyData.Data)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Any("error", errService))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultWeb.Status = 1
		resultWeb.Msg = "Thành công"
	}

	// Gửi log về Kibana
	resultLocalConverted := models.RespLocal{
		StatusCode: resultWeb.Status,
		Message:    resultWeb.Msg,
	}
	LocalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultLocalConverted, startTime)

	return ctx.Status(fiber.StatusOK).JSON(resultWeb)
}
func (tk *programdirectHandlers) UpdateProgramDirect(ctx *fiber.Ctx) error {
	funcName := "UpdateProgramDirect"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Lấy thông tin request
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	// Định nghĩa struct cho dữ liệu đầu vào
	type InputProgramDirect struct {
		Data models.ProgramDirectScreenMobileTb `json:"data"`
	}

	// Parse body
	bodyData := &InputProgramDirect{}
	if err := ctx.BodyParser(bodyData); err != nil {
		internal.Log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}

	// Lấy thông tin người cập nhật
	updateBy, ok := ctx.Locals("email").(string)
	if !ok {
		internal.Log.Error("Failed to get update_by from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    internal.MSG_SYSTEM_ERROR,
		})
	}

	// Gán thông tin người cập nhật
	bodyData.Data.UpdateBy = updateBy

	// Validate dữ liệu đầu vào
	if validateError := models.Validate.Struct(bodyData.Data); validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}

	// Gọi service để cập nhật program direct
	errService := tk.svc.UpdateProgramDirect(&bodyData.Data)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Any("error", errService))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultWeb.Status = 1
		resultWeb.Msg = "Thành công"
	}

	// Gửi log về Kibana
	resultLocalConverted := models.RespLocal{
		StatusCode: resultWeb.Status,
		Message:    resultWeb.Msg,
	}
	LocalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultLocalConverted, startTime)

	return ctx.Status(fiber.StatusOK).JSON(resultWeb)
}

func (tk *programdirectHandlers) DeleteProgramDirect(ctx *fiber.Ctx) error {
	funcName := "DeleteProgramDirect"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultWeb := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Lấy thông tin request
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	// Định nghĩa struct cho dữ liệu đầu vào
	type InputDelete struct {
		Data models.DeleteProgramDirectRequest `json:"data"`
	}

	// Parse body
	bodyData := &InputDelete{}
	if err := ctx.BodyParser(bodyData); err != nil {
		internal.Log.Error("Failed to parse request body", zap.Error(err))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}

	// Lấy thông tin người thực hiện
	updateBy, ok := ctx.Locals("email").(string)
	if !ok {
		internal.Log.Error("Failed to get update_by from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}

	// // Gán thông tin người cập nhật
	// bodyData.Data.UpdateBy = updateBy

	// Validate dữ liệu đầu vào
	if validateError := models.Validate.Struct(bodyData.Data); validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultWeb)
	}

	// Gọi service để xóa dữ liệu
	errService := tk.svc.ProgramDirectService.DeleteProgramDirect(&bodyData.Data, updateBy)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.Any("error", errService))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultWeb.Status = 1
		resultWeb.Msg = "Thành công"
	}

	// Gửi log về Kibana
	resultLocalConverted := models.RespLocal{
		StatusCode: resultWeb.Status,
		Message:    resultWeb.Msg,
	}
	LocalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultLocalConverted, startTime)

	return ctx.Status(fiber.StatusOK).JSON(resultWeb)
}
