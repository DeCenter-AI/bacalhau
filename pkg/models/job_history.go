//go:generate stringer -type=JobHistoryType --trimprefix=JobHistoryType --output job_history_string.go
package models

import (
	"strings"
	"time"
)

type JobHistoryType int

const (
	JobHistoryTypeUndefined JobHistoryType = iota
	JobHistoryTypeJobLevel
	JobHistoryTypeExecutionLevel
)

func (s JobHistoryType) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *JobHistoryType) UnmarshalText(text []byte) (err error) {
	name := string(text)
	for typ := JobHistoryTypeUndefined; typ <= JobHistoryTypeExecutionLevel; typ++ {
		if strings.EqualFold(typ.String(), strings.TrimSpace(name)) {
			*s = typ
			return
		}
	}
	return
}

// StateChange represents a change in state of one of the state types.
type StateChange[StateType any] struct {
	Previous StateType `json:"Previous,omitempty"`
	New      StateType `json:"New,omitempty"`
}

// JobHistory represents a single event in the history of a job. An event can be
// at the job level, or execution (node) level.
//
// {Job,Event}State fields will only be present if the Type field is of
// the matching type.
type JobHistory struct {
	SeqNum      uint64         `json:"SeqNum"`
	Type        JobHistoryType `json:"Type"`
	JobID       string         `json:"JobID"`
	JobVersion  uint64         `json:"JobVersion"`
	ExecutionID string         `json:"ExecutionID,omitempty"`
	Event       Event          `json:"Event,omitempty"`
	Time        time.Time      `json:"Time"`

	// TODO: remove with v1.5
	// Deprecated: Left for backward compatibility with v1.4.x clients
	JobState *StateChange[JobStateType] `json:"JobState,omitempty"`
	// Deprecated: Left for backward compatibility with v1.4.x clients
	ExecutionState *StateChange[ExecutionStateType] `json:"ExecutionState,omitempty"`
}

// Occurred returns when the action that triggered an update to job history
// actually occurred.
//
// The Time field represents the moment that the JobHistory item was recorded,
// i.e. it is almost always set to time.Now() when creating the object. This is
// different to the Event.Timestamp which represents when the source of the
// history update actually occurred.
func (jh JobHistory) Occurred() time.Time {
	if !jh.Event.Timestamp.Equal(time.Time{}) {
		return jh.Event.Timestamp
	}
	return jh.Time
}

// IsJobLevel returns true if the JobHistory is at the job level.
func (jh JobHistory) IsJobLevel() bool {
	return jh.Type == JobHistoryTypeJobLevel
}

// IsExecutionLevel returns true if the JobHistory is at the execution level.
func (jh JobHistory) IsExecutionLevel() bool {
	return jh.Type == JobHistoryTypeExecutionLevel
}
