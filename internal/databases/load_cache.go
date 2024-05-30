package databases

import (
	"database/sql"
	"log"
	"module_0/internal/models"
	"time"

	"github.com/maxchagin/go-memorycache-example"
)

func LoadInCache(cache *memorycache.Cache) {
	db, err := sql.Open("postgres", "dbname=wb_db sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, order_uid, track_number, entry, local, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM information_order")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Orders = models.Orders{}
		var id int
		err := rows.Scan(&id, &order.OrderUid, &order.TrackNumber, &order.Entry, &order.Local, &order.InternalSignature,
			&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		if err != nil {
			log.Fatalln(err)
		}

		selectDelivery(db, &order, &id)
		selectPayment(db, &order, &id)
		selectItems(db, &order, &id)

		cache.Set(order.OrderUid, order, 1*time.Hour)
	}
}

func selectDelivery(db *sql.DB, order *models.Orders, id *int) {
	rows_delivery, err := db.Query("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_id = $1", id)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows_delivery.Close()
	rows_delivery.Next()
	err = rows_delivery.Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		log.Fatalln(err)
	}
}

func selectPayment(db *sql.DB, order *models.Orders, id *int) {
	rows_payment, err := db.Query("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_id = $1", id)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows_payment.Close()
	rows_payment.Next()
	err = rows_payment.Scan(&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee)
	if err != nil {
		log.Fatalln(err)
	}
}

func selectItems(db *sql.DB, order *models.Orders, id *int) {
	rows_item, err := db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_id = $1", id)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows_item.Close()
	for rows_item.Next() {
		var itm models.Item
		err = rows_item.Scan(&itm.ChrtId, &itm.TrackNumber, &itm.Price, &itm.Rid, &itm.Name, &itm.Sale, &itm.Size, &itm.TotalPrice, &itm.NmID, &itm.Brand, &itm.Status)
		if err != nil {
			log.Fatalln(err)
		}
		order.Items = append(order.Items, itm)
	}
}
