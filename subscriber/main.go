package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		db, _ := sql.Open("postgres", "dbname=wb_db sslmode=disable")
		defer db.Close()

		rows, _ := db.Query("SELECT id, order_uid, track_number, entry, local, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM information_order")
		defer rows.Close()

		var orders []Orders
		for rows.Next() {
			var order Orders = Orders{}
			var id int
			err := rows.Scan(&id, &order.OrderUid, &order.TrackNumber, &order.Entry, &order.Local, &order.InternalSignature,
				&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
			if err != nil {
				log.Fatalln(err)
			}

			rows_delivery, _ := db.Query("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_id = $1", id)
			defer rows_delivery.Close()
			rows_delivery.Next()
			err = rows_delivery.Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
				&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
			if err != nil {
				log.Fatalln(err)
			}

			rows_payment, _ := db.Query("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_id = $1", id)
			defer rows_payment.Close()
			rows_payment.Next()
			err = rows_payment.Scan(&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider,
				&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
				&order.Payment.CustomFee)
			if err != nil {
				log.Fatalln(err)
			}

			rows_item, _ := db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_id = $1", id)
			defer rows_item.Close()
			for rows_item.Next() {
				var itm Item
				err = rows_item.Scan(&itm.ChrtId, &itm.TrackNumber, &itm.Price, &itm.Rid, &itm.Name, &itm.Sale, &itm.Size, &itm.TotalPrice, &itm.NmID, &itm.Brand, &itm.Status)
				if err != nil {
					log.Fatalln(err)
				}
				order.Items = append(order.Items, itm)
			}

			orders = append(orders, order)
		}

		js, _ := json.Marshal(orders)

		fmt.Println(string(js))

	}()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl, _ := template.ParseFiles("home.html")
			tmpl.Execute(w, "")
		})
		fmt.Println("Server is listening...")
		http.ListenAndServe("localhost:8181", nil)
	}()

	go func() {
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

		nc.Subscribe("a", func(msg *nats.Msg) {
			var orders Orders
			err = json.Unmarshal(msg.Data, &orders)

			_, _ = db.Exec(`INSERT INTO information_order(id, order_uid, track_number, entry, local, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
			VALUES (COALESCE((SELECT MAX(id) FROM information_order), 0) + 1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, orders.OrderUid, orders.TrackNumber, orders.Entry, orders.Local,
				orders.InternalSignature, orders.CustomerId, orders.DeliveryService, orders.Shardkey, orders.SmId, orders.DateCreated, orders.OofShard)

			_, _ = db.Exec(`INSERT INTO delivery(id, order_id, name, phone, zip, city, address, region, email)
			VALUES (COALESCE((SELECT MAX(id) FROM delivery), 0) + 1, (SELECT MAX(id) FROM information_order), $1, $2, $3, $4, $5, $6, $7)`, orders.Delivery.Name, orders.Delivery.Phone,
				orders.Delivery.Zip, orders.Delivery.City, orders.Delivery.Address, orders.Delivery.Region, orders.Delivery.Email)

			_, _ = db.Exec(`INSERT INTO payment(id, order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
				VALUES (COALESCE((SELECT MAX(id) FROM payment), 0) + 1, (SELECT MAX(id) FROM information_order), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, orders.Payment.Transaction,
				orders.Payment.RequestId, orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount, orders.Payment.PaymentDt, orders.Payment.Bank, orders.Payment.DeliveryCost,
				orders.Payment.GoodsTotal, orders.Payment.CustomFee)

			for _, value := range orders.Items {
				_, err = db.Exec(`INSERT INTO items(id, order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
				VALUES (COALESCE((SELECT MAX(id) FROM items), 0) + 1, (SELECT MAX(id) FROM information_order), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, value.ChrtId,
					value.TrackNumber, value.Price, value.Rid, value.Name, value.Sale, value.Size, value.TotalPrice, value.NmID, value.Brand, value.Status)

				if err != nil {
					log.Fatalln(err)
				}
			}

		})
		wg.Wait()
	}()

	wg.Wait()
}
