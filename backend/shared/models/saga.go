package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SagaStatus string
type SagaStepStatus string
type CompensationStrategy string
type ExecutionMode string

const (
	// Saga Status
	SagaStatusCreated     SagaStatus = "created"
	SagaStatusRunning     SagaStatus = "running"
	SagaStatusCompleted   SagaStatus = "completed"
	SagaStatusFailed      SagaStatus = "failed"
	SagaStatusCompensating SagaStatus = "compensating"
	SagaStatusCompensated SagaStatus = "compensated"
	SagaStatusCancelled   SagaStatus = "cancelled"
	SagaStatusTimeout     SagaStatus = "timeout"

	// Saga Step Status
	SagaStepStatusPending      SagaStepStatus = "pending"
	SagaStepStatusRunning      SagaStepStatus = "running"
	SagaStepStatusCompleted    SagaStepStatus = "completed"
	SagaStepStatusFailed       SagaStepStatus = "failed"
	SagaStepStatusCompensating SagaStepStatus = "compensating"
	SagaStepStatusCompensated  SagaStepStatus = "compensated"
	SagaStepStatusSkipped      SagaStepStatus = "skipped"
	SagaStepStatusTimeout      SagaStepStatus = "timeout"

	// Compensation Strategies
	CompensationStrategyNone       CompensationStrategy = "none"
	CompensationStrategyAutomatic  CompensationStrategy = "automatic"
	CompensationStrategyManual     CompensationStrategy = "manual"
	CompensationStrategyBestEffort CompensationStrategy = "best_effort"

	// Execution Modes
	ExecutionModeSequential ExecutionMode = "sequential"
	ExecutionModeParallel   ExecutionMode = "parallel"
	ExecutionModeMixed      ExecutionMode = "mixed"
)

// Common Saga Types for Southeast Asian E-commerce Platform
const (
	SagaTypeOrderProcessing       = "order_processing"
	SagaTypePaymentProcessing     = "payment_processing"
	SagaTypeUserRegistration     = "user_registration"
	SagaTypeKYCVerification      = "kyc_verification"
	SagaTypeWalletTopup          = "wallet_topup"
	SagaTypeWalletWithdrawal     = "wallet_withdrawal"
	SagaTypeShopOnboarding       = "shop_onboarding"
	SagaTypeProductListing       = "product_listing"
	SagaTypeOrderFulfillment     = "order_fulfillment"
	SagaTypeRefundProcessing     = "refund_processing"
	SagaTypeAccountClosure       = "account_closure"
	SagaTypeDataMigration        = "data_migration"
	SagaTypeSystemMaintenance    = "system_maintenance"
)

type SagaContext struct {
	UserID          uuid.UUID      `json:"user_id,omitempty"`
	OrderID         uuid.UUID      `json:"order_id,omitempty"`
	PaymentID       uuid.UUID      `json:"payment_id,omitempty"`
	ShopID          uuid.UUID      `json:"shop_id,omitempty"`
	ProductID       uuid.UUID      `json:"product_id,omitempty"`
	WalletID        uuid.UUID      `json:"wallet_id,omitempty"`
	TransactionID   uuid.UUID      `json:"transaction_id,omitempty"`
	Country         string         `json:"country,omitempty"`
	Currency        string         `json:"currency,omitempty"`
	Amount          float64        `json:"amount,omitempty"`
	Locale          string         `json:"locale,omitempty"`
	Region          string         `json:"region,omitempty"`
	Environment     string         `json:"environment,omitempty"`
	TraceID         string         `json:"trace_id,omitempty"`
	CorrelationID   string         `json:"correlation_id,omitempty"`
	RequestID       string         `json:"request_id,omitempty"`
	SessionID       string         `json:"session_id,omitempty"`
	Data            map[string]any `json:"data,omitempty"`
}

type StepConfiguration struct {
	ServiceName          string                `json:"service_name"`
	EndpointURL          string                `json:"endpoint_url"`
	Method               string                `json:"method"`
	Headers              map[string]string     `json:"headers"`
	Timeout              time.Duration         `json:"timeout"`
	RetryPolicy          RetryPolicy           `json:"retry_policy"`
	CompensationStrategy CompensationStrategy  `json:"compensation_strategy"`
	CompensationEndpoint string                `json:"compensation_endpoint,omitempty"`
	Dependencies         []string              `json:"dependencies,omitempty"` // Step names this depends on
	Conditions           map[string]any        `json:"conditions,omitempty"`   // Conditions for execution
	Metadata             map[string]any        `json:"metadata,omitempty"`
}

