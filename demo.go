package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

func getYAddr(xAddr int) int {
	yDoubleAddr := math.Sin(float64(xAddr) / 30 * (math.Pi))
	return int(yDoubleAddr*30+0.5) + 30
}

func main() {
	// Connect to the sql database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:4000)/test?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection succedssed")
	defer db.Close()

	// Create 61 tables named testTableX
	orderPrefix := "CREATE TABLE IF NOT EXISTS test.testTable"
	orderSuffix := "(id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL, data varchar(50))"
	var wg sync.WaitGroup
	for i := 0; i <= 60; i++ {
		order := orderPrefix + strconv.Itoa(i) + orderSuffix
		wg.Add(1)
		go func() {
			_, err := db.Exec(order)
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Set a 1 min time ticker
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	// Draw sin function
	whichTime := 0
	yAddr := 30
	dataValue := 0

	for whichTime <= 60 {
		select {
		case <-t.C:
			{
				// switch to next block
				dataValue = 0
				yAddr = getYAddr(whichTime)
				whichTime++
				fmt.Println("location (%d,%d)", whichTime, yAddr)
				break
			}
		default:
			{
				// Insert
				dataValue++
				insertOrderPrefix := "INSERT INTO test.testTable"
				insertOrderSuffix := "(data) VALUES (?)"
				insertOrder := insertOrderPrefix + strconv.Itoa(yAddr) + " " + insertOrderSuffix
				_, err := db.Exec(insertOrder, strconv.Itoa(dataValue))
				if err != nil {
					log.Fatalln(err)
				}

				// Query
				queryOrderPrefix := "SELECT * FROM test.testTable"
				queryOrder := queryOrderPrefix + strconv.Itoa(yAddr)
				row, err := db.Query(queryOrder)
				if err != nil {
					log.Fatalln(err)
				}
				row.Close()
				break
			}
		}
	}
}
