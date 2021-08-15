package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	//imported gorrila mux to handle the router
	"github.com/gorilla/mux"
	//imported bson because the mongo db data is
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "packaging/rootHandle"
)

type Todo struct {
	Task string `json: "Field Task"`
	Done bool   `json: "Field Done"`
}

//declaring the struc to update the database
//this struc has only the task data
type UpdateData struct {
	UpdateTask string `json: "Field Task"`
}

func getAllData() ([]bson.M, error) {

	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientOptions type:", reflect.TypeOf(clientOptions))

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Access a MongoDB collection through a database
	col := client.Database("ToDo_Database").Collection("Todo_Collection")
	fmt.Println("Collection type:", reflect.TypeOf(col))

	// Call the collection's Find() method to return Cursor obj
	// with all of the col's documents
	cursor, err := col.Find(context.TODO(), bson.M{})

	// Find() method raised an error
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(ctx)
		return nil, err
		// If the API call was a success
	} else {

		var datas []bson.M
		if err = cursor.All(ctx, &datas); err != nil {
			log.Fatal(err)
		}
		fmt.Println(datas)
		return datas, nil

		// //returning the data as a response
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(episodes)
	}

}

//this function accepts response writer and the response of http
//the w means write and r means the read of the response
func rootHandleFunc(w http.ResponseWriter, r *http.Request) {
	// Declare host and port options to pass to the Connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientOptions type:", reflect.TypeOf(clientOptions))

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Access a MongoDB collection through a database
	col := client.Database("ToDo_Database").Collection("Todo_Collection")
	fmt.Println("Collection type:", reflect.TypeOf(col))

	// Call the collection's Find() method to return Cursor obj
	// with all of the col's documents
	cursor, err := col.Find(context.TODO(), bson.M{})

	// Find() method raised an error
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(ctx)
		// If the API call was a success
	} else {

		var episodes []bson.M
		if err = cursor.All(ctx, &episodes); err != nil {
			log.Fatal(err)
		}
		fmt.Println(episodes)

		//returning the data as a response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(episodes)

		// // iterate over docs using Next()
		// // declare a result BSON object
		// // var result bson.M
		// len := cursor.RemainingBatchLength()
		// // fmt.Println(reflect.TypeOf(cursor.RemainingBatchLength()))

		// var result [100]bson.M

		// count := 0
		// for cursor.Next(ctx) {

		// 	err := cursor.Decode(&result[count])

		// 	count += 1
		// 	// If there is a cursor.Decode error
		// 	if err != nil {
		// 		fmt.Println("cursor.Next() error:", err)
		// 		os.Exit(1)

		// 		// If there are no cursor.Decode errors
		// 	} else {
		// 		fmt.Println("\nresult type:", reflect.TypeOf(result))
		// 		fmt.Println("result:", result)
		// 		// result = append(result, result...)

		// 	}
		// }
		// // //returning the data as a response
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(result)

	}
}