type RetryPolicy struct {
	MaxRetries    int           `json:"max_retries"`
	Interval      time.Duration `json:"interval"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxInterval   time.Duration `json:"max_interval"`
	RetryOn       []string      `json:"retry_on"` // Error codes/types to retry on
}

type ExecutionResult struct {
	Success        bool              `json:"success"`
	StatusCode     int               `json:"status_code,omitempty"`
	ResponseBody   string            `json:"response_body,omitempty"`
	ErrorMessage   string            `json:"error_message,omitempty"`
	ErrorCode      string            `json:"error_code,omitempty"`
	ExecutedAt     time.Time         `json:"executed_at"`
	Duration       time.Duration     `json:"duration"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
	OutputData     map[string]any    `json:"output_data,omitempty"`
}

type CompensationResult struct {
	Success        bool              `json:"success"`
	StatusCode     int               `json:"status_code,omitempty"`
	ResponseBody   string            `json:"response_body,omitempty"`
	ErrorMessage   string            `json:"error_message,omitempty"`
	ErrorCode      string            `json:"error_code,omitempty"`
	CompensatedAt  time.Time         `json:"compensated_at"`
	Duration       time.Duration     `json:"duration"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
}

type Saga struct {
	ID             uuid.UUID      `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name           string         `json:"name" gorm:"type:varchar(255);not null"`
	Type           string         `json:"type" gorm:"type:varchar(100);not null;index"`
	Status         SagaStatus     `json:"status" gorm:"type:varchar(20);default:'created';index"`
	ExecutionMode  ExecutionMode  `json:"execution_mode" gorm:"type:varchar(20);default:'sequential'"`

	// Context and Data
	Context        SagaContext    `json:"context" gorm:"type:json"`
	InputData      json.RawMessage `json:"input_data" gorm:"type:json"`
	OutputData     json.RawMessage `json:"output_data" gorm:"type:json"`

	// Configuration
	TimeoutDuration time.Duration  `json:"timeout_duration" gorm:"default:1800000000000"` // 30 minutes in nanoseconds
	MaxRetries      int            `json:"max_retries" gorm:"default:3"`
	RetryInterval   time.Duration  `json:"retry_interval" gorm:"default:30000000000"`    // 30 seconds in nanoseconds

	// Execution Information
	CurrentStepIndex int           `json:"current_step_index" gorm:"default:0"`
	StartedAt        *time.Time    `json:"started_at,omitempty" gorm:"index"`
	CompletedAt      *time.Time    `json:"completed_at,omitempty" gorm:"index"`
	FailedAt         *time.Time    `json:"failed_at,omitempty" gorm:"index"`
	ExpiresAt        *time.Time    `json:"expires_at,omitempty" gorm:"index"`

	// Error Handling
	LastError        string        `json:"last_error,omitempty" gorm:"type:text"`
	RetryCount       int           `json:"retry_count" gorm:"default:0"`
	LastRetryAt      *time.Time    `json:"last_retry_at,omitempty"`

	// Audit Fields
	CreatedBy        uuid.UUID     `json:"created_by" gorm:"type:varchar(36);index"`
	CreatedAt        time.Time     `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt        time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type SagaStep struct {
	ID             uuid.UUID          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	SagaID         uuid.UUID          `json:"saga_id" gorm:"type:varchar(36);not null;index"`
	Name           string             `json:"name" gorm:"type:varchar(255);not null"`
	StepIndex      int                `json:"step_index" gorm:"not null"`
	Status         SagaStepStatus     `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	Configuration  StepConfiguration  `json:"configuration" gorm:"type:json"`

	// Execution Results
	ExecutionResults    []ExecutionResult    `json:"execution_results" gorm:"type:json"`
	CompensationResults []CompensationResult `json:"compensation_results" gorm:"type:json"`

	// Timing
	StartedAt       *time.Time        `json:"started_at,omitempty"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	FailedAt        *time.Time        `json:"failed_at,omitempty"`
	CompensatedAt   *time.Time        `json:"compensated_at,omitempty"`

	// Error Handling
	RetryCount      int               `json:"retry_count" gorm:"default:0"`
	LastError       string            `json:"last_error,omitempty" gorm:"type:text"`
	LastRetryAt     *time.Time        `json:"last_retry_at,omitempty"`

	// Input/Output Data
	InputData       json.RawMessage   `json:"input_data" gorm:"type:json"`
	OutputData      json.RawMessage   `json:"output_data" gorm:"type:json"`

	// Audit Fields
	CreatedAt       time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

type SagaDefinition struct {
	ID             uuid.UUID           `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name           string              `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Type           string              `json:"type" gorm:"type:varchar(100);not null;index"`
	Version        string              `json:"version" gorm:"type:varchar(20);not null"`
	Description    string              `json:"description" gorm:"type:text"`
	ExecutionMode  ExecutionMode       `json:"execution_mode" gorm:"type:varchar(20);default:'sequential'"`
	StepDefinitions []StepDefinition   `json:"step_definitions" gorm:"type:json"`
	DefaultConfig   SagaConfig         `json:"default_config" gorm:"type:json"`
	IsActive       bool                `json:"is_active" gorm:"default:true"`
	Tags           []string            `json:"tags" gorm:"type:json"`

	// Audit Fields
	CreatedBy      uuid.UUID           `json:"created_by" gorm:"type:varchar(36);index"`
	CreatedAt      time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
}

type StepDefinition struct {
	Name           string            `json:"name"`
	ServiceName    string            `json:"service_name"`
	EndpointURL    string            `json:"endpoint_url"`
	Method         string            `json:"method"`
	CompensationURL string           `json:"compensation_url,omitempty"`
	CompensationStrategy CompensationStrategy `json:"compensation_strategy"`
	Dependencies   []string          `json:"dependencies,omitempty"`
	Timeout        time.Duration     `json:"timeout"`
	RetryPolicy    RetryPolicy       `json:"retry_policy"`
	Conditions     map[string]any    `json:"conditions,omitempty"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
}

