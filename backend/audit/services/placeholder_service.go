package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tchat/backend/audit/models"
	"gorm.io/gorm"
)

type PlaceholderService struct {
	db *gorm.DB
}

type PlaceholderFilters struct {
	ServiceID  string
	Platform   string
	Status     string
	Priority   string
	AssignedTo string
	Regional   string
	SortBy     string
	SortOrder  string
	Limit      int
	Offset     int
}

type PlaceholderStats struct {
	TotalCount         int64                  `json:"totalCount"`
	ByStatus          map[string]int64       `json:"byStatus"`
	ByPriority        map[string]int64       `json:"byPriority"`
	ByPlatform        map[string]int64       `json:"byPlatform"`
	ByPerformanceImpact map[string]int64     `json:"byPerformanceImpact"`
	ByRegion          map[string]int64       `json:"byRegion"`
	AverageAge        float64                `json:"averageAge"`
	CompletionRate    float64                `json:"completionRate"`
	EstimatedHours    float64                `json:"estimatedHours"`
	CriticalCount     int64                  `json:"criticalCount"`
	OverdueCount      int64                  `json:"overdueCount"`
}

func NewPlaceholderService(db *gorm.DB) *PlaceholderService {
	return &PlaceholderService{db: db}
}

func (s *PlaceholderService) GetPlaceholders(ctx context.Context, filters *PlaceholderFilters) ([]models.PlaceholderItem, int64, error) {
	var placeholders []models.PlaceholderItem
	var total int64

	// Build query with filters
	query := s.db.WithContext(ctx).Model(&models.PlaceholderItem{})

	// Apply filters
	if filters.ServiceID != "" {
		query = query.Where("service_id = ?", filters.ServiceID)
	}
	if filters.Platform != "" {
		query = query.Where("platform = ?", filters.Platform)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Priority != "" {
		query = query.Where("priority = ?", filters.Priority)
	}
	if filters.AssignedTo != "" {
		query = query.Where("assigned_to = ?", filters.AssignedTo)
	}
	if filters.Regional != "" {
		query = query.Where("regional_context = ?", filters.Regional)
	}

	// Get total count for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count placeholders: %w", err)
	}

	// Apply sorting
	switch filters.SortBy {
	case "priority":
		if filters.SortOrder == "asc" {
			query = query.Order("priority ASC")
		} else {
			query = query.Order("CASE priority WHEN 'CRITICAL' THEN 4 WHEN 'HIGH' THEN 3 WHEN 'MEDIUM' THEN 2 WHEN 'LOW' THEN 1 ELSE 0 END DESC")
		}
	case "created_at":
		if filters.SortOrder == "asc" {
			query = query.Order("created_at ASC")
		} else {
			query = query.Order("created_at DESC")
		}
	case "updated_at":
		if filters.SortOrder == "asc" {
			query = query.Order("updated_at ASC")
		} else {
			query = query.Order("updated_at DESC")
		}
	case "performance_impact":
		if filters.SortOrder == "asc" {
			query = query.Order("performance_impact ASC")
		} else {
			query = query.Order("CASE performance_impact WHEN 'CRITICAL' THEN 5 WHEN 'HIGH' THEN 4 WHEN 'MEDIUM' THEN 3 WHEN 'LOW' THEN 2 WHEN 'NONE' THEN 1 ELSE 0 END DESC")
		}
	default:
		// Default to priority DESC
		query = query.Order("CASE priority WHEN 'CRITICAL' THEN 4 WHEN 'HIGH' THEN 3 WHEN 'MEDIUM' THEN 2 WHEN 'LOW' THEN 1 ELSE 0 END DESC")
	}

	// Apply pagination
	query = query.Limit(filters.Limit).Offset(filters.Offset)

	// Execute query
	if err := query.Find(&placeholders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve placeholders: %w", err)
	}

	return placeholders, total, nil
}

func (s *PlaceholderService) CreatePlaceholder(ctx context.Context, placeholder *models.PlaceholderItem) (*models.PlaceholderItem, error) {
	// Generate UUID for new placeholder
	placeholder.ID = uuid.New().String()
	placeholder.DetectionDate = time.Now()
	placeholder.CreatedAt = time.Now()
	placeholder.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Create(placeholder).Error; err != nil {
		return nil, fmt.Errorf("failed to create placeholder: %w", err)
	}

	return placeholder, nil
}

