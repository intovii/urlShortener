# urlShortener
## Описание
Сервис для сокращения ссылок. 

## API Endpoints (Задание)
`/post` - 

`/get/{short_url}` - 

## Инструкция по запуску
Склонировать репозиторий и перейти в рабочую директорию.
``` bash
git clone https://github.com/intovii/
cd 
```
Запуск проекта с базой данных PostgreSQL.
``` bash
make psql
```
Запуск проекта с in-memory базой данных.
``` bash
make in_memory
```
Остановка проекта и удаление соответсвующего тома.
``` bash
make remove
```
Запуск тестов.
``` bash
make tests
```
