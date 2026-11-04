/*
    Аудит изменений прав пользователей приложения.

    Фиксирует, кто и когда изменил роли, доступы к доменам или разделы
    пользователя, а также полную очистку прав. Записи добавляются только
    после успешного изменения; аудит не влияет на саму операцию.

    Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не изменит схему.
*/

SET XACT_ABORT ON;
GO

IF OBJECT_ID(N'dbo.PermissionAudit', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.PermissionAudit
    (
        id          BIGINT IDENTITY(1, 1) NOT NULL
            CONSTRAINT PK_PermissionAudit PRIMARY KEY,
        changed_at  DATETIME2(0) NOT NULL
            CONSTRAINT DF_PermissionAudit_changed_at DEFAULT SYSUTCDATETIME(),
        actor_upn   NVARCHAR(256) NOT NULL,
        target_upn  NVARCHAR(256) NOT NULL,
        action      NVARCHAR(50) NOT NULL,
        object_name NVARCHAR(256) NULL,
        details     NVARCHAR(500) NULL
    );
END
GO

IF OBJECT_ID(N'dbo.PermissionAudit', N'U') IS NOT NULL
   AND NOT EXISTS (
       SELECT 1 FROM sys.indexes
       WHERE name = N'IX_PermissionAudit_changed_at'
         AND object_id = OBJECT_ID(N'dbo.PermissionAudit')
   )
BEGIN
    CREATE INDEX IX_PermissionAudit_changed_at
        ON dbo.PermissionAudit (changed_at DESC, id DESC);
END
GO

IF OBJECT_ID(N'dbo.PermissionAudit', N'U') IS NOT NULL
   AND NOT EXISTS (
       SELECT 1 FROM sys.indexes
       WHERE name = N'IX_PermissionAudit_target'
         AND object_id = OBJECT_ID(N'dbo.PermissionAudit')
   )
BEGIN
    CREATE INDEX IX_PermissionAudit_target
        ON dbo.PermissionAudit (target_upn, changed_at DESC);
END
GO
