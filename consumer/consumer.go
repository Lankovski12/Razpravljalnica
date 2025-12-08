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
	razp "razpravljalnica/razpravljalnica"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
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

	newUser, err := grpcClient.CreateUser(ctx, &razp.CreateUserRequest{Name: "UserRequest"})
	fmt.Printf("Username: %s, Id: %d\n", newUser.Name, newUser.Id)
}
