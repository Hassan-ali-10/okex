package main

import (
	
	routes "okex/routes"
	connectionhelper "okex/db"
	// "fmt"
	// config "okex/config"
	
)

func main() {
	//db.ConnectMongo()
	connectionhelper.GetMongoClient()
	routes.SetRouters()
	
}