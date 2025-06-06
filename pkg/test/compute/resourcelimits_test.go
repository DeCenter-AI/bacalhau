//go:build integration || !unit

package compute

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/bacalhau-project/bacalhau/pkg/config"
	"github.com/bacalhau-project/bacalhau/pkg/config/types"
	"github.com/bacalhau-project/bacalhau/pkg/devstack"
	noop_executor "github.com/bacalhau-project/bacalhau/pkg/executor/noop"
	"github.com/bacalhau-project/bacalhau/pkg/logger"
	"github.com/bacalhau-project/bacalhau/pkg/models"
	"github.com/bacalhau-project/bacalhau/pkg/node"
	"github.com/bacalhau-project/bacalhau/pkg/orchestrator"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/apimodels"
	clientv2 "github.com/bacalhau-project/bacalhau/pkg/publicapi/client/v2"
	"github.com/bacalhau-project/bacalhau/pkg/system"
	"github.com/bacalhau-project/bacalhau/pkg/test/scenario"
	"github.com/bacalhau-project/bacalhau/pkg/test/teststack"
	nodeutils "github.com/bacalhau-project/bacalhau/pkg/test/utils/node"
)

type ComputeNodeResourceLimitsSuite struct {
	suite.Suite
}

func TestComputeNodeResourceLimitsSuite(t *testing.T) {
	suite.Run(t, new(ComputeNodeResourceLimitsSuite))
}

func (suite *ComputeNodeResourceLimitsSuite) SetupTest() {
	logger.ConfigureTestLogging(suite.T())
}

type SeenJobRecord struct {
	Id          string
	CurrentJobs int
	MaxJobs     int
	Start       int64
	End         int64
}

type TotalResourceTestCaseCheck struct {
	name    string
	handler func(seenJobs []SeenJobRecord) (bool, error)
}

type TotalResourceTestCase struct {
	// the total list of jobs to throw at the cluster all at the same time
	jobs        []*models.ResourcesConfig
	totalLimits *models.ResourcesConfig
	wait        TotalResourceTestCaseCheck
	checkers    []TotalResourceTestCaseCheck
}

