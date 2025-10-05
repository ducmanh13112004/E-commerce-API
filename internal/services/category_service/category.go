package categoryservice

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"ecom_promotion_v2/internal/utils_call"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var BotTelegram utils_call.BotTelegramService

type CategoryService interface {
	CreateCategory(*models.CategoryTb) *internal.SystemStatus
	UpdateCategory(category *models.CategoryTb, updateBy string) *internal.SystemStatus
	DeleteCategory(category *models.CategoryTb, updateBy, funcName string) *internal.SystemStatus
	GetCategoryListService(filter *models.FilterCategory) ([]map[string]interface{}, *internal.SystemStatus)
	GetCategoryListNameService() ([]map[string]interface{}, *internal.SystemStatus)
}

type categoryService struct {
	repo *repositories.Repositories
}

func NewCategoryService(
	repo *repositories.Repositories,
) CategoryService {
	return &categoryService{
		repo: repo,
	}
}
func (s *categoryService) CreateCategory(category *models.CategoryTb) *internal.SystemStatus {
	funcName := "CreateCategory"

	// Ghi log khi bắt đầu tạo danh mục
	internal.Log.Info("Creating new category",
		zap.String("funcName", funcName),
		zap.Any("input", category),
	)

	// Cập nhật thời gian hiện tại
	now := utils.GetTimeUTC7()
	category.TUpdate = &now

	// Gọi hàm tạo category
	err := s.repo.Categories.CreateCategory(category)
	if err != nil {
		internal.Log.Error("CreateCategory failed",
			zap.String("funcName", funcName),
			zap.Any("input", category),
			zap.Error(err),
		)

		// Kiểm tra cơ sở dữ liệu
		if strings.Contains(err.Error(), "database") {
			internal.Log.Error("Database error occurred",
				zap.String("funcName", funcName),
				zap.String("category_name", category.CategoryName),
				zap.Error(err),
			)
			go utils_call.SendStatusMessage(
				fmt.Sprintf("❌ Lỗi cơ sở dữ liệu khi tạo danh mục: %s \n User: %s, ID: %d", category.CategoryName, category.UpdateBy, category.CategoryId),
				3)
			return &internal.SystemStatus{
				Status: internal.CODE_SYSTEM_ERROR,
				Msg:    "Database error occurred",
				Detail: err.Error(),
			}
		}

		// Thông báo lỗi tạo danh mục thất bại
		go utils_call.SendStatusMessage(
			fmt.Sprintf("⚠️ Tạo danh mục thất bại: %v \n User: %s, ID: %d", err, category.UpdateBy, category.CategoryId),
			2,
		)
		return internal.SysStatus.SystemError
	}

	// Ghi log khi tạo thành công
	internal.Log.Info("Category created successfully",
		zap.String("funcName", funcName),
		zap.Any("input", category),
	)
	go utils_call.SendStatusMessage(
		fmt.Sprintf("✅ Danh mục mới vừa được thêm: %s \n User: %s, ID: %d", category.CategoryName, category.UpdateBy, category.CategoryId),
		1,
	)
	return nil
}

func (s *categoryService) UpdateCategory(category *models.CategoryTb, updateBy string) *internal.SystemStatus {
	funcName := "UpdateCategory"
	// Kiểm tra xem danh mục có tồn tại không
	existingCategory, err := s.repo.Categories.GetCategoryByID(category.CategoryId)
	if err != nil {
		internal.Log.Error("GetCategoryByID: category not found",
			zap.String("funcName", funcName),
			zap.Int("category_id", category.CategoryId),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError

	}

	// Cập nhật thông tin danh mục
	now := utils.GetTimeUTC7()
	category.TUpdate = &now

	existingCategory.CategoryName = category.CategoryName
	existingCategory.CategoryPriority = category.CategoryPriority
	existingCategory.CategorySrcImg = category.CategorySrcImg
	existingCategory.State = category.State
	existingCategory.UpdateBy = category.UpdateBy
	existingCategory.TUpdate = category.TUpdate
	// Gọi hàm cập nhật danh mục
	err = s.repo.Categories.UpdateCategory(existingCategory)
	if err != nil {
		internal.Log.Error("UpdateCategory failed",
			zap.String("funcName", funcName),
			zap.Any("input", category),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	return nil
}

func (s *categoryService) DeleteCategory(category *models.CategoryTb, updateBy, funcName string) *internal.SystemStatus {

	if category.State == nil {
		defaultState := 0
		category.State = &defaultState
	}

	// Lấy thông tin category từ DB để cập nhật
	categorydb, err := s.repo.Categories.GetCategoryByID(category.CategoryId)
	if err != nil {
		internal.Log.Error("Categories.GetCategoryByID",
			zap.String("funcName", funcName),
			zap.Int("category_id", category.CategoryId),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// Cập nhật các trường cần thiết
	now := utils.GetTimeUTC7()
	categorydb.UpdateBy = updateBy
	categorydb.TUpdate = &now
	categorydb.State = category.State

	// Ghi log trước khi thực hiện xóa
	internal.Log.Info("Before Deleting category",
		zap.String("funcName", funcName),
		zap.Any("input", categorydb),
	)

	// Gọi hàm xóa category
	err = s.repo.Categories.DeleteCategory(categorydb) //*
	if err != nil {
		internal.Log.Error("DeleteCategory failed",
			zap.Any("funcName", funcName),
			zap.Any("input", categorydb),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	return nil
}

func (p *categoryService) GetCategoryListService(filter *models.FilterCategory) ([]map[string]interface{}, *internal.SystemStatus) {
	funcName := "GetCategoryListService"

	// Lấy danh sách danh mục từ repository
	categoryList, err := p.repo.Categories.GetCategoryList(filter)
	if err != nil {
		internal.Log.Error(funcName, zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// Tạo slice danh sách danh mục
	listCategories := make([]map[string]interface{}, 0, len(categoryList))

	for _, category := range categoryList {
		// Định dạng lại thời gian cập nhật
		var formattedTime string
		if category.TUpdate != nil {
			formattedTime = category.TUpdate.Format("02/01/2006 15:04:05")
		} else {
			formattedTime = "" // Nếu NULL trả về chuỗi rỗng
		}

		listCategories = append(listCategories, map[string]interface{}{
			"category_id":       category.CategoryId,
			"category_name":     category.CategoryName,
			"category_priority": category.CategoryPriority,
			"category_src_img":  category.CategorySrcImg,
			"state":             category.State,
			"t_update":          formattedTime,
			"update_by":         category.UpdateBy,
		})
	}

	return listCategories, nil
}

func (p *categoryService) GetCategoryListNameService() ([]map[string]interface{}, *internal.SystemStatus) {
	funcName := "GetCategoryListNameService"

	// Lấy danh sách danh mục từ repository
	categoryList, err := p.repo.Categories.GetCategoryListName()
	if err != nil {
		internal.Log.Error(funcName, zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// Tạo slice danh sách danh mục
	listCategories := make([]map[string]interface{}, 0, len(categoryList))

	for _, category := range categoryList {
		listCategories = append(listCategories, map[string]interface{}{
			"category_id":   category.CategoryId,
			"category_name": category.CategoryName,
		})
	}

	return listCategories, nil
}
