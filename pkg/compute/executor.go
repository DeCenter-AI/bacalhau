package compute

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bacalhau-project/bacalhau/pkg/bacerrors"
	"github.com/bacalhau-project/bacalhau/pkg/executor"
	"github.com/bacalhau-project/bacalhau/pkg/models"
	"github.com/bacalhau-project/bacalhau/pkg/telemetry"

	"github.com/bacalhau-project/bacalhau/pkg/compute/store"
	"github.com/bacalhau-project/bacalhau/pkg/publisher"
	"github.com/bacalhau-project/bacalhau/pkg/storage"
	"github.com/bacalhau-project/bacalhau/pkg/system"
)

const (
	StorageDirectoryPerms     = 0o755
	executionRootCleanupDelay = 1 * time.Hour
)

type BaseExecutorParams struct {
	ID                     string
	Store                  store.ExecutionStore
	Storages               storage.StorageProvider
	StorageDirectory       string
	Executors              executor.ExecProvider
	ResultsPath            ResultsPath
	Publishers             publisher.PublisherProvider
	FailureInjectionConfig models.FailureInjectionConfig
	EnvResolver            EnvVarResolver
	PortAllocator          PortAllocator

	// TODO: this is a temporary solution and should be replaced with a more generic
	//  solution to populate jobs with default resources and network config.
	//  Most likely at a higher level in the stack and before queueing the execution.
	DefaultNetworkType models.Network
}

// BaseExecutor is the base implementation for backend service.
// All operations are executed asynchronously, and a callback is used to notify the caller of the result.
type BaseExecutor struct {
	ID                 string
	store              store.ExecutionStore
	Storages           storage.StorageProvider
	storageDirectory   string
	executors          executor.ExecProvider
	publishers         publisher.PublisherProvider
	resultsPath        ResultsPath
	failureInjection   models.FailureInjectionConfig
	envResolver        EnvVarResolver
	portAllocator      PortAllocator
	defaultNetworkType models.Network
}

func NewBaseExecutor(params BaseExecutorParams) *BaseExecutor {
	return &BaseExecutor{
		ID:                 params.ID,
		store:              params.Store,
		Storages:           params.Storages,
		storageDirectory:   params.StorageDirectory,
		executors:          params.Executors,
		publishers:         params.Publishers,
		failureInjection:   params.FailureInjectionConfig,
		resultsPath:        params.ResultsPath,
		envResolver:        params.EnvResolver,
		portAllocator:      params.PortAllocator,
		defaultNetworkType: params.DefaultNetworkType,
	}
}

func (e *BaseExecutor) prepareInputVolumes(
	ctx context.Context,
	execution *models.Execution,
) ([]storage.PreparedStorage, func(context.Context) error, error) {
	inputVolumes, err := storage.ParallelPrepareStorage(
		ctx, e.Storages, e.storageDirectory, execution, execution.Job.Task().InputSources...)
	if err != nil {
		return nil, nil, err
	}
	return inputVolumes, func(ctx context.Context) error {
		return storage.ParallelCleanStorage(ctx, e.Storages, inputVolumes)
	}, nil
}

// InputCleanupFn is a function type that defines the contract for cleaning up
// resources associated with input volume data after the job execution has either completed
// or failed to start. The function is expected to take a context.Context as an argument,
// which can be used for timeout and cancellation signals. It returns an error if
// the cleanup operation fails.
//
// For example, an InputCleanupFn might be responsible for deallocating storage used
// for input volumes, or deleting temporary input files that were created as part of the
// job's execution. The nature of it operation depends on the storage provided by `storageProvider` and
// input sources of the jobs associated tasks.
type InputCleanupFn = func(context.Context) error

