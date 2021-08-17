package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"
	"os"
	"reflect"
	"strings"
	"time"

	//imported gorrila mux to handle the router
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	//imported bson because the mongo db data is
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "packaging/rootHandle"

	"github.com/joho/godotenv"
)

//declaring the global variables to read the environment variables
var dataBaseURL string
var dataBaseName string
var dataBaseCollection string

//declaring the logger globally so that it can be used through out the program
var logger *zap.Logger

//declaring all the variables globally that are required to connect to our mongodb
var clientOptions *options.ClientOptions
var client *mongo.Client
var col *mongo.Collection

type Todo struct {
	Task string `json: "Field Task"`
	Done bool   `json: "Field Done"`
}

//declaring the struc to update the database
//this struc has only the task data
type UpdateData struct {
	UpdateTask string `json: "Field Task"`
}

//this is the function to get all the data
func getAllData() ([]bson.M, error) {

	//function to assign all the global variable value required for the mongo db connections
	assignMongoDBValue()

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second) //don't ignore the error

	// Call the collection's Find() method to return Cursor obj
	// with all of the col's documents
	cursor, err := col.Find(context.TODO(), bson.M{})

	// if Find() method raised an error
	if err != nil {
		logger.Error(
			"Error while getting the data",
			zap.Error(err))
		// fmt.Println("Finding all documents ERROR:", err)
		cursor.Close(ctx)
		return nil, err //if error occurs returning the nil data and the err msg

	}
	var datas []bson.M
	if err = cursor.All(ctx, &datas); err != nil {
		// log.Fatal(err)
		logger.Error(
			"Error while fetching all the data",
			zap.Error(err))
		cursor.Close(ctx)
		return nil, err
	}
	// //now disconnecting from the mongodb database
	// err := client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(datas)
	cursor.Close(ctx)

	return datas, nil
}

