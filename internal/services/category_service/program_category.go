package categoryservice

import (
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/utils"
	"sort"

	"go.uber.org/zap"
)

type ProgramCategoryService interface {
	GetProgramCategoryService() ([]map[string]interface{}, *internal.SystemStatus)
	GetProgramCategoryListAllService(filter *models.FilterProgramcategory) ([]map[string]interface{}, *internal.SystemStatus)
	CreateProgramCategories(categories []models.ProgramCategory, updateBy string) *internal.SystemStatus
	UpdateProgramCategories(categories []models.ProgramCategory, updateBy string) *internal.SystemStatus
	DeleteProgramCategory(category *models.ProgramCategory, updateby string) *internal.SystemStatus
	OnOffProgramCategory(category *models.UpdateOnOffProgramCategory, updateBy string) *internal.SystemStatus
}

type programcategoryService struct {
	repo *repositories.Repositories
}

func NewProgramCategoryService(repo *repositories.Repositories) ProgramCategoryService {
	return &programcategoryService{repo: repo}
}

// Sắp bắt đầu: 1
// Sắp kết thúc: 2
// Đang diễn ra: 3
// Đã kết thúc: 4
// Tạm ẩn: 5
func (p *programcategoryService) setStateTime(stateTime *int, category *models.ProgramCategory) {
	now := utils.GetTimeUTC7()
	nowFormat := now.Format("2006-01-02")
	isDurationActive := category.ProgramPromotionTb.BeginUsable.Before(now) && category.ProgramPromotionTb.EndUsable.After(now)
	isTimeInDay := category.ProgramPromotionTb.BeginUsable.Format("2006-01-02") == nowFormat || category.ProgramPromotionTb.EndUsable.Format("2006-01-02") == nowFormat
	if isDurationActive || isTimeInDay {
		if category.ProgramPromotionTb.State == 1 {
			// Chương trình đang diễn ra
			*stateTime = 3
			if category.ProgramPromotionTb.EndUsable.Sub(now).Hours() <= 24 {
				//Sắp kết thúc trong 24h
				*stateTime = 2
			}
		} else if category.ProgramPromotionTb.State == 0 {
			// Chương trình tạm ẩn
			*stateTime = 5
		}

	} else {
		if category.ProgramPromotionTb.EndUsable.Before(now) {
			// Chương trình đã kết thúc
			*stateTime = 4
		} else if category.ProgramPromotionTb.BeginUsable.Sub(now).Hours() <= 24 {
			// Chương trình sắp bắt đầu
			*stateTime = 1
		}
	}

}

