package api

import (
	"encoding/json"
	"context"
	//"log"
	"fmt"
	"net/http"
	//"time"
	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	 connectionhelper "okex/db"
	 //config "okex/config"
	 //constants "okex/constants"
	 helpers "okex/common"
	//"io/ioutil"
	"sync"
)
var (
	counter int32          // counter is a variable incremented by all goroutines.
	wg      sync.WaitGroup // wg is used to wait for the program to finish.
	mutex   sync.Mutex     // mutex is used to define a critical section of code.
)
	//A mutex is used to create a critical section around code that ensures only one goroutine at a time can execute that code section.
	
func Home(response http.ResponseWriter, request *http.Request) {
response.Header().Set("content-type", "application/json")
//Perform Find operation & validate against the error.
collectionName:="employees2" // for office	
	//collectionName:="hssn1sss"  // for home
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	var results []bson.M
	//Get MongoDB collection using connectionhelper.
	collection, err := connectionhelper.GetMongoDbCollection(collectionName)
	if err != nil {
		defaultJSON := bson.M{
			"success": false,
			"message": err,
		}
		json.NewEncoder(response).Encode(defaultJSON)
		return
	}

	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
			defaultJSON := bson.M{
			"success": false,
			"message": findError,
		}
		json.NewEncoder(response).Encode(defaultJSON)
		return
	}
	//Map result to slice
	cur.All(context.TODO(), &results)
	// once exhausted, close the cursor  it will free mongodb server resource ...
	cur.Close(context.TODO())
	// return response
	// fmt.Println(results)
	if len(results) == 0 {
	    defaultJSON := bson.M{
			"success": true,
			"message": "Nothing Found",
		}
		json.NewEncoder(response).Encode(defaultJSON)
	} else {
		json.NewEncoder(response).Encode(results)
	}
//json.NewEncoder(response).Encode("WELCOME HOME")
}
func ExecuteOrdersPostRequest(response http.ResponseWriter, request *http.Request) {
	
	// validate headers ....

 	// ok := request.Header.Get("rulesRequest")
 	// fmt.Println("353vsvc")
 	// fmt.Println(ok)
	// if !ok  {  // this request can not be identified please return...
	//        result := bson.M{
	// 		"success": false,
	// 		"message": "Invalid Request",
	// 	}
	// 	json.NewEncoder(response).Encode(result)
	// 	return;
 //    }
    
	response.Header().Set("content-type", "application/json")
	var wg sync.WaitGroup
	orders := make(chan bool)
	
	 // Declare a unbuffered channel
	wg.Add(1)
	var payload map[string]interface{}
	//Decode Incoming Payload By Mapping it on payload i.e. Struct Instance in Models
	_ = json.NewDecoder(request.Body).Decode(&payload)
	go helpers.PickParentsAndMakeChilds(payload, orders,&wg)
	go func() {
		fmt.Println("channel is closed")
		wg.Wait()	// this blocks the goroutine until WaitGroup counter is zero
		close(orders)    // Channels need to be closed,
	}()    // This calls itself
	 defaultJSON := bson.M{
			"success": true,
			"message": "Creating orders",
		}
		fmt.Println("HASSAN")
		json.NewEncoder(response).Encode(defaultJSON)
		return
	
	
	
}
// Note : For Okex we are only Implementing it to execute  Live Orders ....
// func executeAutoOrder(response http.ResponseWriter, request *http.Request) {

// 	//Response Header Definition
// 	response.Header().Set("content-type", "application/json")
// 	// validate Headers ....
// 	headers := request.Header
	
// 	val, ok := headers["rulesRequest"]
// 	if !ok || val != "yes" {  // this request can not be identified please return...
// 	       result := bson.M{
// 			"success": false,
// 			"message": "Invalid Request",
// 		}
// 		json.NewEncoder(response).Encode(result)
// 		return;
//     }
// 	OrderMode:= constants.GLOBALMODE
	
// 	var payload map[string]interface{}
// 	var mode string
// 	trigger_type := "barrier_percentile_trigger"
	
// 	_ = json.NewDecoder(request.Body).Decode(&payload)


// 	//Changes interface int types to float64 type
// 	payload = helpers.AlterInterfaceTypesToFloat(payload)

// 	liveTradingBool := payload["enable_buy_barrier_percentile"]
// 	testTradingBool := payload["enable_test_buy_barrier_percentile"]

// 	switch {
//     case OrderMode == "live":
//         orderpickType := constants.LIVEORDERSTYPE
//     case OrderMode == "test":
//         orderpickType := constants.TESTORDERSTYPE
//     case OrderMode == "both":
//         orderpickType := constants.TESTANDLIVEBOTHORDERSTYPE   
//     }
//     // we will pick and execute orders on the basis of our filter results .....
//     // if our rules live is true and our exchange also handles any one either that is live or both
//     if(liveTradingBool == true && (orderpickType == 1 || orderpickType == 3)){
//     	mode:="live"
//     	buyOrderExecutionAckChannel := make(chan bool)
//     	//GoRoutine
// 		go func() {
// 			buyOrderExecutionAckChannel <- helpers.ExecuteAutoBuyOrders(trigger_type, mode, payload)
// 		}()

// 		//Use Channel Value
// 		 <-buyOrderExecutionAckChannel
// 		//Close Channel
// 		close(buyOrderExecutionAckChannel)
//     }
//     // if our rules test is true and our exchange also handles any one either that is test or both
//     if(testTradingBool == true && (orderpickType == 2 || orderpickType == 3)){
//     	mode:="test"
//     	buyOrderExecutionAckChannel := make(chan bool)
//     	//GoRoutine
// 		go func() {
// 			buyOrderExecutionAckChannel <- helpers.ExecuteAutoBuyOrders(trigger_type, mode, payload)
// 		}()

// 		//Use Channel Value
// 		buyOrderExecutionAck := <-buyOrderExecutionAckChannel

// 		_ = buyOrderExecutionAck

// 		//Close Channel
// 		close(buyOrderExecutionAckChannel)
    	
//     }
	
// 	result := bson.M{
// 		"success": true,
// 		"message": "Creating orders",
// 	}
// 	json.NewEncoder(response).Encode(result)

// }