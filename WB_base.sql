CREATE DATABASE wb_db;

CREATE ROLE adm
NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
DROP ROLE adm;

CREATE USER maksim;
DROP USER maksim;

DROP TABLE information_order;

CREATE TABLE IF NOT EXISTS Delivery(
    id BIGINT PRIMARY KEY NOT NULL,
    Name    varchar,
	Phone   varchar,
	Zip     varchar,
	City    varchar,
	Address varchar,
	Region  varchar,
	Email   varchar
);

CREATE TABLE IF NOT EXISTS Payment(
    id BIGINT PRIMARY KEY NOT NULL,
    Transaction  varchar,
	RequestId    varchar,
	Currency     varchar,
	Provider     varchar,
	Amount       int,
	PaymentDt    int,
	Bank         varchar,
	DeliveryCost int,
	GoodsTotal   int,
	CustomFee    int
);

CREATE TABLE IF NOT EXISTS Items(
    id BIGINT PRIMARY KEY NOT NULL,
    ChrtId      int,
	TrackNumber varchar,
	Price       int,
	Rid         varchar,
	Name        varchar,
	Sale        int,
	Size        varchar,
	TotalPrice  int,
	NmID        int,
	Brand       varchar,
	Status      int
);

CREATE TABLE IF NOT EXISTS information_order(
    id BIGINT PRIMARY KEY NOT NULL,
    OrderUid  varchar,
    TrackNumber varchar,
    Entry varchar,
    Delivery BIGINT REFERENCES Delivery(id),
    Payment BIGINT REFERENCES Payment(id),
    Items BIGINT REFERENCES Items(id),
    Local             varchar,
	InternalSignature varchar,
	CustomerId        varchar,
	DeliveryService   varchar,
	Shardkey          varchar,
	SmId              int,
	DateCreated       time,
	OofShard          varchar
);



