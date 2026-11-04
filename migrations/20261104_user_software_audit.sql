/*
    Аудит изменений записей «пользователь — система» (dbo.UserSoftwares).

    Фиксирует создание, изменение и удаление записей, включая сведения о
    пользователе и системе. Записи добавляются только после успешной
    операции; аудит не влияет на саму операцию.

    Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не изменит схему.
*/

SET XACT_ABORT ON;
GO

IF OBJECT_ID(N'dbo.UserSoftwareAudit', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.UserSoftwareAudit
    (
        id               BIGINT IDENTITY(1, 1) NOT NULL
            CONSTRAINT PK_UserSoftwareAudit PRIMARY KEY,
        changed_at       DATETIME2(0) NOT NULL
            CONSTRAINT DF_UserSoftwareAudit_changed_at DEFAULT SYSUTCDATETIME(),
        actor_upn        NVARCHAR(256) NOT NULL,
        action           NVARCHAR(20) NOT NULL,
        user_software_id BIGINT NULL,
        id_soft          INT NULL,
        soft_name        NVARCHAR(256) NULL,
        user_name        NVARCHAR(256) NULL,
        user_login       NVARCHAR(256) NULL,
        fio              NVARCHAR(256) NULL,
        details          NVARCHAR(500) NULL
    );
END
GO

IF OBJECT_ID(N'dbo.UserSoftwareAudit', N'U') IS NOT NULL
   AND NOT EXISTS (
       SELECT 1 FROM sys.indexes
       WHERE name = N'IX_UserSoftwareAudit_changed_at'
         AND object_id = OBJECT_ID(N'dbo.UserSoftwareAudit')
   )
BEGIN
    CREATE INDEX IX_UserSoftwareAudit_changed_at
        ON dbo.UserSoftwareAudit (changed_at DESC, id DESC);
END
GO

IF OBJECT_ID(N'dbo.UserSoftwareAudit', N'U') IS NOT NULL
   AND NOT EXISTS (
       SELECT 1 FROM sys.indexes
       WHERE name = N'IX_UserSoftwareAudit_user'
         AND object_id = OBJECT_ID(N'dbo.UserSoftwareAudit')
   )
BEGIN
    CREATE INDEX IX_UserSoftwareAudit_user
        ON dbo.UserSoftwareAudit (user_name, changed_at DESC);
END
GO
