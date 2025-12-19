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
	"sync"

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
	topics              = make([]*razp.Topic, 0, 10)     //
	topicMessages       = make([][]*razp.Message, 0, 10) //[topicId][messageId]
	likes               = make([][][]bool, 10)           // [topicId][messageId][userId]
	first         bool  = true
	topicSubsCh         = make([]map[int]chan *razp.MessageEvent, 0, 10)
	userMu        sync.RWMutex
	topicMu       sync.RWMutex
	messageMu     sync.RWMutex
	likeMu        sync.RWMutex
	subMu         sync.RWMutex
	nextSubID     int
	seqNumber     int64 = 0
)

type server struct {
	razp.UnimplementedMessageBoardServer
}

func lockWrite(l ...*sync.RWMutex) {
	for _, mu := range l {
		mu.Lock()
	}
}

func unlockWrite(l ...*sync.RWMutex) {
	for i := len(l) - 1; i >= 0; i-- {
		l[i].Unlock()
	}
}

//support functions for mutex (writes)

func lockRead(l ...*sync.RWMutex) {
	for _, mu := range l {
		mu.RLock()
	}
}

func unlockRead(l ...*sync.RWMutex) {
	for i := len(l) - 1; i >= 0; i-- {
		l[i].RUnlock()
	}
}

//support functions for mutex (reads)

