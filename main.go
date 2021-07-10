package main

import (
	"strings"
	"net/http"
  "log"
  "io/ioutil"
  //"fmt"
  "github.com/PuerkitoBio/goquery"
  "github.com/Syfaro/telegram-bot-api"
  "time"
  "github.com/go-co-op/gocron"
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

func parsePage(html string) string{
	dom,err:=goquery.NewDocumentFromReader(strings.NewReader(html))
	if err!=nil{
		log.Fatalln(err)
	}
  result := dom.Find(".pi_text").First().Text()
  return result
}

func formatter(content string) string {
  content = strings.Replace(content, "Показать полностью...", "", -1)
  formatted := strings.Replace(content, "#", "\n\n#",-1)
  return formatted
}

func sendAnnounce(announceText string) string{
  bot, err := tgbotapi.NewBotAPI("<token>")
    if err != nil {
      log.Panic(err)
    }
  bot.Debug = true
  log.Printf("Authorized on account %s", bot.Self.UserName)
  msg := tgbotapi.NewMessage(<chatId>, announceText)
  bot.Send(msg)
  return "Message sent"
}

func cronTask() {
  sendAnnounce(formatter(parsePage(makeRequest())))
}

func main() {
  scheduler := gocron.NewScheduler(time.UTC)
  scheduler.Every(1).Friday().At("16:00").Do(cronTask)
  scheduler.StartBlocking()
}