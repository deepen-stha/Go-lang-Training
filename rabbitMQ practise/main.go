package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load the env file : Error = %s", err.Error())
		return
	}
}
func main() {
	fmt.Println("Rabbit MQ tutorial")

	conn, err := amqp.Dial("amqp://" + os.Getenv("USER_NAME") + ":" + os.Getenv("PASSWORD") + "@10.20.30.25:5672/")
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:6072/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to our RabbitMq")

	ch,err := conn.Channel()
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err!=nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(q)

	err =ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:	[]byte("Hello World"),
		},
	)

	if err!=nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully Published Message to Queue")
}
