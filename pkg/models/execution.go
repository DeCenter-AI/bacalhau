//go:generate stringer -type=ExecutionStateType --trimprefix=ExecutionState --output execution_state_string.go
//go:generate stringer -type=ExecutionDesiredStateType --trimprefix=ExecutionDesiredState --output execution_desired_state_string.go
package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/bacalhau-project/bacalhau/pkg/lib/validate"
)

// ExecutionStateType The state of an execution. An execution represents a single attempt to execute a job on a node.
// A compute node can have multiple executions for the same job due to retries, but there can only be a single active execution
// per node at any given time.
type ExecutionStateType int

// TODO: change states to reflect non-bidding scheduling
const (
	ExecutionStateUndefined ExecutionStateType = iota
	// ExecutionStateNew The execution has been created, but not pushed to a compute node yet.
	ExecutionStateNew
	// ExecutionStateAskForBid A node has been selected to execute a job, and is being asked to bid on the job.
	ExecutionStateAskForBid
	// ExecutionStateAskForBidAccepted compute node has rejected the ask for bid.
	ExecutionStateAskForBidAccepted
	// ExecutionStateAskForBidRejected compute node has rejected the ask for bid.
	ExecutionStateAskForBidRejected
	// ExecutionStateBidAccepted requester has accepted the bid, and the execution is expected to be running on the compute node.
	ExecutionStateBidAccepted // aka running
	// ExecutionStateRunning The execution is running on the compute node.
	ExecutionStateRunning
	// ExecutionStatePublishing The execution has completed, and the result is being published.
	ExecutionStatePublishing
	// ExecutionStateBidRejected requester has rejected the bid.
	ExecutionStateBidRejected
	// ExecutionStateCompleted The execution has been completed, and the result has been published.
	ExecutionStateCompleted
	// ExecutionStateFailed The execution has failed.
	ExecutionStateFailed
	// ExecutionStateCancelled The execution has been canceled by the user
	ExecutionStateCancelled
)

func ExecutionStateTypes() []ExecutionStateType {
	var res []ExecutionStateType
	for typ := ExecutionStateUndefined; typ <= ExecutionStateCancelled; typ++ {
		res = append(res, typ)
	}
	return res
}

// IsUndefined returns true if the execution state is undefined
func (s ExecutionStateType) IsUndefined() bool {
	return s == ExecutionStateUndefined
}

func (s ExecutionStateType) IsTerminal() bool {
	return s == ExecutionStateBidRejected ||
		s == ExecutionStateCompleted ||
		s == ExecutionStateFailed ||
		s == ExecutionStateCancelled ||
		s == ExecutionStateAskForBidRejected
}

// IsExecuting returns true if the execution is running in the backend
func (s ExecutionStateType) IsExecuting() bool {
	return s == ExecutionStateBidAccepted || s == ExecutionStateRunning || s == ExecutionStatePublishing
}

type ExecutionDesiredStateType int

const (
	ExecutionDesiredStatePending ExecutionDesiredStateType = iota
	ExecutionDesiredStateRunning
	ExecutionDesiredStateStopped
)