func (suite *ComputeNodeResourceLimitsSuite) TestTotalResourceLimits() {
	// for this test we use the transport so the compute_node is calling
	// the executor in a go-routine and we can test what jobs
	// look like over time - this test leave each job running for X seconds
	// and consuming Y resources
	// we will have set a total amount of resources on the compute_node
	// and we want to see that the following things are true:
	//
	//  * all jobs ran eventually (because there is no per job limit and no one job is bigger than the total limit)
	//  * at no time - the total job resource usage exceeds the configured total
	//  * we submit all the jobs at the same time so we prove that compute_nodes "back bid"
	//    * i.e. a job that was seen 20 seconds ago we now have space to run so let's bid on it now
	//
	runTest := func(
		testCase TotalResourceTestCase,
	) {
		ctx := context.Background()
		epochSeconds := time.Now().Unix()

		var seenJobs []SeenJobRecord
		var seenJobsMutex sync.Mutex

		addSeenJob := func(job SeenJobRecord) {
			seenJobsMutex.Lock()
			defer seenJobsMutex.Unlock()
			seenJobs = append(seenJobs, job)
		}

		currentJobCount := 0
		maxJobCount := 0

		// our function that will "execute the job"
		// record time stamps of start and end
		// sleep for a bit to simulate real work happening
		jobHandler := func(_ context.Context, execContext noop_executor.ExecutionContext) (*models.RunCommandResult, error) {
			currentJobCount++
			if currentJobCount > maxJobCount {
				maxJobCount = currentJobCount
			}
			seenJob := SeenJobRecord{
				Id:          execContext.JobID,
				Start:       time.Now().Unix() - epochSeconds,
				CurrentJobs: currentJobCount,
				MaxJobs:     maxJobCount,
			}
			time.Sleep(time.Second * 1)
			currentJobCount--
			seenJob.End = time.Now().Unix() - epochSeconds
			addSeenJob(seenJob)
			return &models.RunCommandResult{}, nil
		}

		getVolumeSizeHandler := func(ctx context.Context, volume models.InputSource) (uint64, error) {
			size, err := datasize.ParseString(volume.Target)
			return size.Bytes(), err
		}

		cfg, err := config.NewTestConfigWithOverrides(types.Bacalhau{
			Compute: types.Compute{
				AllocatedCapacity: types.ResourceScalerFromModelsResourceConfig(*testCase.totalLimits),
			},
		})
		require.NoError(suite.T(), err)

		stack := teststack.Setup(ctx, suite.T(),
			devstack.WithNumberOfHybridNodes(1),
			devstack.WithBacalhauConfigOverride(cfg),
			teststack.WithNoopExecutor(noop_executor.ExecutorConfig{
				ExternalHooks: noop_executor.ExecutorConfigExternalHooks{
					JobHandler:    jobHandler,
					GetVolumeSize: getVolumeSizeHandler,
				},
			}, cfg.Engines),
		)

		for _, jobResources := range testCase.jobs {
			// what the job is doesn't matter - it will only end up
			sanitizedName := strings.Replace(suite.T().Name(), "/", "-", -1)
			jobName := fmt.Sprintf("%s-%d", sanitizedName, rand.Intn(10001))
			j := &models.Job{
				Name:  jobName,
				Type:  models.JobTypeBatch,
				Count: 1,
				Tasks: []*models.Task{
					{
						Name: suite.T().Name(),
						Engine: &models.SpecConfig{
							Type: models.EngineNoop,
						},
						ResourcesConfig: jobResources,
					},
				},
			}
			j.Normalize()
			client := clientv2.New(fmt.Sprintf("http://%s:%d", stack.Nodes[0].APIServer.Address, stack.Nodes[0].APIServer.Port))
			_, err := client.Jobs().Put(ctx, &apimodels.PutJobRequest{
				Job: j,
			})
			require.NoError(suite.T(), err)

			// sleep a bit here to simulate jobs being submitted over time
			time.Sleep((10 + time.Duration(rand.Intn(10))) * time.Millisecond)
		}

		// wait for all the jobs to have completed
		// we can check the seenJobs because that is easier
		waiter := &system.FunctionWaiter{
			Name:        "wait for jobs",
			MaxAttempts: 30,
			Delay:       time.Second * 1,
			Handler: func() (bool, error) {
				return testCase.wait.handler(seenJobs)
			},
		}

		err = waiter.Wait(ctx)
		require.NoError(suite.T(), err, fmt.Sprintf("there was an error in the wait function: %s", testCase.wait.name))

		if err != nil {
			fmt.Printf("error waiting for jobs to have been seen\n")
			spew.Dump(seenJobs)
		}

		checkOk := true
		failingCheckMessage := ""

		for _, checker := range testCase.checkers {
			innerCheck, err := checker.handler(seenJobs)
			errorMessage := ""
			if err != nil {
				errorMessage = fmt.Sprintf("there was an error in the check function: %s %s", checker.name, err.Error())
			}
			require.NoError(suite.T(), err, errorMessage)
			if !innerCheck {
				checkOk = false
				failingCheckMessage = fmt.Sprintf("there was an fail in the check function: %s", checker.name)
			}
		}

		require.True(suite.T(), checkOk, failingCheckMessage)

		if !checkOk {
			fmt.Printf("error checking results on seen jobs\n")
			spew.Dump(seenJobs)
		}
	}

	waitUntilSeenAllJobs := func(expected int) TotalResourceTestCaseCheck {
		return TotalResourceTestCaseCheck{
			name: fmt.Sprintf("waitUntilSeenAllJobs: %d", expected),
			handler: func(seenJobs []SeenJobRecord) (bool, error) {
				return len(seenJobs) >= expected, nil
			},
		}
	}

	checkMaxJobs := func(max int) TotalResourceTestCaseCheck {
		return TotalResourceTestCaseCheck{
			name: fmt.Sprintf("checkMaxJobs: %d", max),
			handler: func(seenJobs []SeenJobRecord) (bool, error) {
				seenMax := 0
				for _, seenJob := range seenJobs {
					if seenJob.MaxJobs > seenMax {
						seenMax = seenJob.MaxJobs
					}
				}
				return seenMax <= max, nil
			},
		}
	}

	// 2 jobs at a time
	// we should end up with 2 groups of 2 in terms of timing
	// and the highest number of jobs at one time should be 2
	suite.Run("2 jobs at a time", func() {
		runTest(
			TotalResourceTestCase{
				jobs: getResourcesArray([][]string{
					{"1", "500Mb", ""},
					{"1", "500Mb", ""},
					{"1", "500Mb", ""},
					{"1", "500Mb", ""},
				}),
				totalLimits: getResources("2", "1Gb", "1Gb"),
				wait:        waitUntilSeenAllJobs(4),
				checkers: []TotalResourceTestCaseCheck{
					// there should only have ever been 2 jobs at one time
					checkMaxJobs(2),
				},
			},
		)
	})

	// test disk space
	// we have a 1Gb disk
	// and 2 jobs each with 900Mb disk space requirements
	// we should only see 1 job at a time
	suite.Run("test disk space", func() {
		runTest(
			TotalResourceTestCase{
				jobs: getResourcesArray([][]string{
					{"100m", "100Mb", "600Mb"},
					{"100m", "100Mb", "600Mb"},
				}),
				totalLimits: getResources("2", "1Gb", "1Gb"),
				wait:        waitUntilSeenAllJobs(2),
				checkers: []TotalResourceTestCaseCheck{
					// there should only have ever been 1 job at one time
					checkMaxJobs(1),
				},
			},
		)
	})

}

