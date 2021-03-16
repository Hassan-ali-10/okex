package helpers

import (
	 "sync"
	 "time"
	 "context"
	 connectionhelper "okex/db"
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

func PickParentsAndMakeChilds( payload map[string]interface{},orders chan bool, wg *sync.WaitGroup){
	defer (*wg).Done()
	
	 pair, _ := payload["coinSymbol"].(string)
	// level, _ := payload["orderLevel"].(string)
	// openLimit, ok := payload["trade_limit"].(float64)

	// if !ok {
	// 	//fmt.Println("Payload Error: trade_limit for openLimit is not defined in payload.")
	// 	orders <- false
	// }
	currentMarketPrice := GetCurrentMarketPrice(pair)
	if currentMarketPrice == 0.0 {
		//fmt.Println("Current Market Price is 0 - Trading must be stopped")
		orders <- false

	}

}
//GetCurrentMarketPrice "Gets Current Market Price Of The Given pair"
func GetCurrentMarketPrice(pair string) float64 {
	var mongoResult map[string]interface{}
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
		"created_date": bson.M{
			"$gte": time5MinsBefore,
		},
	}

	marketPriceCursor, err := collection.Find(context.TODO(), searchFilter, findOpts)

	if err != nil {

		return 0.0

	}
	//Map result to slice
	marketPriceCursor.All(context.TODO(), &mongoResult)
	// once exhausted, close the cursor  it will free mongodb server resource ...
	marketPriceCursor.Close(context.TODO())
	
	
	
	
	if len(mongoResult) > 0 {
		price, _ := mongoResult["price"].(float64)
		return price
	} else {
		return 0.0
	}
}