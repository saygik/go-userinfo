package usecase

import (
	"fmt"
	"time"
)

func (u *UseCase) GetScheduleTasksNotifications() error {
	currentTime := time.Now().Add(time.Minute * 10)
	tasks, err := u.repo.GetScheduleTasksNotifications(currentTime.Format(time.RFC3339))
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MSSQL: %s", err.Error()))
	}
	_ = tasks
	for _, task := range tasks {
		if task.Tip == 3 && task.Status == 3 {
			u.repo.UpdateScheduleTaskNotification(task.Id)
			continue
		}
		// oldTime, _ := time.Parse(time.RFC3339, "1910-11-12T11:45:26.371Z")
		// taskTime, err := time.Parse(time.RFC3339, task.End)
		// if err != nil {
		// 	fmt.Println("Error parsing date:", err)
		// 	continue
		// }
		// if taskTime.Before(oldTime) || task.Tip == 1 {
		// 	taskTime, err = time.Parse(time.RFC3339, task.Start)
		// 	if err != nil {
		// 		fmt.Println("Error parsing date:", err)
		// 		continue
		// 	}
		// }
		taskTime, err := time.Parse(time.RFC3339, task.Start)
		if err != nil {
			fmt.Println("Ошибка определения даты события:", err)
			continue
		}
		taskTimeChek := taskTime
		if task.AllDay {
			taskTimeChek = taskTime.Add(8 * time.Hour)
		}
		tb := taskTimeChek.Before(currentTime)
		te := taskTimeChek.Equal(currentTime)
		_ = tb
		_ = te
		if taskTimeChek.Before(currentTime) || taskTimeChek.Equal(currentTime) {
			schedule, err := u.repo.GetSchedule(task.Idc)
			if err != nil {
				continue
			}
			title := ""
			if len(schedule.Mattermost) > 0 {
				switch task.Tip {
				case 1:
					title = "Напоминание календаря о графике пользователя " + task.Title
				case 2:
					title = "Напоминание календаря о событии: " + task.Title
				case 3:
					title = "Напоминание календаря о задаче: " + task.Title
				default:
					title = "Напоминание календаря о событии: " + task.Title
				}
				taskTimeFormatted := taskTime.Format(time.DateTime)
				if task.AllDay {
					taskTimeFormatted = taskTime.Format(time.DateOnly)
				}
				err = u.matt.SendPost(
					schedule.Mattermost,
					title,
					task.Comment,
					schedule.Name,
					fmt.Sprintf("https://userinfo.brnv.rw/schedules/schedule/%d", schedule.Id),
					"Дата события: "+taskTimeFormatted,
					false)

				if err != nil {
					u.log.Error(fmt.Sprintf("Error sending post for calendar task %d to Mattermost channel %s. Error: %v", task.Id, schedule.Mattermost, err))
				}
				u.repo.UpdateScheduleTaskNotification(task.Id)
			}

		}
	}
	return nil
}