// Execution is used to allocate the placement of a task group to a node.
type Execution struct {
	// ID of the execution (UUID)
	ID string `json:"ID"`

	// Namespace is the namespace the execution is created in
	Namespace string `json:"Namespace"`

	// ID of the evaluation that generated this execution
	EvalID string `json:"EvalID"`

	// Name is a logical name of the execution.
	Name string `json:"Name"`

	// NodeID is the node this is being placed on
	NodeID string `json:"NodeID"`

	// Job is the parent job of the task being allocated.
	// This is copied at execution time to avoid issues if the job
	// definition is updated.
	JobID string `json:"JobID"`
	// TODO: evaluate using a copy of the job instead of a pointer
	Job *Job `json:"Job,omitempty"`

	// JobVersion version of the job that this execution was created in
	JobVersion uint64 `json:"JobVersion"`

	// AllocatedResources is the total resources allocated for the execution tasks.
	AllocatedResources *AllocatedResources `json:"AllocatedResources"`

	// DesiredState of the execution on the compute node
	DesiredState State[ExecutionDesiredStateType] `json:"DesiredState"`

	// ComputeState observed state of the execution on the compute node
	ComputeState State[ExecutionStateType] `json:"ComputeState"`

	// the published results for this execution
	PublishedResult *SpecConfig `json:"PublishedResult"`

	// RunOutput is the output of the run command
	// TODO: evaluate removing this from execution spec in favour of calling `bacalhau job logs`
	RunOutput *RunCommandResult `json:"RunOutput"`

	// PreviousExecution is the execution that this execution is replacing
	PreviousExecution string `json:"PreviousExecution"`

	// NextExecution is the execution that this execution is being replaced by
	NextExecution string `json:"NextExecution"`

	// FollowupEvalID captures a follow up evaluation created to handle a failed execution
	// that can be rescheduled in the future
	FollowupEvalID string `json:"FollowupEvalID"`

	// PartitionIndex is the index of this execution in the job's total partitions (0-based)
	// Only relevant when Job.Count > 1
	PartitionIndex int `json:"PartitionIndex,omitempty"`

	// Revision is increment each time the execution is updated.
	Revision uint64 `json:"Revision"`

	// CreateTime is the time the execution has finished scheduling and been
	// verified by the plan applier.
	CreateTime int64 `json:"CreateTime"`
	// ModifyTime is the time the execution was last updated.
	ModifyTime int64 `json:"ModifyTime"`
}

func (e *Execution) String() string {
	return e.ID
}

func (e *Execution) JobNamespacedID() NamespacedID {
	return NewNamespacedID(e.JobID, e.Namespace)
}

// GetCreateTime returns the creation time
func (e *Execution) GetCreateTime() time.Time {
	return time.Unix(0, e.CreateTime).UTC()
}

// GetModifyTime returns the modify time
func (e *Execution) GetModifyTime() time.Time {
	return time.Unix(0, e.ModifyTime).UTC()
}

// IsExpired returns true if the execution is still running beyond the expiration time
// We return true if the execution is in the bid accepted state (i.e. running)
// and the modify time is older than the expiration time
func (e *Execution) IsExpired(expirationTime time.Time) bool {
	return e.ComputeState.StateType.IsExecuting() && e.GetModifyTime().Before(expirationTime)
}

// Normalize Allocation to ensure fields are initialized to the expectations
// of this version of Bacalhau. Should be called when restoring persisted
// Executions or receiving Executions from Bacalhau clients potentially on an
// older version of Bacalhau.
func (e *Execution) Normalize() {
	if e == nil {
		return
	}
	if e.AllocatedResources == nil {
		e.AllocatedResources = &AllocatedResources{
			Tasks: make(map[string]*Resources),
		}
	}
	if e.PublishedResult == nil {
		e.PublishedResult = &SpecConfig{}
	}
	if e.RunOutput == nil {
		e.RunOutput = &RunCommandResult{}
	}
	e.Job.Normalize()
}

// Copy provides a copy of the allocation and deep copies the job
func (e *Execution) Copy() *Execution {
	if e == nil {
		return nil
	}
	na := new(Execution)
	*na = *e

	na.Job = na.Job.Copy()
	na.AllocatedResources = na.AllocatedResources.Copy()
	na.PublishedResult = na.PublishedResult.Copy()
	na.RunOutput = na.RunOutput.Copy()
	return na
}

