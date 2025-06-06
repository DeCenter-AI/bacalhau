package apimodels

import (
	"strconv"
	"strings"

	"github.com/bacalhau-project/bacalhau/pkg/bacerrors"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/bacalhau-project/bacalhau/pkg/models"
)

type PutJobRequest struct {
	BasePutRequest
	Job   *models.Job `json:"Job"`
	Force bool        `json:"Force"`
}

// Validate is used to validate fields in the PutJobRequest.
func (r *PutJobRequest) Validate() error {
	return r.Job.ValidateSubmission()
}

type PutJobResponse struct {
	BasePutResponse
	JobID        string   `json:"JobID"`
	EvaluationID string   `json:"EvaluationID"`
	Warnings     []string `json:"Warnings"`
}

type DiffJobRequest struct {
	BasePutRequest
	Job *models.Job `json:"Job"`
}

// Validate is used to validate fields in the DiffJobRequest.
func (r *DiffJobRequest) Validate() error {
	// Ensure that the job name is not empty. This validation applies to Job Diffs.
	if strings.TrimSpace(r.Job.Name) == "" {
		return bacerrors.Newf("invalid job spec: name is missing").
			WithCode(bacerrors.ValidationError).
			WithHint("Job name must be provided in the job spec for diff operations")
	}

	return r.Job.ValidateSubmission()
}

type DiffJobResponse struct {
	BasePutResponse
	Diff     string   `json:"Diff"`
	Warnings []string `json:"Warnings"`
}

type GetJobRequest struct {
	BaseGetRequest
	JobIDOrName string
	JobVersion  uint64 `query:"job_version" validate:"omitempty"`
	Include     string `query:"include" validate:"omitempty,oneof=history executions"`
	Limit       uint32 `query:"limit"`
}

// ToHTTPRequest is used to convert the request to an HTTP request
func (o *GetJobRequest) ToHTTPRequest() *HTTPRequest {
	r := o.BaseGetRequest.ToHTTPRequest()

	if o.JobVersion != 0 {
		r.Params.Set("job_version", strconv.FormatUint(o.JobVersion, 10))
	}
	if o.Include != "" {
		r.Params.Set("include", o.Include)
	}
	if o.Limit > 0 {
		r.Params.Set("limit", strconv.FormatUint(uint64(o.Limit), 10))
	}
	return r
}

type GetJobResponse struct {
	BaseGetResponse
	Job        *models.Job                `json:"Job"`
	History    *ListJobHistoryResponse    `json:"History,omitempty"`
	Executions *ListJobExecutionsResponse `json:"Executions,omitempty"`
}

// Normalize is used to canonicalize fields in the GetJobResponse.
func (r *GetJobResponse) Normalize() {
	r.BaseGetResponse.Normalize()
	if r.Job != nil {
		r.Job.Normalize()
	}
	if r.History != nil {
		r.History.Normalize()
	}
	if r.Executions != nil {
		r.Executions.Normalize()
	}
}

type ListJobsRequest struct {
	BaseListRequest
	Labels []labels.Requirement `query:"-"` // don't auto bind as it requires special handling
}

// ToHTTPRequest is used to convert the request to an HTTP request
func (o *ListJobsRequest) ToHTTPRequest() *HTTPRequest {
	r := o.BaseListRequest.ToHTTPRequest()

	for _, v := range o.Labels {
		r.Params.Add("labels", v.String())
	}
	return r
}

type ListJobsResponse struct {
	BaseListResponse
	Items []*models.Job `json:"Items"`
}

// Normalize is used to canonicalize fields in the ListJobsResponse.
func (r *ListJobsResponse) Normalize() {
	r.BaseListResponse.Normalize()
	for _, job := range r.Items {
		job.Normalize()
	}
}

type ListJobHistoryRequest struct {
	BaseVersionedListRequest
	JobIDOrName string `query:"-"`
	Since       int64  `query:"since" validate:"min=0"`
	EventType   string `query:"event_type" validate:"omitempty,oneof=all job execution"`
	ExecutionID string `query:"execution_id" validate:"omitempty"`
}

// ToHTTPRequest is used to convert the request to an HTTP request
func (o *ListJobHistoryRequest) ToHTTPRequest() *HTTPRequest {
	r := o.BaseVersionedListRequest.ToHTTPRequest()

	if o.Since != 0 {
		r.Params.Set("since", strconv.FormatInt(o.Since, 10))
	}
	if o.EventType != "" {
		r.Params.Set("event_type", o.EventType)
	}
	if o.ExecutionID != "" {
		r.Params.Set("execution_id", o.ExecutionID)
	}
	return r
}

type ListJobHistoryResponse struct {
	BaseListResponse
	Items []*models.JobHistory `json:"Items"`
}

type ListJobExecutionsRequest struct {
	BaseVersionedListRequest
	JobIDOrName string `query:"-"`
}

type ListJobExecutionsResponse struct {
	BaseListResponse
	Items []*models.Execution `json:"Items"`
}

type ListJobVersionsRequest struct {
	BaseListRequest
	JobIDOrName string `query:"-"`
}

type ListJobVersionsResponse struct {
	BaseListResponse
	Items []*models.Job `json:"Items"`
}

type ListJobResultsRequest struct {
	BaseListRequest
	JobID string `query:"-"`
}

type ListJobResultsResponse struct {
	BaseListResponse
	Items []*models.SpecConfig `json:"Items"`
}

type StopJobRequest struct {
	BasePutRequest
	JobID  string `json:"-"`
	Reason string `json:"reason"`
}

type StopJobResponse struct {
	BasePutResponse
	EvaluationID string `json:"EvaluationID"`
}

type RerunJobRequest struct {
	BasePutRequest
	JobIDOrName string `json:"-"`
	JobVersion  uint64 `json:"JobVersion"`
	Reason      string `json:"Reason"`
}

type RerunJobResponse struct {
	BasePutResponse
	EvaluationID string   `json:"EvaluationID"`
	JobID        string   `json:"JobID"`
	JobVersion   uint64   `json:"JobVersion"`
	Warnings     []string `json:"Warnings"`
}

type GetLogsRequest struct {
	BaseGetRequest
	JobID          string `query:"-"`
	JobVersion     uint64 `query:"job_version" validate:"omitempty"`
	AllJobVersions bool   `query:"all_job_versions" validate:"omitempty"`
	ExecutionID    string `query:"execution_id" validate:"omitempty"`
	Tail           bool   `query:"tail"`
	Follow         bool   `query:"follow"`
}

// ToHTTPRequest is used to convert the request to an HTTP request
func (o *GetLogsRequest) ToHTTPRequest() *HTTPRequest {
	r := o.BaseGetRequest.ToHTTPRequest()

	if o.ExecutionID != "" {
		r.Params.Set("execution_id", o.ExecutionID)
	}
	if o.Tail {
		r.Params.Set("tail", "true")
	}
	if o.Follow {
		r.Params.Set("follow", "true")
	}
	if o.JobVersion != 0 {
		r.Params.Set("job_version", strconv.FormatUint(o.JobVersion, 10))
	}
	if o.AllJobVersions {
		r.Params.Set("all_job_versions", strconv.FormatBool(o.AllJobVersions))
	}
	return r
}
