package services

import (
	"ecom_promotion_v2/internal/repositories"

	categoryservice "ecom_promotion_v2/internal/services/category_service"
	programcategoryservice "ecom_promotion_v2/internal/services/category_service"
	healthcheckservice "ecom_promotion_v2/internal/services/health_check_service"
	loginserervice "ecom_promotion_v2/internal/services/login_service"
	programDirect "ecom_promotion_v2/internal/services/program_direct_service"
	programproductservice "ecom_promotion_v2/internal/services/program_product_service"
	programPromotionService "ecom_promotion_v2/internal/services/program_promotion_service"

	promotions "ecom_promotion_v2/internal/services/promotions_service"

	"github.com/go-redis/redis"
)

type AppServices struct {
	healthcheckservice.HealthCheckService
	loginserervice.LoginService
	promotions.PromotionsService
	categoryservice.CategoryService //
	programPromotionService.ProgramPromotionService
	programcategoryservice.ProgramCategoryService
	programproductservice.ProgramProductService
	programDirect.ProgramDirectService

	programDirect.ProgramRedirectService
}

func NewAppServices(
	repo *repositories.Repositories,
	rdbToken *redis.Client,
	rdbCache *redis.Client,
) *AppServices {
	jwtSecret := "my-secret-key" //mẫu , nên tạo 1 env
	return &AppServices{

		healthcheckservice.NewHealthCheckService(repo),
		loginserervice.NewLoginService(repo, jwtSecret),
		promotions.NewPromotionsService(repo),
		categoryservice.NewCategoryService(repo), // Kiểm tra kiểu trả về
		programPromotionService.NewProgramPromotionService(repo),
		programcategoryservice.NewProgramCategoryService(repo), // Kiểm tra kiểu trả về
		programproductservice.NewProgramProductService(repo),
		programDirect.NewProgramDirectService(repo),
		programDirect.NewProgramRedirectService(repo),
	}
}
