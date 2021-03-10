package api

import (
	"encoding/json"
	"context"
	//"log"
	//"fmt"
	"net/http"
	//"time"
	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	 connectionhelper "okex/db"
	 config "okex/config"
	 constants "okex/constants"
	 helpers "okex/common"
)

	
	

func executeOrders(response http.ResponseWriter, request *http.Request) {
	collectionName:="employees"
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	var results []bson.M
	//Get MongoDB connection using connectionhelper.
	client, err := connectionhelper.GetMongoClient()
	if err != nil {
		json.NewEncoder(response).Encode(err)
		return
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(config.DBNAME).Collection(collectionName)
	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		json.NewEncoder(response).Encode(findError)
		return
	}
	//Map result to slice
	cur.All(context.TODO(), &results)
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	
	json.NewEncoder(response).Encode(results)
	// collectionName:="employees"
	// collection, err := db.GetMongoDbCollection(collectionName)
	// if err != nil {
	// 	json.NewEncoder(response).Encode(err)
		
	// 	return
	// }

	// var filter bson.M = bson.M{}
	// var results []bson.M
	// cur, err := collection.Find(context.Background(), filter)
	// defer cur.Close(context.Background())

	// if err != nil {
	// 	json.NewEncoder(response).Encode(err)
		
	// 	return
	// }

	// cur.All(context.Background(), &results)

	// if results == nil {
	// 	json.NewEncoder(response).Encode(404)
	// 	return
	// }
	// json.NewEncoder(response).Encode(results)
}
// Note : For Okex we are only Implementing it to execute  Live Orders ....
func executeAutoOrder(response http.ResponseWriter, request *http.Request) {

	//Response Header Definition
	response.Header().Set("content-type", "application/json")
	// validate Headers ....
	headers := request.Header
	
	val, ok = headers["rulesRequest"]
	if !ok || val != "yes" {  // this request can not be identified please return...
	       result := bson.M{
			"success": false,
			"message": "Invalid Request",
		}
		json.NewEncoder(response).Encode(result)
		return;
    }
	OrderMode:= constants.GLOBALMODE
	
	var payload map[string]interface{}
	var mode string
	trigger_type := "barrier_percentile_trigger"
	
	_ = json.NewDecoder(request.Body).Decode(&payload)


	//Changes interface int types to float64 type
	payload = helpers.AlterInterfaceTypesToFloat(payload)

	liveTradingBool := payload["enable_buy_barrier_percentile"]
	testTradingBool := payload["enable_test_buy_barrier_percentile"]

	switch {
    case OrderMode == "live":
        orderpickType := constants.LIVEORDERSTYPE
    case OrderMode == "test":
        orderpickType := constants.TESTORDERSTYPE
    case OrderMode == "both":
        orderpickType := constants.TESTANDLIVEBOTHORDERSTYPE   
    }
    // we will pick and execute orders on the basis of our filter results .....
    // if our rules live is true and our exchange also handles any one either that is live or both
    if(liveTradingBool == true && (orderpickType == 1 || orderpickType == 3)){
    	mode:="live"
    	buyOrderExecutionAckChannel := make(chan bool)
    	//GoRoutine
		go func() {
			buyOrderExecutionAckChannel <- helpers.ExecuteAutoBuyOrders(trigger_type, mode, payload)
		}()

		//Use Channel Value
		 <-buyOrderExecutionAckChannel
		//Close Channel
		close(buyOrderExecutionAckChannel)
    }
    // if our rules test is true and our exchange also handles any one either that is test or both
    if(testTradingBool == true && (orderpickType == 2 || orderpickType == 3)){
    	mode:="test"
    	buyOrderExecutionAckChannel := make(chan bool)
    	//GoRoutine
		go func() {
			buyOrderExecutionAckChannel <- helpers.ExecuteAutoBuyOrders(trigger_type, mode, payload)
		}()

		//Use Channel Value
		buyOrderExecutionAck := <-buyOrderExecutionAckChannel

		_ = buyOrderExecutionAck

		//Close Channel
		close(buyOrderExecutionAckChannel)
    	
    }
	
	result := bson.M{
		"success": true,
		"message": "Creating orders",
	}
	json.NewEncoder(response).Encode(result)

}