//this is the function to handle the adding of the todo application
func addHandlerFunc(w http.ResponseWriter, r *http.Request) {

	//code to get the data from the request Body
	decoder := json.NewDecoder(r.Body)
	var t Todo
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Task)

	//now forming the data that we want to insert into our database
	newData := Todo{
		Task: t.Task,
		Done: t.Done,
	}

	//declaring host, options and port number to pass to the connect() method
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	fmt.Println("Client Type : ", reflect.TypeOf(clientOptions))

	//now connecting to the mongo database
	client, err := mongo.Connect(context.TODO(), clientOptions)

	//checking if the errrors occurs or not while connecting to the mongo db
	if err != nil {
		fmt.Println("Error while connecting to the mongoDB ", err)
		os.Exit(1)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	fmt.Println(reflect.TypeOf(ctx))

	//access a mongodb collection through a database
	col := client.Database("ToDo_Database").Collection("Todo_Collection")

	fmt.Println("the type of the collection is : ", reflect.TypeOf(col))

	fmt.Println("The type of the doc is :", reflect.TypeOf(newData))

	result, insert_err := col.InsertOne(ctx, newData)
	if insert_err != nil {
		fmt.Println("Insertion Error ", insert_err)
		os.Exit(1)
	} else {
		fmt.Println("InsertOne() type : ", reflect.TypeOf(result))
		fmt.Println(result)

		newID := result.InsertedID
		fmt.Println(newID)
		fmt.Println(reflect.TypeOf(newID))
	}

	//now getting all the data after doing addition of the new data
	data, _ := getAllData()
	//returning the data as a response of the api call
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

//this is the function to delete the todo list from the database
func deleteHandlerFunc(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL)

	//spliting the url to the particular id that we want to delete
	url := r.URL.String()
	split := strings.Split(url, "/")
	fmt.Println(split)
	fmt.Println(split[2])

	//now connecting to mongodb
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	//if error occurs while connecting to the database
	if err != nil {
		log.Fatal("This is an error")
		os.Exit(1)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	//getting all the columns form the database
	col := client.Database("ToDo_Database").Collection("Todo_Collection")

	//split of the 2 holds the data id which is unique for each of our data
	idPrimitive, err := primitive.ObjectIDFromHex(split[2])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		//calling the delete one method to delete the data
		res, err := col.DeleteOne(ctx, bson.M{"_id": idPrimitive})

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		} else {
			//checking if the response is nil
			if res.DeletedCount == 0 {
				fmt.Println("The data is not found")
			} else {
				fmt.Println("Deleteone Result : ", res)
			}
		}
	}

	//now getting all the data after doing deletion of the data
	data, _ := getAllData()
	//returning the data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

//this is the function to update the database
func updateHandlerFunc(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL)

	url := r.URL.String()
	id := strings.Split(url, "/")
	//getting the id to be updated
	fmt.Println(id[2])

	//now getting the data from the request body that we want to update
	//decoding the json Body
	//code to get the data(task to update) from the request Body
	decoder := json.NewDecoder(r.Body)
	var t UpdateData
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t.UpdateTask)

	//declaring the host and the port number to connect to the database
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	//Conenct to teh MOngoDB and return the client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//getting the mongodb col
	col := client.Database("ToDo_Database").Collection("Todo_Collection")

	objId, err := primitive.ObjectIDFromHex(id[2])

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//using the _id to get the specific document from the mongodb database
	filter := bson.M{"_id": bson.M{"$eq": objId}}

	//now updating the data
	update := bson.M{"$set": bson.M{"task": t.UpdateTask}}
	//calling the UpdateOne method and pass filter and update to it
	result, err := col.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	//checking for the error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Update One result :", result)
		fmt.Println("UpdateOne match count: ", result.MatchedCount)
		fmt.Println("UpdateOne modified count: ", result.ModifiedCount)
		fmt.Println("UpdateOne result upsertedId: ", result.UpsertedID)
	}

	//now getting all the data after doing updation of the data
	data, _ := getAllData()
	//returning the updated data as a response of the api call
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//this is the main function which is initiated at the beggining
func main() {

	//declaring the gorilla mux to make use of our routes
	r := mux.NewRouter()
	fmt.Println("Starting the application")

	//if it is the root path then at beggining it will only show the todo list data
	r.HandleFunc("/", rootHandleFunc)

	//this endpoint is called if user want to add the data in the mongo db
	// http.HandleFunc("/new", addHandlerFunc)
	r.HandleFunc("/new", addHandlerFunc)

	//this endpoint is called if user want to delete the todo list from the mongodb
	//this endpoint accepts the regular expression that is used to form the id that we want to delete
	r.HandleFunc("/delete/{id:[a-z0-9]+}", deleteHandlerFunc)

	//declaring the endpoint to update the data of the mongodb database
	r.HandleFunc("/update/{id:[a-z0-9]+}", updateHandlerFunc)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