func (e *BaseExecutor) PrepareRunArguments(
	ctx context.Context,
	execution *models.Execution,
	executionDir string,
) (*executor.RunCommandRequest, InputCleanupFn, error) {
	var cleanupFuncs []func(context.Context) error

	inputVolumes, inputCleanup, err := e.prepareInputVolumes(ctx, execution)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, inputCleanup)

	// Allocate ports
	portMappings, err := e.portAllocator.AllocatePorts(execution)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, func(ctx context.Context) error {
		e.portAllocator.ReleasePorts(execution)
		return nil
	})

	// Update execution with allocated ports
	execution.AllocatePorts(portMappings)

	env, err := GetExecutionEnvVars(execution, e.envResolver)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	networkConfig := execution.Job.Task().Network
	if networkConfig.Type == models.NetworkDefault {
		networkConfig.Type = e.defaultNetworkType
	}

	return &executor.RunCommandRequest{
			JobID:        execution.Job.ID,
			ExecutionID:  execution.ID,
			Resources:    execution.TotalAllocatedResources(),
			Network:      networkConfig,
			Outputs:      execution.Job.Task().ResultPaths,
			Inputs:       inputVolumes,
			ExecutionDir: executionDir,
			EngineParams: execution.Job.Task().Engine,
			Env:          env,
			OutputLimits: executor.OutputLimits{
				MaxStdoutFileLength:   system.MaxStdoutFileLength,
				MaxStdoutReturnLength: system.MaxStdoutReturnLength,
				MaxStderrFileLength:   system.MaxStderrFileLength,
				MaxStderrReturnLength: system.MaxStderrReturnLength,
			},
		}, func(ctx context.Context) error {
			var cleanupErr error
			for _, cleanupFunc := range cleanupFuncs {
				if err := cleanupFunc(ctx); err != nil {
					cleanupErr = errors.Join(cleanupErr, err)
				}
			}
			return cleanupErr
		}, nil
}

type StartResult struct {
	cleanup InputCleanupFn
	Err     error
}

func (r *StartResult) Cleanup(ctx context.Context) error {
	if r.cleanup != nil {
		return r.cleanup(ctx)
	}
	return nil
}

func (e *BaseExecutor) Start(ctx context.Context, execution *models.Execution) *StartResult {
	result := new(StartResult)
	jobExecutor, err := e.executors.Get(ctx, execution.Job.Task().Engine.Type)
	if err != nil {
		result.Err = fmt.Errorf("getting executor %s: %w", execution.Job.Task().Engine, err)
		return result
	}

	executionDir, err := e.resultsPath.PrepareExecutionOutputDir(execution.ID)
	if err != nil {
		result.Err = fmt.Errorf("preparing results path: %w", err)
		return result
	}

	args, cleanup, err := e.PrepareRunArguments(ctx, execution, executionDir)
	result.cleanup = cleanup
	if err != nil {
		result.Err = fmt.Errorf("preparing arguments: %w", err)
		return result
	}

	if err = e.store.UpdateExecutionState(ctx, store.UpdateExecutionRequest{
		ExecutionID: execution.ID,
		Condition: store.UpdateExecutionCondition{
			ExpectedStates: []models.ExecutionStateType{
				models.ExecutionStateBidAccepted,
				models.ExecutionStateRunning, // allow retries during node restarts
			},
		},
		NewValues: models.Execution{
			ComputeState: models.NewExecutionState(models.ExecutionStateRunning),
		},
	}); err != nil {
		result.Err = fmt.Errorf("updating execution state from expected: %s to: %s",
			models.ExecutionStateBidAccepted, models.ExecutionStateRunning)
		return result
	}

	log.Ctx(ctx).Debug().Msg("starting execution")

	if e.failureInjection.IsBadActor {
		result.Err = fmt.Errorf("i am a bad node. i failed execution %s", execution.ID)
		return result
	}

	if err := jobExecutor.Start(ctx, args); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to start execution")
		result.Err = err
	}

	return result
}

