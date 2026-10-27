/*
    Переход от одного срока/оповещения в UserSoftwares к отдельной таблице
    полномочий. У пользователя может быть несколько полномочий на систему:
    основное и дополнительные, каждое со своим сроком и признаком оповещения.

    Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не создаст дубликаты.
*/

SET XACT_ABORT ON;
GO

IF OBJECT_ID(N'dbo.UserSoftwaresPermissions', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.UserSoftwaresPermissions
    (
        id                  INT IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_UserSoftwaresPermissions PRIMARY KEY,
        id_user_software    INT NOT NULL,
        name                NVARCHAR(100) NOT NULL,
        enddate             DATETIMEOFFSET(7) NULL,
        sended              BIT NOT NULL
            CONSTRAINT DF_UserSoftwaresPermissions_sended DEFAULT (0),
        CONSTRAINT FK_UserSoftwaresPermissions_UserSoftwares
            FOREIGN KEY (id_user_software)
            REFERENCES dbo.UserSoftwares (id)
            ON DELETE CASCADE
    );
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE object_id = OBJECT_ID(N'dbo.UserSoftwaresPermissions')
      AND name = N'IX_UserSoftwaresPermissions_User'
)
BEGIN
    CREATE INDEX IX_UserSoftwaresPermissions_User
        ON dbo.UserSoftwaresPermissions (id_user_software);
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE object_id = OBJECT_ID(N'dbo.UserSoftwaresPermissions')
      AND name = N'IX_UserSoftwaresPermissions_EOL'
)
BEGIN
    CREATE INDEX IX_UserSoftwaresPermissions_EOL
        ON dbo.UserSoftwaresPermissions (enddate, sended)
        INCLUDE (id_user_software)
        WHERE enddate IS NOT NULL;
END
GO

/*
    Переносим существующий единственный срок в основное полномочие.
    Если старые колонки уже удалены, блок пропускается.
*/
IF COL_LENGTH(N'dbo.UserSoftwares', N'enddate') IS NOT NULL
BEGIN
    INSERT INTO dbo.UserSoftwaresPermissions (id_user_software, name, enddate, sended)
    SELECT u.id,
           N'Основные полномочия',
           u.enddate,
           CASE
               WHEN u.enddate IS NULL
                    OR u.enddate <= CONVERT(datetimeoffset(7), '1900-01-02T00:00:00+00:00')
                   THEN 0
               ELSE ISNULL(u.sended, 0)
           END
    FROM dbo.UserSoftwares u
    WHERE NOT EXISTS (
        SELECT 1
        FROM dbo.UserSoftwaresPermissions p
        WHERE p.id_user_software = u.id
    );
END
GO

/*
    Старые колонки больше не используются. Если в вашей базе есть внешние
    отчёты/процедуры, которые читают UserSoftwares.enddate или sended,
    их нужно перевести на UserSoftwaresPermissions до выполнения этого блока.
*/
IF COL_LENGTH(N'dbo.UserSoftwares', N'enddate') IS NOT NULL
BEGIN
    ALTER TABLE dbo.UserSoftwares DROP COLUMN enddate;
END
GO

IF COL_LENGTH(N'dbo.UserSoftwares', N'sended') IS NOT NULL
BEGIN
    ALTER TABLE dbo.UserSoftwares DROP COLUMN sended;
END
GO
