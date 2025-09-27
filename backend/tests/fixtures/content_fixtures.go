package fixtures

import (
	"github.com/google/uuid"
	contentModels "tchat.dev/content/models"
)

// ContentFixtures provides test data for Content models
type ContentFixtures struct {
	*BaseFixture
}

// NewContentFixtures creates a new content fixtures instance
func NewContentFixtures(seed ...int64) *ContentFixtures {
	return &ContentFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// BasicContent creates basic content item for testing
func (c *ContentFixtures) BasicContent(category string, country ...string) *contentModels.ContentItem {
	countryCode := "TH"
	if len(country) > 0 {
		countryCode = c.CountryCode(country[0])
	}

	// Create content value based on type
	contentValue := contentModels.ContentValue{
		"text": c.SEAContent(countryCode, "product"),
		"html": "<p>" + c.SEAContent(countryCode, "product") + "</p>",
	}

	// Create metadata
	metadata := contentModels.ContentMetadata{
		"author":      "Test Author",
		"version":     "1.0",
		"language":    c.Locale(countryCode),
		"region":      countryCode,
		"target_audience": "general",
		"seo_tags":    []string{"test", "content", countryCode},
	}

	return &contentModels.ContentItem{
		ID:        c.UUID("content-" + category + "-" + countryCode),
		Category:  category,
		Type:      contentModels.ContentTypeText,
		Value:     contentValue,
		Metadata:  metadata,
		Status:    contentModels.ContentStatusDraft,
		Tags:      []string{"test", "fixture", countryCode},
		Notes:     nil,
		CreatedAt: c.PastTime(120), // Created 2 hours ago
		UpdatedAt: c.PastTime(60),  // Updated 1 hour ago
	}
}

// PublishedContent creates published content for testing
func (c *ContentFixtures) PublishedContent(category string, country ...string) *contentModels.ContentItem {
	content := c.BasicContent(category, country...)
	content.ID = c.UUID("published-content-" + category)
	content.Status = contentModels.ContentStatusPublished
	content.UpdatedAt = c.PastTime(30) // Published 30 minutes ago
	return content
}

// ArchivedContent creates archived content for testing
func (c *ContentFixtures) ArchivedContent(category string, country ...string) *contentModels.ContentItem {
	content := c.PublishedContent(category, country...)
	content.ID = c.UUID("archived-content-" + category)
	content.Status = contentModels.ContentStatusArchived
	content.UpdatedAt = c.PastTime(15) // Archived 15 minutes ago
	return content
}

// MultilingualContent creates content for multiple languages
func (c *ContentFixtures) MultilingualContent(category string) []*contentModels.ContentItem {
	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}
	contents := make([]*contentModels.ContentItem, 0, len(countries))

	for _, country := range countries {
		content := c.BasicContent(category, country)
		content.ID = c.UUID("multilingual-" + category + "-" + country)
		content.Metadata["language"] = c.Locale(country)
		content.Metadata["region"] = country
		content.Value = contentModels.ContentValue{
			"text": c.SEAContent(country, "greeting"),
			"html": "<p>" + c.SEAContent(country, "greeting") + "</p>",
		}
		contents = append(contents, content)
	}

	return contents
}

// RichContent creates content with rich media for testing
func (c *ContentFixtures) RichContent(category string) *contentModels.ContentItem {
	content := c.BasicContent(category)
	content.ID = c.UUID("rich-content-" + category)
	content.Type = contentModels.ContentTypeHTML

	// Rich content with media
	content.Value = contentModels.ContentValue{
		"text":  c.LoremText(50),
		"html":  "<p>" + c.LoremText(50) + "</p>",
		"image": "https://example.com/images/test-image.jpg",
		"video": "https://example.com/videos/test-video.mp4",
		"audio": "https://example.com/audio/test-audio.mp3",
	}

	content.Metadata["media_count"] = 3
	content.Metadata["has_image"] = true
	content.Metadata["has_video"] = true
	content.Metadata["has_audio"] = true

	return content
}