func (s *PlaceholderService) UpdatePlaceholder(ctx context.Context, placeholderID string, updates map[string]interface{}) (*models.PlaceholderItem, error) {
	var placeholder models.PlaceholderItem

	// First, get the existing placeholder
	if err := s.db.WithContext(ctx).Where("id = ?", placeholderID).First(&placeholder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("placeholder not found")
		}
		return nil, fmt.Errorf("failed to retrieve placeholder: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Handle status changes
	if status, exists := updates["status"]; exists && status == string(models.StatusCompleted) {
		completedAt := time.Now()
		updates["completed_at"] = &completedAt
	}

	// Update the placeholder
	if err := s.db.WithContext(ctx).Model(&placeholder).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update placeholder: %w", err)
	}

	// Retrieve updated placeholder
	if err := s.db.WithContext(ctx).Where("id = ?", placeholderID).First(&placeholder).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve updated placeholder: %w", err)
	}

	return &placeholder, nil
}

func (s *PlaceholderService) GetPlaceholderStats(ctx context.Context, serviceID, platform string) (*PlaceholderStats, error) {
	var stats PlaceholderStats

	// Build base query
	query := s.db.WithContext(ctx).Model(&models.PlaceholderItem{})
	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	// Get total count
	if err := query.Count(&stats.TotalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get statistics by status
	stats.ByStatus = make(map[string]int64)
	var statusRows []struct {
		Status string
		Count  int64
	}

	statusQuery := query.Session(&gorm.Session{}).Select("status, COUNT(*) as count").Group("status")
	if err := statusQuery.Scan(&statusRows).Error; err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}

	for _, row := range statusRows {
		stats.ByStatus[row.Status] = row.Count
	}

	// Get statistics by priority
	stats.ByPriority = make(map[string]int64)
	var priorityRows []struct {
		Priority string
		Count    int64
	}

	priorityQuery := query.Session(&gorm.Session{}).Select("priority, COUNT(*) as count").Group("priority")
	if err := priorityQuery.Scan(&priorityRows).Error; err != nil {
		return nil, fmt.Errorf("failed to get priority stats: %w", err)
	}

	for _, row := range priorityRows {
		stats.ByPriority[row.Priority] = row.Count
	}

	// Get statistics by platform
	stats.ByPlatform = make(map[string]int64)
	var platformRows []struct {
		Platform string
		Count    int64
	}

	platformQuery := query.Session(&gorm.Session{}).Select("platform, COUNT(*) as count").Group("platform")
	if err := platformQuery.Scan(&platformRows).Error; err != nil {
		return nil, fmt.Errorf("failed to get platform stats: %w", err)
	}

	for _, row := range platformRows {
		stats.ByPlatform[row.Platform] = row.Count
	}

	// Get statistics by performance impact
	stats.ByPerformanceImpact = make(map[string]int64)
	var impactRows []struct {
		PerformanceImpact string
		Count             int64
	}

	impactQuery := query.Session(&gorm.Session{}).Select("performance_impact, COUNT(*) as count").Group("performance_impact")
	if err := impactQuery.Scan(&impactRows).Error; err != nil {
		return nil, fmt.Errorf("failed to get performance impact stats: %w", err)
	}

	for _, row := range impactRows {
		stats.ByPerformanceImpact[row.PerformanceImpact] = row.Count
	}

	// Get statistics by region (Southeast Asian markets)
	stats.ByRegion = make(map[string]int64)
	var regionRows []struct {
		RegionalContext string
		Count           int64
	}

	regionQuery := query.Session(&gorm.Session{}).Where("regional_context IS NOT NULL").Select("regional_context, COUNT(*) as count").Group("regional_context")
	if err := regionQuery.Scan(&regionRows).Error; err != nil {
		return nil, fmt.Errorf("failed to get region stats: %w", err)
	}

	for _, row := range regionRows {
		stats.ByRegion[row.RegionalContext] = row.Count
	}

	// Calculate average age in days
	var avgAgeResult struct {
		AvgAge float64
	}

	avgAgeQuery := query.Session(&gorm.Session{}).Select("AVG(EXTRACT(EPOCH FROM (NOW() - detection_date))/86400) as avg_age")
	if err := avgAgeQuery.Scan(&avgAgeResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average age: %w", err)
	}
	stats.AverageAge = avgAgeResult.AvgAge

	// Calculate completion rate
	completedCount := stats.ByStatus[string(models.StatusCompleted)]
	if stats.TotalCount > 0 {
		stats.CompletionRate = float64(completedCount) / float64(stats.TotalCount) * 100
	}

	// Calculate total estimated hours remaining
	var estimatedHoursResult struct {
		TotalHours float64
	}

	hoursQuery := query.Session(&gorm.Session{}).Where("status != ? AND estimated_hours IS NOT NULL", string(models.StatusCompleted)).Select("SUM(estimated_hours) as total_hours")
	if err := hoursQuery.Scan(&estimatedHoursResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate estimated hours: %w", err)
	}
	stats.EstimatedHours = estimatedHoursResult.TotalHours

	// Get critical count
	stats.CriticalCount = stats.ByPriority[string(models.PriorityCritical)]

	// Calculate overdue count (items older than 30 days and not completed)
	var overdueResult struct {
		Count int64
	}

	overdueQuery := query.Session(&gorm.Session{}).Where("status != ? AND detection_date < ?", string(models.StatusCompleted), time.Now().AddDate(0, 0, -30)).Select("COUNT(*) as count")
	if err := overdueQuery.Scan(&overdueResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate overdue count: %w", err)
	}
	stats.OverdueCount = overdueResult.Count

	return &stats, nil
}

