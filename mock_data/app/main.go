package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	host                  = "postgres"
	port                  = 5432
	user                  = "postgres"
	password              = "postgres"
	dbname                = "order_db"
	customer_count        = 200
	order_update_delay_ms = 250
	order_update_count    = 50
	new_order_delay_ms    = 15000
)

type Order struct {
	OrderId     int
	OrderAmount float64
	OrderStatus int
	CustomerId  int
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	go generateOrders(db)
	time.Sleep(1 * time.Second)
	for {
		updateRandomOrder(db)
		time.Sleep(order_update_delay_ms * time.Millisecond)
	}
}

func createNewOrder(db *sql.DB, customerId int) error {
	query := `INSERT INTO order_schema.order (order_amount, order_status_id, customer_id) VALUES ($1, $2, $3);`
	_, err := db.Exec(query, float64(rand.Intn(18001)+2000)/100, 1, customerId)
	if err != nil {
		return err
	}
	return nil
}

func generateOrders(db *sql.DB) {
	// generate initial order concurrently
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				customerId := rand.Intn(customer_count) + 1
				err := createNewOrder(db, customerId)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Finished generating initial order")
	// generate new order every 15s
	for {
		customerId := rand.Intn(customer_count) + 1
		err := createNewOrder(db, customerId)
		if err != nil {
			panic(err)
		}
		time.Sleep(new_order_delay_ms * time.Millisecond)
	}
}

func updateRandomOrder(db *sql.DB) {
	// get the number of order
	var count int
	var wg sync.WaitGroup
	err := db.QueryRow("SELECT MAX(order_id) FROM order_schema.order;").Scan(&count)
	if err != nil {
		panic(err)
	}
	// randomly select 50 order to update
	for i := 0; i < order_update_count; i++ {
		wg.Add(1)
		go func() {
			orderId := rand.Intn(count) + 1
			var order Order
			err := db.QueryRow("SELECT order_id,order_amount,order_status_id,customer_id FROM order_schema.order WHERE order_id = $1;", orderId).Scan(&order.OrderId, &order.OrderAmount, &order.OrderStatus, &order.CustomerId)
			if err != nil {
				panic(err)
			}
			// increment or cancel the order status if it is less than 3
			if order.OrderStatus < 4 {
				cancelOrder := rand.Intn(20)
				if cancelOrder == 0 {
					order.OrderStatus = 5
				} else {
					order.OrderStatus += 1
				}
				// update the order
				_, err = db.Exec("UPDATE order_schema.order SET order_status_id = $1, updated_at = CURRENT_TIMESTAMP WHERE order_id = $2;", order.OrderStatus, order.OrderId)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
