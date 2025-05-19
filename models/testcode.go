package models

import (
	"time"
)

type TestCodeDataSet struct {
	ID                            uint      `gorm:"primaryKey"`
	FileName                      *string   `gorm:"column:file_name"`
	FunctionName                  *string   `gorm:"column:function_name"`
	RepoName                      *string   `gorm:"column:repo_name"`
	OrgName                       *string   `gorm:"column:org_name"`
	SrcCode                       *string   `gorm:"column:src_code"`
	Context                       *string   `gorm:"column:context"`
	Error                         *string   `gorm:"column:error"`
	ErrorCategory                 *string   `gorm:"column:error_category"`
	Iteration                     *int      `gorm:"column:iteration"`
	TestCode                      *string   `gorm:"column:test_code"`
	Runnability                   *bool     `gorm:"column:runnability"`
	Coverage                      *string   `gorm:"column:coverage"`
	Comment                       *string   `gorm:"column:comment"`
	CreatedAt                     time.Time `gorm:"column:created_at"`
	UpdatedAt                     time.Time `gorm:"column:updated_at"`
	IsTestCodeComplete            *bool     `gorm:"column:is_test_code_complete"`
	RequestID                     *string   `gorm:"column:request_id"`
	InstructionToFix              *string   `gorm:"column:instruction_to_fix"`
	ModelUsedForInitialGeneration *string   `gorm:"column:model_used_for_initial_generation"`
	ModelUsedForAnalysis          *string   `gorm:"column:model_used_for_analysis"`
	ModelUsedForCorrection        *string   `gorm:"column:model_used_for_correction"`
	BuildTime                     *float64  `gorm:"column:build_time"`
	TestCodeGenerationTime        *float64  `gorm:"column:test_code_generation_time"`
}

func (TestCodeDataSet) TableName() string {
	return "test_code_data_set"
}
