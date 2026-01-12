package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	razp "razpravljalnica/razpravljalnica"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := razp.NewMessageBoardClient(conn)

	// ============================================================================
	// UPORABNIKI
	// ============================================================================

	users := []struct {
		name     string
		password string
	}{
		{"Alice", "pass123"}, {"Bob", "pass456"}, {"Charlie", "pass789"},
		{"Diana", "pass000"}, {"Eve", "pass111"}, {"Adam", "pass666"},
		{"Mike", "pass3456"}, {"Loti", "pass98765"}, {"Rick", "pass7667"},
		{"Butter", "pass87654"},
	}

	userIDs := make(map[string]int64)
	userList := []string{}

	fmt.Println("\n========== USTVARJANJE UPORABNIKOV ==========")
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
		userList = append(userList, u.name)
		fmt.Printf("✅ Created user: %s (ID: %d)\n", u.name, user.Id)
	}

	// ============================================================================
	// TEME
	// ============================================================================

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

	fmt.Println("\n========== USTVARJANJE TEM ==========")
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

	// ============================================================================
	// SPOROČILA - veliko več sporočil za vsako temo
	// ============================================================================

	messages := map[string][]struct {
		user string
		text string
	}{
		"Programing": {
			{"Alice", "[Go Tutorial] Kako se naučiti Go programskega jezika - začetni vodič"},
			{"Bob", "[Go] Goroutines so res uporabne za concurrent programming!"},
			{"Charlie", "[Python] Best practices za Python razvoj v 2024"},
			{"Diana", "[Python] Django vs Flask - kateri framework izbrati?"},
			{"Eve", "[JavaScript] Kaj je novo v ES2024 - top 5 features"},
			{"Adam", "[JavaScript] React 19 je izšel! Kdo ga že uporablja?"},
			{"Mike", "[Rust] Performance comparison: Rust vs C++ - moji benchmarki"},
			{"Loti", "[Rust] Začenjam se učiti Rust, kakšni so vaši nasveti?"},
			{"Rick", "[Database] SQL optimization tips za velike baze"},
			{"Butter", "[Database] PostgreSQL vs MySQL - primerjava"},
			{"Alice", "[Go] Kako pravilno strukturirati Go projekt?"},
			{"Bob", "[TypeScript] Zakaj bi morali vsi uporabljati TypeScript"},
			{"Charlie", "[DevOps] Docker best practices za produkcijo"},
			{"Diana", "[Cloud] AWS vs Azure vs GCP - moje izkušnje"},
			{"Eve", "[Security] Kako zaščititi API pred napadi"},
		},
		"Games": {
			{"Bob", "[Gaming News] Nova igra leta - nominations so zunaj!"},
			{"Charlie", "[Elden Ring] Best build guide za začetnike"},
			{"Diana", "[Elden Ring] Shadow of the Erdtree DLC je fenomenalen!"},
			{"Eve", "[Fortnite] Sezona 8 je tu! Novi battle pass"},
			{"Adam", "[Chess] Tips za izboljšanje ratinga - od 1000 do 1500"},
			{"Mike", "[Chess] Magnus Carlsen spet dokazal da je najboljši"},
			{"Loti", "[Board Games] Top 10 igre tega leta za družinske večere"},
			{"Rick", "[Board Games] Catan turnir v Ljubljani - kdo gre?"},
			{"Butter", "[PS5] Najboljše ekskluzive za PlayStation"},
			{"Alice", "[Xbox] Game Pass je res vredn svojega denarja"},
			{"Bob", "[PC Gaming] Kako sestaviti gaming PC v 2024"},
			{"Charlie", "[VR] Meta Quest 3 review po 6 mesecih uporabe"},
			{"Diana", "[Nintendo] Zelda TOTK je mojstrovina!"},
			{"Eve", "[Indie Games] Hades 2 early access - prve impresije"},
			{"Adam", "[Retro] Najboljše retro igre vseh časov"},
		},
		"Phones": {
			{"Charlie", "[AI] ChatGPT nova verzija - GPT-5 rumors"},
			{"Diana", "[AI] Lokalni LLM modeli na telefonu - je to možno?"},
			{"Eve", "[Hardware] Best laptops 2024 za programerje"},
			{"Adam", "[Hardware] M3 MacBook Pro vs ThinkPad - primerjava"},
			{"Mike", "[Phone] iPhone 16 review - vredno nadgradnje?"},
			{"Loti", "[Phone] Samsung S24 Ultra camera test"},
			{"Rick", "[IoT] Smart home setup guide - od začetka"},
			{"Butter", "[IoT] Home Assistant vs Google Home"},
			{"Alice", "[5G] Kako deluje 5G omrežje - razlaga"},
			{"Bob", "[Wearables] Apple Watch vs Galaxy Watch"},
			{"Charlie", "[Tablets] iPad Pro M4 - ali nadomesti laptop?"},
			{"Diana", "[Audio] Najboljše brezžične slušalke 2024"},
			{"Eve", "[Photography] Smartphone photography tips"},
			{"Adam", "[Privacy] Kako zaščititi svojo zasebnost na telefonu"},
			{"Mike", "[Apps] Najboljše produktivnostne aplikacije"},
		},
		"Travel": {
			{"Diana", "[Europe] Best destinations v Evropi za poletje"},
			{"Eve", "[Europe] Interrail potovanje - moje izkušnje"},
			{"Adam", "[Asia] Backpacking through Southeast Asia - 3 meseci"},
			{"Mike", "[Asia] Japan travel guide - what to know"},
			{"Loti", "[America] Road trip planing tips - Route 66"},
			{"Rick", "[America] NYC v 5 dneh - itinerary"},
			{"Butter", "[Budget] Cheap travel hacks ki jih morate poznati"},
			{"Alice", "[Budget] Kako potovati s 30€ na dan"},
			{"Bob", "[Visa] Travel visa guide za Slovence"},
			{"Charlie", "[Visa] Digital nomad vize - katere države jih nudijo"},
			{"Diana", "[Adventure] Hiking v Alpah - najboljše poti"},
			{"Eve", "[Beach] Top 10 plaž v Evropi"},
			{"Adam", "[City Break] Weekend v Barceloni"},
			{"Mike", "[Winter] Skiing v Avstriji - primerjava smučišč"},
			{"Loti", "[Culture] Kulturne znamenitosti Italije"},
		},
		"Food": {
			{"Eve", "[Recipe] Kako narediti perfect pasta carbonara"},
			{"Adam", "[Recipe] Domač kruh - enostaven recept"},
			{"Mike", "[Cooking] Top 5 kuharskih trikov profesionalcev"},
			{"Loti", "[Cooking] Air fryer recepti - best of"},
			{"Rick", "[Diet] Healthy meal prep ideas za cel teden"},
			{"Butter", "[Diet] Intermittent fasting - moje izkušnje"},
			{"Alice", "[Restaurant] Best Italian restaurants v Ljubljani"},
			{"Bob", "[Restaurant] Michelin restavracije v Sloveniji"},
			{"Charlie", "[Baking] Homemade bread tutorial za začetnike"},
			{"Diana", "[Baking] Torte za posebne priložnosti"},
			{"Eve", "[Vegan] Najboljši veganski recepti"},
			{"Adam", "[BBQ] Kako pripraviti perfect steak"},
			{"Mike", "[Asian] Domači sushi - vodič"},
			{"Loti", "[Dessert] Čokoladna torta recept"},
			{"Rick", "[Drinks] Domači koktejli za poletje"},
		},
		"Animals": {
			{"Adam", "[Dogs] Najboljše pasme psov za družine"},
			{"Mike", "[Dogs] Kako vzgojiti mladička"},
			{"Loti", "[Cats] Mačke vs psi - večna debata"},
			{"Rick", "[Cats] Najboljša hrana za mačke"},
			{"Butter", "[Exotic] Eksotične živali kot hišni ljubljenčki"},
			{"Alice", "[Birds] Papige - kako jih pravilno skrbeti"},
			{"Bob", "[Fish] Akvarij za začetnike - setup guide"},
			{"Charlie", "[Wildlife] Safari v Afriki - izkušnje"},
			{"Diana", "[Pets] Kako izbrati pravega hišnega ljubljenčka"},
			{"Eve", "[Health] Veterinarski nasveti za zdrave živali"},
			{"Adam", "[Training] Kako naučiti psa trike"},
			{"Mike", "[Adoption] Posvojitev živali iz zavetišča"},
			{"Loti", "[Funny] Smešni posnetki živali"},
			{"Rick", "[Nature] Divje živali v Sloveniji"},
			{"Butter", "[Care] Nega starejših hišnih ljubljenčkov"},
		},
		"Movies": {
			{"Mike", "[Review] Dune Part 2 - najboljši film leta?"},
			{"Loti", "[Review] Oppenheimer - Nolanovo mojstrovina"},
			{"Rick", "[Marvel] MCU prihodnost - kaj nas čaka"},
			{"Butter", "[Marvel] Deadpool 3 - prve reakcije"},
			{"Alice", "[Horror] Najboljši horror filmi 2024"},
			{"Bob", "[Comedy] Top komedije za filmski večer"},
			{"Charlie", "[Sci-Fi] Blade Runner je še vedno najboljši"},
			{"Diana", "[Drama] Oscar nominacije - napovedi"},
			{"Eve", "[Animation] Miyazaki in Studio Ghibli"},
			{"Adam", "[Series] Najboljše serije na Netflixu"},
			{"Mike", "[Series] House of the Dragon vs Rings of Power"},
			{"Loti", "[Documentary] Dokumentarci ki jih morate videti"},
			{"Rick", "[Classic] Klasični filmi ki jih vsak mora videti"},
			{"Butter", "[Streaming] Kateri streaming service izbrati"},
			{"Alice", "[Cinema] IMAX vs Dolby Cinema - primerjava"},
		},
		"Funny": {
			{"Loti", "[Meme] Najboljši memi tega tedna 😂"},
			{"Rick", "[Meme] Programming humor - samo programerji razumejo"},
			{"Butter", "[Joke] Zakaj programerji ne marajo narave? Preveč bugov!"},
			{"Alice", "[Fail] Epic fail compilation"},
			{"Bob", "[Animals] Smešne živali - garantiran smeh"},
			{"Charlie", "[Dad Jokes] Najboljše očetovske šale"},
			{"Diana", "[Puns] Besedne igre ki te nasmejijo"},
			{"Eve", "[Story] Smešna zgodba iz službe"},
			{"Adam", "[Reddit] Najboljši subredditi za humor"},
			{"Mike", "[Video] Smešni videi ki jih morate videti"},
			{"Loti", "[Comic] Webcomics ki jim sledim"},
			{"Rick", "[Sarcasm] Sarkazem na najvišji ravni"},
			{"Butter", "[Wholesome] Wholesome content za dober dan"},
			{"Alice", "[Cringe] Tako cringe da je smešno"},
			{"Bob", "[Office] Smešne situacije iz pisarne"},
		},
		"Hobbies": {
			{"Rick", "[Photography] Kako začeti s fotografijo"},
			{"Butter", "[Photography] Lightroom vs Photoshop za editing"},
			{"Alice", "[Music] Kako se naučiti kitaro - za začetnike"},
			{"Bob", "[Music] Najboljši DAW za produkcijo glasbe"},
			{"Charlie", "[Art] Digital art za začetnike - orodja"},
			{"Diana", "[Art] Procreate tips in triki"},
			{"Eve", "[DIY] DIY projekti za dom"},
			{"Adam", "[DIY] 3D printing kot hobi"},
			{"Mike", "[Sports] Kako začeti teči - couch to 5k"},
			{"Loti", "[Sports] Kolesarstvo v Sloveniji - najboljše poti"},
			{"Rick", "[Reading] Knjige ki so mi spremenile življenje"},
			{"Butter", "[Reading] E-reader vs fizične knjige"},
			{"Alice", "[Gaming] Retro gaming kot hobi"},
			{"Bob", "[Collecting] Zbirateljstvo - kaj zbirate?"},
			{"Charlie", "[Gardening] Vrtnarjenje za začetnike"},
		},
	}

	// Shranimo ID-je sporočil za kasnejše všečkanje
	type MessageInfo struct {
		TopicID   int64
		MessageID int64
		TopicName string
	}
	postedMessages := []MessageInfo{}

	fmt.Println("\n========== OBJAVLJANJE SPOROČIL ==========")
	messageCount := 0
	for topicName, msgs := range messages {
		topicID := topicIDs[topicName]

		for _, msg := range msgs {
			userID := userIDs[msg.user]

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			postedMsg, err := client.PostMessage(ctx, &razp.PostMessageRequest{
				TopicId: topicID,
				UserId:  userID,
				Text:    msg.text,
			})
			cancel()

			if err != nil {
				log.Printf("❌ Failed to post message: %v", err)
				continue
			}

			postedMessages = append(postedMessages, MessageInfo{
				TopicID:   topicID,
				MessageID: postedMsg.Id,
				TopicName: topicName,
			})

			messageCount++
			fmt.Printf("✅ [%s] %s: %.50s...\n", topicName, msg.user, msg.text)
		}
	}

	// ============================================================================
	// POVZETEK
	// ============================================================================

	fmt.Println("\n========== POVZETEK ==========")
	fmt.Printf("✅ Ustvarjenih uporabnikov: %d\n", len(userIDs))
	fmt.Printf("✅ Ustvarjenih tem: %d\n", len(topicIDs))
	fmt.Printf("✅ Objavljenih sporočil: %d\n", messageCount)
	fmt.Println("==============================")
}
