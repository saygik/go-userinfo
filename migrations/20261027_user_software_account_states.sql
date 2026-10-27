/*
    Состояние учётной записи в конкретной системе.
    Не путать с полномочиями (UserSoftwaresPermissions) и со статусом
    учётной записи в AD: здесь хранится именно состояние доступа в системе.

    Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не создаст дубликаты.
*/

SET XACT_ABORT ON;
GO

IF OBJECT_ID(N'dbo.UserSoftwaresAccountStates', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.UserSoftwaresAccountStates
    (
        id                  INT IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_UserSoftwaresAccountStates PRIMARY KEY,
        id_user_software    INT NOT NULL,
        status              NVARCHAR(20) NOT NULL
            CONSTRAINT DF_UserSoftwaresAccountStates_status DEFAULT (N'active'),
        status_changed_at   DATETIMEOFFSET(7) NULL,
        status_changed_by   NVARCHAR(100) NULL,
        status_reason       NVARCHAR(500) NULL,
        source              NVARCHAR(50) NOT NULL
            CONSTRAINT DF_UserSoftwaresAccountStates_source DEFAULT (N'manual'),
        last_sync_at        DATETIMEOFFSET(7) NULL,

        CONSTRAINT UQ_UserSoftwaresAccountStates_User
            UNIQUE (id_user_software),

        CONSTRAINT FK_UserSoftwaresAccountStates_UserSoftwares
            FOREIGN KEY (id_user_software)
            REFERENCES dbo.UserSoftwares (id)
            ON DELETE CASCADE,

        CONSTRAINT CK_UserSoftwaresAccountStates_Status
            CHECK (status IN (N'active', N'disabled', N'blocked', N'deleted'))
    );
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE object_id = OBJECT_ID(N'dbo.UserSoftwaresAccountStates')
      AND name = N'IX_UserSoftwaresAccountStates_Status'
)
BEGIN
    CREATE INDEX IX_UserSoftwaresAccountStates_Status
        ON dbo.UserSoftwaresAccountStates (status);
END
GO

/*
    Для всех существующих пользователей систем создаём состояние "активна".
    Если запись уже есть, она не меняется.
*/
INSERT INTO dbo.UserSoftwaresAccountStates
    (id_user_software, status, status_changed_at, status_changed_by, status_reason, source)
SELECT u.id,
       N'active',
       NULL,
       NULL,
       NULL,
       N'migration'
FROM dbo.UserSoftwares u
WHERE NOT EXISTS (
    SELECT 1
    FROM dbo.UserSoftwaresAccountStates s
    WHERE s.id_user_software = u.id
);
GO
