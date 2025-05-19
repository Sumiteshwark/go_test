package services

import (
	"backend/models"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a test database with some sample data
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the model
	if err := db.AutoMigrate(&models.TestCodeDataSet{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestGetAllOrgs(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	service := NewTestCodeService(db)

	// Create test data
	// Add your test data setup here

	// Test cases
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty database",
			want:    []string{},
			wantErr: false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetAllOrgs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllOrgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetAllOrgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAvgTestGenerationTimeByOrgName(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	service := NewTestCodeService(db)

	// Create test data
	// Add your test data setup here

	// Test cases
	tests := []struct {
		name     string
		orgNames []string
		want     float64
		wantErr  bool
	}{
		{
			name:     "empty org names",
			orgNames: []string{},
			want:     0,
			wantErr:  false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetAvgTestGenerationTimeByOrgName(tt.orgNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvgTestGenerationTimeByOrgName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAvgTestGenerationTimeByOrgName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRunnabilityAndCoverageByOrgName(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	service := NewTestCodeService(db)

	// Create test data
	// Add your test data setup here

	// Test cases
	tests := []struct {
		name     string
		orgNames []string
		want     struct {
			totalCount       int
			runnableCount    int
			nonRunnableCount int
			avgCoverage      float64
		}
		wantErr bool
	}{
		{
			name:     "empty org names",
			orgNames: []string{},
			want: struct {
				totalCount       int
				runnableCount    int
				nonRunnableCount int
				avgCoverage      float64
			}{
				totalCount:       0,
				runnableCount:    0,
				nonRunnableCount: 0,
				avgCoverage:      0,
			},
			wantErr: false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, runnable, nonRunnable, coverage, err := service.GetRunnabilityAndCoverageByOrgName(tt.orgNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRunnabilityAndCoverageByOrgName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if total != tt.want.totalCount || runnable != tt.want.runnableCount ||
				nonRunnable != tt.want.nonRunnableCount || coverage != tt.want.avgCoverage {
				t.Errorf("GetRunnabilityAndCoverageByOrgName() = %v, %v, %v, %v, want %v, %v, %v, %v",
					total, runnable, nonRunnable, coverage,
					tt.want.totalCount, tt.want.runnableCount, tt.want.nonRunnableCount, tt.want.avgCoverage)
			}
		})
	}
}
