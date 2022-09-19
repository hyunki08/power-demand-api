package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type powerDemand struct {
	coll     *mongo.Collection
	model    mongo.IndexModel
	hasModel bool
	Meta     meta
}

type meta struct {
	MinDate string
	MaxDate string
}

const (
	DBName             string = "PDDB"
	collectionName     string = "PowerDemand"
	metaCollectionName string = "PowerDemandMeta"
)

func (pd *powerDemand) setMetaCollection() {
	if !pd.checkStatus() {
		panic("DB Connection status isbad")
	}

	collections, err := client.Database(DBName).ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	hasMetaColl := false
	for _, coll := range collections {
		if coll == (metaCollectionName) {
			hasMetaColl = true
		}
	}

	metaColl := client.Database(DBName).Collection(metaCollectionName)

	if !hasMetaColl {
		err := client.Database(DBName).CreateCollection(context.TODO(), metaCollectionName)
		if err != nil {
			panic(err)
		}

		projection := bson.D{{"_id", 0}, {"date", 1}}

		sort := bson.D{{"date", 1}}
		opts := options.FindOne().SetSort(sort).SetProjection(projection)
		result := pd.FindOne("", opts)
		minDate := result["date"]

		sort = bson.D{{"date", -1}}
		opts = options.FindOne().SetSort(sort).SetProjection(projection)
		result = pd.FindOne("", opts)
		maxDate := result["date"]

		_, err = metaColl.InsertOne(context.TODO(), bson.D{{"minDate", minDate}, {"maxDate", maxDate}})
		if err != nil {
			panic(err)
		}
	}

	var result map[string]interface{}
	err = metaColl.FindOne(context.TODO(), bson.D{}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No document was found")
	}
	if err != nil {
		panic(err)
	}
	pd.Meta.MaxDate = fmt.Sprintf("%s", result["maxDate"])
	pd.Meta.MinDate = fmt.Sprintf("%s", result["minDate"])
}

func (pd *powerDemand) FindOneByDate(date string) map[string]interface{} {
	return pd.FindOne(date, nil)
}

func (pd *powerDemand) FindOne(date string, opts *options.FindOneOptions) map[string]interface{} {
	if !pd.checkStatus() {
		panic("DB Connection status isbad")
	}

	if opts == nil {
		opts = options.FindOne()
	}

	filter := bson.D{}
	if date != "" {
		filter = bson.D{{"date", date}}
	}

	var result map[string]interface{}
	err := pd.coll.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No document was found")
		return result
	}
	if err != nil {
		panic(err)
	}
	return result
}

func (pd *powerDemand) Find(from string, to string) []map[string]interface{} {
	if !pd.checkStatus() {
		panic("DB Connection status isbad")
	}

	sort := bson.D{{"date", 1}}
	opts := options.Find().SetSort(sort)
	filter := bson.D{{"date", bson.M{"$gte": from, "$lte": to}}}

	var results []map[string]interface{}
	cursor, fErr := pd.coll.Find(context.TODO(), filter, opts)
	if fErr == mongo.ErrNoDocuments {
		fmt.Println("No document was found")
		return results
	}
	if fErr != nil {
		panic(fErr)
	}

	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

/*==== checkStatus Function ====*/
func (pd *powerDemand) checkStatus() bool {
	if !IsConnected() {
		return false
	}
	pd.isSetCollection()
	pd.isSetModel()

	return true
}

func (pd *powerDemand) isSetCollection() bool {
	if pd.coll == nil {
		pd.coll = client.Database(DBName).Collection(collectionName)
	}
	return true
}

func (pd *powerDemand) isSetModel() bool {
	if pd.hasModel == false {
		pd.model = mongo.IndexModel{Keys: bson.D{{"date", "text"}}}
		_, err := pd.coll.Indexes().CreateOne(context.TODO(), pd.model)
		if err != nil {
			panic(err)
		}
		pd.hasModel = true
	}

	return true
}