//this function accepts response writer and the response of http
//the w means write and r means the read of the response
func rootHandleFunc(w http.ResponseWriter, r *http.Request) {

	//function to assign all the global variable value required for the mongo db connections
	assignMongoDBValue()

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second) //do not neglect the error

	// Call the collection's Find() method to return Cursor obj
	// with all of the col's documents
	cursor, err := col.Find(context.TODO(), bson.M{})

	// Find() method raised an error
	if err != nil {
		logger.Error("Error to get all the documents", zap.Error(err))
		// fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(ctx)
		fmt.Println("Connection to MongoDB closed.")
		//returning error if the database is not found
		http.Error(w, "Cannot get the data", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var datas []bson.M
	if err = cursor.All(ctx, &datas); err != nil {
		logger.Error("Eror to get all the data", zap.Error(err))
		cursor.Close(ctx)
		return
	}
	// fmt.Println(episodes)

	logger.Info("Success..",
		zap.Int64("successful get request ", 200))

	//now disconnecting from the mongodb database
	// err := client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cursor.Close(ctx)
	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
	//returning the data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datas)

}

//this is the function to handle the adding of the todo application
func addHandlerFunc(w http.ResponseWriter, r *http.Request) {

	//function to assign all the global variable value required for the mongo db connections
	assignMongoDBValue()

	//code to get the data from the request Body
	decoder := json.NewDecoder(r.Body)
	var t Todo
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error("Error while decoding the body", zap.Error(err))
		//returning error if the database is not found
		http.Error(w, "Error in the id that you have provided", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
		// panic(err)
	}
	// log.Println(t.Task)

	//now forming the data that we want to insert into our database
	newData := Todo{
		Task: t.Task,
		Done: t.Done,
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	// fmt.Println(reflect.TypeOf(ctx))

	result, insert_err := col.InsertOne(ctx, newData)
	if insert_err != nil {
		logger.Error("Error while inserting data to the database", zap.Error(insert_err))
		//returning error if the database is not found
		http.Error(w, "Error cannot enter data into the database", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	newID := result.InsertedID
	fmt.Println(newID)
	// fmt.Println(reflect.TypeOf(newID))
	logger.Info("Successfully added new data..")

	// //now disconnecting to the mongodb database
	// err := client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//now getting all the data after doing addition of the new data
	data, _ := getAllData()
	//returning the data as a response of the api call
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

//this is the function to delete the todo list from the database
func deleteHandlerFunc(w http.ResponseWriter, r *http.Request) {

	//function to assign all the global variable value required for the mongo db connections
	assignMongoDBValue()

	// log.Println(r.URL)
	//spliting the url to the particular id that we want to delete
	url := r.URL.String()
	split := strings.Split(url, "/")

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	//split of the 2 holds the data id which is unique for each of our data
	idPrimitive, err := primitive.ObjectIDFromHex(split[2])
	if err != nil {
		logger.Error("Error while getting the id from the database", zap.Error(err))
		//returning error if the database is not found
		http.Error(w, "Error. Cannot find the id", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	//calling the delete one method to delete the data
	res, err := col.DeleteOne(ctx, bson.M{"_id": idPrimitive})

	if err != nil {
		logger.Error("Error while deleting data from the database", zap.Error(err))
		//returning error if the database is not found
		http.Error(w, "Error. Cannot delete data from the database", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	//checking if the response is nil
	if res.DeletedCount == 0 {
		logger.Info("The data is not found")
		return
		// fmt.Println("The data is not found")
	}
	logger.Info("The data is deleted successfully",
		zap.String("The method called is :", r.Method))

	// //now disconnecting the mongodb database
	// err := client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Deleteone Result : ", res)

	//now getting all the data after doing deletion of the data
	data, _ := getAllData()
	//returning the data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

//this is the function to update the database
func updateHandlerFunc(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(r.URL)

	//function to assign all the global variable value required for the mongo db connections
	assignMongoDBValue()

	url := r.URL.String()
	id := strings.Split(url, "/")
	//getting the id to be updated
	// fmt.Println(id[2])

	//now getting the data from the request body that we want to update
	//decoding the json Body
	//code to get the data(task to update) from the request Body
	decoder := json.NewDecoder(r.Body)
	var t UpdateData
	err := decoder.Decode(&t)
	if err != nil {
		logger.Error("Error while decoding the request body", zap.Error(err))
		// panic(err)
		//returning error if the database is not found
		http.Error(w, "Cannot decode the id ", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	objId, err := primitive.ObjectIDFromHex(id[2])

	if err != nil {
		logger.Error("Error while getting data from the database", zap.Error(err))
		//returning error if the database is not found
		http.Error(w, "Error while getting data from the database", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
		//returning error if the database is not found
		http.Error(w, "Cannot update the data", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// r.Body.Close()
		return
	}
	//using logger info to log the information
	logger.Info("Success..",
		zap.Int64("Update One result: ", result.MatchedCount),
		zap.Int64("Updated modified count", result.ModifiedCount))

	// //now disconnecting to the mongodb database
	// err := client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//now getting all the data after doing updation of the data
	data, _ := getAllData()
	//returning the updated data as a response of the api call
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//function to get the particular data
func getHandlerFunction(w http.ResponseWriter, r *http.Request) {

	//function to assign all the global variable value required for the mongo db connections
	assignMongoDBValue()

	//getting the particular id from the url
	url := r.URL.String()
	id := strings.Split(url, "/")
	//getting the id to be updated
	// fmt.Println(id[2])

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second) //do not neglect the error

	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	docID, err := primitive.ObjectIDFromHex(id[2])
	if err != nil {
		logger.Error("Error ", zap.Error(err))
		//returning error if the database is not found
		http.Error(w, "Cannot get the id", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Declare a struct instance of the MongoDB fields that will contain the document returned
	var result Todo

	// call the collection's Find() method and return Cursor object into result
	// fmt.Println(`bson.M{"_id": docID}:`, bson.M{"_id": docID})
	err = col.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)

	// Check for any errors returned by MongoDB FindOne() method call
	if err != nil {
		fmt.Println("FindOne() ObjectIDFromHex ERROR:", err)
		//returning error if the database is not found
		http.Error(w, "Cannot find the data of particular id", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	logger.Info("Successfly find the documents: ",
		zap.String("result task ", result.Task),
		zap.Bool("result done ", result.Done))

	logger.Info("Success..",
		zap.Int64("successful get request ", 200))

	// //now disconnecting from the mongodb database
	// err := client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// TODO optional you can log your closed MongoDB client
	// fmt.Println("Connection to MongoDB closed.")
	//returning the data as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

//function to initialize the logger
func InitLogger() {
	logger, _ = zap.NewProduction()
}

//function to assign the value to the global mongo db variables
func assignMongoDBValue() {

	//assigning value to the global variable required for conencting to the mongo db
	// Declare host and port options to pass to the Connect() method
	clientOptions = options.Client().ApplyURI(dataBaseURL)
	fmt.Println("clientOptions type:", reflect.TypeOf(clientOptions))
	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	fmt.Println(reflect.TypeOf(client))
	if err != nil {
		fmt.Println("error occurs")
	}
	// Access a MongoDB collection through a database
	col = client.Database(dataBaseName).Collection(dataBaseCollection)
}

//this is the main function which is initiated at the beggining
func main() {

	//we should always load the .env file at first
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")

	if err != nil {
		// log.Fatalf("Error loading .env file")
		logger.Error(
			"Error loading the .env file",
			zap.Error(err))
	}
	// getting env variables
	dataBaseURL = os.Getenv("DB_URL")
	dataBaseName = os.Getenv("DB_NAME")
	dataBaseCollection = os.Getenv("DB_COLLECTION")

	fmt.Printf("godotenv : %s = %s \n", "Site Title", reflect.TypeOf(dataBaseURL))
	fmt.Printf("godotenv : %s = %s \n", "DB Host", dataBaseCollection)

	//calling the initlogger function to initialize the logger
	InitLogger()
	defer logger.Sync()

	//declaring the gorilla mux to make use of our routes
	r := mux.NewRouter()
	logger.Info("Starting the application")

	//if it is the route to get all the documnets
	r.HandleFunc("/todo", rootHandleFunc).Methods("GET")

	//route to get the particular document by id
	r.HandleFunc("/todo/{id:[a-z0-9]+}", getHandlerFunction).Methods("GET")

	//this endpoint is called if user want to add the data in the mongo db
	// http.HandleFunc("/new", addHandlerFunc)
	//post handles the
	r.HandleFunc("/todo", addHandlerFunc).Methods("POST")

	//add the get function with id
	//this endpoint is called if user want to delete the todo list from the mongodb
	//this endpoint accepts the regular expression that is used to form the id that we want to delete
	r.HandleFunc("/todo/{id:[a-z0-9]+}", deleteHandlerFunc).Methods("DELETE")

	//declaring the endpoint to update the data of the mongodb database
	r.HandleFunc("/todo/{id:[a-z0-9]+}", updateHandlerFunc).Methods("PUT")

	http.Handle("/", r)
	//log the error if the server is already in use
	log.Fatal(http.ListenAndServe(":8080", nil))
}
