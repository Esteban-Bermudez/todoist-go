package todoist

type User struct {
	ActivatedUser         bool           `json:"activated_user"`
	AutoReminder          int            `json:"auto_reminder"`
	AvatarBig             string         `json:"avatar_big"`
	AvatarMedium          string         `json:"avatar_medium"`
	AvatarS640            string         `json:"avatar_s640"`
	AvatarSmall           string         `json:"avatar_small"`
	BusinessAccountID     *string        `json:"business_account_id"`
	DailyGoal             int            `json:"daily_goal"`
	DateFormat            int            `json:"date_format"`
	DaysOff               []int          `json:"days_off"`
	DeletedAt             *string        `json:"deleted_at,omitempty"`
	Email                 string         `json:"email"`
	FeatureIdentifier     string         `json:"feature_identifier"`
	Features              map[string]any `json:"features"`
	FullName              string         `json:"full_name"`
	HasMagicNumber        *bool          `json:"has_magic_number,omitempty"`
	HasPassword           bool           `json:"has_password"`
	HasStartedaTrial      *bool          `json:"has_started_trial,omitempty"`
	ID                    string         `json:"id"`
	ImageID               string         `json:"image_id"`
	InboxProjectID        string         `json:"inbox_project_id"`
	IsCelebrationsEnabled bool           `json:"is_celebrations_enabled"`
	IsDeleted             *bool          `json:"is_deleted,omitempty"`
	IsPremium             bool           `json:"is_premium"`
	JoinableWorkspace     *bool          `json:"joinable_workspace"`
	JoinedAt              string         `json:"joined_at"`
	Karma                 float32        `json:"karma"`
	KarmaDisabled         int            `json:"karma_disabled,omitempty"`
	KarmaTrend            string         `json:"karma_trend"`
	Lang                  string         `json:"lang"`
	MfaEnabled            bool           `json:"mfa_enabled"`
	NextWeek              int            `json:"next_week"`
	OnboardingCompleted   *bool          `json:"onboarding_completed,omitempty"`
	OnboardingInitiated   *bool          `json:"onboarding_initiated,omitempty"`
	OnboardingLevel       *string        `json:"onboarding_level,omitempty"`
	OnboardingPersona     *string        `json:"onboarding_persona,omitempty"`
	OnboardingStarted     *bool          `json:"onboarding_started,omitempty"`
	OnboardingTeamMode    *bool          `json:"onboarding_team_mode,omitempty"`
	OnboardingUseCases    []string       `json:"onboarding_use_cases,omitempty"`
	PremiumStatus         string         `json:"premium_status"`
	PremiumUntil          *string        `json:"premium_until"`
	ShardID               *int           `json:"shard_id"`
	ShareLimit            int            `json:"share_limit"`
	SortOrder             int            `json:"sort_order"`
	StartDay              int            `json:"start_day"`
	StartPage             string         `json:"start_page"`
	ThemeID               string         `json:"theme_id"`
	TimeFormat            int            `json:"time_format"`
	Token                 string         `json:"token"`
	TZInfo                map[string]any `json:"tz_info"`
	UniquePrefix          *int           `json:"unique_prefix,omitempty"`
	VerificationStatus    string         `json:"verification_status"`
	WebSocketURL          *string        `json:"web_socket_url,omitempty"`
	WeekendStartDay       int            `json:"weekend_start_day"`
	WeeklyGoal            int            `json:"weekly_goal"`
}

// The UserPlanLimts sync resource type describes the available features and
// limits applicable to the current user plan.
type UserPlanLimits struct {
	Current UserPlanInfo  `json:"current"`
	Next    *UserPlanInfo `json:"next,omitempty"` // This field is optional and will be null if there is no upgrade available for the user.
}

// The UserPlanInfo returned within the current property shows the values that
// are currently applied to the user.
type UserPlanInfo struct {
	ActivityLog              bool   `json:"activity_log"`
	ActivityLogLimit         int    `json:"activity_log_limit"`
	AdvancedPermissions      bool   `json:"advanced_permissions"`
	AutomaticBackups         bool   `json:"automatic_backups"`
	CalendarFeeds            bool   `json:"calendar_feeds"`
	CalendarLayout           bool   `json:"calendar_layout"`
	Comments                 bool   `json:"comments"`
	CompletedTasks           bool   `json:"completed_tasks"`
	CustomizationColor       bool   `json:"customization_color"`
	Deadlines                bool   `json:"deadlines"`
	Durations                bool   `json:"durations"`
	EmailAiForwarding        bool   `json:"email_ai_forwarding,omitempty"`
	EmailForwarding          bool   `json:"email_forwarding"`
	Filters                  bool   `json:"filters"`
	Labels                   bool   `json:"labels"`
	MaxCalendarAccounts      int    `json:"max_calendar_accounts"`
	MaxCollaborators         int    `json:"max_collaborators"`
	MaxFilters               int    `json:"max_filters"`
	MaxFoldersPerWorkspace   int    `json:"max_folders_per_workspace"`
	MaxFreeWorkspacesCreated int    `json:"max_free_workspaces_created"`
	MaxGuestsPerWorkspace    int    `json:"max_guests_per_workspace"`
	MaxLabels                int    `json:"max_labels"`
	MaxProjects              int    `json:"max_projects"`
	MaxProjectsJoined        int    `json:"max_projects_joined"`
	MaxRemindersLocation     int    `json:"max_reminders_location"`
	MaxRemindersTime         int    `json:"max_reminders_time"`
	MaxSections              int    `json:"max_sections"`
	MaxTasks                 int    `json:"max_tasks"`
	MaxUserTemplates         int    `json:"max_user_templates"`
	PlanName                 string `json:"plan_name"`
	Reminders                bool   `json:"reminders"`
	RemindersAtDue           bool   `json:"reminders_at_due"`
	Templates                bool   `json:"templates"`
	UploadLimitMB            int    `json:"upload_limit_mb"`
	Uploads                  bool   `json:"uploads"`
	WeeklyTrends             bool   `json:"weekly_trends"`
}
