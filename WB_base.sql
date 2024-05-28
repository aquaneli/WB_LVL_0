CREATE DATABASE wb_db;

CREATE ROLE adm
NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
DROP ROLE adm;

CREATE USER maksim;
DROP USER maksim;

DROP TABLE delivery;
DROP TABLE payment;
DROP TABLE items;
DROP TABLE information_order;

TRUNCATE TABLE delivery, payment, items, information_order;

CREATE TABLE IF NOT EXISTS information_order(
    id BIGINT PRIMARY KEY NOT NULL,
    order_uid  VARCHAR,
    track_number VARCHAR,
    entry VARCHAR,
    local   VARCHAR,
	internal_signature VARCHAR,
	customer_id        VARCHAR,
	delivery_service   VARCHAR,
	shardkey          VARCHAR,
	sm_id              INT,
	date_created       time,
	oof_shard          VARCHAR
);

CREATE TABLE IF NOT EXISTS delivery(
    id BIGINT PRIMARY KEY NOT NULL,
    order_id BIGINT REFERENCES information_order(id),
    name    VARCHAR,
	phone   VARCHAR,
	zip     VARCHAR,
	city    VARCHAR,
	address VARCHAR,
	region  VARCHAR,
	email   VARCHAR
);

CREATE TABLE IF NOT EXISTS payment(
    id BIGINT PRIMARY KEY NOT NULL,
    order_id BIGINT REFERENCES information_order(id),
    transaction  VARCHAR,
	request_id    VARCHAR,
	currency     VARCHAR,
	provider     VARCHAR,
	amount       INT,
	payment_dt    INT,
	bank         VARCHAR,
	delivery_cost INT,
	goods_total   INT,
	custom_fee    INT
);

CREATE TABLE IF NOT EXISTS items(
    id BIGINT PRIMARY KEY NOT NULL,
    order_id BIGINT REFERENCES information_order(id),
    chrt_id      INT,
	track_number VARCHAR,
	price       INT,
	rid         VARCHAR,
	name        VARCHAR,
	sale        INT,
	size        VARCHAR,
	total_price  INT,
	nm_id        INT,
	brand       VARCHAR,
	status      INT
);

