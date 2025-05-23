package scenario

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bacalhau-project/bacalhau/pkg/downloader"

	"github.com/bacalhau-project/bacalhau/pkg/models"

	"github.com/bacalhau-project/bacalhau/pkg/executor"
	"github.com/bacalhau-project/bacalhau/pkg/executor/noop"
	"github.com/bacalhau-project/bacalhau/pkg/system"
)

func noopScenario(t testing.TB) Scenario {
	return Scenario{
		Stack: &StackConfig{
			ExecutorConfig: noop.ExecutorConfig{
				ExternalHooks: noop.ExecutorConfigExternalHooks{
					JobHandler: func(ctx context.Context, execContext noop.ExecutionContext) (*models.RunCommandResult, error) {
						return executor.WriteJobResults(execContext.ExecutionDir, strings.NewReader("hello, world!\n"), nil, 0, nil, executor.OutputLimits{
							MaxStdoutFileLength:   system.MaxStdoutFileLength,
							MaxStdoutReturnLength: system.MaxStdoutReturnLength,
							MaxStderrFileLength:   system.MaxStderrFileLength,
							MaxStderrReturnLength: system.MaxStderrReturnLength,
						}), nil
					},
				},
			},
		},
		Job: &models.Job{
			Name:  t.Name(),
			Type:  models.JobTypeBatch,
			Count: 1,
			Tasks: []*models.Task{
				{
					Name: t.Name(),
					Engine: &models.SpecConfig{
						Type:   models.EngineNoop,
						Params: make(map[string]interface{}),
					},
				},
			},
		},
		ResultsChecker: FileEquals(downloader.DownloadFilenameStdout, "hello, world!\n"),
		JobCheckers:    WaitUntilSuccessful(1),
	}
}

type NoopTest struct {
	ScenarioRunner
}

func Example_noop() {
	// In a real example, use the testing.T passed to the TestXxx method.
	suite.Run(&testing.T{}, new(NoopTest))
}

func (suite *NoopTest) TestRunNoop() {
	suite.RunScenario(noopScenario(suite.T()))
}
