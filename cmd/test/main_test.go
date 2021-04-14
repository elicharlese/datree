package test

import (
	"path/filepath"
	"testing"

	"github.com/datreeio/datree/bl"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/datreeio/datree/pkg/propertiesExtractor"
	"github.com/stretchr/testify/mock"
)

type mockLocalConfigManager struct {
	mock.Mock
}

func (m *mockLocalConfigManager) GetConfiguration() (localConfig.LocalConfiguration, error) {
	args := m.Called()
	return args.Get(0).(localConfig.LocalConfiguration), args.Error(1)
}

type mockEvaluator struct {
	mock.Mock
}

func (m *mockEvaluator) PrintResults(results *bl.EvaluationResults, cliId string) error {
	m.Called(results, cliId)
	return nil
}

func (m *mockEvaluator) PrintFileParsingErrors(errors []propertiesExtractor.FileError) {
	m.Called(errors)
}

func (m *mockEvaluator) Evaluate(pattern string, cliId string, evaluationConc int) (*bl.EvaluationResults, []propertiesExtractor.FileError, error) {
	args := m.Called(pattern, cliId, evaluationConc)
	return args.Get(0).(*bl.EvaluationResults), args.Get(1).([]propertiesExtractor.FileError), args.Error(2)
}
func TestTestCommand(t *testing.T) {
	evaluator := &mockEvaluator{}
	mockedEvaluateResponse := &bl.EvaluationResults{
		FileNameRuleMapper: map[string]map[int]*bl.Rule{}, Summary: struct {
			RulesCount       int
			TotalFailedRules int
			FilesCount       int
		}{RulesCount: 1, TotalFailedRules: 0, FilesCount: 0},
	}
	evaluator.On("Evaluate", mock.Anything, mock.Anything, mock.Anything).Return(mockedEvaluateResponse, []propertiesExtractor.FileError{}, nil)
	evaluator.On("PrintFileParsingErrors", mock.Anything).Return()
	evaluator.On("PrintResults", mock.Anything, mock.Anything).Return()

	localConfigManager := &mockLocalConfigManager{}
	localConfigManager.On("GetConfiguration").Return(localConfig.LocalConfiguration{CliId: "134kh"}, nil)

	ctx := &TestCommandContext{
		Evaluator:   evaluator,
		LocalConfig: localConfigManager,
	}

	test(ctx, "8/*")
	localConfigManager.AssertCalled(t, "GetConfiguration")

	expectedPath, _ := filepath.Abs("8/*")
	evaluator.AssertCalled(t, "Evaluate", expectedPath, "134kh", 50)
	evaluator.AssertNotCalled(t, "PrintFileParsingErrors")
	evaluator.AssertCalled(t, "PrintResults", mockedEvaluateResponse, "134kh")
}