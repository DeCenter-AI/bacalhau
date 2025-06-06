// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package orchestrator is a generated GoMock package.
package orchestrator

import (
	context "context"
	reflect "reflect"
	time "time"

	models "github.com/bacalhau-project/bacalhau/pkg/models"
	nodes "github.com/bacalhau-project/bacalhau/pkg/orchestrator/nodes"
	gomock "go.uber.org/mock/gomock"
)

// MockEvaluationBroker is a mock of EvaluationBroker interface.
type MockEvaluationBroker struct {
	ctrl     *gomock.Controller
	recorder *MockEvaluationBrokerMockRecorder
}

// MockEvaluationBrokerMockRecorder is the mock recorder for MockEvaluationBroker.
type MockEvaluationBrokerMockRecorder struct {
	mock *MockEvaluationBroker
}

// NewMockEvaluationBroker creates a new mock instance.
func NewMockEvaluationBroker(ctrl *gomock.Controller) *MockEvaluationBroker {
	mock := &MockEvaluationBroker{ctrl: ctrl}
	mock.recorder = &MockEvaluationBrokerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvaluationBroker) EXPECT() *MockEvaluationBrokerMockRecorder {
	return m.recorder
}

// Ack mocks base method.
func (m *MockEvaluationBroker) Ack(evalID, receiptHandle string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ack", evalID, receiptHandle)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ack indicates an expected call of Ack.
func (mr *MockEvaluationBrokerMockRecorder) Ack(evalID, receiptHandle interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ack", reflect.TypeOf((*MockEvaluationBroker)(nil).Ack), evalID, receiptHandle)
}

// Dequeue mocks base method.
func (m *MockEvaluationBroker) Dequeue(types []string, timeout time.Duration) (*models.Evaluation, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dequeue", types, timeout)
	ret0, _ := ret[0].(*models.Evaluation)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Dequeue indicates an expected call of Dequeue.
func (mr *MockEvaluationBrokerMockRecorder) Dequeue(types, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dequeue", reflect.TypeOf((*MockEvaluationBroker)(nil).Dequeue), types, timeout)
}

// Enqueue mocks base method.
func (m *MockEvaluationBroker) Enqueue(evaluation *models.Evaluation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", evaluation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockEvaluationBrokerMockRecorder) Enqueue(evaluation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockEvaluationBroker)(nil).Enqueue), evaluation)
}

// EnqueueAll mocks base method.
func (m *MockEvaluationBroker) EnqueueAll(evaluation map[*models.Evaluation]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnqueueAll", evaluation)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnqueueAll indicates an expected call of EnqueueAll.
func (mr *MockEvaluationBrokerMockRecorder) EnqueueAll(evaluation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueAll", reflect.TypeOf((*MockEvaluationBroker)(nil).EnqueueAll), evaluation)
}

// Inflight mocks base method.
func (m *MockEvaluationBroker) Inflight(evaluationID string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inflight", evaluationID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Inflight indicates an expected call of Inflight.
func (mr *MockEvaluationBrokerMockRecorder) Inflight(evaluationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inflight", reflect.TypeOf((*MockEvaluationBroker)(nil).Inflight), evaluationID)
}

// InflightExtend mocks base method.
func (m *MockEvaluationBroker) InflightExtend(evaluationID, receiptHandle string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InflightExtend", evaluationID, receiptHandle)
	ret0, _ := ret[0].(error)
	return ret0
}

// InflightExtend indicates an expected call of InflightExtend.
func (mr *MockEvaluationBrokerMockRecorder) InflightExtend(evaluationID, receiptHandle interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InflightExtend", reflect.TypeOf((*MockEvaluationBroker)(nil).InflightExtend), evaluationID, receiptHandle)
}

// Nack mocks base method.
func (m *MockEvaluationBroker) Nack(evalID, receiptHandle string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nack", evalID, receiptHandle)
	ret0, _ := ret[0].(error)
	return ret0
}

// Nack indicates an expected call of Nack.
func (mr *MockEvaluationBrokerMockRecorder) Nack(evalID, receiptHandle interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nack", reflect.TypeOf((*MockEvaluationBroker)(nil).Nack), evalID, receiptHandle)
}

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockScheduler) Process(ctx context.Context, eval *models.Evaluation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, eval)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockSchedulerMockRecorder) Process(ctx, eval interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockScheduler)(nil).Process), ctx, eval)
}

