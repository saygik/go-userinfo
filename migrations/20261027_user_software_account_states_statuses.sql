/*
    Расширение статусов учётной записи в системе:
    active / disabled / blocked / deleted.

    Нужен для баз, где миграция
    20261027_user_software_account_states.sql уже применена
    и CHECK-ограничение допускает только active/disabled.

    Перед выполнением сделайте резервную копию базы.
*/

SET XACT_ABORT ON;
GO

IF OBJECT_ID(N'dbo.UserSoftwaresAccountStates', N'U') IS NOT NULL
BEGIN
    IF EXISTS (
        SELECT 1
        FROM sys.check_constraints
        WHERE name = N'CK_UserSoftwaresAccountStates_Status'
          AND parent_object_id = OBJECT_ID(N'dbo.UserSoftwaresAccountStates')
    )
    BEGIN
        ALTER TABLE dbo.UserSoftwaresAccountStates
            DROP CONSTRAINT CK_UserSoftwaresAccountStates_Status;
    END

    ALTER TABLE dbo.UserSoftwaresAccountStates
        ADD CONSTRAINT CK_UserSoftwaresAccountStates_Status
        CHECK (status IN (N'active', N'disabled', N'blocked', N'deleted'));
END
GO
