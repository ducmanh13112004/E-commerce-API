package cmd

import (
	"ecom_promotion_v2/delivery"
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/models"
	"ecom_promotion_v2/internal/repositories"
	"ecom_promotion_v2/internal/services"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	fiber_pprof "github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func FreeOSMemory() {
	for {
		debug.FreeOSMemory()
		internal.Log.Info("FreeOSMemory", zap.Any("NumGoroutine", runtime.NumGoroutine()), zap.Any("NumCPU", runtime.NumCPU()))
		time.Sleep(10 * time.Second)
	}
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		// init timezone
		os.Setenv("TZ", "Asia/Ho_Chi_Minh")

		go FreeOSMemory()
		// app setting
		// utils.StartTelegramService()
		// go utils.SendTelegramMessage("Start "+internal.Envs.HostName, 1)
		// utils_call.UploadFramedImage("test.jpg", "9291918.png", "output/output.png", "miniotests")

		repositories := repositories.NewRepositories(internal.Db.Debug())

		fmt.Println("INIT SERVICE...")
		// Init service
		// Compose into one app service
		appSvcs := services.NewAppServices(repositories)
		// Init handler
		appHandler := delivery.NewAppHandlers(appSvcs)

		models.Validate = validator.New()
		// Hẹn lịch lại các deal sau khi chạy lại service
		AppServer := fiber.New()
		AppServer.Use(cors.New(cors.Config{
			AllowOrigins: "*",
		}))
		AppServer.Use(logger.New(logger.Config{
			Format: "\nEND API\n",
			Done: func(c *fiber.Ctx, logString []byte) {
				var body interface{}
				var response interface{}
				c.BodyParser(&body)
				json.Unmarshal(c.Response().Body(), &response)
				internal.Log.InfoNoConsole("Request processed",
					zap.Any("method", c.Method()), zap.Any("path", c.Path()),
					zap.Any("header", c.GetReqHeaders()), zap.Any("body", body), zap.Any("param_query", c.Queries()),
					zap.Any("latency", time.Since(c.Context().Time())), zap.Any("status", c.Response().StatusCode()), zap.Any("response_body", response))
			},
		}))
		loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
		s := gocron.NewScheduler(loc)
		s.StartAsync()
		fmt.Println("INIT ROUTE")
		// * Serve APIs

		AppServer.Get("/hi-ecom-promotion-v2-api/metrics", monitor.New())
		AppServer.Use(fiber_pprof.New(fiber_pprof.Config{Prefix: "/hi-ecom-promotion-v2-api"}))
		AppServer.Post("/hi-ecom-promotion-v2-api/v1/test-api", TestAPI)
		//
		AppServer.Post("/hi-ecom-promotion-v2-api/test/bottelegram",
			appHandler.RequireXKeyTelegramPortal, appHandler.BotTelegramHandlers)

		AppServer.Post("/hi-ecom-promotion-v2-api/v1/test/lib/telegram", appHandler.RequireXKeyTelegramPortal, appHandler.BotTelegramLibraryHandlers)
		// AppServer.Post("/hi-ecom-promotion-v2-api/v1/image/merge", appHandler.MergeImageHandler)

		AppServer = PromotionsGroup(AppServer, appHandler)
		//
		AppServer = ProgramDirectGroup(AppServer, appHandler)
		AppServer = ProgramDirectUpdateGroup(AppServer, appHandler)
		AppServer = CategoryGroup(AppServer, appHandler)
		AppServer = CategoryProgramGroup(AppServer, appHandler)
		AppServer = ProgramPromotionGroup(AppServer, appHandler)
		AppServer = ProgramProductGroup(AppServer, appHandler)

		fmt.Println("INIT ROUTE SUCCESS")
		if err := AppServer.Listen(":3000"); err != nil {
			fmt.Println("Fiber server got error ", err)
		}
		fmt.Println("Error config!!!!!!! ")
	},
}

func CategoryProgramGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/on-off-program-category", appHandler.RequireTokenLocal, appHandler.OnOffProgramCategoryListAllHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/program-category-list-all", appHandler.RequireTokenLocal, appHandler.GetProgramCategoryListAllHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/program-category-list-on", appHandler.RequireTokenLocal, appHandler.GetProgramCategoryHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/program-category-add", appHandler.RequireTokenLocal, appHandler.CreateProgramCategory)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/program-category-update", appHandler.RequireTokenLocal, appHandler.UpdateProgramCategory)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/program-category-delete", appHandler.RequireTokenLocal, appHandler.DeleteProgramCategory)
	return AppServer
}

func CategoryGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	//khai báo CreateCategory trong app_handlers để sử dụng

	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/category-list", appHandler.RequireTokenLocal, appHandler.GetCategoryListHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/category-add", appHandler.RequireTokenLocal, appHandler.CreateCategory)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/category-update", appHandler.RequireTokenLocal, appHandler.UpdateCategory)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/category-delete", appHandler.RequireTokenLocal, appHandler.DeleteCategory)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/category/category-list-name", appHandler.RequireTokenLocal, appHandler.GetCategoryListNameHandler)
	return AppServer
}
func ProgramDirectGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/ProgramDirect/ProgramDirect-add", appHandler.RequireTokenLocal, appHandler.CreateProgramDirect)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/ProgramDirect/ProgramDirect-list", appHandler.RequireTokenLocal, appHandler.GetProgramDirectListAllHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/ProgramDirect/ProgramDirect-list-ID", appHandler.RequireTokenLocal, appHandler.GetProgramDirectIDHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/ProgramDirect/ProgramDirect-Update", appHandler.RequireTokenLocal, appHandler.UpdateProgramDirect)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/ProgramDirect/ProgramDirect-On-Off-State", appHandler.RequireTokenLocal, appHandler.DeleteProgramDirect)
	return AppServer
}

// update api theo bảng db mới!
func ProgramDirectUpdateGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/program-direct-list-name", appHandler.RequireTokenLocal, appHandler.GetProgramRedirectListNameAllHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/program-direct-add", appHandler.RequireTokenLocal, appHandler.CreateProgramRedirectHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/program-direct-list", appHandler.RequireTokenLocal, appHandler.GetProgramRedirectListAllHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/program-direct-list-id", appHandler.RequireTokenLocal, appHandler.GetProgramRedirectIDHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/program-direct-update", appHandler.RequireTokenLocal, appHandler.UpdateProgramRedirectHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/program-direct-on-off-state", appHandler.RequireTokenLocal, appHandler.DeleteProgramRedirectHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/get-program-list-navigation", appHandler.RequireTokenLocal, appHandler.GetProgramListNavigationHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/getprogramlistnavigationcreate", appHandler.RequireTokenLocal, appHandler.CreateProgramNavigationHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-direct-update-group/getprogramlistnavigationupdate", appHandler.RequireTokenLocal, appHandler.UpdateProgramNavigationHandler)
	return AppServer
}
func PromotionsGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	// API Cấp voucher cho user
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/promotion/receive-voucher-fr-program-id", appHandler.RequireTokenLocal, appHandler.LocalReceiveVoucherFrProgramId)

	// Danh sách chương trình theo danh mục

	return AppServer
}

func ProgramPromotionGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/get-program-promotion-list", appHandler.RequireTokenLocal, appHandler.GetProgramPromotionListHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/get-program-promotion-by-id", appHandler.RequireTokenLocal, appHandler.GetProgramPromotionByIDHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/create-program-promotion", appHandler.RequireTokenLocal, appHandler.CreateProgramPromotionHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/update-program-promotion", appHandler.RequireTokenLocal, appHandler.UpdateProgramPromotionHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/on-off-program-promotion", appHandler.RequireTokenLocal, appHandler.OnOffProgramPromotionHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/delete-program-promotion", appHandler.RequireTokenLocal, appHandler.DeleteProgramPromotionHandler)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/get-report-program-promotion", appHandler.RequireTokenLocal, appHandler.GetReportProgramPromotionHandler)
	//Promotion code
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/promotion/create-promotion", appHandler.RequireTokenLocal, appHandler.CreatePromotion)
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/promotion/get-promotion", appHandler.RequireTokenLocal, appHandler.GetPromotionByProgramId)

	//Call local
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-promotion/get-report-program-promotion-local", appHandler.RequireTokenLocal, appHandler.GetReportProgramPromotionLocalHandler)

	// Call program priority
	return AppServer
}

func ProgramProductGroup(AppServer *fiber.App, appHandler delivery.AppHandlers) *fiber.App {
	AppServer.Post("/hi-ecom-promotion-v2-api/v1/program-product/create-list-sku", appHandler.RequireTokenLocal, appHandler.CreateListSkuHandler)
	return AppServer
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
