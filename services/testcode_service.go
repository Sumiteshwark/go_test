package services

import (
	"backend/models"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"

	"strconv"
	"strings"
)

type TestCodeService struct {
	db *gorm.DB
}

func NewTestCodeService(db *gorm.DB) *TestCodeService {
	return &TestCodeService{db: db}
}

// func (s *TestCodeService) GetAllOrgs() ([]string, error) {
// 	var orgs []models.TestCodeDataSet
// 	var orgNames []string
// 	orgMap := make(map[string]struct{})
// 	if err := s.db.Select("org_name").Find(&orgs).Error; err != nil {
// 		return nil, fmt.Errorf("failed to get all orgs: %v", err)
// 	}
// 	for _, dataset := range orgs {
// 		if dataset.OrgName != nil {
// 			orgName := *dataset.OrgName
// 			_, exists := orgMap[orgName]
// 			if !exists {
// 				orgMap[orgName] = struct{}{}
// 				orgNames = append(orgNames, orgName)
// 			}
// 		}
// 	}

// 	return orgNames, nil
// }

func (s *TestCodeService) GetAllOrgs() ([]string, error) {
	var orgNames []string
	if err := s.db.Model(&models.TestCodeDataSet{}).Distinct("org_name").Pluck("org_name", &orgNames).Error; err != nil {
		return nil, fmt.Errorf("failed to get all orgs: %v", err)
	}

	return orgNames, nil
}

func (s *TestCodeService) GetOrgDataByName(orgNames []string) ([]models.TestCodeDataSet, error) {
	var orgs []models.TestCodeDataSet
	query := s.db.Where("iteration = -1")
	if len(orgNames) > 0 {
		query = query.Where("org_name IN (?)", orgNames)
	}
	if err := query.Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("could not find data for organizations: %v", err)
	}
	return orgs, nil
}

func (s *TestCodeService) GetAvgTestGenerationTimeByOrgName(orgNames []string) (float64, error) {
	// For each request_id, find the created_at time for iteration 0 and iteration -1
	// Then average time is (sum of times of ('iteration -1' - 'iteration 0') for each request_id / total number of request_id)
	type RequestTime struct {
		RequestID    string
		IterZero     *time.Time
		IterMinusOne *time.Time
	}

	var requestTimes []RequestTime
	var query string
	var args []interface{}

	if len(orgNames) == 0 {
		query = `
			SELECT request_id,
			MIN(CASE WHEN iteration = 0 THEN created_at END) AS iter_zero,
			MAX(CASE WHEN iteration = -1 THEN created_at END) AS iter_minus_one
			FROM test_code_data_set
			GROUP BY request_id
		`
	} else {
		query = `
			SELECT request_id,
			MIN(CASE WHEN iteration = 0 THEN created_at END) AS iter_zero,
			MAX(CASE WHEN iteration = -1 THEN created_at END) AS iter_minus_one
			FROM test_code_data_set
			WHERE org_name IN (?)
			GROUP BY request_id
		`
		args = append(args, orgNames)
	}

	if err := s.db.Raw(query, args...).Scan(&requestTimes).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch request times: %v", err)
	}

	if len(requestTimes) == 0 {
		return 0, nil
	}

	var totalTime float64
	var requestCount int

	for _, req := range requestTimes {
		if req.IterZero != nil && req.IterMinusOne != nil {
			totalTime += req.IterMinusOne.Sub(*req.IterZero).Seconds()
			requestCount++
		}
	}

	if requestCount == 0 {
		return 0, nil
	}

	averageTime := totalTime / float64(requestCount)
	averageTime = math.Round(averageTime*100) / 100

	return averageTime, nil
}

func (s *TestCodeService) GetRunnabilityAndCoverageByOrgName(orgNames []string) (int, int, int, float64, error) {
	// Total number of request_id == Total number of -1	iteration
	// Total number of runnable request_id == Total number of runnable -1 iteration
	// Total number of non-runnable request_id == Total number of non-runnable -1 iteration
	// Average coverage == Average coverage of runnable -1 iteration (i.e. sum of coverage of runnable -1 iteration / total number of runnable -1 iteration)
	var orgs []models.TestCodeDataSet
	var query string
	var args []interface{}

	if len(orgNames) == 0 {
		query = "SELECT * FROM test_code_data_set WHERE iteration = -1"
	} else {
		query = "SELECT * FROM test_code_data_set WHERE iteration = -1 AND org_name IN (?)"
		args = append(args, orgNames)
	}

	if err := s.db.Raw(query, args...).Find(&orgs).Error; err != nil {
		return 0, 0, 0, 0, fmt.Errorf("could not fetch organizations: %v", err)
	}

	iterationMinusOneCount := len(orgs)
	runnableCount, nonRunnableCount := 0, 0
	var totalCoverage float64

	for _, org := range orgs {
		runnable := org.Runnability != nil && *org.Runnability

		if runnable {
			runnableCount++
			if org.Coverage != nil {
				coverageStr := strings.TrimSuffix(*org.Coverage, "%")
				if percentageFloat, err := strconv.ParseFloat(coverageStr, 64); err == nil {
					totalCoverage += percentageFloat
				} else {
					fmt.Printf("Warning: Invalid coverage value for request_id %s\n", *org.RequestID)
				}
			}
		} else {
			nonRunnableCount++
		}
	}
	averageCoverage := 0.0
	if runnableCount > 0 {
		averageCoverage = totalCoverage / float64(runnableCount)
		averageCoverage = math.Round(averageCoverage*100) / 100
	}

	return iterationMinusOneCount, runnableCount, nonRunnableCount, averageCoverage, nil
}

