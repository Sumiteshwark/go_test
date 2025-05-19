package routes

import (
	"backend/controllers"
	// "backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) {
	// db.AutoMigrate(&models.TestCodeDataSet{})

	userService := services.NewTestCodeService(db)
	userController := controllers.NewTestCodeController(userService)

	userRoutes := r.Group("api/orgs")
	{
		userRoutes.GET("/", userController.GetAllOrgs)
		userRoutes.POST("/", userController.GetOrgDataByName)
		userRoutes.POST("/avgtime", userController.GetAvgTestGenerationTimeByOrgName)
		userRoutes.POST("/runnability_coverage", userController.GetRunnabilityAndCoverageByOrgName)
		userRoutes.POST("/runnability_lines", userController.GetTotalRunnableAndNonRunnableLines)
		userRoutes.POST("/request_counts", userController.GetTotalTestCaseGeneratedInIntervalsByOrgName)
	}
}
