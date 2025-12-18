package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	razp "razpravljalnica/razpravljalnica"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr         = flag.String("addr", "localhost:50051", "the address to connect to")
	userId int64 = -1
)

func main() {
	flag.Parse()
	fmt.Printf("gRPC client connecting to %v\n", *addr)
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Print("Write name and pass: ")
	var name, pass string
	fmt.Scan(&name, &pass)
	fmt.Printf("Name: %s Pass: %s\n", name, pass)

	grpcClient := razp.NewMessageBoardClient(conn)

	ctx, cancel := context.WithCancel(context.Background()) //idk kak deluje ctx tocn tak da to mal poglej
	defer cancel()
	newUser, err := grpcClient.CreateUser(ctx, &razp.CreateUserRequest{Name: name, Password: pass})
	fmt.Printf("Username: %s, Id: %d\n", newUser.Name, newUser.Id)

	newTopic, err := grpcClient.CreateTopic(ctx, &razp.CreateTopicRequest{Name: "TestTopic"})
	fmt.Printf("TopicName: %s, Id: %d\n", newTopic.Name, newTopic.Id)

	newTopic, err = grpcClient.CreateTopic(ctx, &razp.CreateTopicRequest{Name: "TestTopic2"})
	fmt.Printf("TopicName: %s, Id: %d\n", newTopic.Name, newTopic.Id)

	newTopic, err = grpcClient.CreateTopic(ctx, &razp.CreateTopicRequest{Name: "TestTopic3"})
	fmt.Printf("TopicName: %s, Id: %d\n", newTopic.Name, newTopic.Id)

	newMessage, err := grpcClient.PostMessage(ctx, &razp.PostMessageRequest{TopicId: 1, UserId: 1, Text: "Test"})
	fmt.Printf("NewMessage %s\n", newMessage.Text)

	newMessage, err = grpcClient.PostMessage(ctx, &razp.PostMessageRequest{TopicId: 1, UserId: 1, Text: "Test2"})
	fmt.Printf("NewMessage %s\n", newMessage.Text)

	newMessage, err = grpcClient.PostMessage(ctx, &razp.PostMessageRequest{TopicId: 1, UserId: 1, Text: "Test3"})
	fmt.Printf("NewMessage %s\n", newMessage.Text)

	newMessagesList, err := grpcClient.GetMessages(ctx, &razp.GetMessagesRequest{TopicId: 1, FromMessageId: 1, Limit: 3})

	for _, message := range newMessagesList.Messages {
		fmt.Printf("%v\n", message)
	}

	newMessage, err = grpcClient.UpdateMessage(ctx, &razp.UpdateMessageRequest{TopicId: 1, UserId: 1, MessageId: 1, Text: "Edited Message"})

	if err != nil {
		fmt.Print(err)
	}

	_, err = grpcClient.DeleteMessage(ctx, &razp.DeleteMessageRequest{TopicId: 1, UserId: 1, MessageId: 1})

	newMessage, err = grpcClient.LikeMessage(ctx, &razp.LikeMessageRequest{TopicId: 1, UserId: 1, MessageId: 2})
	// _, err = grpcClient.LikeMessage(ctx, &razp.LikeMessageRequest{TopicId: 1, UserId: 1, MessageId: 3})

	newMessagesList, err = grpcClient.GetMessages(ctx, &razp.GetMessagesRequest{TopicId: 1, FromMessageId: 1, Limit: 3})

	for _, message := range newMessagesList.Messages {
		// if message.TopicId == -1 {
		// 	fmt.Print("DELETED MESSAGE")
		// } else {
		fmt.Printf("%v with %d likes\n", message, message.Likes)
		// }
	}

	stream, err := grpcClient.SubscribeTopic(context.Background(),
		&razp.SubscribeTopicRequest{TopicId: 1, UserId: newUser.Id},
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			ev, err := stream.Recv()
			if err != nil {
				fmt.Println("subscription ended:", err)
				return
			}
			fmt.Printf("[LIVE] topic=%d: %s\n", ev.Message.TopicId, ev.Message.Text)
		}
	}()

	fmt.Println("Subscribed. Press ENTER to quit.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
