/*
    Связь задач календаря с сообщениями Mattermost.

    Позволяет обновлять исходный пост, когда по задаче срабатывает
    напоминание или когда задача закрывается.

    Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не изменит схему.
*/

SET XACT_ABORT ON;
GO

IF OBJECT_ID(N'dbo.ScheduleTaskMattermostLinks', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.ScheduleTaskMattermostLinks
    (
        task_id          INT NOT NULL
            CONSTRAINT PK_ScheduleTaskMattermostLinks PRIMARY KEY,
        channel_id       NVARCHAR(64) NOT NULL,
        post_id          NVARCHAR(64) NOT NULL,
        source_type      NVARCHAR(20) NOT NULL,
        source_id        BIGINT NULL,
        reminder_sent_at DATETIMEOFFSET(7) NULL,
        closed_at        DATETIMEOFFSET(7) NULL,
        closed_by        NVARCHAR(256) NULL,
        close_reason     NVARCHAR(500) NULL,
        created_at       DATETIMEOFFSET(7) NOT NULL
            CONSTRAINT DF_ScheduleTaskMattermostLinks_created_at DEFAULT SYSDATETIMEOFFSET()
    );
END
GO

IF OBJECT_ID(N'dbo.ScheduleTaskMattermostLinks', N'U') IS NOT NULL
   AND NOT EXISTS (
       SELECT 1 FROM sys.indexes
       WHERE name = N'IX_ScheduleTaskMattermostLinks_post'
         AND object_id = OBJECT_ID(N'dbo.ScheduleTaskMattermostLinks')
   )
BEGIN
    CREATE INDEX IX_ScheduleTaskMattermostLinks_post
        ON dbo.ScheduleTaskMattermostLinks (post_id);
END
GO
