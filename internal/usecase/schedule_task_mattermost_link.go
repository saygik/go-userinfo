package usecase

import "github.com/saygik/go-userinfo/internal/entity"

const (
	schedulePostStatusReminderSent       = "Напоминание отправлено"
	schedulePostStatusClosedInCalendar   = "Задача выполнена в календаре"
	schedulePostStatusClosedInMattermost = "Задача отменена в Mattermost"
	schedulePostColorClosed              = "#808080"
	scheduleLinkSourceHRP                = "hrp"
)

// saveScheduleTaskMattermostLink сохраняет связь задачи календаря с постом.
func (u *UseCase) saveScheduleTaskMattermostLink(taskID int, channelID, postID, sourceType string, sourceID int64) {
	if taskID <= 0 || channelID == "" || postID == "" {
		return
	}
	link := entity.ScheduleTaskMattermostLink{
		TaskID:     taskID,
		ChannelID:  channelID,
		PostID:     postID,
		SourceType: sourceType,
		SourceID:   sourceID,
	}
	if err := u.repo.SaveScheduleTaskMattermostLink(link); err != nil && u.log != nil {
		u.log.WithError(err).Warn("не удалось сохранить связь задачи календаря с сообщением Mattermost")
	}
}

// getScheduleTaskPostLink возвращает связь задачи с постом, если она есть.
func (u *UseCase) getScheduleTaskPostLink(taskID int) (entity.ScheduleTaskMattermostLink, bool) {
	if taskID <= 0 {
		return entity.ScheduleTaskMattermostLink{}, false
	}
	link, err := u.repo.GetScheduleTaskMattermostLink(int64(taskID))
	if err != nil {
		if u.log != nil {
			u.log.WithError(err).Warn("не удалось получить связь задачи календаря с сообщением Mattermost")
		}
		return entity.ScheduleTaskMattermostLink{}, false
	}
	if link.PostID == "" {
		return entity.ScheduleTaskMattermostLink{}, false
	}
	return link, true
}

// updateScheduleTaskLinkPostStatus добавляет статус к исходному посту.
func (u *UseCase) updateScheduleTaskLinkPostStatus(link entity.ScheduleTaskMattermostLink, status, color string) {
	if link.PostID == "" || (status == "" && color == "") {
		return
	}
	if err := u.matt.UpdatePostStatus(link.PostID, status, color); err != nil && u.log != nil {
		u.log.WithError(err).Warn("не удалось обновить сообщение Mattermost по задаче календаря")
	}
}

// updateScheduleTaskPostStatus добавляет статус к исходному посту Mattermost.
func (u *UseCase) updateScheduleTaskPostStatus(taskID int, status string) {
	link, ok := u.getScheduleTaskPostLink(taskID)
	if !ok {
		return
	}
	u.updateScheduleTaskLinkPostStatus(link, status, "")
}

// markScheduleTaskPostReminderSent обновляет исходный пост и фиксирует
// отправку напоминания по задаче.
func (u *UseCase) markScheduleTaskPostReminderSent(taskID int) {
	u.updateScheduleTaskPostStatus(taskID, schedulePostStatusReminderSent)
	if err := u.repo.MarkScheduleTaskMattermostReminderSent(int64(taskID)); err != nil && u.log != nil {
		u.log.WithError(err).Warn("не удалось отметить отправку напоминания по задаче календаря")
	}
}

// markScheduleTaskPostClosed обновляет исходный пост и фиксирует закрытие задачи.
// Повторное закрытие не меняет уже выставленный статус.
func (u *UseCase) markScheduleTaskPostClosed(taskID int, actor, reason, status string) {
	link, ok := u.getScheduleTaskPostLink(taskID)
	if !ok || link.ClosedAt != nil {
		return
	}
	u.updateScheduleTaskLinkPostStatus(link, status, schedulePostColorClosed)
	if err := u.repo.MarkScheduleTaskMattermostClosed(int64(taskID), actor, reason); err != nil && u.log != nil {
		u.log.WithError(err).Warn("не удалось отметить закрытие задачи календаря")
	}
}
