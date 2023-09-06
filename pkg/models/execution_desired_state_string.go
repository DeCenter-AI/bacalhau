// Code generated by "stringer -type=ExecutionDesiredStateType --trimprefix=ExecutionDesiredState --output execution_desired_state_string.go"; DO NOT EDIT.

package models

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ExecutionDesiredStatePending-0]
	_ = x[ExecutionDesiredStateRunning-1]
	_ = x[ExecutionDesiredStateStopped-2]
}

const _ExecutionDesiredStateType_name = "PendingRunningStopped"

var _ExecutionDesiredStateType_index = [...]uint8{0, 7, 14, 21}

func (i ExecutionDesiredStateType) String() string {
	if i < 0 || i >= ExecutionDesiredStateType(len(_ExecutionDesiredStateType_index)-1) {
		return "ExecutionDesiredStateType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ExecutionDesiredStateType_name[_ExecutionDesiredStateType_index[i]:_ExecutionDesiredStateType_index[i+1]]
}