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
	//"bufio"
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	//"os"

	razp "razpravljalnica/razpravljalnica"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	app := tview.NewApplication()
	app.EnableMouse(true)

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	grpcClient := razp.NewMessageBoardClient(conn)

	pages := tview.NewPages()

	status := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	title := tview.NewTextView().
		SetText("Welcome to Razpravljalnici").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorPink)

	createBtn := tview.NewButton("Create Account")
	loginBtn := tview.NewButton("Login")

	choiceFlex := tview.NewFlex(). //zacetna stran
					SetDirection(tview.FlexRow).
					AddItem(title, 2, 1, false).
					AddItem(tview.NewBox(), 1, 1, false).
					AddItem(createBtn, 3, 1, true).
					AddItem(loginBtn, 3, 1, false)
	choiceFlex.SetBorder(true).SetTitle(" Welcome ")

	usernameField := tview.NewInputField().SetLabel("Username: ").SetFieldWidth(20)
	passwordField := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*').SetFieldWidth(20)
	createAccountButton := tview.NewButton("Create Account")

	topicList := tview.NewList().
		ShowSecondaryText(false)

	createTopicBtn := tview.NewButton("Create Topic")
	refreshBtn := tview.NewButton("Refresh Topics")

	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(createTopicBtn, 0, 1, false).
		AddItem(refreshBtn, 0, 1, false)

	homeFlex := tview.NewFlex(). //spremeni v reddit like
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(topicList, 0, 6, true).
		AddItem(status, 1, 1, false)

	homeFlex.SetBorder(true).SetTitle(" Topics ")

	newTopicInput := tview.NewInputField().
		SetLabel("Topic name: ").
		SetFieldWidth(30)

	createTopicModal := tview.NewModal(). //spremeni
						AddButtons([]string{"Create", "Cancel"}).
						SetText("Enter topic name").
						SetBackgroundColor(tcell.ColorBlack)

	createAccount := func() { //doda novega userja
		name := usernameField.GetText()
		pass := passwordField.GetText()
		if name == "" || pass == "" { //potrebno dodati funkcijo ki preverja kako korekten je input
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Username and password required")
			return
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //idk je pomembno sam nevem kak more bit
			defer cancel()

			_, err := grpcClient.CreateUser(ctx, &razp.CreateUserRequest{
				Name:     name,
				Password: pass,
			})
			app.QueueUpdateDraw(func() {
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText(err.Error())
				} else {
					status.SetTextColor(tcell.ColorGreen)
					status.SetText("Account created successfully")
				}
			})
		}()
	}

	loadTopics := func() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.ListTopics(ctx, &emptypb.Empty{})
			app.QueueUpdateDraw(func() {
				topicList.Clear()
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText(err.Error())
					return
				}

				for _, t := range resp.Topics {
					topicList.AddItem(t.Name, "", 0, nil)
				}

				status.SetTextColor(tcell.ColorGreen)
				status.SetText("Topics loaded")
			})
		}()
	}

	createTopic := func(name string) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &razp.CreateTopicRequest{Name: name}

			resp, err := grpcClient.CreateTopic(ctx, req)
			app.QueueUpdateDraw(func() {
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText(err.Error())
					return
				}

				status.SetTextColor(tcell.ColorGreen)
				status.SetText(fmt.Sprintf("Topic '%s' created", resp.Name))

				loadTopics()
			})
		}()
	}

	createTopicModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Create" {
			name := newTopicInput.GetText()
			name = strings.TrimSpace(name)
			if name != "" {
				createTopic(name)
			} else {
				status.SetTextColor(tcell.ColorRed)
				status.SetText("Topic name cannot be empty")
			}
		}

		// Zapri modal in se vrni na home
		pages.HidePage("createTopicModal")
		pages.SwitchToPage("home")
		app.SetFocus(topicList)
	})

	refreshBtn.SetSelectedFunc(loadTopics)

	createTopicBtn.SetSelectedFunc(func() {
		newTopicInput.SetText("")

		// modal mora vsebovati input field
		modalFlex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(newTopicInput, 2, 1, true).
			AddItem(createTopicModal, 0, 3, false)

		pages.AddPage("createTopicModal", modalFlex, true, true)
		app.SetFocus(newTopicInput)
	})

	createAccountButton.SetSelectedFunc(createAccount)

	usernameField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			app.SetFocus(passwordField)
		}
	})
	passwordField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			createAccount()
		}
	})

	signupFlex := tview.NewFlex(). //kjer si ustvaris nov account
					SetDirection(tview.FlexRow).
					AddItem(usernameField, 3, 1, true).
					AddItem(passwordField, 3, 1, false).
					AddItem(createAccountButton, 2, 1, false).
					AddItem(status, 2, 1, false)
	signupFlex.SetBorder(true).SetTitle(" Sign Up ")

	loginUsername := tview.NewInputField().SetLabel("Username: ").SetFieldWidth(20)
	loginPassword := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*').SetFieldWidth(20)
	loginButton := tview.NewButton("Login")

	loginFunc := func() { //potrebno dodat funkcijo ki preveri ce je user ze not
		name := loginUsername.GetText()
		pass := loginPassword.GetText()
		if name == "" || pass == "" {
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Username and password required")
			return
		} else {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //idk je pomembno sam nevem kak more bit
				defer cancel()

				_, err := grpcClient.FindUser(ctx, &razp.FindUserRequest{
					Name:     name,
					Password: pass,
				})
				app.QueueUpdateDraw(func() {
					if err != nil {
						status.SetTextColor(tcell.ColorRed)
						status.SetText(err.Error())
					} else {
						status.SetTextColor(tcell.ColorGreen)
						status.SetText("Login successful")

						pages.SwitchToPage("home")
						loadTopics()
						app.SetFocus(topicList)
					}
				})
			}()
		}
	}

	loginButton.SetSelectedFunc(loginFunc)

	loginUsername.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			app.SetFocus(loginPassword)
		}
	})
	loginPassword.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			loginFunc()
		}
	})

	loginFlex := tview.NewFlex(). //kjer se loggas in
					SetDirection(tview.FlexRow).
					AddItem(loginUsername, 3, 1, true).
					AddItem(loginPassword, 3, 1, false).
					AddItem(loginButton, 2, 1, false).
					AddItem(status, 2, 1, false)
	loginFlex.SetBorder(true).SetTitle(" Login ")

	createBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("signup")
		app.SetFocus(usernameField)
	})
	loginBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("login")
		app.SetFocus(loginUsername)
	})

	pages.AddPage("choice", choiceFlex, true, true) //kako je na zacetku
	pages.AddPage("signup", signupFlex, true, false)
	pages.AddPage("login", loginFlex, true, false)
	pages.AddPage("home", homeFlex, true, false)

	app.SetRoot(pages, true).SetFocus(createBtn)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
