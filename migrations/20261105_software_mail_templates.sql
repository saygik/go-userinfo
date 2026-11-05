-- Шаблоны писем об истечении полномочий для отдельных систем.
-- Если строки для системы нет, приложение использует общий шаблон по умолчанию.

IF OBJECT_ID(N'dbo.SoftwareMailTemplates', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.SoftwareMailTemplates
    (
        id_soft int NOT NULL
            CONSTRAINT PK_SoftwareMailTemplates PRIMARY KEY,
        subject nvarchar(300) NOT NULL,
        body nvarchar(max) NOT NULL,
        updated_at datetimeoffset(7) NOT NULL
            CONSTRAINT DF_SoftwareMailTemplates_updated_at DEFAULT SYSDATETIMEOFFSET(),
        updated_by nvarchar(256) NOT NULL
            CONSTRAINT DF_SoftwareMailTemplates_updated_by DEFAULT (N''),
        CONSTRAINT CK_SoftwareMailTemplates_IdSoft CHECK (id_soft > 0)
    );
END;
