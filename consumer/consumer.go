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

	userMap = make(map[int]string)

	statusTicker *time.Ticker

	loadAllMessages      func()
	loadMessagesForTopic func(topicIndex int)
	updateProfileDisplay func()
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

func showStatusMessage(app *tview.Application, status *tview.TextView, message string, color tcell.Color) {
	if statusTicker != nil {
		statusTicker.Stop()
	}

	status.SetTextColor(color)
	status.SetText(message)

	statusTicker = time.NewTicker(2 * time.Second)
	go func() {
		<-statusTicker.C
		app.QueueUpdateDraw(func() {
			status.SetText("")
		})
		statusTicker.Stop()
		statusTicker = nil
	}()
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

	messagesTable := tview.NewList().
		ShowSecondaryText(true)
	messagesTable.SetBackgroundColor(tcell.ColorDefault)

	messagesDetailAll := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messagesTable, 0, 1, true).
		AddItem(status, 1, 1, false).
		AddItem(subsribeTopic, 3, 1, false).
		AddItem(newMessageBtn, 3, 1, false)
	messagesDetailAll.SetBorder(true).SetTitle(" Home ")

	contentFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(topicsPanel, 0, 1, false).
		AddItem(messagesDetailAll, 0, 3, false)

	//SCREEN 7: NEW MESSAGE

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
		SetOptions([]string{"Pink", "Green", "Orange", "Blue", "Red", "Violet"}, nil)
	dropdown.SetLabelColor(tcell.ColorWhite)
	dropdown.SetFieldBackgroundColor(tcell.ColorBlack)
	dropdown.SetFieldTextColor(tcell.ColorWhite)

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

	// Use List for profile messages (same style as main page)
	profileMessagesList := tview.NewList().
		ShowSecondaryText(true)
	profileMessagesList.SetBackgroundColor(tcell.ColorDefault)

	profileMessagesPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(profileMessagesList, 0, 1, true)
	profileMessagesPanel.SetBorder(true).SetTitle(" My Messages ")

	profileFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(profilePanel, 0, 1, false).
		AddItem(profileMessagesPanel, 0, 2, true)
	profileFlex.SetBorder(true).SetTitle(" Profile ")

	//FUNKCIJE

	//za profil
	updateProfileDisplay := func() {
		profileUsernameDisplay.SetText("Username: " + currentUsername)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.GetUserMessages(ctx, &razp.GetUserMessagesRequest{
				UserId: currentUserID,
			})
			app.QueueUpdateDraw(func() {
				profileMessagesList.Clear()

				if err != nil {
					showStatusMessage(app, status, "Failed to load messages", tcell.ColorRed)
					return
				}

				if len(resp.Messages) == 0 {
					profileMessagesList.AddItem("No messages yet.", "", 0, nil)
					return
				}

				for _, msg := range resp.Messages {
					// Header with topic name
					header := fmt.Sprintf("[:/%s]", msg.TopicName)

					// Parse title from text if present
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
							header = fmt.Sprintf("[:/%s | %s]", msg.TopicName, title)
						}
					}

					// Message item (header + text)
					profileMessagesList.AddItem(header, text, 0, nil)

					// Time and likes
					profileMessagesList.AddItem(fmt.Sprintf("%s • %d likes", msg.CreatedAt.AsTime().Format("02.01.2006 15:04"), msg.Likes), "", 0, nil)

					// Edit button - clickable
					topicID := msg.TopicId
					messageID := msg.MessageId
					msgText := msg.Text
					profileMessagesList.AddItem("[Edit]", "", 0, func() {
						// Create edit modal
						editTextarea := tview.NewTextArea()
						editTextarea.SetText(msgText, true)
						editTextarea.SetBorder(true).SetTitle(" Edit Message ")

						saveBtn := tview.NewButton("Save")
						cancelBtn := tview.NewButton("Cancel")

						editGrid := tview.NewGrid().
							SetRows(0, 3).
							SetColumns(0, 15, 15, 0).
							SetBorders(false)
						editGrid.AddItem(editTextarea, 0, 0, 1, 4, 0, 0, true)
						editGrid.AddItem(saveBtn, 1, 1, 1, 1, 0, 0, false)
						editGrid.AddItem(cancelBtn, 1, 2, 1, 1, 0, 0, false)

						editFlex := tview.NewFlex().
							SetDirection(tview.FlexRow).
							AddItem(editGrid, 0, 1, true)
						editFlex.SetBorder(true).SetTitle(" Edit Message ")

						pages.AddPage("editMessage", editFlex, true, true)

						saveBtn.SetSelectedFunc(func() {
							newText := editTextarea.GetText()
							go func() {
								ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
								defer cancel()

								_, err := grpcClient.UpdateMessage(ctx, &razp.UpdateMessageRequest{
									TopicId:   topicID,
									MessageId: messageID + 1,
									UserId:    currentUserID,
									Text:      newText,
								})
								app.QueueUpdateDraw(func() {
									pages.RemovePage("editMessage")
									if err != nil {
										showStatusMessage(app, status, "Failed to update message", tcell.ColorRed)
										return
									}
									showStatusMessage(app, status, "Message updated!", tcell.ColorGreen)
									updateProfileDisplay()
								})
							}()
						})

						cancelBtn.SetSelectedFunc(func() {
							pages.RemovePage("editMessage")
						})

						app.SetFocus(editTextarea)
					})

					// Delete button - clickable
					profileMessagesList.AddItem("[Delete]", "", 0, func() {
						// Confirm delete
						modal := tview.NewModal().
							SetText("Are you sure you want to delete this message?").
							AddButtons([]string{"Delete", "Cancel"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								if buttonLabel == "Delete" {
									go func() {
										ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
										defer cancel()

										_, err := grpcClient.DeleteMessage(ctx, &razp.DeleteMessageRequest{
											TopicId:   topicID,
											MessageId: messageID + 1,
											UserId:    currentUserID,
										})
										app.QueueUpdateDraw(func() {
											pages.RemovePage("deleteConfirm")
											if err != nil {
												showStatusMessage(app, status, "Failed to delete message", tcell.ColorRed)
												return
											}
											showStatusMessage(app, status, "Message deleted!", tcell.ColorGreen)
											updateProfileDisplay()
										})
									}()
								} else {
									pages.RemovePage("deleteConfirm")
								}
							})
						pages.AddPage("deleteConfirm", modal, true, true)
					})

					// Separator
					profileMessagesList.AddItem(" ", "", 0, nil)
				}
			})
		}()
	}

	//za topics na levi
	loadTopics := func() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.ListTopics(ctx, &emptypb.Empty{})
			app.QueueUpdateDraw(func() {
				topicList.Clear()
				if err != nil {
					showStatusMessage(app, status, "Failed to load topics", tcell.ColorRed)
					return
				}

				topicList.AddItem("Home", "", 0, nil)

				for _, t := range resp.Topics {
					topicList.AddItem(":/ "+t.Name, "", 0, nil)
				}

			})
		}()
	}

	//sharni userje
	loadUsers := func() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.ListUsers(ctx, &emptypb.Empty{})
			if err != nil {
				return
			}

			for _, user := range resp.Users {
				userMap[int(user.Id)] = user.Name
			}
		}()
	}

	//FUNKCIJA KI USTVARI NOVEGA UPORABNIKA
	createAccount := func() { //doda novega userja
		name := usernameField.GetText()
		pass := passwordField.GetText()
		if name == "" || pass == "" { //potrebno dodati funkcijo ki preverja kako korekten je input
			showStatusMessage(app, status, "Username and password required", tcell.ColorRed)
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
					showStatusMessage(app, status, "Account with that username already exist", tcell.ColorRed)
				} else {
					currentUsername = name
					currentUserID = raz.Id
					pages.SwitchToPage("home")
					loadTopics()
					loadUsers()
					selectedTopicIndex = -1
					loadAllMessages() // Load all messages as feed
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
					showStatusMessage(app, status, "Topic with same name already exists, please choose a different name", tcell.ColorRed)
					return
				}

				showStatusMessage(app, status, "Topic '"+resp.Name+"' created successfully", tcell.ColorGreen)

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
			showStatusMessage(app, status, "Username and password required", tcell.ColorRed)
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
						showStatusMessage(app, status, "Password does not match the username", tcell.ColorRed)
					} else {
						currentUserID = raz.Id
						currentUsername = name
						pages.SwitchToPage("home")
						loadTopics()
						loadUsers()
						selectedTopicIndex = -1
						loadAllMessages() // Load all messages as feed on login
					}
				})
			}()
		}
	}

	// FUNKCIJA ZA NALAGANJE SPOROČIL IZ TOPIC
	loadMessagesForTopic = func(topicIndex int) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := grpcClient.GetMessages(ctx, &razp.GetMessagesRequest{
				TopicId:       int64(topicIndex) + 1,
				FromMessageId: 0,
				Limit:         1000,
			})

			topicsResp, topicErr := grpcClient.ListTopics(ctx, &emptypb.Empty{})
			topicName := ""
			if topicErr == nil && topicIndex < len(topicsResp.Topics) {
				topicName = topicsResp.Topics[topicIndex].Name
			}

			app.QueueUpdateDraw(func() {
				messagesTable.Clear()

				if err != nil {
					showStatusMessage(app, status, "Failed to load messages", tcell.ColorRed)
					return
				}

				if len(resp.Messages) == 0 {
					messagesTable.AddItem("No messages yet.", "", 0, nil)
					return
				}

				for _, msg := range resp.Messages {
					// Skip deleted messages
					if msg.Text == "message deleted by user" {
						continue
					}

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

					username := userMap[int(msg.UserId)]

					header := fmt.Sprintf("[:/%s | %s]", topicName, title)

					// topic | naslov + sporocilo
					messagesTable.AddItem("[white]"+header+"[-]", "[gray]"+fmt.Sprintf("%s on %s", username, msg.CreatedAt.AsTime().Format("01.02.2025 15:04"))+"[-]", 0, nil)

					messagesTable.AddItem(" ", "", 0, nil)

					// ime in ura
					messagesTable.AddItem("[white]"+text+"[-]", "", 0, nil)

					// lajki - clickable
					msgID := msg.Id
					tIdx := topicIndex
					messagesTable.AddItem("[gray]"+fmt.Sprintf("%d likes", msg.Likes)+"[-]", "", 0, func() {
						go func() {
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()

							_, err := grpcClient.LikeMessage(ctx, &razp.LikeMessageRequest{
								TopicId:   int64(tIdx) + 1,
								MessageId: msgID + 1,
								UserId:    currentUserID,
							})
							app.QueueUpdateDraw(func() {
								if err != nil {
									showStatusMessage(app, status, "Failed to like", tcell.ColorRed)
									return
								}
								showStatusMessage(app, status, "Liked!", tcell.ColorGreen)
								loadMessagesForTopic(tIdx)
							})
						}()
					})

					// Empty separator row
					messagesTable.AddItem(" ", "", 0, nil)
				}
			})
		}()
	}

	// FUNKCIJA ZA NALAGANJE VSEH SPOROČIL IZ VSEH TOPICOV
	loadAllMessages = func() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			topicsResp, err := grpcClient.ListTopics(ctx, &emptypb.Empty{})
			if err != nil {
				app.QueueUpdateDraw(func() {
					showStatusMessage(app, status, "Failed to load topics", tcell.ColorRed)
				})
				return
			}

			type messageWithTopic struct {
				msg       *razp.Message
				topicName string
				topicIdx  int
			}
			var allMessages []messageWithTopic

			for i, topic := range topicsResp.Topics {
				msgResp, err := grpcClient.GetMessages(ctx, &razp.GetMessagesRequest{
					TopicId:       topic.Id,
					FromMessageId: 0,
					Limit:         1000,
				})
				if err != nil {
					continue
				}

				for _, msg := range msgResp.Messages {
					if msg.Text != "message deleted by user" {
						allMessages = append(allMessages, messageWithTopic{
							msg:       msg,
							topicName: topic.Name,
							topicIdx:  i,
						})
					}
				}
			}

			app.QueueUpdateDraw(func() {
				messagesTable.Clear()

				if len(allMessages) == 0 {
					messagesTable.AddItem("No messages yet.", "", 0, nil)
					return
				}

				for _, item := range allMessages {
					msg := item.msg
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

					header := fmt.Sprintf("[:/%s | %s]", item.topicName, title)

					messagesTable.AddItem("[white]"+header+"[-]", "[gray]"+fmt.Sprintf("%s on %s", userMap[int(msg.UserId)], msg.CreatedAt.AsTime().Format("01. 02. 2025 15:04"))+"[-]", 0, nil)

					messagesTable.AddItem(" ", "", 0, nil)

					// ime in ura
					messagesTable.AddItem("[white]"+text+"[-]", "", 0, nil)

					// lajki - clickable
					msgID := msg.Id
					tIdx := item.topicIdx
					messagesTable.AddItem("[gray]"+fmt.Sprintf("%d likes", msg.Likes)+"[-]", "", 0, func() {
						go func() {
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()

							_, err := grpcClient.LikeMessage(ctx, &razp.LikeMessageRequest{
								TopicId:   int64(tIdx) + 1,
								MessageId: msgID + 1,
								UserId:    currentUserID,
							})
							app.QueueUpdateDraw(func() {
								if err != nil {
									showStatusMessage(app, status, "Failed to like", tcell.ColorRed)
									return
								}
								showStatusMessage(app, status, "Liked!", tcell.ColorGreen)
								loadAllMessages()
							})
						}()
					})

					// Empty separator row
					messagesTable.AddItem(" ", "", 0, nil)
				}
			})
		}()
	}

	// FUNKCIJA ZA POŠILJANJE NOVEGA SPOROČILA
	postMessage := func(topicIndex int, messageText string) {
		messageTitle := newMessageTitle.GetText()

		if messageText == "" {
			showStatusMessage(app, status, "Message cannot be empty", tcell.ColorRed)
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
					showStatusMessage(app, status, "Failed to post message", tcell.ColorRed)
					return
				}

				showStatusMessage(app, status, "Message posted successfully", tcell.ColorGreen)
				newMessageTextarea.SetText("", false)
				newMessageTitle.SetText("")

				loadMessagesForTopic(topicIndex)
				pages.SwitchToPage("home")
				app.SetFocus(messagesTable)
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
	refreshBtn.SetSelectedFunc(func() {
		loadTopics()
		loadUsers()
		if selectedTopicIndex == -1 {
			loadAllMessages()
		} else {
			loadMessagesForTopic(selectedTopicIndex)
		}
	})

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
			showStatusMessage(app, status, "Topic name cannot be empty", tcell.ColorRed)
		}
		createTopic(topicName)
	})

	// FUNKCIJA ZA KLIK NA TOPIC
	topicList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index == 0 {
			// "All Topics" selected - show all messages
			selectedTopicIndex = -1
			loadAllMessages()
		} else {
			// Specific topic selected (offset by 1 because of "All Topics" at index 0)
			selectedTopicIndex = index - 1
			loadMessagesForTopic(selectedTopicIndex)
		}
	})

	newMessageBtn.SetSelectedFunc(func() {
		if selectedTopicIndex == -1 {
			showStatusMessage(app, status, "Please select a specific topic first!", tcell.ColorRed)
			return
		}
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
			showStatusMessage(app, status, "Both passwords required", tcell.ColorRed)
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
					showStatusMessage(app, status, "Old password is incorrect", tcell.ColorRed)
					return
				}

				_, err1 := grpcClient.ChangePass(ctx, &razp.ChangeUserRequest{
					Id:       currentUserID,
					Password: newPass,
				})

				if err1 != nil {
					showStatusMessage(app, status, "Could't change password", tcell.ColorRed)
					return
				}

				showStatusMessage(app, status, "Password was changed", tcell.ColorGreen)

				profileOldPassword.SetText("")
				profileNewPassword.SetText("")
			})
		}()
	})

	subsribeTopic.SetSelectedFunc(func() {
		if selectedTopicIndex == -1 {
			showStatusMessage(app, status, "Please select a specific topic first!", tcell.ColorRed)
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
					showStatusMessage(app, status, "Subscribe error", tcell.ColorRed)
				})
				return
			}

			app.QueueUpdateDraw(func() {
				showStatusMessage(app, status, "Subscribed!", tcell.ColorGreen)
			})

			for {
				ev, err := stream.Recv()
				if err != nil {
					app.QueueUpdateDraw(func() {
						showStatusMessage(app, status, "Subscription ended", tcell.ColorYellow)
					})
					break
				}

				app.QueueUpdateDraw(func() {
					loadMessagesForTopic(selectedTopicIndex)
					showStatusMessage(app, status, "New message received!", tcell.ColorGreen)
				})

				_ = ev
			}
		}()
	})

	dropdown.SetSelectedFunc(func(text string, index int) {
		colorMap := map[string]tcell.Color{
			"Pink":   tcell.ColorLightPink,
			"Blue":   tcell.ColorBlue,
			"Red":    tcell.ColorRed,
			"Green":  tcell.ColorGreen,
			"Orange": tcell.ColorOrange,
			"Violet": tcell.ColorPurple,
		}

		themeColorSec = colorMap[text]
		updateThemeColors(app, createBtn, loginBtn, createAccountButton, backFromSignupBtn, loginButton, backFromLoginBtn, createTopicBtn, refreshBtn, newMessageBtn, subsribeTopic, postNewMessageBtn, cancelNewMessageBtn, changePasswordBtn, backFromProfileBtn, profileBtn, themeColorSec)
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
				showStatusMessage(app, status, "Topic name cannot be empty", tcell.ColorRed)
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
