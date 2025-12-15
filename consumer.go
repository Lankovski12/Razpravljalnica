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

				_, err := grpcClient.FindUser(ctx, &razp.CreateUserRequest{
					Name:     name,
					Password: pass,
				})
				app.QueueUpdateDraw(func() {
					if err != nil {
						status.SetTextColor(tcell.ColorRed)
						status.SetText(err.Error())
					} else {
						status.SetTextColor(tcell.ColorYellow)
						status.SetText("Logging in...")
					}
				})
			}()
		}

		// pokazi poste, pojdi na drug page
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

	app.SetRoot(pages, true).SetFocus(createBtn)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
