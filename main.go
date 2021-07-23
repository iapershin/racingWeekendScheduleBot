package main

import (
  "strings"
  "net/http"
  "log"
  "io/ioutil"
  "os"
  "crypto/tls"
  //"fmt"
  "github.com/PuerkitoBio/goquery"
  "github.com/syfaro/telegram-bot-api"
  //"time"
  //"github.com/go-co-op/gocron"
)


func makeRequest() string {
  tr := &http.Transport{
          TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
      }
  client := &http.Client{Transport: tr}
  resp, err := client.Get("https://vk.com/world_of_speed")
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

func sendAnnounce(announceText string, bot *tgbotapi.BotAPI) string{
  msg := tgbotapi.NewMessage(198952278, announceText)
  bot.Send(msg)
  return "Message sent"
}

func cronTask(bot *tgbotapi.BotAPI) {
  sendAnnounce(formatter(parsePage(makeRequest())),bot)
}

func main() {
  bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
    if err != nil {
      log.Panic(err)
    }
  bot.Debug = true
  log.Printf("Authorized on account %s", bot.Self.UserName)
  //scheduler := gocron.NewScheduler(time.UTC)
  //scheduler.Every(1).Friday().At("16:00").Do(cronTask)
  //scheduler.StartBlocking()
  cronTask(bot)
  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60
  updates, err := bot.GetUpdatesChan(u)
  for update := range updates {
  	if update.Message == nil { // ignore any non-Message Updates
  		continue
  	}
  	if update.Message.Text == "/subscribe" {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are subscribed!")
        bot.Send(msg)
        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
    }
    if update.Message.Text == "/unsubscribe" {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are unsubscribed")
        bot.Send(msg)
        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
    }
  }

}