// MockSchedulerProvider is a mock of SchedulerProvider interface.
type MockSchedulerProvider struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerProviderMockRecorder
}

// MockSchedulerProviderMockRecorder is the mock recorder for MockSchedulerProvider.
type MockSchedulerProviderMockRecorder struct {
	mock *MockSchedulerProvider
}

// NewMockSchedulerProvider creates a new mock instance.
func NewMockSchedulerProvider(ctrl *gomock.Controller) *MockSchedulerProvider {
	mock := &MockSchedulerProvider{ctrl: ctrl}
	mock.recorder = &MockSchedulerProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchedulerProvider) EXPECT() *MockSchedulerProviderMockRecorder {
	return m.recorder
}

// EnabledSchedulers mocks base method.
func (m *MockSchedulerProvider) EnabledSchedulers() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnabledSchedulers")
	ret0, _ := ret[0].([]string)
	return ret0
}

// EnabledSchedulers indicates an expected call of EnabledSchedulers.
func (mr *MockSchedulerProviderMockRecorder) EnabledSchedulers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnabledSchedulers", reflect.TypeOf((*MockSchedulerProvider)(nil).EnabledSchedulers))
}

// Scheduler mocks base method.
func (m *MockSchedulerProvider) Scheduler(jobType string) (Scheduler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scheduler", jobType)
	ret0, _ := ret[0].(Scheduler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scheduler indicates an expected call of Scheduler.
func (mr *MockSchedulerProviderMockRecorder) Scheduler(jobType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scheduler", reflect.TypeOf((*MockSchedulerProvider)(nil).Scheduler), jobType)
}

// MockPlanner is a mock of Planner interface.
type MockPlanner struct {
	ctrl     *gomock.Controller
	recorder *MockPlannerMockRecorder
}

// MockPlannerMockRecorder is the mock recorder for MockPlanner.
type MockPlannerMockRecorder struct {
	mock *MockPlanner
}

// NewMockPlanner creates a new mock instance.
func NewMockPlanner(ctrl *gomock.Controller) *MockPlanner {
	mock := &MockPlanner{ctrl: ctrl}
	mock.recorder = &MockPlannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlanner) EXPECT() *MockPlannerMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockPlanner) Process(ctx context.Context, plan *models.Plan) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, plan)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockPlannerMockRecorder) Process(ctx, plan interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockPlanner)(nil).Process), ctx, plan)
}

// MockNodeDiscoverer is a mock of NodeDiscoverer interface.
type MockNodeDiscoverer struct {
	ctrl     *gomock.Controller
	recorder *MockNodeDiscovererMockRecorder
}

// MockNodeDiscovererMockRecorder is the mock recorder for MockNodeDiscoverer.
type MockNodeDiscovererMockRecorder struct {
	mock *MockNodeDiscoverer
}

// NewMockNodeDiscoverer creates a new mock instance.
func NewMockNodeDiscoverer(ctrl *gomock.Controller) *MockNodeDiscoverer {
	mock := &MockNodeDiscoverer{ctrl: ctrl}
	mock.recorder = &MockNodeDiscovererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeDiscoverer) EXPECT() *MockNodeDiscovererMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockNodeDiscoverer) List(ctx context.Context, filter ...nodes.NodeStateFilter) ([]models.NodeState, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range filter {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].([]models.NodeState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockNodeDiscovererMockRecorder) List(ctx interface{}, filter ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, filter...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNodeDiscoverer)(nil).List), varargs...)
}

// MockNodeRanker is a mock of NodeRanker interface.
type MockNodeRanker struct {
	ctrl     *gomock.Controller
	recorder *MockNodeRankerMockRecorder
}

