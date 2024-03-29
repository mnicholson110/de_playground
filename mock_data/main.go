package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const (
	host = "0.0.0.0"
	port = 5432
	user = "postgres"
)

type Order struct {
	OrderId     int
	OrderAmount float64
	OrderStatus int
	CustomerId  int
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable",
		host, port, user)
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
		time.Sleep(100 * time.Millisecond)
	}
}

func createNewOrder(db *sql.DB, customerId int) error {
	query := `INSERT INTO order_schema.orders (order_amount, order_status_id, customer_id) VALUES ($1, $2, $3);`
	_, err := db.Exec(query, float64(rand.Intn(18001)+2000)/100, 1, customerId)
	if err != nil {
		return err
	}
	return nil
}

func generateOrders(db *sql.DB) {
	// generate initial 1000 orders
	for i := 0; i < 1000; i++ {
		customerId := rand.Intn(50) + 1
		err := createNewOrder(db, customerId)
		if err != nil {
			panic(err)
		}
	}
	// generate new orders every 500ms
	for {
		customerId := rand.Intn(50) + 1
		err := createNewOrder(db, customerId)
		if err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func updateRandomOrder(db *sql.DB) {
	// get the number of orders
	var count int
	err := db.QueryRow("SELECT MAX(order_id) FROM order_schema.orders;").Scan(&count)
	if err != nil {
		panic(err)
	}
	// randomly select 5 orders to update
	for i := 0; i < 5; i++ {
		orderId := rand.Intn(count) + 1
		var order Order
		err := db.QueryRow("SELECT * FROM order_schema.orders WHERE order_id = $1;", orderId).Scan(&order.OrderId, &order.OrderAmount, &order.OrderStatus, &order.CustomerId)
		if err != nil {
			panic(err)
		}
		// increment or cancel the order status if it is less than 3
		if order.OrderStatus < 4 {
			cancelOrder := rand.Intn(10)
			if cancelOrder == 0 {
				order.OrderStatus = 5
			} else {
				order.OrderStatus += 1
			}
			// update the order
			_, err = db.Exec("UPDATE order_schema.orders SET order_status_id = $1 WHERE order_id = $2;", order.OrderStatus, order.OrderId)
			if err != nil {
				panic(err)
			}
		}
	}
}
