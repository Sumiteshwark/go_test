package controllers

import (
	"backend/services"
	"backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TestCodeController struct {
	service *services.TestCodeService
}

func NewTestCodeController(service *services.TestCodeService) *TestCodeController {
	return &TestCodeController{service: service}
}

func (c *TestCodeController) GetAllOrgs(ctx *gin.Context) {
	orgNames, err := c.service.GetAllOrgs()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, orgNames)
}

type OrgRequests struct {
	OrgNames []string `json:"org_names" binding:"required"`
}

func (c *TestCodeController) GetOrgDataByName(ctx *gin.Context) {
	var req OrgRequests
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	orgs, err := c.service.GetOrgDataByName(req.OrgNames)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, orgs)
}

func (c *TestCodeController) GetAvgTestGenerationTimeByOrgName(ctx *gin.Context) {
	var req OrgRequests
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	avgTime, err := c.service.GetAvgTestGenerationTimeByOrgName(req.OrgNames)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, avgTime)
}

func (c *TestCodeController) GetRunnabilityAndCoverageByOrgName(ctx *gin.Context) {
	var req OrgRequests
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	iterationMinusOneCount, runnableCount, nonRunnableCount, coverage, err := c.service.GetRunnabilityAndCoverageByOrgName(req.OrgNames)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"total_tests":              iterationMinusOneCount,
		"total_runnable_tests":     runnableCount,
		"total_non_runnable_tests": nonRunnableCount,
		"coverage":                 coverage,
	})
}

func (c *TestCodeController) GetTotalRunnableAndNonRunnableLines(ctx *gin.Context) {
	var req OrgRequests
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	runnableLines, nonRunnableLines, err := c.service.GetTotalRunnableAndNonRunnableLines(req.OrgNames)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"runnable_lines":     runnableLines,
		"non_runnable_lines": nonRunnableLines,
	})
}

type dateInput struct {
	StartDate string   `json:"start_date" binding:"required"`
	EndDate   string   `json:"end_date" binding:"required"`
	OrgNames  []string `json:"org_names" binding:"required"`
}

func (c *TestCodeController) GetTotalTestCaseGeneratedInIntervalsByOrgName(ctx *gin.Context) {
	var input dateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	startDateStr := input.StartDate
	endDateStr := input.EndDate
	OrgNames := input.OrgNames

	startDate, err1 := time.Parse("2006-01-02", startDateStr)
	endDate, err2 := time.Parse("2006-01-02", endDateStr)

	if err1 != nil || err2 != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid date format")
		return
	}

	results, err := c.service.GetTotalTestCaseGeneratedInIntervalsByOrgName(startDate, endDate, OrgNames)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, results)
}
