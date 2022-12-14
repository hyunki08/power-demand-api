package db

import (
	"context"
	"fmt"
	"strings"

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

func (pd *powerDemand) FindByPipeline(pipeline mongo.Pipeline) []map[string]interface{} {
	if !pd.checkStatus() {
		panic("DB Connection status isbad")
	}

	cursor, err := pd.coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func (pd *powerDemand) FindDemandedYearly(year string) []map[string]interface{} {
	matchStage := bson.D{{"$match", bson.D{{"date", bson.M{"$gte": strings.Join([]string{year, "00", "00"}, "-"), "$lte": strings.Join([]string{year, "99", "99"}, "-")}}}}}
	setStage := bson.D{
		{"$set", bson.D{
			{"sum", bson.D{{"$add", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}},
			{"avg", bson.D{{"$avg", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}}},
		},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", ""},
			{"avg_year", bson.D{{"$avg", "$avg"}}},
			{"sum_year", bson.D{{"$sum", "$sum"}}},
			{"max_year", bson.D{{"$max", "$sum"}}},
			{"min_year", bson.D{{"$min", "$sum"}}}},
		},
	}
	unsetStage := bson.D{{"$unset", bson.A{"_id"}}}

	return pd.FindByPipeline(mongo.Pipeline{matchStage, setStage, groupStage, unsetStage})
}

func (pd *powerDemand) FindDemandedMonthly(year string, month string) []map[string]interface{} {
	matchStage := bson.D{{"$match", bson.D{{"date", bson.M{"$gte": strings.Join([]string{year, month, "00"}, "-"), "$lte": strings.Join([]string{year, month, "99"}, "-")}}}}}
	setStage := bson.D{
		{"$set", bson.D{
			{"_id", 0},
			{"sum", bson.D{{"$add", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}},
			{"avg", bson.D{{"$avg", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}}},
		},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", ""},
			{"avg_month", bson.D{{"$avg", "$avg"}}},
			{"sum_month", bson.D{{"$sum", "$sum"}}},
			{"max_month", bson.D{{"$max", "$sum"}}},
			{"min_month", bson.D{{"$min", "$sum"}}}},
		},
	}
	unsetStage := bson.D{{"$unset", bson.A{"_id"}}}

	return pd.FindByPipeline(mongo.Pipeline{matchStage, setStage, groupStage, unsetStage})
}

func (pd *powerDemand) FindDemandedDailyByRange(from string, to string) []map[string]interface{} {
	matchStage := bson.D{{"$match", bson.D{{"date", bson.M{"$gte": from, "$lte": to}}}}}
	setStage := bson.D{
		{"$set", bson.D{
			{"sum", bson.D{{"$add", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}},
			{"avg", bson.D{{"$avg", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}}},
		},
	}
	unsetStage := bson.D{{"$unset", bson.A{"_id"}}}

	return pd.FindByPipeline(mongo.Pipeline{matchStage, setStage, unsetStage})
}

func (pd *powerDemand) FindDemandedDaily(date string) []map[string]interface{} {
	matchStage := bson.D{{"$match", bson.D{{"date", date}}}}
	setStage := bson.D{
		{"$set", bson.D{
			{"sum", bson.D{{"$add", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}},
			{"avg", bson.D{{"$avg", bson.A{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8", "$9", "$10", "$11", "$12", "$13", "$14", "$15", "$16", "$17", "$18", "$19", "$20", "$21", "$22", "$23", "$24"}}}}},
		},
	}
	unsetStage := bson.D{{"$unset", bson.A{"_id"}}}

	return pd.FindByPipeline(mongo.Pipeline{matchStage, setStage, unsetStage})
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
