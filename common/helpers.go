package helpers

import (
	 "sync"
	 "time"
	 "fmt"
	 "context"
	  "encoding/json"
	 //"reflect"
	 connectionhelper "okex/db"
	 "go.mongodb.org/mongo-driver/bson/primitive"
	// //"go.mongodb.org/mongo-driver/mongo"
	 "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)
//AlterInterfaceTypes "Alters interface types for int types converted to float and wraps it up into a map again with same keys"
func AlterInterfaceTypesToFloat(mapp map[string]interface{}) map[string]interface{} {

	var alteredMapp = make(map[string]interface{})

	for key, value := range mapp {

		switch v := value.(type) {
		case nil:
		case int:
			alteredMapp[key] = float64(v)
		case int16:
			alteredMapp[key] = float64(v)
		case int32:
			alteredMapp[key] = float64(v)
		case int64:
			alteredMapp[key] = float64(v)
		default:
			alteredMapp[key] = v
		}
	}

	return alteredMapp

}

func PickParentsAndMakeChilds( payload map[string]interface{},orderpickType string,orders chan bool, wg *sync.WaitGroup){
	defer (*wg).Done()
	
	 coin, _ := payload["coinSymbol"].(string)
	level, _ := payload["orderLevel"].(string)
	openLimit, ok := payload["trade_limit"].(float64)
	
	if !ok {
		fmt.Println("Payload Error: trade_limit for openLimit is not defined in payload.")
		orders <- false
	}
	currentMarketPrice := GetCurrentMarketPrice(coin)
	fmt.Println(currentMarketPrice)
	if currentMarketPrice == 0.0 {
		//fmt.Println("Current Market Price is 0 - Trading must be stopped")
		orders <- false // returning channels 

	}
	switch {
    case orderpickType == "live":
   		go RunChildOrders("live",currentMarketPrice,openLimit,coin,level,payload);
   		orders <- true // returning channels 
    case orderpickType == "test":
        go RunChildOrders("test",currentMarketPrice,openLimit,coin,level,payload);
        orders <- true // returning channels 
    case orderpickType == "both":
    	go RunChildOrders("live",currentMarketPrice,openLimit,coin,level,payload);
    	go RunChildOrders("test",currentMarketPrice,openLimit,coin,level,payload);
    	orders <- true  // returning channels 
    default :
    	orders <- false // returning channels 
    }

}

func RunChildOrders(mode string,price float64,limit float64,coin string,level string,data map[string]interface{}){
	ordersData:=ListParentActiveOrders(coin,level,mode,"barrier_percentile_trigger",limit)
	fmt.Println("coming here",ordersData)
}

//ListParentActiveOrders "Pick Active Parents"
func ListParentActiveOrders(coin string, level string, mode string, trigger_type string, limit float64) []primitive.M {
	
	var mongoResult []bson.M
	
	
	searchFilter := bson.M{
			"symbol": coin,
			"order_mode": mode,
			"parent_status": "parent",
			"order_level": level,
			"status": "new",
			"trigger_type": trigger_type,
			"pause_status": "play",
			"pick_parent": "yes",
			//"admin_id":      "5c0912b7fc9aadaac61dd072", //for testing only

		}
	 prettyJSON, err := json.MarshalIndent(searchFilter, "", "    ")
    if err != nil {
        //log.Fatal("Failed to generate json", err)
    }
    fmt.Printf("%s\n", string(prettyJSON))	
    limitInt := int64(limit)
	findOpts := options.Find().SetLimit(limitInt)	
	collectionName:="BuyOrders_okex" 
	collection, err := connectionhelper.GetMongoDbCollection(collectionName)
	if err != nil {
		return mongoResult
	}
	ordersCursor, err := collection.Find(context.TODO(), searchFilter, findOpts)

	if err != nil {

		return mongoResult

	}
	//Map result to slice
	ordersCursor.All(context.TODO(), &mongoResult)
	// once exhausted, close the cursor  it will free mongodb server resource ...
	ordersCursor.Close(context.TODO())
	fmt.Println(len(mongoResult))
	return mongoResult
	

}
//GetCurrentMarketPrice "Gets Current Market Price Of The Given pair"
func GetCurrentMarketPrice(pair string) float64 {
	fmt.Println("bvasgasasgwq647g",pair)
	
	var mongoResult []bson.M
	collectionName:="MarketPrices" 
	collection, err := connectionhelper.GetMongoDbCollection(collectionName)
	if err != nil {
		return 0.0
	}
	time5MinsBefore := time.Now()
	time5MinsBefore = time5MinsBefore.Add(-60 * time.Minute)
	//Find Options Like Sort And Limit Etc
	findOpts := options.Find()
	findOpts.SetLimit(1)
	findOpts.SetSort(bson.M{"_id": -1})

	searchFilter := bson.M{
		"coin": pair,
		// "created_date": bson.M{
		// 	"$gte": time5MinsBefore,
		// },
	}
	 // prettyJSON, err := json.MarshalIndent(searchFilter, "", "    ")
  //   if err != nil {
  //       //log.Fatal("Failed to generate json", err)
  //   }
  //   fmt.Printf("%s\n", string(prettyJSON))
	marketPriceCursor, err := collection.Find(context.TODO(), searchFilter, findOpts)

	if err != nil {

		return 0.0

	}

	//Map result to slice
	marketPriceCursor.All(context.TODO(), &mongoResult)
	// once exhausted, close the cursor  it will free mongodb server resource ...
	marketPriceCursor.Close(context.TODO())
	
	fmt.Println(mongoResult)
	
	
	if len(mongoResult) > 0 {
		price, _ := mongoResult[0]["price"].(float64)
		return price
	} else {
		return 0.0
	}
}