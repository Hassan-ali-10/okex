package main

import (
	
	routes "okex/routes"
	
	// "fmt"
     crons "okex/crons"
	
)

func main() {
	
	routes.SetRouters()
	crons.AllCrons();
	
}