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

var (
	host                     = os.Getenv("HOST")
	port                     = os.Getenv("PORT")
	user                     = os.Getenv("USER")
	password                 = os.Getenv("PASSWORD")
	dbname                   = os.Getenv("DBNAME")
	customer_count           = 200
	order_count              = 0
	initial_order_count_k, _ = strconv.Atoi(os.Getenv("INITIAL_ORDER_COUNT_K"))
	order_update_count_k, _  = strconv.Atoi(os.Getenv("ORDER_UPDATE_COUNT_K"))
	new_order_count_k, _     = strconv.Atoi(os.Getenv("NEW_ORDER_COUNT_K"))
)

type Order struct {
	OrderId     int
	OrderAmount float64
	OrderStatus int
	CustomerId  int
}

func main() {
	fmt.Println("INITIAL_ORDER_COUNT_K: ", initial_order_count_k)
	fmt.Println("ORDER_UPDATE_COUNT_K: ", order_update_count_k)
	fmt.Println("NEW_ORDER_COUNT_K: ", new_order_count_k)

	start := time.Now()
	generateOrders(initial_order_count_k)
	fmt.Println("Finished generating initial orders in", time.Since(start).Seconds(), "seconds.")
	order_count = initial_order_count_k * 1000

	for {
		start = time.Now()
		updateOrders()
		fmt.Println("Finished updating orders in", time.Since(start).Seconds(), "seconds.")
		time.Sleep(5 * time.Second)

		start = time.Now()
		generateOrders(new_order_count_k)
		fmt.Println("Finished generating new orders in", time.Since(start).Seconds(), "seconds.")
		order_count += new_order_count_k * 1000

	}
}

func generateOrders(order_count int) {
	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
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

	// generate orders concurrently
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < order_count; i++ {
				query := `INSERT INTO order_schema.order (order_amount, order_status_id, customer_id) VALUES ($1, $2, $3);`
				_, err := db.Exec(query, float64(rand.Intn(18001)+2000)/100, 1, rand.Intn(customer_count)+1)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func updateOrders() {
	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
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

	var wg sync.WaitGroup
	// randomly select orders to update
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < order_update_count_k; i++ {
				orderId := rand.Intn(order_count) + 1
				cancelOrder := rand.Intn(20)
				if cancelOrder == 0 {
					_, err = db.Exec(`UPDATE order_schema.order
                            SET order_status_id = 5, updated_at = CURRENT_TIMESTAMP
                            WHERE order_id = $1 AND order_status_id NOT IN (4,5);`,
						orderId)
					if err != nil {
						panic(err)
					}
				} else {
					// update the order
					_, err = db.Exec(`UPDATE order_schema.order
                            SET order_status_id = order_status_id + 1, updated_at = CURRENT_TIMESTAMP
                            WHERE order_id = $1 AND order_status_id NOT IN (4,5);`,
						orderId)
					if err != nil {
						panic(err)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