func (e *BaseExecutor) Wait(ctx context.Context, execution *models.Execution) (*models.RunCommandResult, error) {
	jobExecutor, err := e.executors.Get(ctx, execution.Job.Task().Engine.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get executor %s: %w", execution.Job.Task().Engine, err)
	}

	waitC, errC := jobExecutor.Wait(ctx, execution.ID)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-waitC:
		return res, nil
	case err := <-errC:
		log.Ctx(ctx).Error().Err(err).Msg("failed to wait on execution")
		return nil, err
	}
}

// Run the execution after it has been accepted, and propose a result to the requester to be verified.
//
//nolint:funlen
func (e *BaseExecutor) Run(ctx context.Context, execution *models.Execution) (err error) {
	ctx = log.Ctx(ctx).With().
		Str("job", execution.Job.ID).
		Str("execution", execution.ID).
		Logger().WithContext(ctx)

	stopwatch := telemetry.Timer(ctx, jobDurationMilliseconds, execution.Job.MetricAttributes()...)
	topic := EventTopicExecutionRunning
	defer func() {
		if err != nil {
			if !bacerrors.IsErrorWithCode(err, executor.ExecutionAlreadyCancelled) {
				e.handleFailure(ctx, execution, err, topic)
			}
		}
		dur := stopwatch()
		log.Ctx(ctx).Debug().
			Dur("duration", dur).
			Str("jobID", execution.JobID).
			Str("executionID", execution.ID).
			Msg("run complete")
	}()

	res := e.Start(ctx, execution)

	defer func(executionOutputDir string) {
		// For now execution results are removed after a fixed delay.
		go func() {
			log.Debug().Str("path", executionOutputDir).Msg("scheduled execution results dir removal")
			ticker := time.NewTicker(executionRootCleanupDelay)
			defer ticker.Stop()
			for range ticker.C {
				err = os.RemoveAll(executionOutputDir)
				if err != nil {
					log.Error().Err(err).Str("path", executionOutputDir).Msg("failed to remove execution results dir")
					return
				}
				log.Debug().Str("path", executionOutputDir).Msg("removed execution results dir")
			}
		}()

		// The rest can be cleaned up immediately after completion
		if err := res.Cleanup(ctx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to clean up start arguments")
		}
	}(e.resultsPath.ExecutionOutputDir(execution.ID))

	if err := res.Err; err != nil {
		if bacerrors.IsErrorWithCode(err, executor.ExecutionAlreadyStarted) {
			// by not returning this error to the caller when the execution has already been started/is already running
			// we allow duplicate calls to `Run` to be idempotent and fall through to the below `Wait` call.
			log.Ctx(ctx).Warn().Err(err).Str("execution", execution.ID).
				Msg("execution is already running processing to wait on execution")
		} else {
			// We don't consider the job failed if the execution is already running or has already completed.
			// TODO(forrest): [correctness] do we really want to record a job failed metric if (one of) its execution(s)
			// failed to start? Perhaps it would be better to have metrics for execution failures here and job failures
			// higher up the call stack?
			jobsFailed.Add(ctx, 1)
			return err
		}
	}

	result, err := e.Wait(ctx, execution)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			// TODO(forrest) [correctness]:
			// The ExecutorBuffer is using a context with a timeout to signal an execution has timed out and should end.
			//
			// We don't handle context.Canceled here as it means the node is shutting down. Still we should do a
			// better job at gracefully shutting down the execution and either reporting that to the requester
			// or retrying the execution during startup.
			//
			// Moving forward we must avoid canceling executions via the context.Context. When pluggable executors
			// become the default since canceling the context will simply result in the RPC connection closing (I think)
			// The general solution here is to stop using contexts for canceling jobs and to instead make explicit calls
			// the an executors `Cancel` method.
			return NewErrExecTimeout(execution.Job.Task().Timeouts.GetExecutionTimeout())
		}
		return err
	}
	if result.ErrorMsg != "" {
		return fmt.Errorf("%s", result.ErrorMsg)
	}
	jobsCompleted.Add(ctx, 1)

	expectedState := models.ExecutionStateRunning
	var publishedResult *models.SpecConfig

	// publish if the job has a publisher defined
	if !execution.Job.Task().Publisher.IsEmpty() {
		topic = EventTopicExecutionPublishing
		if err = e.store.UpdateExecutionState(ctx, store.UpdateExecutionRequest{
			ExecutionID: execution.ID,
			Condition: store.UpdateExecutionCondition{
				ExpectedStates: []models.ExecutionStateType{expectedState},
			},
			NewValues: models.Execution{
				ComputeState: models.NewExecutionState(models.ExecutionStatePublishing),
				RunOutput:    result,
			},
		}); err != nil {
			return err
		}

		expectedState = models.ExecutionStatePublishing

		resultsDir := ExecutionResultsDir(e.resultsPath.ExecutionOutputDir(execution.ID))
		publishedResult, err = e.publish(ctx, execution, resultsDir)
		if err != nil {
			return err
		}

		defer func() {
			// cleanup execution results
			log.Ctx(ctx).Debug().
				Str("execution", execution.ID).
				Str("path", resultsDir).
				Msg("cleaning up execution results")
			err = os.RemoveAll(resultsDir)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf("failed to remove results directory at %s", resultsDir)
			}
		}()
	}

	// mark the execution as completed
	if err = e.store.UpdateExecutionState(ctx, store.UpdateExecutionRequest{
		ExecutionID: execution.ID,
		Condition: store.UpdateExecutionCondition{
			ExpectedStates: []models.ExecutionStateType{expectedState},
		},
		NewValues: models.Execution{
			ComputeState:    models.NewExecutionState(models.ExecutionStateCompleted),
			PublishedResult: publishedResult,
			RunOutput:       result,
		},
		Events: []*models.Event{ExecCompletedEvent()},
	}); err != nil {
		return err
	}

	return err
}