func (s *PlaceholderService) BulkUpdatePlaceholders(ctx context.Context, placeholderIDs []string, updates map[string]interface{}) (int64, error) {
	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Handle status changes for bulk update
	if status, exists := updates["status"]; exists && status == string(models.StatusCompleted) {
		completedAt := time.Now()
		updates["completed_at"] = &completedAt
	}

	result := s.db.WithContext(ctx).Model(&models.PlaceholderItem{}).Where("id IN ?", placeholderIDs).Updates(updates)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to bulk update placeholders: %w", result.Error)
	}

	return result.RowsAffected, nil
}

func (s *PlaceholderService) GetPlaceholdersByRegion(ctx context.Context, region string, limit, offset int) ([]models.PlaceholderItem, int64, error) {
	var placeholders []models.PlaceholderItem
	var total int64

	// Build query for regional placeholders
	query := s.db.WithContext(ctx).Model(&models.PlaceholderItem{}).Where("regional_context = ?", region)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count regional placeholders: %w", err)
	}

	// Apply sorting (priority-based for regional optimization)
	query = query.Order("CASE priority WHEN 'CRITICAL' THEN 4 WHEN 'HIGH' THEN 3 WHEN 'MEDIUM' THEN 2 WHEN 'LOW' THEN 1 ELSE 0 END DESC")

	// Apply pagination
	query = query.Limit(limit).Offset(offset)

	// Execute query
	if err := query.Find(&placeholders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve regional placeholders: %w", err)
	}

	return placeholders, total, nil
}

func (s *PlaceholderService) ArchivePlaceholder(ctx context.Context, placeholderID, archiveReason, reviewedBy string) error {
	var placeholder models.PlaceholderItem

	// Get the placeholder
	if err := s.db.WithContext(ctx).Where("id = ?", placeholderID).First(&placeholder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("placeholder not found")
		}
		return fmt.Errorf("failed to retrieve placeholder: %w", err)
	}

	// Check if placeholder is completed
	if placeholder.Status != string(models.StatusCompleted) {
		return fmt.Errorf("placeholder not completed")
	}

	// Update with archive information
	now := time.Now()
	updates := map[string]interface{}{
		"status":            string(models.StatusCancelled), // Using cancelled as archived status
		"resolution_notes":  archiveReason,
		"reviewed_by":       reviewedBy,
		"reviewed_at":       &now,
		"updated_at":        now,
	}

	if err := s.db.WithContext(ctx).Model(&placeholder).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to archive placeholder: %w", err)
	}

	return nil
}