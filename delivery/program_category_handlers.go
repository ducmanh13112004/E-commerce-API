package delivery

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/services"
	"ecom_promotion_v2/internal/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Định nghĩa interface và struct cho handler
type ProgramCategoryHandlers interface {
	GetProgramCategoryHandler(*fiber.Ctx) error
	CreateProgramCategory(*fiber.Ctx) error
	GetProgramCategoryListAllHandler(*fiber.Ctx) error
	UpdateProgramCategory(*fiber.Ctx) error
	DeleteProgramCategory(*fiber.Ctx) error
	OnOffProgramCategoryListAllHandler(*fiber.Ctx) error
}

type programcategoryHandlers struct {
	svc *services.AppServices
}

func NewProgramCategoryHandlers(
	appService *services.AppServices,
) ProgramCategoryHandlers {
	return &programcategoryHandlers{
		svc: appService,
	}
}

func (p *programcategoryHandlers) GetProgramCategoryListAllHandler(ctx *fiber.Ctx) error {
	funcName := "GetProgramCategoryListAllHandler"
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
		Data models.FilterProgramcategory `json:"data"`
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

	resp, errService := p.svc.ProgramCategoryService.GetProgramCategoryListAllService(&bodyData.Data)
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
				"list_all_programcategories": resp,
			},
		}
	}

	// Trả về JSON response
	return ctx.Status(http.StatusOK).JSON(resultFe)
}

func (p *programcategoryHandlers) OnOffProgramCategoryListAllHandler(ctx *fiber.Ctx) error {
	funcName := "OnOffProgramCategoryListAllHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}
	var body interface{}
	ctx.BodyParser(&body)
	type InputGetPromotion struct {
		Data models.UpdateOnOffProgramCategory `json:"data"`
	}
	bodyData := &InputGetPromotion{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultFe)
	}
	update_by := ctx.Locals("email").(string)
	errService := p.svc.OnOffProgramCategory(&bodyData.Data, update_by)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programcategoryHandlers) GetProgramCategoryHandler(ctx *fiber.Ctx) error {
	funcName := "GetProgramCategoryHandler"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)
	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	startTime := time.Now()

	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))
	listCategories, errService := p.svc.ProgramCategoryService.GetProgramCategoryService() //
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("error", errService))
		resultWeb := models.RespWeb{
			Status: errService.Status,
			Msg:    errService.Msg,
			Detail: errService.Detail,
		}
		resultLocalConverted := models.RespLocal{
			StatusCode: resultWeb.Status,
			Message:    resultWeb.Msg,
		}
		LocalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultLocalConverted, startTime)

		return ctx.Status(http.StatusInternalServerError).JSON(resultWeb)
	}
	resultWeb := models.RespWeb{
		Status: 1,
		Msg:    "Thành công",
		Detail: map[string]interface{}{
			"list_active_programcategories": listCategories,
		},
	}
	resultLocalConverted := models.RespLocal{
		StatusCode: resultWeb.Status,
		Message:    resultWeb.Msg,
	}
	LocalSendKibana(ctx, funcName, uri, tokenAuth, nil, resultLocalConverted, startTime)
	return ctx.Status(http.StatusOK).JSON(resultWeb)
}
func (p *programcategoryHandlers) CreateProgramCategory(ctx *fiber.Ctx) error {
	funcName := "CreateProgramCategory"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Đọc body từ request
	var body interface{}
	ctx.BodyParser(&body)

	// Định nghĩa cấu trúc input
	type InputGetCategory struct {
		Data struct {
			ListCategories []models.ProgramCategory `json:"list_categories"`
		} `json:"data"`
	}

	// Parse body thành InputGetCategory
	bodyData := &InputGetCategory{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	// Lấy email từ context
	updateBy, ok := ctx.Locals("email").(string)
	if !ok {
		internal.Log.Error("Failed to get update_by from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}

	for i := range bodyData.Data.ListCategories {
		// Kiểm tra giá trị active
		if bodyData.Data.ListCategories[i].Active != 0 && bodyData.Data.ListCategories[i].Active != 1 {
			internal.Log.Error("Invalid active value", zap.Int("category_id", bodyData.Data.ListCategories[i].CategoryId), zap.Int("active", bodyData.Data.ListCategories[i].Active))
			resultFe.Msg = "Giá trị 'active' phải là 0 hoặc 1"
			return ctx.Status(http.StatusBadRequest).JSON(resultFe)
		}
		// Gán thông tin người cập nhật
		bodyData.Data.ListCategories[i].UpdateBy = updateBy
	}
	// Kiểm tra tính hợp lệ của dữ liệu
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultFe)
	}

	// Gọi dịch vụ để tạo các danh mục
	errService := p.svc.ProgramCategoryService.CreateProgramCategories(bodyData.Data.ListCategories, updateBy)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	// Trả về kết quả
	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programcategoryHandlers) UpdateProgramCategory(ctx *fiber.Ctx) error {
	funcName := "UpdateProgramCategory"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	// Nhận dữ liệu từ body
	var body interface{}
	ctx.BodyParser(&body)

	// Định nghĩa struct để nhận dữ liệu
	type InputUpdateCategory struct {
		Data struct {
			ListCategories []models.ProgramCategory `json:"list_categories"`
		} `json:"data"`
	}

	bodyData := &InputUpdateCategory{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	// Lấy thông tin người cập nhật từ context
	updateBy, ok := ctx.Locals("email").(string)
	if !ok {
		internal.Log.Error(funcName, zap.String("error", "Failed to get update_by from context"))
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_SYSTEM_ERROR,
			Msg:    internal.MSG_SYSTEM_ERROR,
		})
	}

	for i := range bodyData.Data.ListCategories {
		// Kiểm tra giá trị active
		if bodyData.Data.ListCategories[i].Active != 0 && bodyData.Data.ListCategories[i].Active != 1 {
			internal.Log.Error("Invalid active value", zap.Int("category_id", bodyData.Data.ListCategories[i].CategoryId), zap.Int("active", bodyData.Data.ListCategories[i].Active))
			resultFe.Msg = "Giá trị 'active' phải là 0 hoặc 1"
			return ctx.Status(http.StatusBadRequest).JSON(resultFe)
		}
		// Gán thông tin người cập nhật
		bodyData.Data.ListCategories[i].UpdateBy = updateBy
	}

	// Kiểm tra tính hợp lệ của dữ liệu
	validateError := models.Validate.Struct(bodyData)
	if validateError != nil {
		internal.Log.Error("Validation error", zap.Error(validateError))
		return ctx.Status(http.StatusBadRequest).JSON(resultFe)
	}

	// Gọi service để cập nhật danh sách danh mục
	errService := p.svc.ProgramCategoryService.UpdateProgramCategories(bodyData.Data.ListCategories, updateBy)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}

