package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Presence represents a user's online presence and activity status
type Presence struct {
	ID           uuid.UUID        `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID       uuid.UUID        `json:"user_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	Status       PresenceStatus   `json:"status" gorm:"type:varchar(20);not null;default:'offline'"`
	LastSeen     *time.Time       `json:"last_seen,omitempty" gorm:"index"`
	IsOnline     bool             `json:"is_online" gorm:"default:false;index"`
	Platform     Platform         `json:"platform" gorm:"type:varchar(20)"`
	DeviceInfo   DeviceInfo       `json:"device_info" gorm:"type:json"`
	Location     *UserLocation    `json:"location,omitempty" gorm:"type:json"`
	Activity     ActivityStatus   `json:"activity" gorm:"type:varchar(30);default:'idle'"`
	CustomStatus *string          `json:"custom_status,omitempty" gorm:"type:varchar(100)"`
	Privacy      PresencePrivacy  `json:"privacy" gorm:"type:json"`
	Metadata     PresenceMetadata `json:"metadata" gorm:"type:json"`
	CreatedAt    time.Time        `json:"created_at" gorm:"not null"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"not null"`
}

// PresenceStatus represents the overall presence status of a user
type PresenceStatus string

const (
	PresenceStatusOnline       PresenceStatus = "online"        // Actively using the app
	PresenceStatusAway         PresenceStatus = "away"          // Away from keyboard/device
	PresenceStatusBusy         PresenceStatus = "busy"          // Do not disturb
	PresenceStatusInvisible    PresenceStatus = "invisible"     // Appear offline to others
	PresenceStatusOffline      PresenceStatus = "offline"       // Not connected
	PresenceStatusIdle         PresenceStatus = "idle"          // Inactive for a period
)

// Platform represents the platform/device type being used
type Platform string

const (
	PlatformWeb     Platform = "web"        // Web browser
	PlatformMobile  Platform = "mobile"     // Mobile app (iOS/Android)
	PlatformDesktop Platform = "desktop"    // Desktop application
	PlatformTablet  Platform = "tablet"     // Tablet application
	PlatformAPI     Platform = "api"        // API access
	PlatformBot     Platform = "bot"        // Bot/automation
)

// ActivityStatus represents what the user is currently doing
type ActivityStatus string

const (
	ActivityStatusIdle     ActivityStatus = "idle"         // Not actively using
	ActivityStatusTyping   ActivityStatus = "typing"       // Currently typing
	ActivityStatusCalling  ActivityStatus = "calling"      // In a voice/video call
	ActivityStatusGaming   ActivityStatus = "gaming"       // Playing games
	ActivityStatusShopping ActivityStatus = "shopping"     // Using commerce features
	ActivityStatusStreaming ActivityStatus = "streaming"   // Streaming content
	ActivityStatusReading  ActivityStatus = "reading"      // Reading messages/content
)

// DeviceInfo represents information about the user's device
type DeviceInfo struct {
	DeviceID     string `json:"device_id,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	OS           string `json:"os,omitempty"`           // iOS, Android, Windows, macOS, Linux
	OSVersion    string `json:"os_version,omitempty"`
	AppVersion   string `json:"app_version,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	Country      string `json:"country,omitempty"`      // ISO country code
	Timezone     string `json:"timezone,omitempty"`     // IANA timezone
	Language     string `json:"language,omitempty"`     // ISO language code
	NetworkType  string `json:"network_type,omitempty"` // wifi, cellular, ethernet
}

// UserLocation represents optional location information
type UserLocation struct {
	Country   string  `json:"country,omitempty"`    // ISO country code
	Region    string  `json:"region,omitempty"`     // State/province
	City      string  `json:"city,omitempty"`       // City name
	Latitude  float64 `json:"latitude,omitempty"`   // Precise coordinates (if permitted)
	Longitude float64 `json:"longitude,omitempty"`  // Precise coordinates (if permitted)
	Accuracy  int     `json:"accuracy,omitempty"`   // Location accuracy in meters
	Timestamp time.Time `json:"timestamp,omitempty"` // When location was captured
}

