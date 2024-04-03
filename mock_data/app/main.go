package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	host           = "postgres"
	port           = 5432
	user           = "postgres"
	password       = "postgres"
	dbname         = "order_db"
	customer_count = 200
)

var (
	initial_order_count_k, _ = strconv.Atoi(os.Getenv("INITIAL_ORDER_COUNT_K"))
	order_update_delay_ms, _ = strconv.Atoi(os.Getenv("ORDER_UPDATE_DELAY_MS"))
	order_update_count_da, _ = strconv.Atoi(os.Getenv("ORDER_UPDATE_COUNT_DA"))
	new_order_delay_ms, _    = strconv.Atoi(os.Getenv("NEW_ORDER_DELAY_MS"))
)

type Order struct {
	OrderId     int
	OrderAmount float64
	OrderStatus int
	CustomerId  int
}

func main() {
	fmt.Println("INITIAL_ORDER_COUNT_K: ", initial_order_count_k)
	fmt.Println("ORDER_UPDATE_DELAY_MS: ", order_update_delay_ms)
	fmt.Println("ORDER_UPDATE_COUNT_DA: ", order_update_count_da)
	fmt.Println("NEW_ORDER_DELAY_MS: ", new_order_delay_ms)

	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(500)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	go generateOrders(db)
	time.Sleep(5 * time.Second) // wait for some initial orders to be generated

	delay := time.Duration(order_update_delay_ms)
	fmt.Println("Start updating orders")
	for {
		updateRandomOrder(db)
		time.Sleep(delay * time.Millisecond)
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
	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < initial_order_count_k*10; i++ {
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

	fmt.Println("Finished generating initial orders in ", time.Since(start).Seconds(), " seconds.")
	// generate new orders
	delay := time.Duration(new_order_delay_ms)
	for {
		customerId := rand.Intn(customer_count) + 1
		err := createNewOrder(db, customerId)
		if err != nil {
			panic(err)
		}
		time.Sleep(delay * time.Millisecond)
	}
}

func updateRandomOrder(db *sql.DB) {
	// get the number of order
	var count int
	var wg sync.WaitGroup
	start := time.Now()
	err := db.QueryRow("SELECT MAX(order_id) FROM order_schema.order;").Scan(&count)
	if err != nil {
		panic(err)
	}
	// randomly select orders to update
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < order_update_count_da; i++ {
				orderId := rand.Intn(count) + 1
				cancelOrder := rand.Intn(20)
				if cancelOrder == 0 {
					_, err = db.Exec("UPDATE order_schema.order SET order_status_id = 5, updated_at = CURRENT_TIMESTAMP WHERE order_id = $1 AND order_status_id != 4;", orderId)
					if err != nil {
						panic(err)
					}
				} else {
					// update the order
					_, err = db.Exec("UPDATE order_schema.order SET order_status_id = order_status_id + 1, updated_at = CURRENT_TIMESTAMP WHERE order_id = $1 AND order_status_id NOT IN (4,5);", orderId)
					if err != nil {
						panic(err)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Finished updating ", fmt.Sprint(order_update_count_da*10), " orders in ", time.Since(start).Seconds(), " seconds.")
}
