/*
    Нормализация прав доступа к календарям.

    Семантика, которую использует приложение:
      UserScheduleUsers.tip
        1 — полный доступ к календарю;
        2 — просмотр и участие в графике дежурств.
      UserScheduleUserGroups.role
        1 — полный доступ к календарю;
        2 — просмотр и участие в графике дежурств.

    Скрипт убирает дубли прав (приоритет у полного доступа), приводит
    недопустимые значения к просмотру и добавляет индексы для выборок
    по календарю. Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не изменит данные и схему.
*/

SET XACT_ABORT ON;
GO

BEGIN TRANSACTION;

-- Недопустимые значения трактуем как просмотр: это безопаснее, чем
-- случайно выдать полный доступ.
UPDATE dbo.UserScheduleUsers
SET tip = 2
WHERE tip NOT IN (1, 2);

UPDATE dbo.UserScheduleUserGroups
SET role = 2
WHERE role NOT IN (1, 2);

-- Один пользователь может иметь только одно право на календарь.
-- При дублях оставляем строку с полным доступом (tip = 1).
;WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY id_shedule, name
               ORDER BY tip ASC, id ASC
           ) AS rn
    FROM dbo.UserScheduleUsers
)
DELETE FROM duplicates
WHERE rn > 1;

-- Одна группа GLPI может иметь только одно право на календарь.
-- При дублях оставляем строку с полным доступом (role = 1).
;WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY id_shedule, tip, name
               ORDER BY role ASC, id ASC
           ) AS rn
    FROM dbo.UserScheduleUserGroups
)
DELETE FROM duplicates
WHERE rn > 1;

COMMIT TRANSACTION;
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.check_constraints
    WHERE name = N'CK_UserScheduleUsers_tip'
      AND parent_object_id = OBJECT_ID(N'dbo.UserScheduleUsers')
)
BEGIN
    ALTER TABLE dbo.UserScheduleUsers
        ADD CONSTRAINT CK_UserScheduleUsers_tip CHECK (tip IN (1, 2));
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.check_constraints
    WHERE name = N'CK_UserScheduleUserGroups_role'
      AND parent_object_id = OBJECT_ID(N'dbo.UserScheduleUserGroups')
)
BEGIN
    ALTER TABLE dbo.UserScheduleUserGroups
        ADD CONSTRAINT CK_UserScheduleUserGroups_role CHECK (role IN (1, 2));
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'IX_UserSchedules_id'
      AND object_id = OBJECT_ID(N'dbo.UserSchedules')
)
BEGIN
    CREATE INDEX IX_UserSchedules_id
        ON dbo.UserSchedules (id);
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'IX_UserScheduleUsers_schedule_name'
      AND object_id = OBJECT_ID(N'dbo.UserScheduleUsers')
)
BEGIN
    CREATE INDEX IX_UserScheduleUsers_schedule_name
        ON dbo.UserScheduleUsers (id_shedule, name, tip);
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'IX_UserScheduleUserGroups_schedule_group'
      AND object_id = OBJECT_ID(N'dbo.UserScheduleUserGroups')
)
BEGIN
    CREATE INDEX IX_UserScheduleUserGroups_schedule_group
        ON dbo.UserScheduleUserGroups (id_shedule, tip, name, role);
END
GO
