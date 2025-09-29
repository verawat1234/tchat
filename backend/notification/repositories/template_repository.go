package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/notification/models"
)

// TemplateRepositoryImpl implements the TemplateRepository interface
type TemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &TemplateRepositoryImpl{db: db}
}

// Create creates a new notification template
func (r *TemplateRepositoryImpl) Create(ctx context.Context, template *models.NotificationTemplate) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	if template.ID == uuid.Nil {
		template.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(template).Error
}

// GetByID retrieves a template by ID
func (r *TemplateRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update updates an existing template
func (r *TemplateRepositoryImpl) Update(ctx context.Context, template *models.NotificationTemplate) error {
	template.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete deletes a template by ID
func (r *TemplateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.NotificationTemplate{}, "id = ?", id).Error
}

// GetByType retrieves a template by notification type
func (r *TemplateRepositoryImpl) GetByType(ctx context.Context, notificationType string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.WithContext(ctx).
		Where("type = ?", notificationType).
		Where("active = ?", true).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByTypeAndLanguage retrieves a template by type and language
func (r *TemplateRepositoryImpl) GetByTypeAndLanguage(ctx context.Context, notificationType, language string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate

	// First try to find template with specific language
	err := r.db.WithContext(ctx).
		Where("type = ?", notificationType).
		Where("language = ?", language).
		Where("active = ?", true).
		First(&template).Error

	if err == gorm.ErrRecordNotFound {
		// Fallback to default language (en)
		err = r.db.WithContext(ctx).
			Where("type = ?", notificationType).
			Where("language = ? OR language IS NULL", "en").
			Where("active = ?", true).
			Order("language DESC"). // Prioritize specific language over NULL
			First(&template).Error
	}

	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByCategory retrieves templates by category with pagination
func (r *TemplateRepositoryImpl) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*models.NotificationTemplate, error) {
	var templates []*models.NotificationTemplate
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&templates).Error
	return templates, err
}

// GetAll retrieves all templates with pagination
func (r *TemplateRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*models.NotificationTemplate, error) {
	var templates []*models.NotificationTemplate
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&templates).Error
	return templates, err
}

// GetActiveTemplates retrieves only active templates with pagination
func (r *TemplateRepositoryImpl) GetActiveTemplates(ctx context.Context, limit, offset int) ([]*models.NotificationTemplate, error) {
	var templates []*models.NotificationTemplate
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&templates).Error
	return templates, err
}

// GetVariables extracts variable names from a template
func (r *TemplateRepositoryImpl) GetVariables(ctx context.Context, templateID uuid.UUID) ([]string, error) {
	template, err := r.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	variables := make(map[string]bool)

	// Extract variables from title template
	if template.TitleTemplate != "" {
		extractVariables(template.TitleTemplate, variables)
	}

	// Extract variables from body template
	if template.BodyTemplate != "" {
		extractVariables(template.BodyTemplate, variables)
	}

	// Convert map to slice
	result := make([]string, 0, len(variables))
	for variable := range variables {
		result = append(result, variable)
	}

	return result, nil
}

// extractVariables extracts variable names from template text
func extractVariables(template string, variables map[string]bool) {
	// Look for {{variable}} patterns
	start := 0
	for {
		startIdx := strings.Index(template[start:], "{{")
		if startIdx == -1 {
			break
		}
		startIdx += start

		endIdx := strings.Index(template[startIdx+2:], "}}")
		if endIdx == -1 {
			break
		}
		endIdx += startIdx + 2

		variable := strings.TrimSpace(template[startIdx+2 : endIdx])
		if variable != "" {
			variables[variable] = true
		}

		start = endIdx + 2
	}
}

// ValidateTemplate validates a notification template
func (r *TemplateRepositoryImpl) ValidateTemplate(ctx context.Context, template *models.NotificationTemplate) error {
	// Check required fields
	if template.Type == "" {
		return ErrTemplateTypeRequired
	}

	if template.TitleTemplate == "" && template.BodyTemplate == "" {
		return ErrTemplateContentRequired
	}

	// Validate template syntax
	if template.TitleTemplate != "" {
		if err := validateTemplateContent(template.TitleTemplate); err != nil {
			return err
		}
	}

	if template.BodyTemplate != "" {
		if err := validateTemplateContent(template.BodyTemplate); err != nil {
			return err
		}
	}

	// Check for duplicate active templates of the same type and language
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.NotificationTemplate{}).
		Where("type = ?", template.Type).
		Where("active = ?", true)

	if template.Language != "" {
		query = query.Where("language = ?", template.Language)
	} else {
		query = query.Where("language IS NULL")
	}

	if template.ID != uuid.Nil {
		query = query.Where("id != ?", template.ID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrDuplicateActiveTemplate
	}

	return nil
}

// validateTemplateContent validates template content syntax
func validateTemplateContent(content string) error {
	// Check for balanced braces
	openBraces := strings.Count(content, "{{")
	closeBraces := strings.Count(content, "}}")

	if openBraces != closeBraces {
		return ErrInvalidTemplateSyntax
	}

	// Check for nested braces (not allowed)
	inBrace := false
	for i := 0; i < len(content)-1; i++ {
		if content[i:i+2] == "{{" {
			if inBrace {
				return ErrInvalidTemplateSyntax
			}
			inBrace = true
			i++ // Skip next character
		} else if content[i:i+2] == "}}" {
			if !inBrace {
				return ErrInvalidTemplateSyntax
			}
			inBrace = false
			i++ // Skip next character
		}
	}

	if inBrace {
		return ErrInvalidTemplateSyntax
	}

	return nil
}

// Template validation errors
var (
	ErrTemplateTypeRequired     = fmt.Errorf("template type is required")
	ErrTemplateContentRequired  = fmt.Errorf("template must have either title or body content")
	ErrInvalidTemplateSyntax    = fmt.Errorf("invalid template syntax")
	ErrDuplicateActiveTemplate  = fmt.Errorf("active template already exists for this type and language")
)