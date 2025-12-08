# go-userinfo

Сервис для агрегирования информации о пользователях и их ресурсах из нескольких систем: Active Directory, GLPI, Mattermost, внутренних расписаний и инвентаризации ПО. Приложение предоставляет HTTP API с аутентификацией через Authentik (OAuth2), обрабатывает интеграционные вебхуки и периодически синхронизирует кэши в Redis.

## Основные возможности
- OAuth2 через Authentik: получение и обновление access/refresh токенов, logout.
- Работа с AD: поиск пользователей и компьютеров, получение активности, делегирования почтовых ящиков, управление группами и ролями приложения.
- Доступы в приложении: перечень ролей, групп, ресурсов, информация о текущем пользователе.
- GLPI: чтение заявок/проблем/отказов/отчетов, статистика, добавление комментариев, решений, пользователей и создание заявок.
- Software инвентаризация: учет систем, пользователей, выдача/изменение/удаление доступов.
- Расписание: управление задачами календарей и выдача расписаний.
- Mattermost интеграции: комментарии в GLPI, отключение уведомлений о задачах календаря, slash-команды.
- IUTM: выдача whitelist-данных.
- Служебные функции: статические страницы, изображения для статусов заявок, метрики Prometheus на `:9100`.

## Конфигурация
- Конфиг читается из `config.json` в корне (см. пример в репозитории). В нем задаются:
  - `app`: окружение, порт HTTP (`9099` по умолчанию), период фоновой проверки ПО.
  - `vault`: адрес Vault и данные approle (`roleid`, `secretid`, `secretpath`) для загрузки секретов.
  ```{
    "app": {
        "env": "local",
        "port": "9099",
        "softwarebottime": 2
    },
    "vault": {
        "server": "https://vault.domain.com:8200/",
        "roleid": "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
        "secretid": "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
        "secretpath": "/secret/data/test/"
    }}
    ```

  - Остальные секции (AD, репозитории, интеграции, Authentik) подтягиваются из Vault по путям, указанным в `secretpath`.
## Секция ad
```{"domains": [
    {
      "adminGLPIGroup": 21,
      "base": "DC=domain,DC=com",
      "bindDN": "CN=read-only-admin,OU=Служебные записи,DC=xxx,DC=rw",
      "bindPassword": "password",
      "computerFilter": "(&(objectCategory=computer)(!(userAccountControl:1.2.840.113556.1.4.803:=2))(!(memberof=CN=hidden computers,OU=_Groups,DC=domain,DC=com)))",
      "dc": "dc.domain.com",
      "domain": "domain",
      "filter": "(&(&(&(objectClass=user)(objectCategory=person)(!(memberof=CN=queryFiltred,OU=_Services,DC=brnv,DC=rw))(!(memberof=CN=services,OU=_Services,DC=domain,DC=com)))))",
      "groupFilter": "(&(&(&(objectClass=user)(objectCategory=person)((memberof=%s))(!(memberof=CN=queryFiltred,OU=_Services,DC=domain,DC=com))(!(userAccountControl:1.2.840.113556.1.4.803:=2)))))",
      "internetGroups": {
        "full": "Интернет Полный",
        "tech": "Интернет Технологический",
        "whitelist": "Интернет Белый список;Интернет ИВЦ2"
      },
      "key": "domain.com",
      "name": "Домен Предприятия",
      "rmsPort": 25650,
      "time": 5
    }
  ]
}
```
## Секция authentik
```
{
  "client-id": "id",
  "client-secret": "secret",
  "log-out-url": "https://sso.domain.com/application/o/userinfo/end-session/",
  "redirect-url": "https://userinfo.domain.com/authentik/callback",
  "scopes": [
    "openid",
    "profile",
    "email",
    "offline",
    "offline_access"
  ],
  "url": "https://sso.domain.rw/application/o/userinfo/"
}
```
## Секция integrations
```
{
  "add-comment-from-api": "https://userinfoapi.domain.com/api/mattermost/glpi/comment",
  "allowed-api": [
    "X-Api-Key"
  ],
  "allowed-hosts": [
    "ip1",
    "ip2"
  ],
  "disable-calendar-task-notification-api": "url",
  "n8n-webhook-ivc2-kaspersky": "https://n8n.domain.com/webhook/xxx"
}
```
## Секция repository
```
{
  "glpi": {
    "dbname": "dbname",
    "password": "password",
    "server": "xxx.xxx.xxx.xxx",
    "user": "dbuser"
  },
  "glpiapi": {
    "server": "https://domain.com/apirest.php/",
    "token": "apptoken",
    "usertoken": "usertoken"
  },
  "iutm": {
    "server": "https://domain.com:8443",
    "token": "user",
    "usertoken": "password"
  },
  "mattermost": {
    "server": "https://domain.com/api/v4",
    "token": "token"
  },
  "mssql": {
    "dbname": "dbname",
    "password": "password",
    "server": "servername",
    "user": "user"
  },
  "redis": {
    "password": "",
    "secret": "secret",
    "server": "127.0.0.1:6379"
  }
}
```