func (s *TestCodeService) GetTotalRunnableAndNonRunnableLines(orgNames []string) (int, int, error) {
	// Total number of request_id == Total number of -1	iteration
	// Total number of runnable lines = Total number of lines in runnable -1 iteration
	// Total number of non-runnable lines = Total number of lines in non-runnable -1 iteration
	type testStruct struct {
		TestCode    string
		Runnability *bool
	}

	var testDataList []testStruct
	var query string
	var args []interface{}

	if len(orgNames) == 0 {
		query = "SELECT test_code, runnability FROM test_code_data_set WHERE iteration = -1"
	} else {
		query = "SELECT test_code, runnability FROM test_code_data_set WHERE iteration = -1 AND org_name IN (?)"
		args = append(args, orgNames)
	}

	if err := s.db.Raw(query, args...).Scan(&testDataList).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch test_code data: %v", err)
	}

	totalRunnableLines := 0
	totalNonRunnableLines := 0

	for _, data := range testDataList {
		var unquotedData string

		if strings.HasPrefix(data.TestCode, `"`) && strings.HasSuffix(data.TestCode, `"`) {
			var err error
			unquotedData, err = strconv.Unquote(data.TestCode)
			if err != nil {
				unquotedData = data.TestCode
			}
		} else {
			unquotedData = data.TestCode
		}

		lineCount := strings.Count(unquotedData, "\n")

		if data.Runnability != nil && *data.Runnability {
			totalRunnableLines += lineCount
		} else {
			totalNonRunnableLines += lineCount
		}
	}

	return totalRunnableLines, totalNonRunnableLines, nil
}

func (s *TestCodeService) GetTotalTestCaseGeneratedInIntervalsByOrgName(startDate time.Time, endDate time.Time, orgNames []string) ([]map[string]interface{}, error) {
	currentTime := time.Now()
	if endDate.After(currentTime) {
		endDate = currentTime
	} else {
		endDate = endDate.Truncate(24 * time.Hour).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	if endDate.After(currentTime) {
		endDate = currentTime
	}

	//// to handle timezone
	// startDate = startDate.Local()
	// endDate = endDate.Local()
	// currentTime := time.Now().Local()
	// adjustedEndDate := endDate.Add(24 * time.Hour)
	// if adjustedEndDate.After(currentTime) {
	// 	endDate = currentTime
	// } else {
	// 	endDate = adjustedEndDate
	// }

	// fmt.Println("startDate: ", startDate)
	// fmt.Println("endTime: ", endDate)
	totalDuration := endDate.Sub(startDate)
	// fmt.Println("totalDuration: ", totalDuration)

	interval := totalDuration / 7

	var results []map[string]interface{}

	for i := 0; i < 7; i++ {
		intervalStart := startDate.Add(time.Duration(i) * interval)
		intervalEnd := startDate.Add(time.Duration(i+1) * interval)
		// var count int64
		// s.db.Model(&models.TestCodeDataSet{}).Where("created_at >= ? AND created_at < ? AND org_name = ?", intervalStart, intervalEnd, orgName).Count(&count)

		// results = append(results, map[string]interface{}{
		// 	"interval_start": intervalStart.Format("2006-01-02 15:04:05"),
		// 	"interval_end":   intervalEnd.Format("2006-01-02 15:04:05"),
		// 	"count":          count,
		// })

		var count int64
		query := s.db.Model(&models.TestCodeDataSet{}).Where("created_at >= ? AND created_at < ?", intervalStart, intervalEnd)

		if len(orgNames) > 0 {
			query = query.Where("org_name IN (?)", orgNames)
		}

		query.Count(&count)

		results = append(results, map[string]interface{}{
			"interval_start": intervalStart.Format("2006-01-02 15:04:05"),
			"interval_end":   intervalEnd.Format("2006-01-02 15:04:05"),
			"count":          count,
		})
	}

	return results, nil
}
