# Grafana dashboards

Дашборды для Prometheus-метрик UserInfo (`:9100/metrics`).

| Файл | UID | Назначение |
|---|---|---|
| `ad-users-account-state.json` | `ad-users-account-state` | состояние учётных записей AD: активные, отключённые, привилегированные, служебная почта, ошибки атрибутов, пароли, компьютеры, интернет |
| `userinfo-service-health.json` | `userinfo-service-health` | общее состояние сервиса: доступность, импорт AD, свежесть кеша, HTTP, кеш списка пользователей, память, CPU, GC |

## Импорт

1. Grafana → **Dashboards** → **New** → **Import**.
2. Загрузить JSON-файл.
3. Выбрать Prometheus datasource (в файлах указан UID `af0lqps3nq58gf` — при необходимости заменить на свой).
4. Нажать **Import**.

## Переменные

`userinfo-service-health.json`:

- `job` — задание Prometheus, которое собирает метрики UserInfo;
- `instance` — конкретный экземпляр;
- `domain` — домен AD для панели последнего успешного импорта.

`ad-users-account-state.json`:

- `domain` — домен AD, можно выбрать несколько или **All**.

## Что должно быть развёрнуто

Дашборды используют метрики, добавленные в коммите `d46a409`:

- `http_requests_total`, `http_request_duration_seconds`, `http_requests_in_flight`, `http_response_size_bytes`;
- `ad_cache_fill_duration_seconds`, `ad_cache_fill_errors_total`, `ad_cache_fill_timeouts_total`, `ad_cache_last_success_timestamp_seconds`, `ad_cache_last_complete_timestamp_seconds`;
- `ad_users_cache_requests_total`, `ad_users_build_duration_seconds`, `ad_users_processed_users`;
- расширенный `ad_users_per_domain` с `type=privileged|service_mail|service|archive|sap_disabled|attribute_errors|password_never_expires|password_cant_change|no_computer|os_unsupported`.

После первого деплоя метрики `ad_cache_*` и новые `ad_users_per_domain{type=...}` появятся только после очередного заполнения кеша AD.

## Полезные алерты

```promql
# кеш AD не обновлялся больше 25 минут
time() - ad_cache_last_success_timestamp_seconds > 1500

# ошибки или таймауты импорта AD
increase(ad_cache_fill_errors_total[15m]) > 0
increase(ad_cache_fill_timeouts_total[15m]) > 0

# доля 5xx выше 1%
sum(rate(http_requests_total{status=~"5.."}[5m]))
/ sum(rate(http_requests_total[5m])) > 0.01

# p95 запроса по endpoint выше 1 секунды
histogram_quantile(0.95,
  sum by (route, le) (rate(http_request_duration_seconds_bucket[5m]))
) > 1
```
