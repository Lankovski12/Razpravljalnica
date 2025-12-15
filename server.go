/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	razp "razpravljalnica/razpravljalnica"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	port                = flag.Int("port", 50051, "The server port")
	numOfUsers    int64 = 0
	numOfTopics   int64 = 0
	users               = make([]*razp.User, 0, 10)
	topics              = make([]*razp.Topic, 0, 10)
	topicMessages       = make([][]*razp.Message, 0, 10)
)

type server struct {
	razp.UnimplementedMessageBoardServer
}

func (s *server) CreateUser(_ context.Context, in *razp.CreateUserRequest) (*razp.User, error) {

	newUserName := in.Name
	newPassword := in.Password
	if numOfUsers != 0 {
		for _, user := range users {
			if user.Name == newUserName {
				return nil, status.Errorf(codes.Aborted, "User with same name already exists, please choose a different name")
				// ne izpise errorja plus poglej ce se da to mogoce nrdit tak da avtomatsko relauncha request za CreateUser verjetno rabim v clientu se enkrat callat
			}
		}
	}
	if newPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}
	//dodat potrebno za sifrirat
	newUser := &razp.User{Id: numOfUsers, Name: newUserName, Password: newPassword}
	numOfUsers += 1
	users = append(users, newUser)
	fmt.Printf("Added New User, Id:%d Name: %s, Password: %s\n", newUser.Id, newUser.Name, newUser.Password)

	return newUser, nil
}

func (s *server) FindUser(_ context.Context, in *razp.CreateUserRequest) (*emptypb.Empty, error) {
	userName := in.Name
	userPass := in.Password

	for _, user := range users {
		if user.Name == userName {
			if user.Password == userPass {
				return &emptypb.Empty{}, nil
			} else {
				return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "Wrong password")
			}
		}
	}
	return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "User does not exist")
}

func (s *server) CreateTopic(_ context.Context, in *razp.CreateTopicRequest) (*razp.Topic, error) {

	newTopicName := in.Name

	if numOfTopics != 0 {
		// fmt.Print("tsets")
		for _, topic := range topics {
			if topic.Name == newTopicName {
				return nil, status.Error(codes.Aborted, "Topic with same name already exists, please choose a different name")
				// ne izpise errorja ampak ga poslje, lahka si ga pol printas ce hoces
				// plus poglej ce se da to mogoce nrdit tak da avtomatsko relauncha request za CreateUser verjetno rabim v clientu se enkrat callat
			}
		}
	}

	newTopic := &razp.Topic{Id: numOfTopics, Name: newTopicName, NumberOfMessage: 0}
	numOfTopics += 1
	topics = append(topics, newTopic)

	topicMessages = append(topicMessages, make([]*razp.Message, 0, 10))

	fmt.Printf("Added New Topic, Id:%d Name: %s\n", newTopic.Id, newTopic.Name)

	return newTopic, nil
}

func (s *server) PostMessage(_ context.Context, in *razp.PostMessageRequest) (*razp.Message, error) {

	topicIdx := in.TopicId - 1
	numOfTopicMessages := topics[topicIdx].NumberOfMessage
	currentTime := timestamppb.Now()
	newMessage := &razp.Message{Id: numOfTopicMessages, TopicId: in.TopicId, UserId: in.UserId, Text: in.Text, CreatedAt: currentTime, Likes: 0}
	topics[topicIdx].NumberOfMessage += 1
	topicMessages[in.TopicId-1] = append(topicMessages[topicIdx], newMessage)

	return newMessage, nil
}

func (s *server) GetMessages(_ context.Context, in *razp.GetMessagesRequest) (*razp.GetMessagesResponse, error) {

	// implementi linked list za messages al pa ne sj nvm probi zdle array/slice pa bomo pol fixal ce bo treba
	topicIdx := in.TopicId - 1
	numOfTopicMessages := topics[topicIdx].NumberOfMessage
	if in.FromMessageId > numOfTopicMessages {
		return nil, status.Error(codes.Aborted, "This messageID does not exist")
	}
	// ce je messageId prevlek ne bo delal
	if in.TopicId > numOfTopics {
		return nil, status.Error(codes.Aborted, "This topicID does not exist")
	}

	var messageInterval []*razp.Message = topicMessages[topicIdx][:]
	var leftInterval int64
	var rightInterval int64
	if in.Limit+int32(in.FromMessageId) > int32(numOfTopicMessages) {
		leftInterval = in.FromMessageId - 1
		rightInterval = numOfTopicMessages // ker slice ig vzame od vkljucno indexa do indexa (ne vkljucno)
	} else {
		leftInterval = in.FromMessageId - 1
		rightInterval = in.FromMessageId + int64(in.Limit)
	}

	messageInterval = messageInterval[leftInterval:rightInterval]
	// fmt.Printf("LeftInterval:%d RightInterval:%d", leftInterval, rightInterval)
	newGetMessagesResponse := &razp.GetMessagesResponse{Messages: messageInterval}

	return newGetMessagesResponse, nil
}

func (s *server) ListTopics(_ context.Context, empty *emptypb.Empty) (*razp.ListTopicsResponse, error) {

	// implementi linked list za messages al pa ne sj nvm probi zdle array/slice pa bomo pol fixal ce bo treba

	var allTopics []*razp.Topic = topics[:]
	newListTopicsResponse := &razp.ListTopicsResponse{Topics: allTopics}

	return newListTopicsResponse, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	razp.RegisterMessageBoardServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
