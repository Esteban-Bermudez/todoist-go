package todoist

type Workspace struct {
	ID                    string         `json:"id"`
	Name                  string         `json:"name"`
	Description           string         `json:"description"`
	Plan                  string         `json:"plan"`
	IsLinkSharingEnabled  bool           `json:"is_link_sharing_enabled"`
	IsGuestAllowed        bool           `json:"is_guest_allowed"`
	InviteCode            *string        `json:"invite_code,omitempty"`
	Role                  string         `json:"role"` // ADMIN, MEMBER, GUEST
	LogoBig               string         `json:"logo_big"`
	LogoMedium            string         `json:"logo_medium"`
	LogoSmall             string         `json:"logo_small"`
	LogoS640              string         `json:"logo_s640"`
	Limits                *any           `json:"limits"` // This can be a complex object, so using `any` for flexibility
	CreatorID             string         `json:"creator_id"`
	CreatedAt             string         `json:"created_at"`
	IsDeleted             bool           `json:"is_deleted"`
	IsCollapsed           bool           `json:"is_collapsed"`
	DomainName            *string        `json:"domain_name,omitempty"`
	DomainDiscovery       *bool          `json:"domain_discovery,omitempty"`
	RestrictEmailDomains  *bool          `json:"restrict_email_domains,omitempty"`
	PendingInvitations    []string       `json:"pending_invitations,omitempty"`
	PendingInvitesByType  map[string]int `json:"pending_invites_by_type,omitempty"`
	MemberCountByType     map[string]int `json:"member_count_by_type,omitempty"`
	CurrentActiveProjects int            `json:"current_active_projects,omitempty"`
	CurrentMemberCount    int            `json:"current_member_count,omitempty"`
	CurrentTemplateCount  int            `json:"current_template_count,omitempty"`
	Properties            *any           `json:"properties,omitempty"` // Flexible properties, can be any type
}

// WorkspaceUsers are not returned in full sync responses, only in incremental
// sync. To keep a list of workspace users up-to-date, clients should first list
// all workspace users, then use incremental
// sync to update that initial list as needed.
// WorkspaceUsers are not the same as collaborators. Two users can be members of
// a common workspace without having a common shared project, so they will both
// “see” each other in workspace_users but not in collaborators.
// Guests will not receive WorkspaceUsers sync events or resources.
type WorkspaceUser struct {
	UserID       string  `json:"user_id"`
	WorkspaceID  string  `json:"workspace_id"`
	UserEmail    string  `json:"user_email"`
	FullName     string  `json:"full_name"`
	Timezone     string  `json:"timezone"`
	AvatarBig    string  `json:"avatar_big,omitempty"`
	AvatarMedium string  `json:"avatar_medium,omitempty"`
	AvatarS640   string  `json:"avatar_s640,omitempty"`
	AvatarSmall  string  `json:"avatar_small,omitempty"`
	ImageID      *string `json:"image_id"`
	Role         string  `json:"role"` // ADMIN, MEMBER, GUEST
	IsDeleted    bool    `json:"is_deleted"`
}