func (s *server) CreateUser(_ context.Context, in *razp.CreateUserRequest) (*razp.User, error) {
	userMu.Lock()
	defer userMu.Unlock()
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

func (s *server) FindUser(_ context.Context, in *razp.FindUserRequest) (*razp.User, error) {
	userMu.RLock()
	defer userMu.RUnlock()
	userName := in.Name
	userPass := in.Password

	for _, user := range users {
		if user.Name == userName {
			if user.Password == userPass {
				return user, nil
			} else {
				return nil, status.Error(codes.InvalidArgument, "Wrong password")
			}
		}
	}
	return nil, status.Error(codes.InvalidArgument, "User does not exist")
}

func (s *server) CreateTopic(_ context.Context, in *razp.CreateTopicRequest) (*razp.Topic, error) {
	lockWrite(&topicMu, &messageMu, &subMu)
	defer unlockWrite(&topicMu, &messageMu, &subMu)
	newTopicName := in.Name
	if numOfTopics != 0 {
		// fmt.Print("tsets")
		for _, topic := range topics {
			if topic.Name == newTopicName {
				return nil, status.Errorf(codes.Aborted, "Topic with same name already exists, please choose a different name")
				// ne izpise errorja ampak ga poslje, lahka si ga pol printas ce hoces
				// plus poglej ce se da to mogoce nrdit tak da avtomatsko relauncha request za CreateUser verjetno rabim v clientu se enkrat callat
			}
		}
	}

	newTopic := &razp.Topic{Id: numOfTopics, Name: newTopicName}
	numOfTopics += 1
	topics = append(topics, newTopic)
	// ustvari topic in ga dodaj v global tabelo in dodaj kanal za subscription
	topicMessages = append(topicMessages, make([]*razp.Message, 0, 10))
	topicSubsCh = append(topicSubsCh, make(map[int]chan *razp.MessageEvent))

	fmt.Printf("Added New Topic, Id:%d Name: %s\n", newTopic.Id, newTopic.Name)

	return newTopic, nil
}

func (s *server) PostMessage(_ context.Context, in *razp.PostMessageRequest) (*razp.Message, error) {
	messageMu.Lock()
	defer messageMu.Unlock()
	topicIdx := in.TopicId - 1
	numOfTopicMessages := int64(len(topicMessages[topicIdx]))
	currentTime := timestamppb.Now()
	newMessage := &razp.Message{Id: numOfTopicMessages, TopicId: in.TopicId, UserId: in.UserId, Text: in.Text, CreatedAt: currentTime, Likes: 0}
	topicMessages[in.TopicId-1] = append(topicMessages[topicIdx], newMessage)

	// dodaj message v globalno tabelo sporocil

	publishToTopic(topicIdx, &razp.MessageEvent{
		SequenceNumber: seqNumber,
		Op:             razp.OpType_OP_POST,
		Message:        newMessage,
		EventAt:        timestamppb.Now(),
	})
	seqNumber += 1

	// funkcija za subscription ki

	return newMessage, nil
}

func (s *server) GetMessages(_ context.Context, in *razp.GetMessagesRequest) (*razp.GetMessagesResponse, error) {
	// implementi linked list za messages al pa ne sj nvm probi zdle array/slice pa bomo pol fixal ce bo treba
	lockRead(&topicMu, &messageMu)
	defer unlockRead(&topicMu, &messageMu)
	topicIdx := in.TopicId - 1
	numOfTopicMessages := int64(len(topicMessages[topicIdx]))
	if topicIdx >= numOfTopics {
		return nil, status.Error(codes.Aborted, "This topicID does not exist")
	}

	if in.FromMessageId > numOfTopicMessages {
		return nil, status.Error(codes.Aborted, "This messageID does not exist")
	}
	// ce je messageId prevlek ne bo delal
	var leftInterval int64
	var rightInterval int64
	if in.Limit+int32(in.FromMessageId) > int32(numOfTopicMessages) {
		leftInterval = in.FromMessageId - 1
		rightInterval = numOfTopicMessages // ker slice ig vzame od vkljucno indexa do indexa (ne vkljucno)
	} else {
		leftInterval = in.FromMessageId - 1
		rightInterval = in.FromMessageId + int64(in.Limit)
	}

	messageInterval := append([]*razp.Message(nil), topicMessages[topicIdx][leftInterval:rightInterval]...)
	// doda message v globalno tabelo sporocil
	// fmt.Printf("LeftInterval:%d RightInterval:%d", leftInterval, rightInterval)
	newGetMessagesResponse := &razp.GetMessagesResponse{Messages: messageInterval}
	// in naredi response zato da ga vrnemo
	return newGetMessagesResponse, nil
}

func (s *server) ListTopics(_ context.Context, empty *emptypb.Empty) (*razp.ListTopicsResponse, error) {
	topicMu.RLock()
	defer topicMu.RUnlock()
	allTopics := append([]*razp.Topic(nil), topics...)
	newListTopicsResponse := &razp.ListTopicsResponse{Topics: allTopics}
	// vrne tabelo vseh topicov
	return newListTopicsResponse, nil
}

func (s *server) UpdateMessage(ctx context.Context, in *razp.UpdateMessageRequest) (*razp.Message, error) {
	lockWrite(&messageMu, &likeMu)
	defer unlockWrite(&messageMu, &likeMu)
	topicIdx := in.TopicId - 1
	messageIdx := in.MessageId - 1

	if topicIdx >= numOfTopics {
		return nil, status.Error(codes.Aborted, "This topicID does not exist")
	}

	if messageIdx >= int64(len(topicMessages[topicIdx])) {
		return nil, status.Error(codes.Aborted, "This messageID does not exist")
	}

	if topicMessages[topicIdx][messageIdx] == nil {
		return nil, status.Error(codes.Aborted, "This message was deleted by the author")
	}

	if topicMessages[topicIdx][messageIdx].UserId != in.UserId {
		return nil, status.Error(codes.Aborted, "This message can only be edited by the author")
	}

	// napake ki nam lahko povzrocijo tezave ali nedovoljene dostope

	topicMessages[topicIdx][messageIdx].Text = in.Text
	topicMessages[topicIdx][messageIdx].CreatedAt = timestamppb.Now()

	// popravek v globalni tabeli

	returnMessage := topicMessages[topicIdx][messageIdx]

	// struct ki ga vrnemo

	return returnMessage, nil
}

func (s *server) DeleteMessage(ctx context.Context, in *razp.DeleteMessageRequest) (*emptypb.Empty, error) {
	lockWrite(&messageMu, &likeMu)
	defer unlockWrite(&messageMu, &likeMu)
	topicIdx := in.TopicId - 1
	messageIdx := in.MessageId - 1

	if topicIdx >= numOfTopics {
		return nil, status.Error(codes.Aborted, "This topicID does not exist")
	}

	if messageIdx >= int64(len(topicMessages[topicIdx])) {
		return nil, status.Error(codes.Aborted, "This messageID does not exist")
	}

	if topicMessages[topicIdx][messageIdx].UserId != in.UserId {
		return nil, status.Error(codes.Aborted, "This message can only be deleted by the author")
	}

	// napake ki nam lahko povzrocijo tezave ali nedovoljene dostope

	topicMessages[topicIdx][messageIdx].TopicId = -1
	topicMessages[topicIdx][messageIdx].UserId = -1
	topicMessages[topicIdx][messageIdx].Likes = -1
	topicMessages[topicIdx][messageIdx].Text = "message deleted by user"

	// fake brisanje, dejansko samo zamenjamo vse na neveljavne vrednosti

	return nil, nil
}

func (s *server) LikeMessage(ctx context.Context, in *razp.LikeMessageRequest) (*razp.Message, error) {
	lockWrite(&messageMu, &likeMu)
	defer unlockWrite(&messageMu, &likeMu)
	topicIdx := in.TopicId - 1
	messageIdx := in.MessageId - 1

	if topicIdx >= numOfTopics {
		return nil, status.Error(codes.Aborted, "This topicID does not exist")
	}

	if messageIdx >= int64(len(topicMessages[topicIdx])) {
		return nil, status.Error(codes.Aborted, "This messageID does not exist")
	}

	if topicMessages[topicIdx][messageIdx].Text == "message deleted by user" {
		return nil, status.Error(codes.Aborted, "This message was deleted by the author")
	}

	// napake ki nam lahko povzrocijo tezave ali nedovoljene dostope

	likes[topicIdx][messageIdx][in.UserId] = !likes[topicIdx][messageIdx][in.UserId]
	if likes[topicIdx][messageIdx][in.UserId] {
		topicMessages[topicIdx][messageIdx].Likes += 1
	} else {
		topicMessages[topicIdx][messageIdx].Likes -= 1
	}

	// dodamo zapis v boolean tabelo likeov in inkrementira vrednost v tabeli sporocil

	returnMessage := topicMessages[topicIdx][messageIdx]
	return returnMessage, nil
}

func (s *server) SubscribeTopic(in *razp.SubscribeTopicRequest, stream grpc.ServerStreamingServer[razp.MessageEvent]) error {

	if in.TopicId <= 0 {
		return status.Error(codes.InvalidArgument, "TopicId must be >= 1")
	}
	topicMu.RLock()
	topicCount := numOfTopics
	topicMu.RUnlock()
	if in.TopicId > topicCount {
		return status.Error(codes.NotFound, "This topicID does not exist")
	}
	t := int(in.TopicId - 1)

	// napake ki nam lahko povzrocijo tezave ali nedovoljene dostope

	// ustvari kanal za subscriberja
	ch := make(chan *razp.MessageEvent, 256)

	subMu.Lock()
	if t >= len(topicSubsCh) {
		subMu.Unlock()
		return status.Error(codes.Aborted, "This topicID does not exist")
	}
	id := nextSubID
	nextSubID++
	topicSubsCh[t][id] = ch
	subMu.Unlock()

	// dodaj ga v map kanalov ki je per topic zato je vseeeno kaksen je subId ker

	defer func() {
		subMu.Lock()
		delete(topicSubsCh[t], id)
		subMu.Unlock()
		close(ch)
	}()
	// zapremo kanal in ga odstranimo iz mapa ko se konca

	ctx := stream.Context()

	messageMu.RLock()
	msgs := append([]*razp.Message(nil), topicMessages[t]...) // snapshot
	messageMu.RUnlock()

	start := int64(0)
	if in.FromMessageId > 0 {
		start = in.FromMessageId - 1
	}
	if start > int64(len(msgs)) {
		start = int64(len(msgs))
	}

	for i := start; i < int64(len(msgs)); i++ {
		m := msgs[i]
		if m == nil {
			continue
		}
		ev := &razp.MessageEvent{
			SequenceNumber: m.Id,
			Op:             razp.OpType_OP_POST,
			Message:        m,
			EventAt:        m.CreatedAt,
		}
		if err := stream.Send(ev); err != nil {
			return err
		}
	}

	// od prejsnega commenta do tu je zato da from message id deluje previlno ker rabimo za nazaj tud poslat po kanalu ko se subscribeamo
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-ch:
			if err := stream.Send(ev); err != nil {
				return err
			}
		}
	}
	// za live posiljanje
}

