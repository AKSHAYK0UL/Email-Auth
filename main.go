package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AKSHAYK0UL/Email_Auth/helper"
	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/AKSHAYK0UL/Email_Auth/routefunc"
)

//	func init() {
//		model.ConnectMongo()
//	}
func main() {
	fmt.Println("NEW Auth Server")
	model.ConnectMongo()
	route := routefunc.RouteTable()
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	go func() {
		route.Run(":" + port)
	}()
	for {
		fmt.Println("Ticker ticked at", time.Now())
		helper.DeleteVcode()
		time.Sleep(2 * time.Minute)
	}
}
