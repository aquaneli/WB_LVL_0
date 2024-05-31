CREATE DATABASE wb_db;

CREATE USER maksim WITH PASSWORD '12345';
GRANT ALL PRIVILEGES ON DATABASE wb_db TO maksim;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO maksim;

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
	sm_id              BIGINT,
	date_created       TIME,
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
	amount       BIGINT,
	payment_dt    BIGINT,
	bank         VARCHAR,
	delivery_cost BIGINT,
	goods_total   BIGINT,
	custom_fee    BIGINT
);

CREATE TABLE IF NOT EXISTS items(
    id BIGINT PRIMARY KEY NOT NULL,
    order_id BIGINT REFERENCES information_order(id),
    chrt_id      BIGINT,
	track_number VARCHAR,
	price       BIGINT,
	rid         VARCHAR,
	name        VARCHAR,
	sale        BIGINT,
	size        VARCHAR,
	total_price  BIGINT,
	nm_id        BIGINT,
	brand       VARCHAR,
	status      BIGINT
);
