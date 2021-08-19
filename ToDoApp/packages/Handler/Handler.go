package Handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"test3/packages/Initializer"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var col *mongo.Collection
var sugarLogger *zap.SugaredLogger

//funcition to initialize the logger it will run at first
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		sugarLogger.Errorf("Failed to load the env file : Error = %s", err.Error())
		return
	}
	sugarLogger = Initializer.InitLogger()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		sugarLogger.Errorf("Failed to connect to mongodb = %s", err.Error())
		return
	}
	if err != nil {
		sugarLogger.Errorf("Failed to ping to mongodb = %s", err.Error())
		return
	}
	col = client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
}

type Todo struct {
	Task string `json: "Field Task"`
	Done bool   `json: "Field Done"`
}

type UpdateData struct {
	UpdateTask string `json: "Field Task"`
}

//this is the function to get all the data
func getAllData() ([]bson.M, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	cursor, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		cursor.Close(ctx)
		return nil, err
	}
	var datas []bson.M
	if err = cursor.All(ctx, &datas); err != nil {
		cursor.Close(ctx)
		return nil, err
	}
	cursor.Close(ctx)
	return datas, nil
}

//this function accepts response writer and the response of http. the w means write and r means the read of the request
func RootHandleFunc(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second) //do not neglect the error
	cursor, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		defer cursor.Close(ctx)
		http.Error(w, "Cannot get the data", http.StatusInternalServerError)
		return
	}
	var datas []bson.M
	if err = cursor.All(ctx, &datas); err != nil {
		cursor.Close(ctx)
		return
	}
	cursor.Close(ctx)
	fmt.Println("Connection to MongoDB closed.")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datas)
}

//this is the function to handle the adding of the todo application
func AddHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//code to get the data from the request Body
	decoder := json.NewDecoder(r.Body)
	var t Todo
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "Error in the id that you have provided", http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//now forming the data that we want to insert into our database
	newData := Todo{
		Task: t.Task,
		Done: t.Done,
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	result, insert_err := col.InsertOne(ctx, newData)
	if insert_err != nil {
		http.Error(w, "Error cannot enter data into the database", http.StatusInternalServerError)
		return
	}
	newID := result.InsertedID
	fmt.Println(newID)
	data, _ := getAllData()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//this is the function to delete the todo list from the database
func DeleteHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var id = getID(r.URL.String())
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Error. Cannot find the id", http.StatusInternalServerError)
		return
	}
	res, err := col.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	if err != nil {
		http.Error(w, "Error. Cannot delete data from the database", http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		return
	}
	data, _ := getAllData()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//this is the function to update the database
func UpdateHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var id = getID(r.URL.String())
	decoder := json.NewDecoder(r.Body)
	var t UpdateData
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "Cannot decode the id ", http.StatusInternalServerError)
		return
	}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Error while getting data from the database", http.StatusInternalServerError)
		return
	}
	//using the _id to get the specific document from the mongodb database
	filter := bson.M{"_id": bson.M{"$eq": objId}}
	update := bson.M{"$set": bson.M{"task": t.UpdateTask}}
	//calling the UpdateOne method and pass filter and update to it
	result, err := col.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	fmt.Println(result)
	if err != nil {
		http.Error(w, "Cannot update the data", http.StatusInternalServerError)
		return
	}
	data, _ := getAllData()
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//function to get the particular data
func GetHandlerFunction(w http.ResponseWriter, r *http.Request) {
	var id = getID(r.URL.String())
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Cannot get the id", http.StatusInternalServerError)
		return
	}
	var result Todo
	err = col.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		http.Error(w, "Cannot find the data of particular id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

//function to return id out of the url
func getID(url string) string {
	id := strings.Split(url, "/")
	return id[2]
}
