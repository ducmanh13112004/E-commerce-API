package healthcheckservice

import (
	"ecom_promotion_v2/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckService interface {
	HealthCheck(ctx *fiber.Ctx) (int, interface{})
}
type healthCheckService struct {
	repo *repositories.Repositories
}

func NewHealthCheckService(repo *repositories.Repositories) HealthCheckService {
	return &healthCheckService{
		repo: repo,
	}
}
