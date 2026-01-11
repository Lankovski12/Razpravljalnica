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

	themeColorSec = tcell.ColorBlue

	currentUserID   int64 = -1
	currentUsername       = ""

	selectedTopicIndex int = -1

	userMap = make(map[int64]string)
)

func updateThemeColors(app *tview.Application, createBtn, loginBtn, createAccountButton, backFromSignupBtn, loginButton, backFromLoginBtn, createTopicBtn, refreshBtn, newMessageBtn, subsribeTopic, postNewMessageBtn, cancelNewMessageBtn, changePasswordBtn, backFromProfileBtn, profileBtn *tview.Button, themeColorSec tcell.Color) {
	createBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	loginBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	createAccountButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	backFromSignupBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	loginButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	backFromLoginBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	createTopicBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	refreshBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	profileBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	newMessageBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	subsribeTopic.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	postNewMessageBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	cancelNewMessageBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	changePasswordBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	backFromProfileBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	app.Draw()
}

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

	// SCREEN 1: CREATE ACCOUNT OR LOG IN

	createBtn := tview.NewButton("Create account")
	createBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	loginBtn := tview.NewButton("Login")
	loginBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	welcomeGrid := tview.NewGrid().
		SetRows(10, 3, 1, 3, 2, 1).
		SetColumns(0, 30, 0).
		SetBorders(false)

	welcomeGrid.AddItem(createBtn, 1, 1, 1, 1, 0, 0, true)
	welcomeGrid.AddItem(loginBtn, 3, 1, 1, 1, 0, 0, false)
	welcomeGrid.AddItem(status, 5, 1, 1, 1, 0, 0, false)

	//GLAVNA STRUKTURA
	choiceFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(welcomeGrid, 0, 1, true)
	choiceFlex.SetBorder(true).SetTitle(" Welcome ")

	// SCREEN 2: CREATE NEW ACCOUNT

	usernameField := tview.NewInputField().SetLabel("Username: ").SetFieldWidth(20)
	usernameField.SetLabelColor(tcell.ColorWhite)

	passwordField := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*').SetFieldWidth(20)
	passwordField.SetLabelColor(tcell.ColorWhite)

	createAccountButton := tview.NewButton("Create account")
	createAccountButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	backFromSignupBtn := tview.NewButton("Back")
	backFromSignupBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	signupGrid := tview.NewGrid().
		SetRows(3, 1, 2, 1, 3, 3, 1, 3, 0, 1).
		SetColumns(0, 50, 0).
		SetBorders(false)

	signupGrid.AddItem(usernameField, 1, 1, 1, 1, 0, 0, false)
	signupGrid.AddItem(passwordField, 3, 1, 1, 1, 0, 0, false)
	signupGrid.AddItem(createAccountButton, 5, 1, 1, 1, 0, 30, false)
	signupGrid.AddItem(backFromSignupBtn, 7, 1, 1, 1, 0, 30, false)
	signupGrid.AddItem(status, 9, 1, 1, 2, 0, 0, false)

	//GLAVNA STRUKTURA
	signupFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(signupGrid, 0, 1, true)
	signupFlex.SetBorder(true).SetTitle(" Sign Up ")

	//SCREEN 3: LOGIN SCREEN

	loginUsername := tview.NewInputField().SetLabel("Username: ").SetFieldWidth(20)
	loginUsername.SetLabelColor(tcell.ColorWhite)

	loginPassword := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*').SetFieldWidth(20)
	loginPassword.SetLabelColor(tcell.ColorWhite)

	loginButton := tview.NewButton("Login")
	loginButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	backFromLoginBtn := tview.NewButton("Back")
	backFromLoginBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	loginGrid := tview.NewGrid().
		SetRows(3, 1, 2, 1, 3, 3, 1, 3, 0, 1).
		SetColumns(0, 50, 0).
		SetBorders(false)

	loginGrid.AddItem(loginUsername, 1, 1, 1, 1, 0, 0, false)
	loginGrid.AddItem(loginPassword, 3, 1, 1, 1, 0, 0, false)
	loginGrid.AddItem(loginButton, 5, 1, 1, 1, 0, 30, false)
	loginGrid.AddItem(backFromLoginBtn, 7, 1, 1, 1, 0, 30, false)
	loginGrid.AddItem(status, 9, 1, 1, 2, 0, 0, false)

	//GLAVNA STRUKTURA
	loginFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginGrid, 0, 1, true)
	loginFlex.SetBorder(true).SetTitle(" Login ")

	//SCREEN 4: MAIN PAGE

	topicList := tview.NewList().ShowSecondaryText(false)

	newTopicInput := tview.NewInputField().
		SetLabel("New topic: ").
		SetFieldWidth(20)
	newTopicInput.SetLabelColor(tcell.ColorWhite)
	newTopicInput.SetFieldTextColor(tcell.ColorWhite)

	createTopicBtn := tview.NewButton("Create")
	createTopicBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	refreshBtn := tview.NewButton("Refresh")
	refreshBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	profileBtn := tview.NewButton("My Profile")
	profileBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	newMessageBtn := tview.NewButton("New Message")
	newMessageBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	subsribeTopic := tview.NewButton("Subscribe")
	subsribeTopic.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	topicsGrid := tview.NewGrid().
		SetRows(1, 19, 10, 1, 1, 3, 1, 3).
		SetColumns(0, 0).
		SetBorders(false)

	topicsGrid.AddItem(topicList, 1, 0, 1, 2, 0, 0, true)
	topicsGrid.AddItem(newTopicInput, 3, 0, 1, 1, 0, 0, false)
	topicsGrid.AddItem(createTopicBtn, 5, 0, 1, 2, 0, 0, false)
	topicsGrid.AddItem(profileBtn, 7, 0, 1, 1, 0, 0, false)
	topicsGrid.AddItem(refreshBtn, 7, 1, 1, 1, 0, 0, false)

	topicsPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topicsGrid, 0, 1, false)
	topicsPanel.SetBorder(true).SetTitle(" Topics ")

	messagesGrid := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(false)

	messagesScroll := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messagesGrid, 0, 1, true)

	messagesDetailAll := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messagesScroll, 0, 1, false).
		AddItem(status, 1, 1, false).
		AddItem(subsribeTopic, 3, 1, false).
		AddItem(newMessageBtn, 3, 1, false)
	messagesDetailAll.SetBorder(true).SetTitle(" Home ")

	contentFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(topicsPanel, 0, 1, false).
		AddItem(messagesDetailAll, 0, 3, false)

	//SCREEN 7: NEW MESSAGE - Pisanje novega sporočila

	newMessageTitle := tview.NewInputField().
		SetLabel("Title: ").
		SetFieldWidth(30)
	newMessageTitle.SetLabelColor(tcell.ColorWhite)
	newMessageTitle.SetFieldTextColor(tcell.ColorWhite)

	newMessageTextarea := tview.NewTextArea()
	newMessageTextarea.SetBorderColor(tcell.ColorWhite)

	postNewMessageBtn := tview.NewButton("Post")
	postNewMessageBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	cancelNewMessageBtn := tview.NewButton("Cancel")
	cancelNewMessageBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	newMessageGrid := tview.NewGrid().
		SetRows(2, 2, 1, 0, 3, 1, 3).
		SetColumns(0, 40, 0).
		SetBorders(false)

	newMessageGrid.AddItem(newMessageTitle, 1, 1, 1, 1, 0, 0, true)
	newMessageGrid.AddItem(tview.NewTextView().SetText("Content:"), 2, 1, 1, 1, 0, 0, false)
	newMessageGrid.AddItem(newMessageTextarea, 3, 1, 1, 1, 0, 0, false)
	newMessageGrid.AddItem(postNewMessageBtn, 4, 1, 1, 1, 0, 0, false)
	newMessageGrid.AddItem(cancelNewMessageBtn, 6, 1, 1, 1, 0, 0, false)

	newMessageFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(newMessageGrid, 0, 1, true)
	newMessageFlex.SetBorder(true).SetTitle(" New Message ")

	//SCREEN 5: PROFILE
	profileUsernameDisplay := tview.NewTextView().
		SetText("Username: " + currentUsername).
		SetTextAlign(tview.AlignLeft)
	profileUsernameDisplay.SetTextColor(tcell.ColorWhite)

	profileOldPassword := tview.NewInputField().
		SetLabel("Old password: ").
		SetMaskCharacter('*').
		SetFieldWidth(15)
	profileOldPassword.SetLabelColor(tcell.ColorWhite)
	profileOldPassword.SetFieldTextColor(tcell.ColorWhite)

	profileNewPassword := tview.NewInputField().
		SetLabel("New password: ").
		SetMaskCharacter('*').
		SetFieldWidth(15)
	profileNewPassword.SetLabelColor(tcell.ColorWhite)
	profileNewPassword.SetFieldTextColor(tcell.ColorWhite)

	changePasswordBtn := tview.NewButton("Change Password")
	changePasswordBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	backFromProfileBtn := tview.NewButton("Back")
	backFromProfileBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	dropdown := tview.NewDropDown().
		SetLabel("Theme: ").
		SetOptions([]string{"Pink", "Green", "Orange", "Blue", "Red", "Violet"}, func(option string, optionIndex int) {
			switch option {
			case "Pink":
				themeColorSec = tcell.ColorRed
			case "Green":
				themeColorSec = tcell.ColorGreen
			case "Orange":
				themeColorSec = tcell.ColorOrange
			case "Blue":
				themeColorSec = tcell.ColorBlue
			case "Red":
				themeColorSec = tcell.ColorRed
			case "Violet":
				themeColorSec = tcell.ColorPaleVioletRed
			}
		})
	//updateThemeColors(app, createBtn, loginBtn, createAccountButton, backFromSignupBtn, loginButton, backFromLoginBtn, createTopicBtn, refreshBtn, profileBtn, newMessageBtn, subsribeTopic, changePasswordBtn, backFromProfileBtn, postNewMessageBtn, cancelNewMessageBtn, themeColorSec)
	dropdown.SetLabelColor(tcell.ColorWhite)

	logoutBtn := tview.NewButton("Log out")
	logoutBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorLightGray))

	exitBtn := tview.NewButton("Exit")
	exitBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorLightGray).Background(tcell.ColorWhite))

	profileGrid := tview.NewGrid().
		SetRows(2, 2, 2, 6, 3, 3, 3, 1, 3, 1, 1, 4, 3, 1, 3).
		SetColumns(0, 30, 0).
		SetBorders(false)

	profileGrid.AddItem(profileUsernameDisplay, 1, 1, 1, 1, 0, 0, false)
	profileGrid.AddItem(dropdown, 2, 1, 1, 1, 0, 0, true)
	profileGrid.AddItem(profileOldPassword, 4, 1, 1, 1, 0, 0, false)
	profileGrid.AddItem(profileNewPassword, 5, 1, 1, 1, 0, 0, false)
	profileGrid.AddItem(changePasswordBtn, 6, 1, 1, 1, 0, 0, false)
	profileGrid.AddItem(backFromProfileBtn, 8, 1, 1, 1, 0, 0, false)
	profileGrid.AddItem(status, 10, 0, 1, 3, 0, 0, false)
	profileGrid.AddItem(logoutBtn, 12, 1, 1, 1, 0, 0, false)
	profileGrid.AddItem(exitBtn, 14, 1, 1, 1, 0, 0, false)

	profilePanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(profileGrid, 0, 1, true)

	profileMessagesGrid := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(false)

	profileMessagesScroll := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(profileMessagesGrid, 0, 1, true)

	profileMessagesPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(profileMessagesScroll, 0, 1, false)

	profileFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(profilePanel, 0, 1, false).
		AddItem(profileMessagesPanel, 0, 2, false)
	profileFlex.SetBorder(true).SetTitle(" Profile ")

	//FUNKCIJE

	//FUNKCIJA KI NALOZI VSE TOPICS
	updateProfileDisplay := func() {
		profileUsernameDisplay.SetText("Username: " + currentUsername)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.GetUserMessages(ctx, &razp.GetUserMessagesRequest{
				UserId: currentUserID,
			})
			app.QueueUpdateDraw(func() {
				profileMessagesGrid.Clear()

				if err != nil {
					errorText := tview.NewTextView().SetText("Error loading messages: " + err.Error())
					errorText.SetTextColor(tcell.ColorRed)
					profileMessagesGrid.AddItem(errorText, 0, 0, 1, 1, 0, 0, false)
					return
				}

				if len(resp.Messages) == 0 {
					emptyText := tview.NewTextView().SetText("No messages yet")
					emptyText.SetTextColor(tcell.ColorGray)
					profileMessagesGrid.AddItem(emptyText, 0, 0, 1, 1, 0, 0, false)
					return
				}

				// Prikaži vsa sporočila
				for i, msg := range resp.Messages {
					messageText := fmt.Sprintf(
						"[%s] Topic: %s (ID: %d)\n%s\n[Likes: %d]",
						msg.CreatedAt.AsTime().Format("2006-01-02 15:04"),
						msg.TopicName,
						msg.MessageId,
						msg.Text,
						msg.Likes,
					)

					msgView := tview.NewTextView().
						SetText(messageText).
						SetDynamicColors(true).
						SetWordWrap(true)
					msgView.SetTextColor(tcell.ColorWhite)
					msgView.SetBorder(true)

					profileMessagesGrid.AddItem(msgView, i, 0, 1, 1, 0, 0, false)
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

			})
		}()
	}

	loadUsers := func() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.ListUsers(ctx, &emptypb.Empty{})
			if err != nil {
				fmt.Printf("Napaka pri učitavanju uporabnikov: %v\n", err)
				return
			}

			// Popuni userMap
			for _, user := range resp.Users {
				userMap[user.Id] = user.Name
			}
		}()
	}

	//FUNKCIJA KI USTVARI NOVEGA UPORABNIKA
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

			raz, err := grpcClient.CreateUser(ctx, &razp.CreateUserRequest{
				Name:     name,
				Password: pass,
			})
			app.QueueUpdateDraw(func() {
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText("Account with that username already exist")
				} else {
					currentUsername = name
					currentUserID = raz.Id
					pages.SwitchToPage("home")
					loadTopics()
					loadUsers()
					app.SetFocus(topicList)
				}
			})
		}()
	}

	//FUNKCIJA KI USTVARI NOV TOPIC
	createTopic := func(name string) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			req := &razp.CreateTopicRequest{Name: name}

			resp, err := grpcClient.CreateTopic(ctx, req)
			app.QueueUpdateDraw(func() {
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText("Topic with same name already exists, please choose a different name")
					return
				}

				status.SetTextColor(tcell.ColorGreen)
				status.SetText("Topic '" + resp.Name + "' created successfully")

				newTopicInput.SetText("")

				loadTopics()
			})
		}()
	}

	//FUNKCIJA KI TE VPISE
	loginFunc := func() {
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

				raz, err := grpcClient.FindUser(ctx, &razp.FindUserRequest{
					Name:     name,
					Password: pass,
				})
				app.QueueUpdateDraw(func() {
					if err != nil {
						status.SetTextColor(tcell.ColorRed)
						status.SetText("Password does not match the username")
					} else {
						currentUserID = raz.Id
						currentUsername = name
						pages.SwitchToPage("home")
						loadTopics()
						loadUsers()
					}
				})
			}()
		}
	}

	// FUNKCIJA ZA NALAGANJE SPOROČIL IZ TOPIC
	loadMessagesForTopic := func(topicIndex int) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.GetMessages(ctx, &razp.GetMessagesRequest{
				TopicId:       int64(topicIndex) + 1,
				FromMessageId: 0,
				Limit:         100,
			})
			app.QueueUpdateDraw(func() {
				messagesGrid.Clear()
				messagesGrid.SetRows(0).SetColumns(0)

				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText("Failed to load messages: " + err.Error())
					return
				}

				row := 0
				for _, msg := range resp.Messages {
					text := msg.Text
					title := ""

					if len(text) > 0 && text[0] == '[' {
						endIdx := -1
						for j, c := range text {
							if c == ']' {
								endIdx = j
								break
							}
						}
						if endIdx > 0 {
							title = text[1:endIdx]
							text = text[endIdx+2:]
						}
					}

					messageView := tview.NewTextView().
						SetDynamicColors(true).
						SetWrap(true).
						SetText(text)
					messageView.SetBorder(true).SetTitle(fmt.Sprintf(" %s ", title)).SetTitleAlign(tview.AlignLeft)
					messageView.SetBackgroundColor(tcell.ColorDefault)

					likeCheckbox := tview.NewCheckbox().
						SetLabel(fmt.Sprintf("%d likes", msg.Likes))
					likeCheckbox.SetChangedFunc(func(checked bool) {
						// Klik na Like - pošalji LikeMessage RPC
						go func() {
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()

							_, err := grpcClient.LikeMessage(ctx, &razp.LikeMessageRequest{
								TopicId:   int64(topicIndex) + 1,
								MessageId: msg.Id + 1,
								UserId:    currentUserID,
							})
							app.QueueUpdateDraw(func() {
								if err != nil {
									status.SetTextColor(tcell.ColorRed)
									status.SetText("Failed to like message: " + err.Error())
									return
								}

							})
						}()
					})

					metaText := fmt.Sprintf("by %s | %s", userMap[msg.UserId], msg.CreatedAt.AsTime().Format("15:04"))
					metaView := tview.NewTextView().
						SetText(metaText).
						SetTextAlign(tview.AlignLeft)
					metaView.SetTextColor(tcell.ColorGray)

					messagesGrid.AddItem(messageView, row, 0, 1, 1, 0, 0, false)
					messagesGrid.AddItem(metaView, row+1, 0, 1, 1, 0, 0, false)
					messagesGrid.AddItem(likeCheckbox, row+2, 0, 1, 1, 0, 0, true)

					row += 3
				}

				rows := make([]int, row)
				for i := range rows {
					if i%2 == 0 {
						rows[i] = 5
					} else {
						rows[i] = 1
					}
				}
				if len(rows) > 0 {
					messagesGrid.SetRows(rows...)
				}

				status.SetTextColor(tcell.ColorGreen)
				status.SetText(fmt.Sprintf("Loaded %d messages", len(resp.Messages)))
			})
		}()
	}

	// FUNKCIJA ZA POŠILJANJE NOVEGA SPOROČILA
	postMessage := func(topicIndex int, messageText string) {
		messageTitle := newMessageTitle.GetText()

		if messageText == "" {
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Message cannot be empty")
			return
		}

		fullMessage := messageText
		if messageTitle != "" {
			fullMessage = "[" + messageTitle + "] " + messageText
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := grpcClient.PostMessage(ctx, &razp.PostMessageRequest{
				TopicId: int64(topicIndex) + 1,
				UserId:  currentUserID,
				Text:    fullMessage,
			})
			app.QueueUpdateDraw(func() {
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText("Failed to post message: " + err.Error())
					return
				}

				status.SetTextColor(tcell.ColorGreen)
				status.SetText("Message posted successfully")
				newMessageTextarea.SetText("", false)
				newMessageTitle.SetText("")

				loadMessagesForTopic(topicIndex)
				pages.SwitchToPage("topicDetail")
				app.SetFocus(messagesGrid)
			})
		}()
	}

	// FUNKCIJE ZA GUMBE

	//CREATE ACCOUNT POJDI NA CREATE ACCOUNT SCREEN
	createBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("signup")
		app.SetFocus(usernameField)
	})

	//LOGIN POJDI NA LOGIN SCREEN
	loginBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("login")
		app.SetFocus(loginUsername)
	})

	//USTVARI RACUN
	createAccountButton.SetSelectedFunc(createAccount)

	//BACK GUMB POJDI NAZAJ NA ZACETEK
	backFromSignupBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("choice")
	})

	backFromLoginBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("choice")
	})

	backFromProfileBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("home")
	})

	//PRIJAVI SE
	loginButton.SetSelectedFunc(loginFunc)

	//OSVEZI TOPICS
	refreshBtn.SetSelectedFunc(loadTopics)

	profileBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("profile")
	})

	profileBtn.SetSelectedFunc(func() {
		updateProfileDisplay()
		pages.SwitchToPage("profile")
	})

	//USTAVRI NOW TOPIC
	createTopicBtn.SetSelectedFunc(func() {
		topicName := newTopicInput.GetText()
		if topicName == "" {
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Topic name cannot be empty")
		}
		createTopic(topicName)
	})

	// FUNKCIJA ZA KLIK NA TOPIC
	topicList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		selectedTopicIndex = index
		loadMessagesForTopic(index)
	})

	newMessageBtn.SetSelectedFunc(func() {
		newMessageTextarea.SetText("", false)
		pages.SwitchToPage("newMessage")
	})

	// BUTTON AKCIJE ZA NEW MESSAGE SCREEN
	postNewMessageBtn.SetSelectedFunc(func() {
		messageText := newMessageTextarea.GetText()
		postMessage(selectedTopicIndex, messageText)
		pages.SwitchToPage("home")
	})

	cancelNewMessageBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("home")
	})

	logoutBtn.SetSelectedFunc(func() {
		pages.SwitchToPage("choice")
		currentUsername = ""
		currentUserID = -1
	})

	exitBtn.SetSelectedFunc(func() {
		app.Stop()
	})

	changePasswordBtn.SetSelectedFunc(func() {
		oldPass := profileOldPassword.GetText()
		newPass := profileNewPassword.GetText()

		if oldPass == "" || newPass == "" {
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Both passwords required")
			return
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := grpcClient.FindUser(ctx, &razp.FindUserRequest{
				Name:     currentUsername,
				Password: oldPass,
			})

			app.QueueUpdateDraw(func() {
				if err != nil {
					status.SetTextColor(tcell.ColorRed)
					status.SetText("Old password is incorrect")
					return
				}

				status.SetTextColor(tcell.ColorGreen)
				status.SetText("Password change: waiting for server implementation")

				profileOldPassword.SetText("")
				profileNewPassword.SetText("")
			})
		}()
	})

	subsribeTopic.SetSelectedFunc(func() {
		if selectedTopicIndex == -1 {
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Please select a topic first!")
			return
		}

		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			stream, err := grpcClient.SubscribeTopic(ctx, &razp.SubscribeTopicRequest{
				TopicId:       int64(selectedTopicIndex) + 1,
				UserId:        currentUserID,
				FromMessageId: 0,
			})

			if err != nil {
				app.QueueUpdateDraw(func() {
					status.SetTextColor(tcell.ColorRed)
					status.SetText("Subscribe error: " + err.Error())
				})
				return
			}

			app.QueueUpdateDraw(func() {
				status.SetTextColor(tcell.ColorGreen)
				status.SetText("Subscribed! Waiting for new messages...")
			})

			for {
				ev, err := stream.Recv()
				if err != nil {
					app.QueueUpdateDraw(func() {
						status.SetTextColor(tcell.ColorYellow)
						status.SetText("Subscription ended")
					})
					break
				}

				app.QueueUpdateDraw(func() {
					loadMessagesForTopic(selectedTopicIndex)
					status.SetTextColor(tcell.ColorGreen)
					status.SetText("New message received!")
				})

				_ = ev
			}
		}()
	})

	//ČE PRITISNEMO ENTER

	//USERNAME->PASSWORD
	usernameField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			app.SetFocus(passwordField)
		}
	})

	//PASSWORD->USTVARI RACUN
	passwordField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			createAccount()
		}
	})

	//USERNAME->PASSWORD PRI LOGIN
	loginUsername.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			app.SetFocus(loginPassword)
		}
	})

	//PASSWORD->PRIJAVI SE
	loginPassword.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			loginFunc()
		}
	})

	newTopicInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			topicName := newTopicInput.GetText()
			if topicName == "" {
				status.SetTextColor(tcell.ColorRed)
				status.SetText("Topic name cannot be empty")
				return
			}
			createTopic(topicName)
		}
	})

	//REGISTTRIRAMO SCREENE
	pages.AddPage("choice", choiceFlex, true, true)
	pages.AddPage("signup", signupFlex, true, false)
	pages.AddPage("login", loginFlex, true, false)
	pages.AddPage("home", contentFlex, true, false)
	pages.AddPage("profile", profileFlex, true, false)
	pages.AddPage("newMessage", newMessageFlex, true, false)

	//ZACNI APLIKACIJO
	app.SetRoot(pages, true).SetFocus(createBtn)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