func (p *programcategoryHandlers) DeleteProgramCategory(ctx *fiber.Ctx) error {
	funcName := "DeleteProgramCategory"
	defer utils.ExecTime(utils.GetTimeUTC7(), funcName, nil)

	resultFe := models.RespWeb{
		Status: internal.SysStatus.WrongParams.Status,
		Msg:    internal.SysStatus.WrongParams.Msg,
	}

	uri := string(ctx.Request().URI().RequestURI())
	tokenAuth := string(ctx.Request().Header.Peek("TOKEN"))
	internal.Log.Info(funcName, zap.Any("uri", uri), zap.Any("auth", tokenAuth))

	var body interface{}
	ctx.BodyParser(&body)
	type InputDeleteCategory struct {
		Data models.ProgramCategory `json:"data"`
	}

	bodyData := &InputDeleteCategory{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("body", body), zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(resultFe)
	}

	updateBy, ok := ctx.Locals("email").(string)
	if !ok {
		internal.Log.Error("Failed to get update_by from context")
		return ctx.Status(http.StatusUnauthorized).JSON(models.RespWeb{
			Status: internal.CODE_WRONG_PARAMS,
			Msg:    internal.MSG_WRONG_PARAMS,
		})
	}

	bodyData.Data.UpdateBy = updateBy

	errService := p.svc.ProgramCategoryService.DeleteProgramCategory(&bodyData.Data, updateBy)
	if errService != nil {
		internal.Log.Error(funcName, zap.Any("input", bodyData), zap.String("error", errService.Msg))
		return ctx.Status(fiber.StatusOK).JSON(errService)
	} else {
		resultFe.Status = 1
		resultFe.Msg = "Thành công"
	}

	return ctx.Status(fiber.StatusOK).JSON(resultFe)
}
