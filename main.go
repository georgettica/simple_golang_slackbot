package main

import (
  "fmt"
  "os"
  "regexp"
  "strings"

  "github.com/nlopes/slack"
)

func getenv(name string) string {
  v := os.Getenv(name)
  if v == "" {
    panic("missing required environment variable " + name)
  }
  return v
}

func main() {
  token := getenv("SLACKTOKEN")
  api := slack.New(token)
  rtm := api.NewRTM()
  go rtm.ManageConnection()

Loop:
  for {
    select {
    case msg := <-rtm.IncomingEvents:
      fmt.Println("Event Received:")
      switch ev := msg.Data.(type) {

      case *slack.MessageEvent:
        info := rtm.GetInfo()

        text := ev.Text
        text = strings.TrimSpace(text)
        text = strings.ToLower(text)

        prefix := ",u,"
        matchedUser, _ := regexp.MatchString("^" + prefix ,text)

        if ev.User != info.User.ID && matchedUser {
          text = text[len(prefix):]
          text = strings.TrimSpace(text)
          splits := strings.Split(text, ":")
          if len(splits) != 2 {
            panic("Not a valid user command")
          }
          user, state := splits[0], splits[1]
          fmt.Printf("%s ::: %s",user, state)
          rtm.SendMessage(rtm.NewOutgoingMessage("\\[T]/ Praise *the* Sun \\[T]/", ev.Channel))
        }

      case *slack.RTMError:
        fmt.Printf("Error: %s\n", ev.Error())

      case *slack.InvalidAuthEvent:
        fmt.Printf("Invalid credentials")
        break Loop

      default:
        // Take no action
      }
    }
  }
}