type SagaConfig struct {
	TimeoutDuration time.Duration `json:"timeout_duration"`
	MaxRetries      int           `json:"max_retries"`
	RetryInterval   time.Duration `json:"retry_interval"`
	AutoCompensate  bool          `json:"auto_compensate"`
	ParallelSteps   [][]string    `json:"parallel_steps,omitempty"` // Groups of steps that can run in parallel
}

func (s *Saga) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Context.CorrelationID == "" {
		s.Context.CorrelationID = uuid.New().String()
	}
	if s.Context.TraceID == "" {
		s.Context.TraceID = uuid.New().String()
	}
	if s.TimeoutDuration == 0 {
		s.TimeoutDuration = 30 * time.Minute
	}
	if s.ExpiresAt == nil {
		expiresAt := time.Now().Add(s.TimeoutDuration)
		s.ExpiresAt = &expiresAt
	}
	return nil
}

func (s *Saga) IsExpired() bool {
	return s.ExpiresAt != nil && time.Now().After(*s.ExpiresAt)
}

func (s *Saga) CanRetry() bool {
	return s.RetryCount < s.MaxRetries && !s.IsExpired()
}

func (s *Saga) Start() {
	s.Status = SagaStatusRunning
	now := time.Now()
	s.StartedAt = &now
}

func (s *Saga) Complete() {
	s.Status = SagaStatusCompleted
	now := time.Now()
	s.CompletedAt = &now
}

func (s *Saga) Fail(errorMessage string) {
	s.Status = SagaStatusFailed
	s.LastError = errorMessage
	now := time.Now()
	s.FailedAt = &now
}

func (s *Saga) StartCompensation() {
	s.Status = SagaStatusCompensating
}

func (s *Saga) CompleteCompensation() {
	s.Status = SagaStatusCompensated
}

func (s *Saga) Cancel() {
	s.Status = SagaStatusCancelled
}

func (s *Saga) Timeout() {
	s.Status = SagaStatusTimeout
	now := time.Now()
	s.FailedAt = &now
}

func (s *Saga) AddRetry(errorMessage string) {
	s.RetryCount++
	s.LastError = errorMessage
	now := time.Now()
	s.LastRetryAt = &now
}

func (s *Saga) UnmarshalInputData(v interface{}) error {
	return json.Unmarshal(s.InputData, v)
}

func (s *Saga) MarshalInputData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.InputData = data
	return nil
}

func (s *Saga) UnmarshalOutputData(v interface{}) error {
	return json.Unmarshal(s.OutputData, v)
}

