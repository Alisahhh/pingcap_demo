package main

import (
	"database/sql"
	"log"
 _ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
	"strconv"
	"sync"
	"math"
)

func getYAddr(xAddr int) int{
	yDoubleAddr:= math.Sin(float64(xAddr) / 30 * (math.Pi))
	return int(yDoubleAddr * 30 + 0.5) + 30
}

func main(){

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:4000)/test?parseTime=true")
    if err != nil{
        log.Fatal(err)
		}
		fmt.Println("connection succedssed")
		defer db.Close()
		
		orderPrefix := "CREATE TABLE IF NOT EXISTS test.testTable";
		orderSuffix := "(id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL, data varchar(50))"
		var wg sync.WaitGroup
		for i:=0;i <= 60;i ++{
			order := orderPrefix + strconv.Itoa(i) + orderSuffix
			wg.Add(1)
			go func(){
				_, err := db.Exec(order)
				if err != nil{
						log.Fatalln(err)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		t := time.NewTicker(30*time.Second)
		defer t.Stop()

		whichTime := 0;
		yAddr := 30
		dataValue := 0

		for whichTime <= 60{
			select{
				case <-t.C:{
					dataValue = 0
					yAddr = getYAddr(whichTime)
					whichTime++
					fmt.Println(yAddr)
					fmt.Println(whichTime)
					break
				}
				default:{
					dataValue ++
					insertOrderPrefix := "INSERT INTO test.testTable"
					insertOrderSuffix := "(data) VALUES (?)"
					insertOrder := insertOrderPrefix +  strconv.Itoa(yAddr) + " " + insertOrderSuffix
					_,err := db.Exec(insertOrder,strconv.Itoa(dataValue))
					if err != nil{
						log.Fatalln(err)
					}
					queryOrderPrefix := "SELECT * FROM test.testTable" 
					queryOrder :=  queryOrderPrefix + strconv.Itoa(yAddr)
					row, err := db.Query(queryOrder)
					if err != nil{
						log.Fatalln(err)
					}
					row.Close()
					break
				}
			}
		}

}
