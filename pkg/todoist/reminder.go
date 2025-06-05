package todoist

type Reminder struct {
	ID        string `json:"id"`
	NotifyUID string `json:"notify_uid"`
	ItemID    string `json:"item_id"`
	Type      string `json:"type"`
	Due       struct {
		Date        string  `json:"date"`
		Timezone    *string `json:"timezone"`
		IsRecurring bool    `json:"is_recurring"`
		String      string  `json:"string"`
		Lang        string  `json:"lang"`
	} `json:"due"`
	MinuteOffset int     `json:"minute_offset"`
	Name         *string `json:"name,omitempty"`
	LocLat       *string `json:"loc_lat,omitempty"`
	LocLong      *string `json:"loc_long,omitempty"`
	LocTrigger   *string `json:"loc_trigger,omitempty"`
	Radius       *int    `json:"radius,omitempty"`
	IsDeleted    bool    `json:"is_deleted"`
}
