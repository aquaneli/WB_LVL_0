package databases

import (
	"database/sql"
	"encoding/json"
	"log"
	"module_0/internal/models"
	"sync"
	"time"

	"github.com/maxchagin/go-memorycache-example"
	"github.com/nats-io/nats.go"
)

func NatsSub(cache *memorycache.Cache, wg *sync.WaitGroup) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()

	db, err := sql.Open("postgres", "dbname=wb_db sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	_, err = nc.Subscribe("orders", func(msg *nats.Msg) {
		var orders models.Orders
		err = json.Unmarshal(msg.Data, &orders)
		if err != nil {
			log.Fatalln(err)
		}
		_, res := cache.Get(orders.OrderUid)
		if !res {
			insertInformationOrder(db, &orders)
			insertDelivery(db, &orders)
			insertPayment(db, &orders)
			insertItems(db, &orders)
			cache.Set(orders.OrderUid, orders, 1*time.Hour)
			log.Println("Message published successfully")
		} else {
			log.Println("This OrderUid already exists")
		}
	})
	if err != nil {
		log.Fatalln(err)
	}
	wg.Wait()
}

func insertInformationOrder(db *sql.DB, orders *models.Orders) {
	_, err := db.Exec(`INSERT INTO information_order(id, order_uid, track_number, entry, local, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
			VALUES (COALESCE((SELECT MAX(id) FROM information_order), 0) + 1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, orders.OrderUid, orders.TrackNumber, orders.Entry, orders.Local,
		orders.InternalSignature, orders.CustomerId, orders.DeliveryService, orders.Shardkey, orders.SmId, orders.DateCreated, orders.OofShard)
	if err != nil {
		log.Fatalln(err)
	}
}

func insertDelivery(db *sql.DB, orders *models.Orders) {
	_, err := db.Exec(`INSERT INTO delivery(id, order_id, name, phone, zip, city, address, region, email)
	VALUES (COALESCE((SELECT MAX(id) FROM delivery), 0) + 1, (SELECT MAX(id) FROM information_order), $1, $2, $3, $4, $5, $6, $7)`, orders.Delivery.Name, orders.Delivery.Phone,
		orders.Delivery.Zip, orders.Delivery.City, orders.Delivery.Address, orders.Delivery.Region, orders.Delivery.Email)
	if err != nil {
		log.Fatalln(err)
	}
}

func insertPayment(db *sql.DB, orders *models.Orders) {
	_, err := db.Exec(`INSERT INTO payment(id, order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
				VALUES (COALESCE((SELECT MAX(id) FROM payment), 0) + 1, (SELECT MAX(id) FROM information_order), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, orders.Payment.Transaction,
		orders.Payment.RequestId, orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount, orders.Payment.PaymentDt, orders.Payment.Bank, orders.Payment.DeliveryCost,
		orders.Payment.GoodsTotal, orders.Payment.CustomFee)
	if err != nil {
		log.Fatalln(err)
	}
}

func insertItems(db *sql.DB, orders *models.Orders) {
	for _, value := range orders.Items {
		_, err := db.Exec(`INSERT INTO items(id, order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES (COALESCE((SELECT MAX(id) FROM items), 0) + 1, (SELECT MAX(id) FROM information_order), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, value.ChrtId,
			value.TrackNumber, value.Price, value.Rid, value.Name, value.Sale, value.Size, value.TotalPrice, value.NmID, value.Brand, value.Status)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