// test that with 10 GPU nodes - that 10 jobs end up being allocated 1 per node
// this is a check of the bidding & capacity manager system
func (suite *ComputeNodeResourceLimitsSuite) TestParallelGPU() {
	nodeCount := 2
	jobsPerNode := 2
	seenJobs := 0
	var jobIds []string

	ctx := context.Background()

	// the job needs to hang for a period of time so the other job will
	// run on another node
	jobHandler := func(_ context.Context, _ noop_executor.ExecutionContext) (*models.RunCommandResult, error) {
		time.Sleep(time.Millisecond * 1000)
		seenJobs++
		return &models.RunCommandResult{}, nil
	}

	cfg, err := config.NewTestConfigWithOverrides(types.Bacalhau{
		Compute: types.Compute{
			AllocatedCapacity: types.ResourceScalerFromModelsResourceConfig(models.ResourcesConfig{
				CPU:    "1",
				Memory: "1Gb",
				Disk:   "1Gb",
				GPU:    "1",
			}),
			Heartbeat: types.Heartbeat{
				Interval: types.Duration(50 * time.Millisecond),
			},
		},
	})
	suite.Require().NoError(err)

	stack := teststack.Setup(ctx, suite.T(),
		devstack.WithNumberOfHybridNodes(1),
		devstack.WithNumberOfComputeOnlyNodes(1),
		devstack.WithBacalhauConfigOverride(cfg),
		devstack.WithSystemConfig(node.SystemConfig{
			OverSubscriptionFactor: 2,
		}),
		teststack.WithNoopExecutor(
			noop_executor.ExecutorConfig{
				ExternalHooks: noop_executor.ExecutorConfigExternalHooks{
					JobHandler: jobHandler,
				},
			}, cfg.Engines),
	)

	// for the requester node to pick up the nodeInfo messages
	nodeutils.WaitForNodeDiscovery(suite.T(), stack.Nodes[0].RequesterNode, nodeCount)

	job := &models.Job{
		Type:  models.JobTypeBatch,
		Count: 1,
		Tasks: []*models.Task{
			{
				Name: suite.T().Name(),
				Engine: &models.SpecConfig{
					Type:   models.EngineNoop,
					Params: make(map[string]interface{}),
				},
				Publisher: &models.SpecConfig{
					Type:   models.PublisherNoop,
					Params: make(map[string]interface{}),
				},
				ResourcesConfig: &models.ResourcesConfig{
					CPU:    "",
					Memory: "",
					Disk:   "",
					GPU:    "1",
				},
			},
		},
	}
	job.Normalize()

	resolver := scenario.NewStateResolverFromStore(stack.Nodes[0].RequesterNode.JobStore)

	sanitizedJobName := strings.Replace(suite.T().Name(), "/", "-", -1)

	for i := 0; i < nodeCount; i++ {
		for j := 0; j < jobsPerNode; j++ {
			job.Name = fmt.Sprintf("%s-%d-%d", sanitizedJobName, i, j)
			submittedJob, err := stack.Nodes[0].RequesterNode.Endpoint.SubmitJob(ctx, &orchestrator.SubmitJobRequest{Job: job})
			require.NoError(suite.T(), err)
			jobIds = append(jobIds, submittedJob.JobID)

			// sleep a bit here to simulate jobs being submitted over time
			// and to give time for compute nodes to accept and run the jobs
			// this needs to be less than the time the job lasts
			// so we are running jobs in parallel
			time.Sleep(200 * time.Millisecond)
		}
	}

	for _, jobId := range jobIds {
		err = resolver.Wait(ctx, jobId, scenario.WaitForSuccessfulCompletion())
		require.NoError(suite.T(), err)
	}

	require.Equal(suite.T(), nodeCount*jobsPerNode, seenJobs)

	allocationMap := map[string]int{}

	for _, jobId := range jobIds {
		jobState, err := resolver.JobState(ctx, jobId)
		require.NoError(suite.T(), err)
		completedExecutionStates := scenario.GetCompletedExecutionStates(jobState)
		require.Equal(suite.T(), 1, len(completedExecutionStates))
		require.Equal(suite.T(), models.ExecutionStateCompleted, completedExecutionStates[0].ComputeState.StateType)
		allocationMap[completedExecutionStates[0].NodeID]++
	}

	// test that each node has 2 job allocated to it
	node1Count, ok := allocationMap[stack.Nodes[0].ID]
	require.True(suite.T(), ok)
	require.Equal(suite.T(), jobsPerNode, node1Count)

	node2Count, ok := allocationMap[stack.Nodes[1].ID]
	require.True(suite.T(), ok)
	require.Equal(suite.T(), jobsPerNode, node2Count)
}
