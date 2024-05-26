package main

// go mod init github.com/lib/pq создался файл .mod
// go get github.com/lib/pq скачивание пакета

import (
	"database/sql"
	"encoding/json"
	_"net"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Orders struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Local             string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	db, err := sql.Open("postgres", "dbname=wb_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	var orders Orders
	nc.Subscribe("a", func(msg *nats.Msg) {
		err = json.Unmarshal(msg.Data, &orders)
		_, _ = db.Exec("insert into Delivery(id, Name,Phone,Zip,City,Address,Region,Email) values (COALESCE((SELECT MAX(id) FROM Delivery), 0) + 1, $1,$2,$3,$4,$5,$6,$7)",
			orders.Delivery.Name, orders.Delivery.Phone, orders.Delivery.Zip, orders.Delivery.City, orders.Delivery.Address, orders.Delivery.Region, orders.Delivery.Email)

		_, _ = db.Exec("insert into Payment(id, Transaction, RequestId , Currency, Provider, Amount, PaymentDt,Bank,DeliveryCost, GoodsTotal, CustomFee) values (COALESCE((SELECT MAX(id) FROM Delivery), 0) + 1, $1,$2,$3,$4,$5,$6,$7)",
			orders.Payment.Transaction, orders.Payment.RequestId, orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount, orders.Payment.PaymentDt, orders.Payment.Bank, orders.Payment.DeliveryCost, orders.Payment.GoodsTotal, orders.Payment.CustomFee)

	})
}
