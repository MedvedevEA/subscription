#!/bin/bash
set -e

# Проверка и создание базы данных с русской локалью
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    -- Проверяем существование локали
    SELECT * FROM pg_collation WHERE collname LIKE 'ru%';
    
    -- Создаем базу данных с русской локалью
    CREATE DATABASE mydb_russian
        WITH 
        OWNER = postgres
        ENCODING = 'UTF8'
        LC_COLLATE = 'ru_RU.UTF-8'
        LC_CTYPE = 'ru_RU.UTF-8'
        TEMPLATE = template0;
EOSQL