// PresencePrivacy represents privacy settings for presence information
type PresencePrivacy struct {
	ShowOnlineStatus  bool `json:"show_online_status"`   // Show online/offline status
	ShowLastSeen      bool `json:"show_last_seen"`       // Show last seen timestamp
	ShowActivity      bool `json:"show_activity"`        // Show current activity
	ShowLocation      bool `json:"show_location"`        // Show location information
	ShowToEveryone    bool `json:"show_to_everyone"`     // Show to all users vs contacts only
	ShowToContacts    bool `json:"show_to_contacts"`     // Show to contacts
	ShowInBusinessHours bool `json:"show_in_business_hours"` // Show during business hours only
}

// PresenceMetadata represents additional metadata and analytics
type PresenceMetadata struct {
	SessionDuration    time.Duration `json:"session_duration,omitempty"`    // Current session duration
	DailyActiveTime    time.Duration `json:"daily_active_time,omitempty"`   // Today's total active time
	TotalSessions      int           `json:"total_sessions,omitempty"`      // Total sessions today
	MessagesSent       int           `json:"messages_sent,omitempty"`       // Messages sent today
	MessagesReceived   int           `json:"messages_received,omitempty"`   // Messages received today
	CallsParticipated  int           `json:"calls_participated,omitempty"`  // Calls participated today
	LastActiveDialog   *uuid.UUID    `json:"last_active_dialog,omitempty"`  // Last dialog user was active in
	FeatureUsage       map[string]int `json:"feature_usage,omitempty"`      // Feature usage counters
	BusinessHours      *BusinessHours `json:"business_hours,omitempty"`     // Business hours for this user
	AutoAwayAfter      *time.Duration `json:"auto_away_after,omitempty"`    // Auto-away timeout
	AutoOfflineAfter   *time.Duration `json:"auto_offline_after,omitempty"` // Auto-offline timeout
}

// BusinessHours represents business hours for presence automation
type BusinessHours struct {
	Enabled   bool   `json:"enabled"`
	Timezone  string `json:"timezone"`  // IANA timezone
	StartTime string `json:"start_time"` // HH:MM format
	EndTime   string `json:"end_time"`   // HH:MM format
	Days      []int  `json:"days"`       // 0=Sunday, 1=Monday, etc.
}

// Value implements the driver.Valuer interface for DeviceInfo
func (di DeviceInfo) Value() (driver.Value, error) {
	return json.Marshal(di)
}

// Scan implements the sql.Scanner interface for DeviceInfo
func (di *DeviceInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into DeviceInfo", value)
	}

	return json.Unmarshal(jsonData, di)
}

// Value implements the driver.Valuer interface for UserLocation
func (ul UserLocation) Value() (driver.Value, error) {
	return json.Marshal(ul)
}

// Scan implements the sql.Scanner interface for UserLocation
func (ul *UserLocation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into UserLocation", value)
	}

	return json.Unmarshal(jsonData, ul)
}

// Value implements the driver.Valuer interface for PresencePrivacy
func (pp PresencePrivacy) Value() (driver.Value, error) {
	return json.Marshal(pp)
}

// Scan implements the sql.Scanner interface for PresencePrivacy
func (pp *PresencePrivacy) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into PresencePrivacy", value)
	}

	return json.Unmarshal(jsonData, pp)
}

// Value implements the driver.Valuer interface for PresenceMetadata
func (pm PresenceMetadata) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

// Scan implements the sql.Scanner interface for PresenceMetadata
func (pm *PresenceMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into PresenceMetadata", value)
	}

	return json.Unmarshal(jsonData, pm)
}