// ImageContent creates image content for testing
func (c *ContentFixtures) ImageContent(category string) *contentModels.ContentItem {
	content := c.BasicContent(category)
	content.ID = c.UUID("image-content-" + category)
	content.Type = contentModels.ContentTypeJSON

	content.Value = contentModels.ContentValue{
		"url":        "https://example.com/images/test-image.jpg",
		"alt_text":   "Test image for " + category,
		"width":      1920,
		"height":     1080,
		"file_size":  1024000,
		"mime_type":  "image/jpeg",
	}

	content.Metadata["image_type"] = "jpeg"
	content.Metadata["dimensions"] = "1920x1080"

	return content
}

// VideoContent creates video content for testing
func (c *ContentFixtures) VideoContent(category string) *contentModels.ContentItem {
	content := c.BasicContent(category)
	content.ID = c.UUID("video-content-" + category)
	content.Type = contentModels.ContentTypeJSON

	content.Value = contentModels.ContentValue{
		"url":       "https://example.com/videos/test-video.mp4",
		"title":     "Test Video for " + category,
		"duration":  120, // 2 minutes
		"width":     1920,
		"height":    1080,
		"file_size": 50000000, // 50MB
		"mime_type": "video/mp4",
		"thumbnail": "https://example.com/thumbnails/test-video.jpg",
	}

	content.Metadata["video_type"] = "mp4"
	content.Metadata["duration_seconds"] = 120
	content.Metadata["quality"] = "1080p"

	return content
}

// JSONContent creates JSON content for testing
func (c *ContentFixtures) JSONContent(category string) *contentModels.ContentItem {
	content := c.BasicContent(category)
	content.ID = c.UUID("json-content-" + category)
	content.Type = contentModels.ContentTypeJSON

	jsonData := map[string]interface{}{
		"api_version": "v1",
		"endpoints": []string{
			"/api/users",
			"/api/content",
			"/api/payments",
		},
		"features": map[string]bool{
			"chat":     true,
			"payments": true,
			"commerce": true,
		},
		"limits": map[string]int{
			"max_file_size": 10485760, // 10MB
			"rate_limit":    1000,
		},
	}

	content.Value = contentModels.ContentValue(jsonData)
	content.Metadata["schema_version"] = "1.0"
	content.Metadata["data_type"] = "configuration"

	return content
}

// ContentCategoryFixtures provides test data for ContentCategory models
type ContentCategoryFixtures struct {
	*BaseFixture
}

