CREATE DATABASE wb_db;

CREATE ROLE adm
NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
DROP ROLE adm;

CREATE USER maksim;
DROP USER maksim;

DROP TABLE information_order;

CREATE TABLE IF NOT EXISTS information_order(
    id BIGINT PRIMARY KEY NOT NULL,
    
);



