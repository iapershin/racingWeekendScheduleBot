package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	//"os"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/syfaro/telegram-bot-api"

	//"time"
	//"github.com/go-co-op/gocron"
	"database/sql"

	_ "github.com/lib/pq"
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
	content = strings.Replace(content, "Показать полностью...", "", -1)
	formatted := strings.Replace(content, "#", "\n\n#", -1)
	return formatted
}

func sendAnnounce(announceText string, bot *tgbotapi.BotAPI, userList []int64) string {
	for _, user := range userList {
		msg := tgbotapi.NewMessage(user, announceText)
		bot.Send(msg)
	}
	return "Message sent"
}

func cronTask(bot *tgbotapi.BotAPI, userList []int64) {
	sendAnnounce(formatter(parsePage(makeRequest())), bot, userList)
}

const (
	db_host     = "//"
	db_port     = 5432
	db_user     = "postgres"
	db_password = "//"
	db_name     = "bot-rwb"
)

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
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_name)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	//scheduler := gocron.NewScheduler(time.UTC)
	//scheduler.Every(1).Friday().At("16:00").Do(cronTask)
	//scheduler.StartBlocking()
	//cronTask(bot, getUsersList(psqlInfo))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Text == "/subscribe" {
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
		}
		if update.Message.Text == "/unsubscribe" {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			deleteUserFromDB(psqlInfo, update.Message.Chat.ID)
			log.Printf("User %d unsubscribed", update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are unsubscribed")
			bot.Send(msg)
		}
	}
}
