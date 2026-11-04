package entity

import "time"

// ScheduleTaskMattermostLink связывает задачу календаря с сообщением Mattermost.
type ScheduleTaskMattermostLink struct {
	TaskID         int        `db:"task_id"`
	ChannelID      string     `db:"channel_id"`
	PostID         string     `db:"post_id"`
	SourceType     string     `db:"source_type"`
	SourceID       int64      `db:"source_id"`
	ReminderSentAt *time.Time `db:"reminder_sent_at"`
	ClosedAt       *time.Time `db:"closed_at"`
	ClosedBy       string     `db:"closed_by"`
	CloseReason    string     `db:"close_reason"`
}
