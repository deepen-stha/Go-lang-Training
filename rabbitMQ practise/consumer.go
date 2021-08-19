//this will consumes the messages that we have published in the queue
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
	fmt.Println("Consumer Application")
	conn, err := amqp.Dial("amqp://" + os.Getenv("USER_NAME") + ":" + os.Getenv("PASSWORD") + "@10.20.30.25:5672/")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {

		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Recived Message: %s\n", d.Body)
		}
	}()

	fmt.Println("successfully connected to the rabbit mq instance")
	fmt.Println("waiting for the messages")
	<-forever

}
