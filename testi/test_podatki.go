package main

import (
	"context"
	"fmt"
	"log"
	"time"

	razp "razpravljalnica/razpravljalnica"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := razp.NewMessageBoardClient(conn)
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users := []struct {
		name     string
		password string
	}{
		{"Alice", "pass123"},
		{"Bob", "pass456"},
		{"Charlie", "pass789"},
		{"Diana", "pass000"},
		{"Eve", "pass111"},
		{"Adam", "pass666"},
		{"Mike", "pass3456"},
		{"Loti", "pass98765"},
		{"Rick", "pass7667"},
		{"Butter", "pass87654"},
	}

	userIDs := make(map[string]int64)
	for _, u := range users {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		user, err := client.CreateUser(ctx, &razp.CreateUserRequest{
			Name:     u.name,
			Password: u.password,
		})
		cancel()
		if err != nil {
			log.Printf("❌ Failed to create user %s: %v", u.name, err)
			continue
		}
		userIDs[u.name] = user.Id
		fmt.Printf("✅ Created user: %s (ID: %d)\n", u.name, user.Id)
	}

	topics := []string{
		"Programing",
		"Games",
		"Phones",
		"Travel",
		"Food",
		"Animals",
		"Movies",
		"Funny",
		"Hobbies",
	}

	topicIDs := make(map[string]int64)
	for _, topicName := range topics {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		topic, err := client.CreateTopic(ctx, &razp.CreateTopicRequest{
			Name: topicName,
		})
		cancel()
		if err != nil {
			log.Printf("❌ Failed to create topic %s: %v", topicName, err)
			continue
		}
		topicIDs[topicName] = topic.Id
		fmt.Printf("✅ Created topic: %s (ID: %d)\n", topicName, topic.Id)
	}

	// Programiranje topic
	messages := map[string][]string{
		"Programing": {
			"[Go Tutorial] Kako se naučiti Go programskega jezika",
			"[Python] Best practices za Python razvoj",
			"[JavaScript] Kaj je novo v ES2024",
			"[Rust] Performance comparison: Rust vs C++",
			"[Database] SQL optimization tips",
		},
		"Games": {
			"[Gaming News] Nova igra leta",
			"[Elden Ring] Best build guide",
			"[Fortnite] Sezona 8 je tu!",
			"[Chess] Tips za izboljšanje ratinga",
			"[Board Games] Top 10 igre tega leta",
		},
		"Phones": {
			"[AI] ChatGPT nova verzija",
			"[Hardware] Best laptops 2024",
			"[Phone] iPhone 16 review",
			"[IoT] Smart home setup guide",
			"[5G] Kako deluje 5G omrežje",
		},
		"Travel": {
			"[Europe] Best destinations v Evropi",
			"[Asia] Backpacking through Southeast Asia",
			"[America] Road trip planing tips",
			"[Budget] Cheap travel hacks",
			"[Visa] Travel visa guide",
		},
		"Food": {
			"[Recipe] Kako narediti Perfect pasta",
			"[Cooking] Top 5 kuharskih trikov",
			"[Diet] Healthy meal prep ideas",
			"[Restaurant] Best Italian restaurants",
			"[Baking] Homemade bread tutorial",
		},
	}

	messageCount := 0
	for topicName, msgs := range messages {
		topicID := topicIDs[topicName]

		// Povežem z različnimi uporabniki
		userList := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}

		for i, msg := range msgs {
			userID := userIDs[userList[i%len(userList)]]

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := client.PostMessage(ctx, &razp.PostMessageRequest{
				TopicId: topicID,
				UserId:  userID,
				Text:    msg,
			})
			cancel()

			if err != nil {
				log.Printf("❌ Failed to post message: %v", err)
				continue
			}
			messageCount++
			fmt.Printf("✅ Posted message in '%s' by %s\n", topicName, userList[i%len(userList)])
		}
	}
}