func publishToTopic(topicIdx int64, ev *razp.MessageEvent) {
	subMu.RLock()
	subs, ok := func() (map[int]chan *razp.MessageEvent, bool) {
		if topicIdx < 0 || topicIdx >= int64(len(topicSubsCh)) {
			return nil, false
		}
		return topicSubsCh[topicIdx], true
	}()
	if !ok {
		subMu.RUnlock()
		return
	}
	// vrne pointer na hrambo kanalov subscriptiona iz globalne tabele povezan z stevilko topica in vsakic ko dodamo novega
	// v topicSubsCh bo pogledal ce je ze notri ker je map (to se zgodi ze v subscribe topic)

	chans := make([]chan *razp.MessageEvent, 0, len(subs))
	for _, ch := range subs {
		chans = append(chans, ch)
	}
	subMu.RUnlock()

	// doda lokalno kopijo vseh kanalov na topicu zato da nam ni treba obdrzati mutex basically doda vse previous v prazen je tist append

	for _, ch := range chans {
		select {
		case ch <- ev:
		default:
			// drop if slow
		}
	}
	// posiljanje v kanal
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	if first {
		for t := 0; t < 10; t++ { // to sm ze zgori inicializiru
			likes[t] = make([][]bool, 50)
			for m := 0; m < 50; m++ {
				likes[t][m] = make([]bool, 20) // initialization za 10 topicov vsak 50 messageov in vsakih 50 msg 20 userev (10 * 50 * 20)
			}
		}
		first = false
	}
	razp.RegisterMessageBoardServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
