# Используем официальный образ PostgreSQL как базовый
FROM postgres:17

# Копируем скрипт инициализации в контейнер
COPY init.sql /docker-entrypoint-initdb.d/

# Устанавливаем переменную окружения для инициализации базы данных
ENV POSTGRES_DB=postgres
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password

EXPOSE 5432