func (s *Saga) MarshalOutputData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.OutputData = data
	return nil
}

func (s *Saga) GetProgressPercentage(totalSteps int) float64 {
	if totalSteps == 0 {
		return 0
	}
	if s.Status == SagaStatusCompleted {
		return 100.0
	}
	return float64(s.CurrentStepIndex) / float64(totalSteps) * 100.0
}

func (ss *SagaStep) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return nil
}

func (ss *SagaStep) Start() {
	ss.Status = SagaStepStatusRunning
	now := time.Now()
	ss.StartedAt = &now
}

func (ss *SagaStep) Complete(result ExecutionResult) {
	ss.Status = SagaStepStatusCompleted
	ss.ExecutionResults = append(ss.ExecutionResults, result)
	now := time.Now()
	ss.CompletedAt = &now
}

func (ss *SagaStep) Fail(result ExecutionResult) {
	ss.Status = SagaStepStatusFailed
	ss.ExecutionResults = append(ss.ExecutionResults, result)
	ss.LastError = result.ErrorMessage
	now := time.Now()
	ss.FailedAt = &now
}

func (ss *SagaStep) StartCompensation() {
	ss.Status = SagaStepStatusCompensating
}

func (ss *SagaStep) CompleteCompensation(result CompensationResult) {
	ss.Status = SagaStepStatusCompensated
	ss.CompensationResults = append(ss.CompensationResults, result)
	now := time.Now()
	ss.CompensatedAt = &now
}

func (ss *SagaStep) Skip() {
	ss.Status = SagaStepStatusSkipped
	now := time.Now()
	ss.CompletedAt = &now
}

func (ss *SagaStep) Timeout() {
	ss.Status = SagaStepStatusTimeout
	now := time.Now()
	ss.FailedAt = &now
}

func (ss *SagaStep) CanRetry() bool {
	return ss.RetryCount < ss.Configuration.RetryPolicy.MaxRetries
}

func (ss *SagaStep) AddRetry(errorMessage string) {
	ss.RetryCount++
	ss.LastError = errorMessage
	now := time.Now()
	ss.LastRetryAt = &now
}

func (ss *SagaStep) GetLatestExecutionResult() *ExecutionResult {
	if len(ss.ExecutionResults) == 0 {
		return nil
	}
	return &ss.ExecutionResults[len(ss.ExecutionResults)-1]
}

func (ss *SagaStep) GetLatestCompensationResult() *CompensationResult {
	if len(ss.CompensationResults) == 0 {
		return nil
	}
	return &ss.CompensationResults[len(ss.CompensationResults)-1]
}

func (ss *SagaStep) UnmarshalInputData(v interface{}) error {
	return json.Unmarshal(ss.InputData, v)
}

func (ss *SagaStep) MarshalInputData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ss.InputData = data
	return nil
}

func (ss *SagaStep) UnmarshalOutputData(v interface{}) error {
	return json.Unmarshal(ss.OutputData, v)
}

func (ss *SagaStep) MarshalOutputData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ss.OutputData = data
	return nil
}

func (ss *SagaStep) ShouldExecute(sagaContext SagaContext) bool {
	if len(ss.Configuration.Conditions) == 0 {
		return true
	}

	// Simple condition evaluation - in production, use a proper expression engine
	for key, expectedValue := range ss.Configuration.Conditions {
		var actualValue interface{}

		switch key {
		case "country":
			actualValue = sagaContext.Country
		case "currency":
			actualValue = sagaContext.Currency
		case "amount_gte":
			if expectedValue, ok := expectedValue.(float64); ok {
				return sagaContext.Amount >= expectedValue
			}
		case "amount_lte":
			if expectedValue, ok := expectedValue.(float64); ok {
				return sagaContext.Amount <= expectedValue
			}
		default:
			if data, exists := sagaContext.Data[key]; exists {
				actualValue = data
			}
		}

		if actualValue != expectedValue {
			return false
		}
	}

	return true
}

func (sd *SagaDefinition) BeforeCreate(tx *gorm.DB) error {
	if sd.ID == uuid.Nil {
		sd.ID = uuid.New()
	}
	return nil
}

