package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bacalhau-project/bacalhau/pkg/models"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/apimodels"
	publisher_local "github.com/bacalhau-project/bacalhau/pkg/publisher/local"
	"github.com/bacalhau-project/bacalhau/pkg/test/mock"

	"github.com/bacalhau-project/bacalhau/pkg/compute"
	"github.com/bacalhau-project/bacalhau/pkg/compute/env"
	"github.com/bacalhau-project/bacalhau/pkg/devstack"
	_ "github.com/bacalhau-project/bacalhau/pkg/logger"
	"github.com/bacalhau-project/bacalhau/pkg/test/scenario"
	testutils "github.com/bacalhau-project/bacalhau/pkg/test/teststack"
)

const testNodeCount = 1

//nolint:funlen,mnd
func RunTestCase(
	t *testing.T,
	testCase scenario.Scenario,
) {
	ctx := context.Background()
	job := testCase.Job

	var devstackOptions []devstack.ConfigOption
	if testCase.Stack != nil && testCase.Stack.DevStackOptions != nil {
		devstackOptions = testCase.Stack.DevStackOptions
	}
	devstackOptions = append(devstackOptions, devstack.WithNumberOfHybridNodes(testNodeCount))

	stack := testutils.Setup(ctx, t, devstackOptions...)
	executor, err := stack.Nodes[0].ComputeNode.Executors.Get(ctx, job.Task().Engine.Type)
	require.NoError(t, err)

	isInstalled, err := executor.IsInstalled(ctx)
	require.NoError(t, err)
	require.True(t, isInstalled)

	prepareStorage := func(getStorage scenario.SetupStorage) []*models.InputSource {
		if getStorage == nil {
			return []*models.InputSource{}
		}

		storageList, stErr := getStorage(ctx)
		require.NoError(t, stErr)

		for _, input := range storageList {
			strger, err := stack.Nodes[0].ComputeNode.Storages.Get(ctx, input.Source.Type)
			require.NoError(t, err)

			hasStorage, stErr := strger.HasStorageLocally(ctx, *input)
			require.NoError(t, stErr)

			require.True(t, hasStorage)
		}

		return storageList
	}

	job.Task().InputSources = append(job.Task().InputSources, prepareStorage(testCase.Inputs)...)
	job.Task().ResultPaths = append(job.Task().ResultPaths, testCase.Outputs...)
	job.Task().Publisher = publisher_local.NewSpecConfig()
	job.Count = testNodeCount

	now := time.Now().UnixNano()
	job.ID = fmt.Sprintf("test-job-%d", now)
	job.Meta = make(map[string]string, 2)
	job.Meta[models.MetaOrchestratorID] = "test-owner"
	job.State = models.NewJobState(models.JobStateTypeUndefined)
	job.Version = 1
	job.Revision = 1
	job.CreateTime = now
	job.ModifyTime = now

	job.Normalize()
	require.NoError(t, job.Validate())

	execution := mock.ExecutionForJob(job)
	execution.AllocateResources(job.Task().Name, models.Resources{})

	resultsPath, err := compute.NewResultsPath(t.TempDir())
	require.NoError(t, err)
	executionDir, err := resultsPath.PrepareExecutionOutputDir(execution.ID)
	require.NoError(t, err)
	storageProvider := stack.Nodes[0].ComputeNode.Storages

	// Create a permissive env resolver for testing
	envResolver := env.NewResolver(env.ResolverParams{
		AllowList: []string{"TEST_*"}, // Allow only TEST_* environment variables to be forwarded
	})
	portAllocator, err := compute.NewPortAllocator(10000, 20000)
	require.NoError(t, err)

	// TODO: use base executor instead of directly using the executor
	baseExecutor := compute.NewBaseExecutor(compute.BaseExecutorParams{
		ID:               "test-executor",
		Storages:         storageProvider,
		StorageDirectory: t.TempDir(),
		EnvResolver:      envResolver,
		PortAllocator:    portAllocator,
	})

	runCommandArguments, cleanup, err := baseExecutor.PrepareRunArguments(
		ctx,
		execution,
		executionDir,
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := cleanup(ctx); err != nil {
			t.Fatal(err)
		}
	})

	result, err := executor.Run(ctx, runCommandArguments)
	if testCase.SubmitChecker != nil {
		// TODO not sure how this behavior should be replicated.
		err = testCase.SubmitChecker(&apimodels.PutJobResponse{
			JobID:        job.ID,
			EvaluationID: "TODO",
			Warnings:     nil,
		}, err)
		require.NoError(t, err)
	}

	if testCase.CommandResultsChecker != nil {
		err = testCase.CommandResultsChecker(result)
		require.NoError(t, err)
	}

	if testCase.ResultsChecker != nil {
		err = testCase.ResultsChecker(compute.ExecutionResultsDir(executionDir))
		require.NoError(t, err)
	}
}
