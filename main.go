package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
	tgbotapi "github.com/syfaro/telegram-bot-api"
)

func makeRequest() string {
	resp, err := http.Get("https://vk.com/world_of_speed")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return sb
}

func parsePage(html string) string {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}
	result := dom.Find(".pi_text").First().Text()
	return result
}

func formatter(content string) string {
	expanderText := [2]string{"Показать полностью...", "See more"}
	for _, text := range expanderText {
		content = strings.ReplaceAll(content, text, "")
	}
	formatted := strings.ReplaceAll(content, "#", "\n\n#")
	return formatted
}

func sendAnnounce(announceText string, bot *tgbotapi.BotAPI, userList []int64) string {
	for _, user := range userList {
		msg := tgbotapi.NewMessage(user, announceText)
		bot.Send(msg)
		log.Println("New announce sent:")
		log.Println(announceText)
	}
	return "Message sent"
}

var (
	db_host     = os.Getenv("DB_HOST")
	db_port     = os.Getenv("DB_PORT")
	db_user     = os.Getenv("DB_USER")
	db_password = os.Getenv("DB_PASSWORD")
	db_name     = os.Getenv("DB_NAME")
)

func checkDBConnection(psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Unable to connect to database with following parameters: " + psqlInfo)
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println("Unable to connect to database with following parameters: " + psqlInfo)
		panic(err)
	}
	log.Println("Successfully connected to Database with following parameters: " + psqlInfo)
}

func checkIfUserExists(psqlInfo string, chatID int64) bool {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	sqlStatement := `SELECT * FROM public.users WHERE id=($1)`
	result := db.QueryRow(sqlStatement, chatID)
	err = result.Scan(&chatID)
	db.Close()
	switch err {
	case sql.ErrNoRows:
		log.Printf("User %d doesn't exists", chatID)
		return false
	case nil:
		log.Printf("User %d exists", chatID)
		return true
	default:
		panic(err)
	}
}

func addUserToDB(psqlInfo string, chatID int64) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	sqlStatement := `
  		INSERT INTO public.users (id)
  		VALUES ($1)`
	db.Exec(sqlStatement, chatID)
	db.Close()
}

func deleteUserFromDB(psqlInfo string, chatID int64) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	sqlStatement := `
  		DELETE FROM public.users
  		WHERE id=($1)`
	db.Exec(sqlStatement, chatID)
	db.Close()
}

func getUsersList(psqlInfo string) []int64 {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	sqlStatement := `SELECT id FROM public.users`
	rows, err := db.Query(sqlStatement)
	userList := []int64{}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		userList = append(userList, id)
	}
	db.Close()
	return userList
}

func main() {
	// establish psql session
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		db_host, db_port, db_user, db_password, db_name)
	checkDBConnection(psqlInfo)

	// bot initialization
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// scheduler block
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Friday().At("18:00").Do(func() {
		sendAnnounce(formatter(parsePage(makeRequest())), bot, getUsersList(psqlInfo))
	})
	scheduler.StartAsync()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// chat handler
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch update.Message.Text {
		case "/subscribe":
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if checkIfUserExists(psqlInfo, update.Message.Chat.ID) == false {
				addUserToDB(psqlInfo, update.Message.Chat.ID)
				log.Printf("New user %d subscribed", update.Message.Chat.ID)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are subscribed!")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are already subscribed")
				bot.Send(msg)
			}
		case "/unsubscribe":
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			deleteUserFromDB(psqlInfo, update.Message.Chat.ID)
			log.Printf("User %d unsubscribed", update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are unsubscribed")
			bot.Send(msg)
		default:
			const title = "Hi! I'm Racing Weekend Schedule Bot\n" +
				"To subscribe on weekly updates type /subscribe.\n" +
				"To unsubscribe from updates type /unsubscribe\n" +
				"Type whatever you want to see this message"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, title)
			bot.Send(msg)
		}
	}
}
