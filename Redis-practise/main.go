package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load the env file : Error = %s", err.Error())
		return
	}
}

func main() {
	fmt.Println("Running")
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CLOUD_ADDR"),
		Password: os.Getenv("PASSWORD"),
		DB:       0,
	})

	//forming the json value
	json, err := json.Marshal(Author{Name: "Deepen", Age: 20})
	if err != nil {
		fmt.Println(err)
	}

	//setting the value
	err = client.Set("id1234", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	val, err := client.Get("id1234").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}