func (sd *SagaDefinition) CreateSaga(context SagaContext, inputData interface{}) (*Saga, []*SagaStep, error) {
	saga := &Saga{
		ID:               uuid.New(),
		Name:             sd.Name,
		Type:             sd.Type,
		Status:           SagaStatusCreated,
		ExecutionMode:    sd.ExecutionMode,
		Context:          context,
		TimeoutDuration:  sd.DefaultConfig.TimeoutDuration,
		MaxRetries:       sd.DefaultConfig.MaxRetries,
		RetryInterval:    sd.DefaultConfig.RetryInterval,
		CurrentStepIndex: 0,
	}

	if inputData != nil {
		if err := saga.MarshalInputData(inputData); err != nil {
			return nil, nil, fmt.Errorf("failed to marshal input data: %w", err)
		}
	}

	steps := make([]*SagaStep, len(sd.StepDefinitions))
	for i, stepDef := range sd.StepDefinitions {
		steps[i] = &SagaStep{
			ID:        uuid.New(),
			SagaID:    saga.ID,
			Name:      stepDef.Name,
			StepIndex: i,
			Status:    SagaStepStatusPending,
			Configuration: StepConfiguration{
				ServiceName:          stepDef.ServiceName,
				EndpointURL:          stepDef.EndpointURL,
				Method:               stepDef.Method,
				Timeout:              stepDef.Timeout,
				RetryPolicy:          stepDef.RetryPolicy,
				CompensationStrategy: stepDef.CompensationStrategy,
				CompensationEndpoint: stepDef.CompensationURL,
				Dependencies:         stepDef.Dependencies,
				Conditions:           stepDef.Conditions,
				Metadata:             stepDef.Metadata,
			},
			ExecutionResults:    make([]ExecutionResult, 0),
			CompensationResults: make([]CompensationResult, 0),
		}
	}

	return saga, steps, nil
}

func (sd *SagaDefinition) ValidateStepDependencies() error {
	stepNames := make(map[string]bool)
	for _, stepDef := range sd.StepDefinitions {
		stepNames[stepDef.Name] = true
	}

	for _, stepDef := range sd.StepDefinitions {
		for _, dependency := range stepDef.Dependencies {
			if !stepNames[dependency] {
				return fmt.Errorf("step '%s' has invalid dependency '%s'", stepDef.Name, dependency)
			}
		}
	}

	return nil
}