// NewContentCategoryFixtures creates a new content category fixtures instance
func NewContentCategoryFixtures(seed ...int64) *ContentCategoryFixtures {
	return &ContentCategoryFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// RootCategory creates a root category for testing
func (c *ContentCategoryFixtures) RootCategory(name string) *contentModels.ContentCategory {
	desc := "Root category for " + name
	return &contentModels.ContentCategory{
		ID:          c.UUID("category-" + name),
		Name:        name,
		Description: &desc,
		ParentID:    nil,
		Parent:      nil,
		Children:    []contentModels.ContentCategory{},
		IsActive:    true,
		CreatedAt:   c.PastTime(180), // Created 3 hours ago
		UpdatedAt:   c.PastTime(90),  // Updated 1.5 hours ago
	}
}

// SubCategory creates a subcategory for testing
func (c *ContentCategoryFixtures) SubCategory(name string, parentID uuid.UUID) *contentModels.ContentCategory {
	desc := "Subcategory for " + name
	return &contentModels.ContentCategory{
		ID:          c.UUID("subcategory-" + name),
		Name:        name,
		Description: &desc,
		ParentID:    &parentID,
		Parent:      nil, // Will be loaded by ORM
		Children:    []contentModels.ContentCategory{},
		IsActive:    true,
		CreatedAt:   c.PastTime(120), // Created 2 hours ago
		UpdatedAt:   c.PastTime(60),  // Updated 1 hour ago
	}
}

// CategoryHierarchy creates a complete category hierarchy for testing
func (c *ContentCategoryFixtures) CategoryHierarchy() []*contentModels.ContentCategory {
	// Root categories
	categories := []*contentModels.ContentCategory{
		c.RootCategory("announcements"),
		c.RootCategory("promotions"),
		c.RootCategory("help"),
		c.RootCategory("features"),
	}

	// Subcategories for announcements
	announcements := categories[0]
	subcats := []*contentModels.ContentCategory{
		c.SubCategory("system-updates", announcements.ID),
		c.SubCategory("maintenance", announcements.ID),
		c.SubCategory("new-features", announcements.ID),
	}
	categories = append(categories, subcats...)

	// Subcategories for promotions
	promotions := categories[1]
	subcats = []*contentModels.ContentCategory{
		c.SubCategory("seasonal", promotions.ID),
		c.SubCategory("special-offers", promotions.ID),
		c.SubCategory("partnerships", promotions.ID),
	}
	categories = append(categories, subcats...)

	return categories
}

// ContentVersionFixtures provides test data for ContentVersion models
type ContentVersionFixtures struct {
	*BaseFixture
}

// NewContentVersionFixtures creates a new content version fixtures instance
func NewContentVersionFixtures(seed ...int64) *ContentVersionFixtures {
	return &ContentVersionFixtures{
		BaseFixture: NewBaseFixture(seed...),
	}
}

// ContentVersion creates a content version for testing
func (c *ContentVersionFixtures) ContentVersion(contentID uuid.UUID, version int) *contentModels.ContentVersion {
	editorUUID := c.UUID("editor-" + string(rune(version)))
	return &contentModels.ContentVersion{
		ID:        c.UUID("version-" + contentID.String() + "-v" + string(rune(version))),
		ContentID: contentID,
		Version:   version,
		Value: contentModels.ContentValue{
			"text": c.LoremText(20) + " (v" + string(rune(version)) + ")",
			"html": "<p>" + c.LoremText(20) + " (v" + string(rune(version)) + ")</p>",
		},
		Metadata: contentModels.ContentMetadata{
			"version": string(rune(version)),
			"editor":  "Test Editor",
		},
		CreatedAt: c.PastTime(180 - (version * 30)), // Older versions created earlier
		CreatedBy: &editorUUID,
	}
}

// ContentVersions creates multiple versions for a content item
func (c *ContentVersionFixtures) ContentVersions(contentID uuid.UUID, count int) []*contentModels.ContentVersion {
	versions := make([]*contentModels.ContentVersion, 0, count)
	for i := 1; i <= count; i++ {
		version := c.ContentVersion(contentID, i)
		versions = append(versions, version)
	}
	return versions
}

// TestContentData creates a comprehensive set of content test data
func (c *ContentFixtures) TestContentData() map[string]interface{} {
	// Create categories
	categoryFixtures := NewContentCategoryFixtures()
	categories := categoryFixtures.CategoryHierarchy()

	// Create various content types
	contents := []*contentModels.ContentItem{
		c.BasicContent("announcements"),
		c.PublishedContent("promotions"),
		c.ArchivedContent("help"),
		c.RichContent("features"),
		c.ImageContent("gallery"),
		c.VideoContent("tutorials"),
		c.JSONContent("configuration"),
	}

	// Add multilingual content
	multilingualContents := c.MultilingualContent("welcome-message")
	contents = append(contents, multilingualContents...)

	// Create versions for some content
	versionFixtures := NewContentVersionFixtures()
	versions := make([]*contentModels.ContentVersion, 0)
	for _, content := range contents[:3] { // Add versions to first 3 content items
		contentVersions := versionFixtures.ContentVersions(content.ID, 3)
		versions = append(versions, contentVersions...)
	}

	return map[string]interface{}{
		"categories": categories,
		"contents":   contents,
		"versions":   versions,
	}
}

// AllContentFixtures creates a complete set of content-related test data
func AllContentFixtures(seed ...int64) (*ContentFixtures, *ContentCategoryFixtures, *ContentVersionFixtures) {
	return NewContentFixtures(seed...), NewContentCategoryFixtures(seed...), NewContentVersionFixtures(seed...)
}