func (p *programcategoryService) GetProgramCategoryListAllService(filter *models.FilterProgramcategory) ([]map[string]interface{}, *internal.SystemStatus) {
	funcName := "GetProgramCategoryListAllService"

	// Lấy danh sách danh mục từ repository
	categoryList, err := p.repo.ProgramCategories.GetProgramCategoryList(filter)
	if err != nil {
		internal.Log.Error(funcName, zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// Tạo danh sách danh mục chỉ với ProgramId và CategoryId
	listCategories := make([]map[string]interface{}, 0, len(categoryList))

	for _, category := range categoryList {
		var stateTime int
		// Gọi setStateTime để tính stateTime từ ProgramCategory
		p.setStateTime(&stateTime, &category)
		if filter.StateTime != nil && *filter.StateTime != stateTime {
			continue
		}
		listCategories = append(listCategories, map[string]interface{}{
			"id":            category.Id,
			"program_id":    category.ProgramId,
			"program_name":  category.ProgramPromotionTb.ProgramName,
			"category_id":   category.CategoryId,
			"category_name": category.CategoryTb.CategoryName,
			"state":         category.ProgramPromotionTb.State,
			"active":        category.Active,
			"update_by":     category.UpdateBy,
			"priority":      category.Priority,
			"t_update":      category.TUpdate.Format("2006-01-02 15:04:05"),
		})
	}

	return listCategories, nil
}

func (p *programcategoryService) OnOffProgramCategory(category *models.UpdateOnOffProgramCategory, updateBy string) *internal.SystemStatus {
	funcName := "OnOffProgramCategory"
	programCategory, err := p.repo.ProgramCategories.GetProgramCategoryByID(category.Id)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", category), zap.Error(err))
		return internal.SysStatus.SystemError
	}

	now := utils.GetTimeUTC7()
	programCategory.Active = category.Active
	programCategory.UpdateBy = updateBy
	programCategory.TUpdate = &now
	err = p.repo.ProgramCategories.UpdateOnOffProgramCategory(programCategory)
	if err != nil {
		internal.Log.Error(funcName, zap.Any("input", category), zap.Error(err))
		return internal.SysStatus.SystemError
	}
	return nil
}

func (p *programcategoryService) GetProgramCategoryService() ([]map[string]interface{}, *internal.SystemStatus) {
	funcName := "GetProgramCategoryService"

	// Lấy danh sách danh mục từ repository
	categoryList, err := p.repo.ProgramCategories.GetProgramCategory()
	if err != nil {
		internal.Log.Error(funcName, zap.Error(err))
		return nil, internal.SysStatus.SystemError
	}

	// Sắp xếp danh sách theo (A → Z)
	sort.Slice(categoryList, func(i, j int) bool {
		return categoryList[i].CategoryName < categoryList[j].CategoryName
	})

	// Chuyển đổi danh sách sang dạng map
	listCategories := make([]map[string]interface{}, 0, len(categoryList))
	for _, category := range categoryList {
		listCategories = append(listCategories, map[string]interface{}{
			"category_id":   category.CategoryId,
			"category_name": category.CategoryName,
		})
	}

	return listCategories, nil
}

func (s *programcategoryService) CreateProgramCategories(categories []models.ProgramCategory, updateBy string) *internal.SystemStatus {
	funcName := "CreateProgramCategories"

	// Ghi log khi bắt đầu tạo danh mục
	internal.Log.Info("Creating new categories",
		zap.String("function", funcName),
		zap.Any("input", categories),
	)

	// Kiểm tra danh sách rỗng
	if len(categories) == 0 {
		return &internal.SystemStatus{
			Status: 400,
			Msg:    "Danh sách danh mục không được để trống",
		}
	}

	// Kiểm tra từng danh mục trong danh sách
	for _, category := range categories {
		exists, err := s.repo.ProgramCategories.IsProgramIDExists(int(category.ProgramId))
		if err != nil || !exists {
			return &internal.SystemStatus{
				Status: 404,
				Msg:    "Chương trình không tồn tại",
			}
		}

		exist, err := s.repo.ProgramCategories.IsCategoryIDExists(int(category.CategoryId))
		if err != nil || !exist {
			return &internal.SystemStatus{
				Status: 404,
				Msg:    "Danh mục không tồn tại",
			}
		}

		existss, err := s.repo.ProgramCategories.IsProgramidCategoryidExists(category.ProgramId, category.CategoryId)
		if err != nil {
			internal.Log.Error("Failed to check if program_id and category_id exist",
				zap.String("function", funcName),
				zap.Int64("program_id", category.ProgramId),
				zap.Int("category_id", category.CategoryId),
				zap.Error(err),
			)
			return internal.SysStatus.SystemError
		}

		if existss {
			internal.Log.Error("Duplicate program_id and category_id detected",
				zap.String("function", funcName),
				zap.Int64("program_id", category.ProgramId),
				zap.Int("category_id", category.CategoryId),
			)
			return &internal.SystemStatus{
				Status: 1062,
				Msg:    "Chương trình và danh mục này đã tồn tại",
			}
		}
	}

	// Cập nhật thời gian hiện tại cho từng danh mục
	now := utils.GetTimeUTC7()
	for i := range categories {
		categories[i].TUpdate = &now
	}

	// Gọi repo để thêm danh sách danh mục
	err := s.repo.ProgramCategories.CreateProgramCategories(categories)
	if err != nil {
		internal.Log.Error("CreateProgramCategories failed",
			zap.String("function", funcName),
			zap.Any("input", categories),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// Ghi log khi tạo thành công
	internal.Log.Info("Program categories created successfully",
		zap.String("function", funcName),
	)

	return nil
}

func (s *programcategoryService) UpdateProgramCategories(categories []models.ProgramCategory, updateBy string) *internal.SystemStatus {
	funcName := "UpdateProgramCategories"

	// Tạo một slice để lưu trữ các danh mục cần cập nhật
	var categoriesToUpdate []models.ProgramCategory

	// Lặp qua từng danh mục trong danh sách categories
	for _, category := range categories {
		// Kiểm tra xem danh mục có tồn tại không
		existingCategory, err := s.repo.ProgramCategories.GetProgramCategoryByID(category.Id)
		if err != nil {
			internal.Log.Error("GetCategoryByID: category not found",
				zap.String("function", funcName),
				zap.Int("id", category.Id),
				zap.Error(err),
			)
			// Không dừng lại mà chỉ ghi log và tiếp tục với các danh mục khác
			continue
		}

		// Cập nhật thông tin danh mục
		now := utils.GetTimeUTC7()
		existingCategory.TUpdate = &now
		existingCategory.ProgramId = category.ProgramId
		existingCategory.CategoryId = category.CategoryId
		existingCategory.Priority = category.Priority
		existingCategory.Active = category.Active
		existingCategory.UpdateBy = updateBy

		// Thêm danh mục đã cập nhật vào danh sách categoriesToUpdate
		categoriesToUpdate = append(categoriesToUpdate, *existingCategory)
	}

	// Nếu không có danh mục nào cần cập nhật, trả về lỗi
	if len(categoriesToUpdate) == 0 {
		internal.Log.Error("No valid categories to update", zap.String("function", funcName))
		return internal.SysStatus.SystemError
	}

	// Gọi hàm cập nhật danh mục với batch update
	err := s.repo.ProgramCategories.UpdateProgramCategories(categoriesToUpdate)
	if err != nil {
		internal.Log.Error("UpdateCategories failed",
			zap.String("function", funcName),
			zap.Any("input", categoriesToUpdate),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// Trả về thành công nếu tất cả các danh mục đều được cập nhật thành công
	return nil
}

func (s *programcategoryService) DeleteProgramCategory(category *models.ProgramCategory, updateBy string) *internal.SystemStatus {
	funcName := "DeleteProgramCategory"

	// Kiểm tra category có tồn tại không
	categorydb, err := s.repo.ProgramCategories.GetProgramCategoryByID(category.Id)
	if err != nil {
		internal.Log.Error("Database error when fetching category",
			zap.String("function", funcName),
			zap.Int("id", category.Id),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	// Ghi log trước khi xóa
	internal.Log.Info("Soft deleting program category",
		zap.String("function", funcName),
		zap.Int("id", category.Id),
		zap.String("updateBy", updateBy),
	)

	// Gọi repo để xóa mềm
	err = s.repo.ProgramCategories.DeleteProgramCategory(*categorydb)
	if err != nil {
		internal.Log.Error("Soft delete failed",
			zap.String("function", funcName),
			zap.Int("id", category.Id),
			zap.Error(err),
		)
		return internal.SysStatus.SystemError
	}

	return nil
}