// Response structures for API
type SagaResponse struct {
	ID               uuid.UUID      `json:"id"`
	Name             string         `json:"name"`
	Type             string         `json:"type"`
	Status           SagaStatus     `json:"status"`
	ExecutionMode    ExecutionMode  `json:"execution_mode"`
	Context          SagaContext    `json:"context"`
	CurrentStepIndex int            `json:"current_step_index"`
	Progress         float64        `json:"progress"`
	StartedAt        *time.Time     `json:"started_at,omitempty"`
	CompletedAt      *time.Time     `json:"completed_at,omitempty"`
	FailedAt         *time.Time     `json:"failed_at,omitempty"`
	ExpiresAt        *time.Time     `json:"expires_at,omitempty"`
	LastError        string         `json:"last_error,omitempty"`
	RetryCount       int            `json:"retry_count"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type SagaStepResponse struct {
	ID               uuid.UUID         `json:"id"`
	SagaID           uuid.UUID         `json:"saga_id"`
	Name             string            `json:"name"`
	StepIndex        int               `json:"step_index"`
	Status           SagaStepStatus    `json:"status"`
	StartedAt        *time.Time        `json:"started_at,omitempty"`
	CompletedAt      *time.Time        `json:"completed_at,omitempty"`
	FailedAt         *time.Time        `json:"failed_at,omitempty"`
	CompensatedAt    *time.Time        `json:"compensated_at,omitempty"`
	RetryCount       int               `json:"retry_count"`
	LastError        string            `json:"last_error,omitempty"`
	ExecutionResults []ExecutionResult `json:"execution_results"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// Request structures for API
type CreateSagaRequest struct {
	Type       string         `json:"type" binding:"required"`
	Name       string         `json:"name"`
	Context    SagaContext    `json:"context" binding:"required"`
	InputData  interface{}    `json:"input_data"`
	Config     *SagaConfig    `json:"config"`
}

type CreateSagaDefinitionRequest struct {
	Name            string           `json:"name" binding:"required"`
	Type            string           `json:"type" binding:"required"`
	Version         string           `json:"version" binding:"required"`
	Description     string           `json:"description"`
	ExecutionMode   ExecutionMode    `json:"execution_mode"`
	StepDefinitions []StepDefinition `json:"step_definitions" binding:"required"`
	DefaultConfig   SagaConfig       `json:"default_config" binding:"required"`
	Tags            []string         `json:"tags"`
}

// Manager for business logic
type SagaManager struct{}

func NewSagaManager() *SagaManager {
	return &SagaManager{}
}

func (sm *SagaManager) ValidateSaga(saga *Saga) error {
	if saga.Name == "" {
		return fmt.Errorf("saga name is required")
	}

	if saga.Type == "" {
		return fmt.Errorf("saga type is required")
	}

	if saga.Context.CorrelationID == "" {
		return fmt.Errorf("correlation ID is required")
	}

	return nil
}

func (sm *SagaManager) GetSupportedSagaTypes() []string {
	return []string{
		SagaTypeOrderProcessing,
		SagaTypePaymentProcessing,
		SagaTypeUserRegistration,
		SagaTypeKYCVerification,
		SagaTypeWalletTopup,
		SagaTypeWalletWithdrawal,
		SagaTypeShopOnboarding,
		SagaTypeProductListing,
		SagaTypeOrderFulfillment,
		SagaTypeRefundProcessing,
		SagaTypeAccountClosure,
		SagaTypeDataMigration,
		SagaTypeSystemMaintenance,
	}
}

func (sm *SagaManager) CalculateExecutionOrder(steps []*SagaStep, executionMode ExecutionMode) ([][]int, error) {
	if executionMode == ExecutionModeSequential {
		// Sequential execution: each step waits for the previous one
		result := make([][]int, len(steps))
		for i := range steps {
			result[i] = []int{i}
		}
		return result, nil
	}

	if executionMode == ExecutionModeParallel {
		// Parallel execution: all steps run simultaneously (if no dependencies)
		indices := make([]int, len(steps))
		for i := range steps {
			indices[i] = i
		}
		return [][]int{indices}, nil
	}

	// Mixed mode: respect dependencies and parallelize where possible
	return sm.calculateMixedExecutionOrder(steps)
}

func (sm *SagaManager) calculateMixedExecutionOrder(steps []*SagaStep) ([][]int, error) {
	stepNameToIndex := make(map[string]int)
	for i, step := range steps {
		stepNameToIndex[step.Name] = i
	}

	// Build dependency graph
	dependsOn := make(map[int][]int)
	dependents := make(map[int][]int)

	for i, step := range steps {
		dependsOn[i] = make([]int, 0)
		for _, depName := range step.Configuration.Dependencies {
			if depIndex, exists := stepNameToIndex[depName]; exists {
				dependsOn[i] = append(dependsOn[i], depIndex)
				if dependents[depIndex] == nil {
					dependents[depIndex] = make([]int, 0)
				}
				dependents[depIndex] = append(dependents[depIndex], i)
			}
		}
	}

	// Topological sort with level grouping
	result := make([][]int, 0)
	completed := make(map[int]bool)
	inProgress := make(map[int]bool)

	for len(completed) < len(steps) {
		currentLevel := make([]int, 0)

		// Find steps that can be executed (no pending dependencies)
		for i := 0; i < len(steps); i++ {
			if completed[i] || inProgress[i] {
				continue
			}

			canExecute := true
			for _, depIndex := range dependsOn[i] {
				if !completed[depIndex] {
					canExecute = false
					break
				}
			}

			if canExecute {
				currentLevel = append(currentLevel, i)
				inProgress[i] = true
			}
		}

		if len(currentLevel) == 0 {
			return nil, fmt.Errorf("circular dependency detected or invalid dependency configuration")
		}

		result = append(result, currentLevel)

		// Mark current level as completed
		for _, stepIndex := range currentLevel {
			completed[stepIndex] = true
			delete(inProgress, stepIndex)
		}
	}

	return result, nil
}

func GetDefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries:    3,
		Interval:      30 * time.Second,
		BackoffFactor: 2.0,
		MaxInterval:   5 * time.Minute,
		RetryOn:       []string{"500", "502", "503", "504", "timeout", "network_error"},
	}
}

func GetDefaultSagaConfig() SagaConfig {
	return SagaConfig{
		TimeoutDuration: 30 * time.Minute,
		MaxRetries:      3,
		RetryInterval:   30 * time.Second,
		AutoCompensate:  true,
		ParallelSteps:   make([][]string, 0),
	}
}