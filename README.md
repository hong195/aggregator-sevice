# Data Aggregation Service

Сервис для эмуляции внешних пакетов, их обработки (поиск максимального значения) и сохранения результатов в БД.  
Поддерживает REST API и gRPC.

## Быстрый старт

### Локальная разработка

```sh
# Postgres
make compose-up
# Запуск приложения и миграций
make run
```

### Весь docker stack с reverse proxy

```sh
make compose-up-all 
```

### Документация API

Документация генерируется автоматически с помощью swaggo/swag.
Команда для обновления:

```sh
make swag-v1
```

####  После генерации доступны файлы:
•	docs/swagger.json
•	docs/swagger.yaml

##### Просмотр
•	Локально (через UI):
http://localhost:8080/swagger/index.html
•	Через онлайн-редактор:
Swagger Editor → загрузите docs/swagger.yaml

##### Примеры эндпоинтов
•	GET /api/v1/packets/:id — получить пакет по ID
•	GET /api/v1/packets?start=...&end=... — получить список пакетов за период


#### gRPC

Прото-файлы лежат в docs/proto/v1/.

Сервис:

```protobuf
service AggregationService {
  rpc FindPacketByID (FindPacketByIDRequest) returns (FindPacketByIDResponse);
  rpc ListPacketsByPeriod (ListPacketsByPeriodRequest) returns (ListPacketsByPeriodResponse);
}
```
Сгенерировать gRPC-код:

```sh
make proto-v1 
```