- Требуемая версия Go: `1.24`.

## Запуск локально
1) Убедитесь, что `config.json` заполнен корректно и у процесса есть доступ к Vault.
2) Выполните `go run .` (или соберите бинарь `go build .` и запустите его).
3) HTTP сервер поднимется на порту из `config.app.port`, метрики — на `:9100`.

## Маршруты
Все API имеют префикс `/api`. В столбце «Авторизация» указано, требуется ли токен (`TokenAuthMiddleware` или извлечение пользователя из токена).

| Секция | Метод | Путь | Авторизация | Описание |
| --- | --- | --- | --- | --- |
| UI | GET | `/` | нет | Главная HTML-страница приложения. |
| UI | GET | `/f` | нет | Страница входа. |
| OAuth2 | GET | `/api/oauth-authentik/login` | нет | Инициировать Authentik login (redirect URL). |
| OAuth2 | POST | `/api/oauth-authentik/logout` | нет | Завершить сессию в Authentik. |
| OAuth2 | GET | `/api/oauth-authentik/token` | нет | Обмен кода на токен. |
| OAuth2 | POST | `/api/oauth-authentik/refresh` | нет | Обновить access токен по refresh. |
| AD | GET | `/api/ad/users` | да | Полный список пользователей домена. |
| AD | GET | `/api/ad/public/users` | да | Сокращенная информация по пользователям. |
| AD | GET | `/api/ad/user/:username` | да | Подробные свойства пользователя. |
| AD | GET | `/api/ad/finduser/:username` | нет | Поиск пользователя (упрощенный профиль). |
| AD | GET | `/api/ad/computers` | да | Список доменных компьютеров. |
| AD | GET | `/api/ad/stats/counts` | да | Базовая статистика по доменам. |
| AD | GET | `/api/ad/activity/user/:username` | да | История активности пользователя в AD. |
| AD | GET | `/api/ad/user-mailbox-delegates/:username` | да | Делегированные почтовые ящики пользователя. |
| AD | PUT | `/api/ad/user/avatar/:username` | да | Установить аватар пользователя. |
| AD | PUT | `/api/ad/user/role/:username` | да | Изменить роль приложения. |
| AD | POST | `/api/ad/user/group/:username` | да | Добавить группу приложения пользователю. |
| AD | DELETE | `/api/ad/user/group/:username` | да | Удалить группу приложения у пользователя. |
| AD | POST | `/api/ad/user/role/:username` | да | Добавить роль приложения пользователю. |
| AD | DELETE | `/api/ad/user/role/:username` | да | Удалить роль приложения у пользователя. |
| AD | GET | `/api/ad/groupusers/:domain/:group` | да | Пользователи указанной группы домена. |
| App | GET | `/api/app/whoami` | да | Текущий пользователь (whoami). |
| App | GET | `/api/app/userresources` | да | Ресурсы, разрешенные текущему пользователю. |
| App | GET | `/api/app/resources` | да | Все ресурсы приложения. |
| App | GET | `/api/app/roles` | да | Все роли приложения. |
| App | GET | `/api/app/groups` | да | Все группы приложения. |
| App | GET | `/api/app/domains` | да | Перечень доменов. |
| App | GET | `/api/app/setip` | нет | Сохранить IP пользователя. |
| App | GET | `/api/app/ip` | нет | Получить IP пользователя. |
| GLPI | GET | `/api/glpi/whoami` | да | Текущий пользователь GLPI. |
| GLPI | GET | `/api/glpi/user/:username` | да | Профиль пользователя GLPI. |
| GLPI | GET | `/api/glpi/nctickets` | да | Нерешенные заявки GLPI. |
| GLPI | GET | `/api/glpi/tickets/mygroups` | да | Незакрытые заявки в группах слежения пользователя. |
| GLPI | GET | `/api/glpi/ticket/:id` | да | Получить заявку. |
| GLPI | GET | `/api/glpi/reports/ticket/:id` | да (по токену пользователя) | Получить заявку для отчетов. |
| GLPI | GET | `/api/glpi/ticket/solutions/:id` | да | Шаблоны решений заявки. |
| GLPI | POST | `/api/glpi/ticket/user/:id` | да | Добавить пользователя к заявке. |
| GLPI | POST | `/api/glpi/ticket` | да | Создать заявку. |
| GLPI | POST | `/api/glpi/comment/ticket/:id` | да | Добавить комментарий к заявке. |
| GLPI | POST | `/api/glpi/solution/ticket/:id` | да | Добавить решение заявки. |
| GLPI | GET | `/api/glpi/problem/:id` | да | Получить проблему. |
| GLPI | GET | `/api/glpi/users` | да | Список пользователей GLPI. |
| GLPI | GET | `/api/glpi/otkazes` | да | Отказы за период. |
| GLPI | GET | `/api/glpi/problems` | да | Проблемы за период. |
| GLPI | GET | `/api/glpi/statistics/tickets` | да | Статистика по заявкам. |
| GLPI | GET | `/api/glpi/statistics/failures` | да | Статистика по отказам. |
| GLPI | GET | `/api/glpi/statistics/period-regions-month-days` | да | Статистика по отказам по регионам/дням. |
| GLPI | GET | `/api/glpi/statistics/statsdays` | да | Статистика заявок по дням. |
| GLPI | GET | `/api/glpi/statistics/top10performers` | да | ТОП-10 исполнителей. |
| GLPI | GET | `/api/glpi/statistics/top10iniciators` | да | ТОП-10 инициаторов. |
| GLPI | GET | `/api/glpi/statistics/top10groups` | да | ТОП-10 групп. |
| GLPI | GET | `/api/glpi/statistics/periodcounts` | да | Количество заявок за период. |
| GLPI | GET | `/api/glpi/statistics/periodrequestypes` | да | ТОП типов запросов за период. |
| GLPI | GET | `/api/glpi/statistics/regions` | да | Статистика по регионам. |
| GLPI | GET | `/api/glpi/statistics/period-org-treemap` | да | Дерево организаций за период. |
| GLPI | GET | `/api/glpi/hrp` | нет | Получение заявок HRP (фоновая интеграция). |
| Img | GET | `/api/img/ticket-status/:id` | нет | Изображение/бейдж статуса заявки. |
| IUTM | GET | `/api/iutm/wlist` | нет | Whitelist-данные IUTM. |
| Manual | GET | `/api/manual/orgcodes` | да | Коды организаций для мануалов. |
| Mattermost | GET | `/api/mattermost/users` | да | Пользователи Mattermost (с сессиями). |
| Mattermost | POST | `/api/mattermost/glpi/comment` | нет | Добавить комментарий HRP/GLPI из Mattermost. |
| Mattermost | POST | `/api/mattermost/schedule/notification` | нет | Отключить уведомление о задаче календаря. |
| Mattermost Commands | POST | `/api/mattermost-commands/glpi` | нет | Slash-команда Mattermost для GLPI. |
| Schedule | GET | `/api/schedule/one/:id` | да (по токену пользователя) | Получить один календарь. |
| Schedule | GET | `/api/schedule/all` | да (по токену пользователя) | Все доступные календари. |
| Schedule | POST | `/api/schedule/task` | да | Добавить задачу календаря. |
| Schedule | DELETE | `/api/schedule/task/:id` | да | Удалить задачу календаря. |
| Schedule | PUT | `/api/schedule/task/:id` | да | Обновить задачу календаря. |
| Software | GET | `/api/software/one/:id` | да | Получить систему. |
| Software | POST | `/api/software/one/:id` | да | Добавить пользователя в систему. |
| Software | PUT | `/api/software/one/:id` | да | Обновить пользователя в системе. |
| Software | GET | `/api/software/one/:id/users` | да | Пользователи системы. |
| Software | GET | `/api/software/user/:username` | да | Все системы пользователя. |
| Software | POST | `/api/software/user/:username` | да | Добавить систему пользователю. |
| Software | DELETE | `/api/software/user/:id` | да | Удалить систему пользователя. |
| Software | GET | `/api/software/all` | да | Все системы. |
| Software | GET | `/api/software/all/users` | да | Все пользователи по системам. |
