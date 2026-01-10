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
)

/*func updateThemeColors(app *tview.Application, createBtn, loginBtn, createAccountButton, backFromSignupBtn, loginButton, backFromLoginBtn, createTopicBtn, refreshBtn, profileBtn *tview.Button, themeColorSec tcell.Color) {
	createAccountButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	backFromSignupBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	loginButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	backFromLoginBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	createTopicBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))
	refreshBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	profileBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	createBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))
	loginBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	app.Draw()
}*/

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

	//FLEX ZA GUMBA
	buttonsColumn := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(createBtn, 3, 0, true).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(loginBtn, 3, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)

	//CENTRIRANO
	buttonsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonsColumn, 50, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)

	//GLAVNA STRUKTURA
	choiceFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(buttonsFlex, 0, 1, true)
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

	//FLEX ZA GUMBA
	signupButtonsRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(createAccountButton, 25, 0, false).
		AddItem(tview.NewBox(), 0, 0, false).
		AddItem(backFromSignupBtn, 25, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)

	//GLAVNA STRUKTURA
	signupFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(usernameField, 3, 1, true).
		AddItem(passwordField, 3, 1, false).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(signupButtonsRow, 3, 0, false).
		AddItem(status, 2, 1, false)
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

	//Flex za gumba
	loginButtonsRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(loginButton, 25, 0, false).
		AddItem(tview.NewBox(), 0, 0, false).
		AddItem(backFromLoginBtn, 25, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)

	//GLAVNA STRUKTURA
	loginFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginUsername, 3, 1, true).
		AddItem(loginPassword, 3, 1, false).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(loginButtonsRow, 3, 0, false).
		AddItem(status, 2, 1, false)
	loginFlex.SetBorder(true).SetTitle(" Login ")

	//SCREEN 4: MAIN PAGE

	topicList := tview.NewList().ShowSecondaryText(false)
	messagesList := tview.NewList().ShowSecondaryText(false)

	newTopicInput := tview.NewInputField().
		SetLabel("New topic: ").
		SetFieldWidth(20)
	newTopicInput.SetLabelColor(tcell.ColorWhite)
	newTopicInput.SetFieldTextColor(tcell.ColorWhite)

	createTopicBtn := tview.NewButton("Create topic")
	createTopicBtn.SetStyle(tcell.StyleDefault.Foreground(themeColorSec).Background(tcell.ColorWhite))

	refreshBtn := tview.NewButton("Refresh page")
	refreshBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	profileBtn := tview.NewButton("My profile")
	profileBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(themeColorSec))

	inputFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(newTopicInput, 0, 6, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(createTopicBtn, 0, 3, false)

	//FLEX ZA GUMBE
	leva := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(profileBtn, 0, 4, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(refreshBtn, 0, 4, false).
		AddItem(tview.NewBox(), 0, 1, false)

	topicPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topicList, 0, 15, false).
		AddItem(inputFlex, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(leva, 0, 1, false)
	topicPanel.SetBorder(true).SetTitle(" Topics ")

	messagesPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Home:").SetTextAlign(tview.AlignCenter), 1, 0, false).
		AddItem(messagesList, 0, 1, false).
		AddItem(status, 2, 0, false)

	contentFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(topicPanel, 0, 1, false).
		AddItem(messagesPanel, 0, 3, false)

	//SCREEN 5: PROFILE

	profileUsernameDisplay := tview.NewTextView().
		SetText("Username: " + currentUsername).
		SetTextAlign(tview.AlignLeft)
	profileUsernameDisplay.SetTextColor(tcell.ColorWhite)

	// PROFILE SCREEN - Change password fields
	profileOldPassword := tview.NewInputField().
		SetLabel("Old password: ").
		SetMaskCharacter('*').
		SetFieldWidth(20)
	profileOldPassword.SetLabelColor(tcell.ColorWhite)
	profileOldPassword.SetFieldTextColor(tcell.ColorWhite)

	profileNewPassword := tview.NewInputField().
		SetLabel("New password: ").
		SetMaskCharacter('*').
		SetFieldWidth(20)
	profileNewPassword.SetLabelColor(tcell.ColorWhite)
	profileNewPassword.SetFieldTextColor(tcell.ColorWhite)

	// PROFILE SCREEN - Theme selector (dropdown)
	dropdown := tview.NewDropDown().
		SetLabel("Theme: ").
		SetOptions([]string{"Pink", "Green", "Orange", "Blue", "Red", "Violet"}, func(option string, optionIndex int) {
			// Ko izbereš temo, nastavi ustrezne barve
			switch option {
			case "Pink":
				themeColorSec = tcell.ColorPink
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
			/*updateThemeColors(app, createBtn, loginBtn, createAccountButton, backFromSignupBtn, loginButton, backFromLoginBtn, createTopicBtn, refreshBtn, profileBtn, themeColorSec)*/
		})
	dropdown.SetLabelColor(tcell.ColorWhite)

	changePasswordBtn := tview.NewButton("Change Password")
	changePasswordBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen))

	// PROFILE SCREEN - Back button
	backFromProfileBtn := tview.NewButton("Back")
	backFromProfileBtn.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed))

	usernameDisplayFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(profileUsernameDisplay, 2, 0, false)

	// PROFILE SCREEN - Buttons flex
	profileButtonsRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(changePasswordBtn, 18, 0, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(backFromProfileBtn, 8, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)

	// PROFILE SCREEN - Main structure
	profileFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("My Profile").SetTextAlign(tview.AlignCenter), 2, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(usernameDisplayFlex, 2, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(dropdown, 2, 0, false).
		AddItem(tview.NewBox(), 2, 0, false).
		AddItem(profileOldPassword, 3, 0, false).
		AddItem(profileNewPassword, 3, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(profileButtonsRow, 3, 0, false).
		AddItem(status, 2, 0, false)
	profileFlex.SetBorder(true).SetTitle(" Profile ")

	//FUNKCIJE

	//FUNKCIJA KI NALOZI VSE TOPICS
	updateProfileDisplay := func() {
		profileUsernameDisplay.SetText("Username: " + currentUsername)
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
					}
				})
			}()
		}
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
		app.SetFocus(dropdown)
	})

	//USTAVRI NOW TOPIC
	createTopicBtn.SetSelectedFunc(func() {
		topicName := newTopicInput.GetText()
		if topicName == "" {
			status.SetTextColor(tcell.ColorRed)
			status.SetText("Topic name cannot be empty")
			return
		}
		createTopic(topicName)
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

	//ZACNI APLIKACIJO
	app.SetRoot(pages, true).SetFocus(createBtn)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