// Validate is used to check a job for reasonable configuration
func (e *Execution) Validate() error {
	err := errors.Join(
		validate.NotBlank(e.ID, "missing execution ID"),
		validate.NoSpaces(e.ID, "execution ID contains a space"),
		validate.NoNullChars(e.ID, "execution ID contains a null character"),
		validate.NotBlank(e.Namespace, "execution must be in a namespace"),
		validate.NotBlank(e.JobID, "missing execution job ID"),
	)
	if e.Job != nil {
		err = errors.Join(err, e.Job.Validate())
		if (e.Job.Type == JobTypeBatch || e.Job.Type == JobTypeService) && e.Job.Count > 1 {
			if e.PartitionIndex < 0 || e.PartitionIndex >= e.Job.Count {
				err = errors.Join(err, fmt.Errorf("partition index must be between 0 and %d", e.Job.Count-1))
			}
		}
	}
	return err
}

// IsTerminalState returns true if the execution desired of observed state is terminal
func (e *Execution) IsTerminalState() bool {
	return e.IsTerminalDesiredState() || e.IsTerminalComputeState()
}

// IsTerminalDesiredState returns true if the execution desired state is terminal
func (e *Execution) IsTerminalDesiredState() bool {
	return e.DesiredState.StateType == ExecutionDesiredStateStopped
}

// IsTerminalComputeState returns true if the execution observed state is terminal
func (e *Execution) IsTerminalComputeState() bool {
	switch e.ComputeState.StateType {
	case ExecutionStateCompleted, ExecutionStateFailed, ExecutionStateCancelled, ExecutionStateAskForBidRejected, ExecutionStateBidRejected:
		return true
	default:
		return false
	}
}

// IsDiscarded returns true if the execution has failed, been cancelled or rejected.
func (e *Execution) IsDiscarded() bool {
	switch e.ComputeState.StateType {
	case ExecutionStateAskForBidRejected, ExecutionStateBidRejected, ExecutionStateCancelled, ExecutionStateFailed:
		return true
	default:
		return false
	}
}

// AllocateResources allocates resources to a task
func (e *Execution) AllocateResources(taskID string, resources Resources) {
	if e.AllocatedResources == nil {
		e.AllocatedResources = &AllocatedResources{
			Tasks: make(map[string]*Resources),
		}
	}
	e.AllocatedResources.Tasks[taskID] = resources.Copy()
}

// AllocatePorts updates the execution's network configuration with the allocated port mappings.
func (e *Execution) AllocatePorts(portMap PortMap) {
	if e.Job == nil || e.Job.Task().Network == nil {
		return
	}

	// Replace the existing port mappings with the allocated ones
	if portMap == nil {
		e.Job.Task().Network.Ports = make(PortMap, 0)
	} else {
		e.Job.Task().Network.Ports = portMap.Copy()
	}
}

// OrchestrationProtocol is the protocol used to orchestrate the execution
func (e *Execution) OrchestrationProtocol() Protocol {
	return e.Job.OrchestrationProtocol()
}

func (e *Execution) TotalAllocatedResources() *Resources {
	return e.AllocatedResources.Total()
}

type RunCommandResult struct {
	// stdout of the run. Yaml provided for `describe` output
	STDOUT string `json:"Stdout"`

	// bool describing if stdout was truncated
	StdoutTruncated bool `json:"StdoutTruncated"`

	// stderr of the run.
	STDERR string `json:"stderr"`

	// bool describing if stderr was truncated
	StderrTruncated bool `json:"StderrTruncated"`

	// exit code of the run.
	ExitCode int `json:"ExitCode"`

	// Runner error
	ErrorMsg string `json:"ErrorMsg"`
}

func NewRunCommandResult() *RunCommandResult {
	return &RunCommandResult{
		STDOUT:          "",    // stdout of the run.
		StdoutTruncated: false, // bool describing if stdout was truncated
		STDERR:          "",    // stderr of the run.
		StderrTruncated: false, // bool describing if stderr was truncated
		ExitCode:        -1,    // exit code of the run.
	}
}

func (r *RunCommandResult) Copy() *RunCommandResult {
	if r == nil {
		return nil
	}

	newRCR := new(RunCommandResult)
	*newRCR = *r

	// Since all fields are simple types (string, bool, int),
	// a shallow copy is sufficient.
	return newRCR
}
