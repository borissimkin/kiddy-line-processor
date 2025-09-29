# todo

[x] протобаф файлы вынести в папку api
[x] cmd добавить вложеность для сервиса
[x] добавить docs
[x] рефактор мейкфайла
[x] рефактор докерфайла
[x] разбить проект на фичи а не по функциональному признаку
[x] вынести интерфейс стораджа из модуля с реализацией
[x] рефактор докер композа (переименование сервиса приложения)
[x] рефактор grpc сервера 
[x] добавить логирование
[] оформить ридми (ссылка на тз, объяснить структуру модулей)
сделать пресеты команд для проверки функционала
// {"sport": "soccer", "sport": "football", "interval": "3s"}
// {"sport": "soccer", "sport": "football", "interval": "1s"}
// {"sport": "baseball", "sport": "football", "interval": "5s"}

тесты

# Kiddy Line Processor

![Go Version](https://img.shields.io/badge/Go-1.23.1-blue)
![Build](https://img.shields.io/badge/build-passing-brightgreen)

Программа для обработки спортивных коэффициентов из внешнего сервиса.

- Получает коэффициенты и сохраняет их в **Redis**.
- Предоставляет **bidirectional streaming gRPC** для подписки на дельты коэффициентов.
- Поднимает **HTTP-сервер** для проверки доступности и факта первой синхронизации.

Подробное ТЗ: [docs/Softpro.Тестовое_задание.pdf](https://github.com/borissimkin/kiddy-line-processor/blob/refactor/docs/Softpro.%D0%A2%D0%B5%D1%81%D1%82%D0%BE%D0%B2%D0%BE%D0%B5_%D0%B7%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5.pdf)

---

## Модули

- **app** — инициализация модулей и сервисов.
- **config** — конфигурация.
- **lineprocessor** — gRPC-сервер подписки на коэффициенты линий.
- **linesprovider** — получение коэффициентов из внешней системы и сохранение во внутреннее хранилище.
- **proto** — сгенерированные сервисы из protobuf-схем.
- **ready** — HTTP-сервис проверки соединения и первой синхронизации.
- **storage** — модуль хранилища (**Redis**).

---

## Запуск

```bash
make run
```

---

## Примеры запросов

### Проверка доступности сервиса
```bash
curl localhost:8080/ready
```
Пример успешного ответа:
```json
{"status":"ok"}
```

### Подписка на коэффициенты (gRPC streaming)
(необходимо [grpcurl](https://github.com/fullstorydev/grpcurl)):

```bash
grpcurl -plaintext -d @ localhost:8081 proto.SportsLinesService.SubscribeOnSportsLines
```

Передача списка спортов и интервала обновления:
```json
{
  "sport": "soccer",
  "sport": "football",
  "interval": "3s"
}
```

Пример ответа от сервиса:
```json
{
  "sport": "soccer",
  "coeff": 1.85,
  "updated_at": "2025-09-29T12:34:56Z"
}
```

---

## Технологии
- Go
- Redis
- gRPC
- HTTP