// Publish the result of an execution after it has been verified.
func (e *BaseExecutor) publish(ctx context.Context, execution *models.Execution,
	resultsDir string,
) (*models.SpecConfig, error) {
	log.Ctx(ctx).Debug().Str("executionID", execution.ID).Str("resultsDir", resultsDir).Msg("Publishing execution results")

	jobPublisher, err := e.publishers.Get(ctx, execution.Job.Task().Publisher.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get publisher %s: %w", execution.Job.Task().Publisher.Type, err)
	}
	publishedResult, err := jobPublisher.PublishResult(ctx, execution, resultsDir)
	if err != nil {
		return nil, bacerrors.Wrap(err, "failed to publish result")
	}
	log.Ctx(ctx).Debug().
		Str("execution", execution.ID).
		Msg("Execution published")

	return &publishedResult, nil
}

// Cancel the execution.
func (e *BaseExecutor) Cancel(ctx context.Context, execution *models.Execution) error {
	log.Ctx(ctx).Debug().Str("Execution", execution.ID).Msg("Canceling execution")
	exe, err := e.executors.Get(ctx, execution.Job.Task().Engine.Type)
	if err != nil {
		return err
	}
	return exe.Cancel(ctx, execution.ID)
}

func (e *BaseExecutor) handleFailure(ctx context.Context, execution *models.Execution, err error, topic models.EventTopic) {
	log.Ctx(ctx).Warn().Err(err).Msgf("%s failed", topic)

	updateError := e.store.UpdateExecutionState(ctx, store.UpdateExecutionRequest{
		ExecutionID: execution.ID,
		NewValues: models.Execution{
			ComputeState: models.NewExecutionState(models.ExecutionStateFailed).WithMessage(err.Error()),
		},
		Events: []*models.Event{models.NewEvent(topic).WithError(err)},
	})

	if updateError != nil {
		var alreadyTerminalError store.ErrExecutionAlreadyTerminal
		if !errors.As(updateError, &alreadyTerminalError) {
			log.Ctx(ctx).Error().Err(updateError).Msgf("Failed to update execution (%s) state to failed: %s", execution.ID, updateError)
		}
	}
}

// compile-time interface check
var _ Executor = (*BaseExecutor)(nil)