// ValidPresenceStatuses returns all supported presence statuses
func ValidPresenceStatuses() []PresenceStatus {
	return []PresenceStatus{
		PresenceStatusOnline,
		PresenceStatusAway,
		PresenceStatusBusy,
		PresenceStatusInvisible,
		PresenceStatusOffline,
		PresenceStatusIdle,
	}
}

// ValidPlatforms returns all supported platforms
func ValidPlatforms() []Platform {
	return []Platform{
		PlatformWeb,
		PlatformMobile,
		PlatformDesktop,
		PlatformTablet,
		PlatformAPI,
		PlatformBot,
	}
}

// ValidActivityStatuses returns all supported activity statuses
func ValidActivityStatuses() []ActivityStatus {
	return []ActivityStatus{
		ActivityStatusIdle,
		ActivityStatusTyping,
		ActivityStatusCalling,
		ActivityStatusGaming,
		ActivityStatusShopping,
		ActivityStatusStreaming,
		ActivityStatusReading,
	}
}

// IsValid validates if the presence status is supported
func (ps PresenceStatus) IsValid() bool {
	for _, valid := range ValidPresenceStatuses() {
		if ps == valid {
			return true
		}
	}
	return false
}

// IsValid validates if the platform is supported
func (p Platform) IsValid() bool {
	for _, valid := range ValidPlatforms() {
		if p == valid {
			return true
		}
	}
	return false
}

