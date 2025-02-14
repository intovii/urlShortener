\c postgres
CREATE EXTENSION IF NOT EXISTS dblink;
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'urls') THEN
      PERFORM dblink_exec('dbname=postgres user=' || current_user, 'CREATE DATABASE urls');
   END IF;
END
$$;
\c urls
DO
$$
    BEGIN
        CREATE TABLE IF NOT EXISTS urls
        (
            short_url   VARCHAR(10) PRIMARY KEY,
            url         TEXT NOT NULL UNIQUE
        );
        
        
        RAISE NOTICE 'Таблицы успешно созданы.';
    EXCEPTION
        WHEN OTHERS THEN
            RAISE EXCEPTION 'Ошибка при создании таблиц: %', SQLERRM;
    END
$$;
