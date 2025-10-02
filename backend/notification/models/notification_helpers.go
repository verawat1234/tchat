package models

// Helper methods for Notification model to provide compatibility with handler expectations

// GetRecipientID returns the UserID as RecipientID for compatibility
func (n *Notification) GetRecipientID() string {
	return n.UserID.String()
}

// GetSubject returns the Title as Subject for compatibility
func (n *Notification) GetSubject() *string {
	if n.Title == "" {
		return nil
	}
	return &n.Title
}

// GetContent returns the Body as Content for compatibility
func (n *Notification) GetContent() string {
	return n.Body
}

// GetIDString returns the ID as string for compatibility
func (n *Notification) GetIDString() string {
	return n.ID.String()
}

// GetTypeString returns the Type as string for compatibility
func (n *Notification) GetTypeString() string {
	return string(n.Type)
}

// GetChannelString returns the Channel as string for compatibility
func (n *Notification) GetChannelString() string {
	return string(n.Channel)
}

// GetStatusString returns the Status as string for compatibility
func (n *Notification) GetStatusString() string {
	return string(n.Status)
}

// GetPriorityString returns the Priority as string for compatibility
func (n *Notification) GetPriorityString() string {
	return string(n.Priority)
}

// NotificationSubscription helpers for handler compatibility
func (ns *NotificationSubscription) GetIDString() string {
	return ns.ID.String()
}

func (ns *NotificationSubscription) GetUserIDString() string {
	return ns.UserID.String()
}

func (ns *NotificationSubscription) GetChannelString() string {
	return string(ns.Channel)
}

func (ns *NotificationSubscription) GetType() string {
	// Map platform to type for compatibility
	return string(ns.Channel)
}

func (ns *NotificationSubscription) GetEndpoint() string {
	// Return DeviceToken as endpoint for compatibility
	return ns.DeviceToken
}

func (ns *NotificationSubscription) GetEnabled() bool {
	return ns.IsActive
}

// NotificationTemplate helpers for handler compatibility
func (nt *NotificationTemplate) GetIDString() string {
	return nt.ID.String()
}

func (nt *NotificationTemplate) GetTypeString() string {
	return string(nt.Type)
}

func (nt *NotificationTemplate) GetCategoryString() string {
	return string(nt.Category)
}

func (nt *NotificationTemplate) GetSubjectPtr() *string {
	if nt.Subject == "" {
		return nil
	}
	return &nt.Subject
}

func (nt *NotificationTemplate) GetContent() string {
	// Return Body as Content for compatibility
	return nt.Body
}

func (nt *NotificationTemplate) GetLocales() map[string]interface{} {
	// Convert LocalizedVersions to Locales format for compatibility
	locales := make(map[string]interface{})
	for _, version := range nt.LocalizedVersions {
		locales[version.Language] = map[string]interface{}{
			"title":      version.Title,
			"body":       version.Body,
			"image_url":  version.ImageURL,
			"action_url": version.ActionURL,
		}
	}
	return locales
}

func (nt *NotificationTemplate) GetMetadata() map[string]interface{} {
	// Return empty metadata for compatibility (model doesn't have this field)
	return make(map[string]interface{})
}

func (nt *NotificationTemplate) GetActive() bool {
	return nt.IsActive
}

// NotificationPreferences helpers for handler compatibility
func (np *NotificationPreferences) GetUserIDString() string {
	return np.UserID.String()
}

func (np *NotificationPreferences) GetEmailEnabled() bool {
	if np.Channels == nil {
		return true
	}
	enabled, exists := np.Channels["email"]
	if !exists {
		return true // Default to enabled
	}
	return enabled
}

func (np *NotificationPreferences) GetSMSEnabled() bool {
	if np.Channels == nil {
		return true
	}
	enabled, exists := np.Channels["sms"]
	if !exists {
		return true
	}
	return enabled
}

func (np *NotificationPreferences) GetPushEnabled() bool {
	if np.Channels == nil {
		return true
	}
	enabled, exists := np.Channels["push"]
	if !exists {
		return true
	}
	return enabled
}

func (np *NotificationPreferences) GetInAppEnabled() bool {
	if np.Channels == nil {
		return true
	}
	enabled, exists := np.Channels["in_app"]
	if !exists {
		return true
	}
	return enabled
}

func (np *NotificationPreferences) GetQuietHoursMap() map[string]interface{} {
	return map[string]interface{}{
		"enabled": np.QuietHours.Enabled,
		"start":   np.QuietHours.Start,
		"end":     np.QuietHours.End,
	}
}

func (np *NotificationPreferences) GetLanguages() []string {
	// Model doesn't have Languages field, return empty slice
	return []string{}
}

// NotificationAnalytics helpers for handler compatibility
func (na *NotificationAnalytics) GetPeriod(period string) string {
	// Return the period passed in, as model doesn't store it
	return period
}

func (na *NotificationAnalytics) GetTotalOpened() int64 {
	// Use TotalRead as opened count
	return na.TotalRead
}

func (na *NotificationAnalytics) GetTotalClicked() int64 {
	// Model doesn't track clicks separately, return 0
	return 0
}

func (na *NotificationAnalytics) GetDeliveryRate() float64 {
	if na.TotalSent == 0 {
		return 0
	}
	return float64(na.TotalDelivered) / float64(na.TotalSent) * 100
}

func (na *NotificationAnalytics) GetOpenRate() float64 {
	if na.TotalDelivered == 0 {
		return 0
	}
	return float64(na.TotalRead) / float64(na.TotalDelivered) * 100
}

func (na *NotificationAnalytics) GetClickRate() float64 {
	// Model doesn't track clicks, return 0
	return 0
}

func (na *NotificationAnalytics) GetByType() map[string]int64 {
	// Model doesn't have ByType, return empty map
	return make(map[string]int64)
}

func (na *NotificationAnalytics) GetByChannelMap() map[string]int64 {
	// Convert NotificationChannel keys to strings
	result := make(map[string]int64)
	for channel, count := range na.ByChannel {
		result[string(channel)] = count
	}
	return result
}

func (na *NotificationAnalytics) GetTopCategories() []map[string]interface{} {
	// Convert ByCategory to top categories format
	categories := make([]map[string]interface{}, 0, len(na.ByCategory))
	for category, count := range na.ByCategory {
		categories = append(categories, map[string]interface{}{
			"category": string(category),
			"count":    count,
		})
	}
	return categories
}

func (na *NotificationAnalytics) GetRecentActivity() []map[string]interface{} {
	// Model doesn't have recent activity, return empty slice
	return []map[string]interface{}{}
}
