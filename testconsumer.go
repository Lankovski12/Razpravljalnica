package main

import (
	"context"
	"flag"
	"fmt"
	razp "razpravljalnica/razpravljalnica"
	"time"

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

	grpcClient := razp.NewMessageBoardClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second) //idk kak deluje ctx tocn tak da to mal poglej
	defer cancel()

	newUser, err := grpcClient.CreateUser(ctx, &razp.CreateUserRequest{Name: "Neza", Password: "Test"})
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
}