// IsValid validates if the activity status is supported
func (as ActivityStatus) IsValid() bool {
	for _, valid := range ValidActivityStatuses() {
		if as == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of PresenceStatus
func (ps PresenceStatus) String() string {
	return string(ps)
}

// String returns the string representation of Platform
func (p Platform) String() string {
	return string(p)
}

// String returns the string representation of ActivityStatus
func (as ActivityStatus) String() string {
	return string(as)
}

// IsVisibleToOthers checks if presence should be visible to other users
func (ps PresenceStatus) IsVisibleToOthers() bool {
	return ps != PresenceStatusInvisible
}

// RequiresAuthentication checks if this platform requires authentication
func (p Platform) RequiresAuthentication() bool {
	return p != PlatformAPI && p != PlatformBot
}

// IsActiveStatus checks if this is an active status
func (as ActivityStatus) IsActiveStatus() bool {
	return as != ActivityStatusIdle
}

// Validate performs comprehensive validation on the Presence model
func (p *Presence) Validate() error {
	var errs []string

	// User ID validation
	if p.UserID == uuid.Nil {
		errs = append(errs, "user_id is required")
	}

	// Status validation
	if !p.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid presence status: %s", p.Status))
	}

	// Platform validation
	if p.Platform != "" && !p.Platform.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid platform: %s", p.Platform))
	}

	// Activity validation
	if p.Activity != "" && !p.Activity.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid activity status: %s", p.Activity))
	}

	// Custom status validation
	if p.CustomStatus != nil {
		if len(*p.CustomStatus) > 100 {
			errs = append(errs, "custom_status cannot exceed 100 characters")
		}
		if strings.TrimSpace(*p.CustomStatus) == "" {
			p.CustomStatus = nil // Set to nil if empty
		}
	}

	// Location validation
	if err := p.validateLocation(); err != nil {
		errs = append(errs, fmt.Sprintf("invalid location: %v", err))
	}

	// Device info validation
	if err := p.validateDeviceInfo(); err != nil {
		errs = append(errs, fmt.Sprintf("invalid device info: %v", err))
	}

	// Metadata validation
	if err := p.validateMetadata(); err != nil {
		errs = append(errs, fmt.Sprintf("invalid metadata: %v", err))
	}

	// Business logic validation
	if err := p.validateBusinessLogic(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validateLocation validates the location information
func (p *Presence) validateLocation() error {
	if p.Location == nil {
		return nil
	}

	loc := p.Location

	// Validate coordinates if provided
	if loc.Latitude != 0 || loc.Longitude != 0 {
		if loc.Latitude < -90 || loc.Latitude > 90 {
			return errors.New("latitude must be between -90 and 90")
		}
		if loc.Longitude < -180 || loc.Longitude > 180 {
			return errors.New("longitude must be between -180 and 180")
		}
	}

	// Validate accuracy
	if loc.Accuracy < 0 {
		return errors.New("accuracy cannot be negative")
	}

	// Validate country code
	if loc.Country != "" && len(loc.Country) != 2 {
		return errors.New("country must be a 2-letter ISO code")
	}

	return nil
}

// validateDeviceInfo validates the device information
func (p *Presence) validateDeviceInfo() error {
	device := p.DeviceInfo

	// Validate country code
	if device.Country != "" && len(device.Country) != 2 {
		return errors.New("device country must be a 2-letter ISO code")
	}

	// Validate language code
	if device.Language != "" && len(device.Language) < 2 {
		return errors.New("language must be a valid language code")
	}

	// Validate IP address format (basic check)
	if device.IPAddress != "" && len(device.IPAddress) < 7 {
		return errors.New("invalid IP address format")
	}

	return nil
}

// validateMetadata validates the metadata information
func (p *Presence) validateMetadata() error {
	meta := p.Metadata

	// Validate counters
	if meta.TotalSessions < 0 {
		return errors.New("total_sessions cannot be negative")
	}
	if meta.MessagesSent < 0 {
		return errors.New("messages_sent cannot be negative")
	}
	if meta.MessagesReceived < 0 {
		return errors.New("messages_received cannot be negative")
	}
	if meta.CallsParticipated < 0 {
		return errors.New("calls_participated cannot be negative")
	}

	// Validate durations
	if meta.SessionDuration < 0 {
		return errors.New("session_duration cannot be negative")
	}
	if meta.DailyActiveTime < 0 {
		return errors.New("daily_active_time cannot be negative")
	}

	// Validate auto-timeout durations
	if meta.AutoAwayAfter != nil && *meta.AutoAwayAfter < time.Minute {
		return errors.New("auto_away_after must be at least 1 minute")
	}
	if meta.AutoOfflineAfter != nil && *meta.AutoOfflineAfter < 5*time.Minute {
		return errors.New("auto_offline_after must be at least 5 minutes")
	}

	// Validate business hours
	if meta.BusinessHours != nil {
		if err := p.validateBusinessHours(meta.BusinessHours); err != nil {
			return fmt.Errorf("invalid business hours: %v", err)
		}
	}

	return nil
}

// validateBusinessHours validates business hours configuration
func (p *Presence) validateBusinessHours(bh *BusinessHours) error {
	if !bh.Enabled {
		return nil
	}

	// Validate timezone
	if bh.Timezone == "" {
		return errors.New("timezone is required when business hours are enabled")
	}

	// Validate time format (HH:MM)
	if bh.StartTime == "" || bh.EndTime == "" {
		return errors.New("start_time and end_time are required")
	}

	// Basic time format validation
	if len(bh.StartTime) != 5 || len(bh.EndTime) != 5 {
		return errors.New("time format must be HH:MM")
	}

	// Validate days (0-6)
	if len(bh.Days) == 0 {
		return errors.New("at least one day must be specified")
	}
	for _, day := range bh.Days {
		if day < 0 || day > 6 {
			return errors.New("days must be between 0 (Sunday) and 6 (Saturday)")
		}
	}

	return nil
}

// validateBusinessLogic validates business rules and constraints
func (p *Presence) validateBusinessLogic() error {
	// Online status consistency
	if p.Status == PresenceStatusOnline && !p.IsOnline {
		return errors.New("status is online but is_online is false")
	}
	if p.Status == PresenceStatusOffline && p.IsOnline {
		return errors.New("status is offline but is_online is true")
	}

	// Last seen validation
	if p.IsOnline && p.LastSeen != nil {
		// If user is currently online, last_seen should be recent
		if time.Since(*p.LastSeen) > 5*time.Minute {
			return errors.New("last_seen is too old for online user")
		}
	}

	// Activity validation for online users
	if p.IsOnline && p.Activity == "" {
		p.Activity = ActivityStatusIdle // Set default activity
	}

	return nil
}

// BeforeCreate sets up the presence before database creation
func (p *Presence) BeforeCreate() error {
	// Generate UUID if not set
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	// Set default values
	p.setDefaults()

	// Validate before creation
	return p.Validate()
}

// BeforeUpdate sets up the presence before database update
func (p *Presence) BeforeUpdate() error {
	// Update timestamp
	p.UpdatedAt = time.Now().UTC()

	// Update last seen if going online
	if p.IsOnline {
		now := time.Now().UTC()
		p.LastSeen = &now
	}

	// Validate before update
	return p.Validate()
}

// setDefaults sets default values for the presence
func (p *Presence) setDefaults() {
	// Set default privacy settings
	if p.Privacy == (PresencePrivacy{}) {
		p.Privacy = PresencePrivacy{
			ShowOnlineStatus:     true,
			ShowLastSeen:         true,
			ShowActivity:         false,
			ShowLocation:         false,
			ShowToEveryone:       false,
			ShowToContacts:       true,
			ShowInBusinessHours:  false,
		}
	}

	// Initialize metadata if empty
	if p.Metadata == (PresenceMetadata{}) {
		autoAway := 10 * time.Minute
		autoOffline := 30 * time.Minute
		p.Metadata = PresenceMetadata{
			AutoAwayAfter:    &autoAway,
			AutoOfflineAfter: &autoOffline,
			FeatureUsage:     make(map[string]int),
		}
	}

	// Set default activity if online
	if p.IsOnline && p.Activity == "" {
		p.Activity = ActivityStatusIdle
	}
}

// UpdateActivity updates the user's current activity
func (p *Presence) UpdateActivity(activity ActivityStatus) error {
	if !activity.IsValid() {
		return fmt.Errorf("invalid activity status: %s", activity)
	}

	p.Activity = activity
	p.UpdatedAt = time.Now().UTC()

	// Update last seen if setting an active status
	if activity.IsActiveStatus() {
		now := time.Now().UTC()
		p.LastSeen = &now
	}

	return nil
}

// SetOnline sets the user as online with optional platform and device info
func (p *Presence) SetOnline(platform Platform, deviceInfo *DeviceInfo) error {
	if platform != "" && !platform.IsValid() {
		return fmt.Errorf("invalid platform: %s", platform)
	}

	p.IsOnline = true
	p.Status = PresenceStatusOnline
	p.Platform = platform
	if deviceInfo != nil {
		p.DeviceInfo = *deviceInfo
	}

	now := time.Now().UTC()
	p.LastSeen = &now
	p.UpdatedAt = now

	// Start new session
	p.Metadata.TotalSessions++

	return nil
}

// SetOffline sets the user as offline
func (p *Presence) SetOffline() error {
	p.IsOnline = false
	p.Status = PresenceStatusOffline
	p.Activity = ActivityStatusIdle
	p.UpdatedAt = time.Now().UTC()

	return nil
}

// SetAway sets the user as away
func (p *Presence) SetAway() error {
	p.Status = PresenceStatusAway
	p.Activity = ActivityStatusIdle
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// SetBusy sets the user as busy/do not disturb
func (p *Presence) SetBusy(customStatus *string) error {
	p.Status = PresenceStatusBusy
	p.CustomStatus = customStatus
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// SetInvisible sets the user as invisible (appear offline)
func (p *Presence) SetInvisible() error {
	p.Status = PresenceStatusInvisible
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateLocation updates the user's location if privacy allows
func (p *Presence) UpdateLocation(location *UserLocation) error {
	if !p.Privacy.ShowLocation {
		return errors.New("location sharing is disabled")
	}

	if location != nil {
		location.Timestamp = time.Now().UTC()
	}

	p.Location = location
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// IncrementFeatureUsage increments a feature usage counter
func (p *Presence) IncrementFeatureUsage(feature string) {
	if p.Metadata.FeatureUsage == nil {
		p.Metadata.FeatureUsage = make(map[string]int)
	}
	p.Metadata.FeatureUsage[feature]++
	p.UpdatedAt = time.Now().UTC()
}

// UpdateSessionMetrics updates session-related metrics
func (p *Presence) UpdateSessionMetrics(messagesSent, messagesReceived, callsParticipated int) {
	p.Metadata.MessagesSent += messagesSent
	p.Metadata.MessagesReceived += messagesReceived
	p.Metadata.CallsParticipated += callsParticipated
	p.UpdatedAt = time.Now().UTC()
}

// ShouldShowToUser checks if presence should be shown to a specific user
func (p *Presence) ShouldShowToUser(requestingUserID uuid.UUID, isContact bool) bool {
	// Don't show invisible status
	if p.Status == PresenceStatusInvisible {
		return false
	}

	// Show to the user themselves
	if p.UserID == requestingUserID {
		return true
	}

	// Check privacy settings
	if p.Privacy.ShowToEveryone {
		return true
	}

	if p.Privacy.ShowToContacts && isContact {
		return true
	}

	return false
}

// GetVisiblePresence returns presence information visible to a specific user
func (p *Presence) GetVisiblePresence(requestingUserID uuid.UUID, isContact bool) map[string]interface{} {
	if !p.ShouldShowToUser(requestingUserID, isContact) {
		return map[string]interface{}{
			"user_id": p.UserID,
			"status":  PresenceStatusOffline,
			"is_online": false,
		}
	}

	response := map[string]interface{}{
		"user_id":    p.UserID,
		"status":     p.Status,
		"is_online":  p.IsOnline,
		"platform":   p.Platform,
		"updated_at": p.UpdatedAt,
	}

	// Include last seen if privacy allows
	if p.Privacy.ShowLastSeen && p.LastSeen != nil {
		response["last_seen"] = p.LastSeen
	}

	// Include activity if privacy allows
	if p.Privacy.ShowActivity {
		response["activity"] = p.Activity
	}

	// Include custom status
	if p.CustomStatus != nil {
		response["custom_status"] = p.CustomStatus
	}

	// Include location if privacy allows
	if p.Privacy.ShowLocation && p.Location != nil {
		response["location"] = p.Location
	}

	return response
}

// IsInBusinessHours checks if current time is within business hours
func (p *Presence) IsInBusinessHours() bool {
	if p.Metadata.BusinessHours == nil || !p.Metadata.BusinessHours.Enabled {
		return false
	}

	// This is a simplified check - in production, you'd use proper timezone handling
	now := time.Now()
	weekday := int(now.Weekday())

	// Check if current day is in business days
	for _, day := range p.Metadata.BusinessHours.Days {
		if day == weekday {
			// TODO: Implement proper time range checking with timezone
			return true
		}
	}

	return false
}

// ShouldAutoAway checks if user should be automatically set to away
func (p *Presence) ShouldAutoAway() bool {
	if p.Metadata.AutoAwayAfter == nil || p.LastSeen == nil {
		return false
	}

	return time.Since(*p.LastSeen) > *p.Metadata.AutoAwayAfter
}

// ShouldAutoOffline checks if user should be automatically set to offline
func (p *Presence) ShouldAutoOffline() bool {
	if p.Metadata.AutoOfflineAfter == nil || p.LastSeen == nil {
		return false
	}

	return time.Since(*p.LastSeen) > *p.Metadata.AutoOfflineAfter
}

// ToPublicPresence returns a sanitized version for public API responses
func (p *Presence) ToPublicPresence(forUserID uuid.UUID, isContact bool) map[string]interface{} {
	return p.GetVisiblePresence(forUserID, isContact)
}