// MockNodeRankerMockRecorder is the mock recorder for MockNodeRanker.
type MockNodeRankerMockRecorder struct {
	mock *MockNodeRanker
}

// NewMockNodeRanker creates a new mock instance.
func NewMockNodeRanker(ctrl *gomock.Controller) *MockNodeRanker {
	mock := &MockNodeRanker{ctrl: ctrl}
	mock.recorder = &MockNodeRankerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeRanker) EXPECT() *MockNodeRankerMockRecorder {
	return m.recorder
}

// RankNodes mocks base method.
func (m *MockNodeRanker) RankNodes(ctx context.Context, job models.Job, nodes []models.NodeInfo) ([]NodeRank, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RankNodes", ctx, job, nodes)
	ret0, _ := ret[0].([]NodeRank)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RankNodes indicates an expected call of RankNodes.
func (mr *MockNodeRankerMockRecorder) RankNodes(ctx, job, nodes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RankNodes", reflect.TypeOf((*MockNodeRanker)(nil).RankNodes), ctx, job, nodes)
}

// MockNodeSelector is a mock of NodeSelector interface.
type MockNodeSelector struct {
	ctrl     *gomock.Controller
	recorder *MockNodeSelectorMockRecorder
}

// MockNodeSelectorMockRecorder is the mock recorder for MockNodeSelector.
type MockNodeSelectorMockRecorder struct {
	mock *MockNodeSelector
}

// NewMockNodeSelector creates a new mock instance.
func NewMockNodeSelector(ctrl *gomock.Controller) *MockNodeSelector {
	mock := &MockNodeSelector{ctrl: ctrl}
	mock.recorder = &MockNodeSelectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeSelector) EXPECT() *MockNodeSelectorMockRecorder {
	return m.recorder
}

// AllNodes mocks base method.
func (m *MockNodeSelector) AllNodes(ctx context.Context) ([]models.NodeInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllNodes", ctx)
	ret0, _ := ret[0].([]models.NodeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllNodes indicates an expected call of AllNodes.
func (mr *MockNodeSelectorMockRecorder) AllNodes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllNodes", reflect.TypeOf((*MockNodeSelector)(nil).AllNodes), ctx)
}

// MatchingNodes mocks base method.
func (m *MockNodeSelector) MatchingNodes(ctx context.Context, job *models.Job) ([]NodeRank, []NodeRank, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MatchingNodes", ctx, job)
	ret0, _ := ret[0].([]NodeRank)
	ret1, _ := ret[1].([]NodeRank)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// MatchingNodes indicates an expected call of MatchingNodes.
func (mr *MockNodeSelectorMockRecorder) MatchingNodes(ctx, job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MatchingNodes", reflect.TypeOf((*MockNodeSelector)(nil).MatchingNodes), ctx, job)
}

// MockRetryStrategy is a mock of RetryStrategy interface.
type MockRetryStrategy struct {
	ctrl     *gomock.Controller
	recorder *MockRetryStrategyMockRecorder
}

// MockRetryStrategyMockRecorder is the mock recorder for MockRetryStrategy.
type MockRetryStrategyMockRecorder struct {
	mock *MockRetryStrategy
}

// NewMockRetryStrategy creates a new mock instance.
func NewMockRetryStrategy(ctrl *gomock.Controller) *MockRetryStrategy {
	mock := &MockRetryStrategy{ctrl: ctrl}
	mock.recorder = &MockRetryStrategyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRetryStrategy) EXPECT() *MockRetryStrategyMockRecorder {
	return m.recorder
}

// ShouldRetry mocks base method.
func (m *MockRetryStrategy) ShouldRetry(ctx context.Context, request RetryRequest) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldRetry", ctx, request)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ShouldRetry indicates an expected call of ShouldRetry.
func (mr *MockRetryStrategyMockRecorder) ShouldRetry(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldRetry", reflect.TypeOf((*MockRetryStrategy)(nil).ShouldRetry), ctx, request)
}
