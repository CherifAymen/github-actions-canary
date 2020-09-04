package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Data : the format of saving objects in DB
type Data struct {
	ID       primitive.ObjectID `bson:"_id"`
	Sentence string             `bson:"sentence"`
	Polarity float32            `bson:"polarity"`
}

//Sentiment : A struct that represents our json object received from Flusk Server
type Sentiment struct {
	Sentence string  `json:"sentence"`
	Polarity float32 `json:"polarity"`
	Version  string  `json:"version"`
}

var collection *mongo.Collection
var ctx = context.TODO()
var uri string

func init() {
	envVar := "DB"
	if len(os.Getenv("DB_ENV")) > 0 {
		envVar = os.Getenv("DB_ENV")
	}
	db := os.Getenv(envVar)
	if len(db) == 0 {
		uri = "mongodb://localhost:27017/"
	} else {
		uri = "mongodb://" + db + ":27017/"
	}
	log.Println(uri)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("AS").Collection("Data")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")

	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleRequests() {
	fmt.Println("* Running on http://localhost:8080/ (Press CTRL+C to quit)")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/testHealth", HealthHandler)
	myRouter.HandleFunc("/data", allData)
	myRouter.HandleFunc("/sentiment", sentimentHandler).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func sentimentHandler(w http.ResponseWriter, r *http.Request) {
	saLogicAPIURL := os.Getenv("URL")
	if saLogicAPIURL == "" {
		saLogicAPIURL = "http://localhost:5000"
	}
	fullURL := saLogicAPIURL + "/analyse/sentiment"

	enableCors(&w)

	// get the body of our POST request and send it to saLogic
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if len(reqBody) != 0 {
		//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		t := time.Now()
		clock := fmt.Sprintf("[%d-%02d-%02d  %02d:%02d:%02d]",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(req.Host+" -- "+clock+" \""+req.Method, r.URL.Path, " ", resp.StatusCode, " - ")

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		var sentiment Sentiment

		json.Unmarshal(body, &sentiment)
		sentiment.Version = "go"
		data := &Data{
			ID:       primitive.NewObjectID(),
			Sentence: sentiment.Sentence,
			Polarity: sentiment.Polarity,
		}
		newData(data)
		json.NewEncoder(w).Encode(sentiment)
	}
}

func allData(w http.ResponseWriter, r *http.Request) {

	// passing bson.D{{}} matches all documents in the collection
	filter := primitive.D{{}}

	datas, err := filterTasks(filter)

	if err != nil {
		fmt.Fprintf(w, "No Data found")
		return
	}

	json.NewEncoder(w).Encode(datas)
}

func newData(data *Data) *mongo.InsertOneResult {
	id, _ := collection.InsertOne(ctx, data)
	return id
}

func filterTasks(filter interface{}) ([]*Data, error) {
	// A slice of tasks for storing the decoded documents
	var datas []*Data

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return datas, err
	}

	for cur.Next(ctx) {
		var d Data
		err := cur.Decode(&d)
		if err != nil {
			return datas, err
		}

		datas = append(datas, &d)
	}

	if err := cur.Err(); err != nil {
		return datas, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(datas) == 0 {
		return datas, mongo.ErrNoDocuments
	}

	return datas, nil
}

// HealthHandler returns a succesful status and a message.
// For use by Consul or other processes that need to verify service health.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, you've hit %s\n", r.URL.Path)
}

func main() {
	handleRequests()
}
