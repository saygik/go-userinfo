/*
    Необязательный ID заявки GLPI, по которой пользователю системы
    предоставлен доступ. Хранится только номер заявки; ссылка
    формируется на клиенте.

    Перед выполнением сделайте резервную копию базы.
    Скрипт идемпотентен: повторный запуск не изменит схему.
*/

SET XACT_ABORT ON;
GO

IF COL_LENGTH(N'dbo.UserSoftwares', N'glpi_ticket_id') IS NULL
BEGIN
    ALTER TABLE dbo.UserSoftwares
        ADD glpi_ticket_id NVARCHAR(50) NULL;
END
GO
