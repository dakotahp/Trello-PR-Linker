package main

import (
  "log"
  "net/http"
  "os"
  "fmt"
  "regexp"

  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "github.com/adlio/trello"
  "gopkg.in/rjz/githubhook.v0"
)

type Head struct {
  Ref string `json:"ref"`
}

type PullRequest struct {
  HtmlUrl string `json:"html_url"`
  Title string `json:"title"`
  Head `json:"head"`
}

type Payload struct {
  Action string `json:"action"`
  PullRequest `json:"pull_request"`
}

func main() {
  port := os.Getenv("PORT")
  secret := []byte(os.Getenv("SECRET_TOKEN"))

  if port == "" {
    log.Fatal("$PORT must be set")
  }

  router := gin.New()
  router.Use(gin.Logger())
  router.LoadHTMLGlob("templates/*.tmpl.html")
  router.Static("/static", "static")

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl.html", nil)
  })

  router.POST("/webhook", func(c *gin.Context) {
    var pr Payload

    hook, err := githubhook.Parse(secret, c.Request)

    if err != nil {
      fmt.Print("Error with secure webhook:", err)
    }


    if pr.Action == "opened" {
      c.BindJSON(&pr)

      fmt.Println("Running on PR:", pr.PullRequest.HtmlUrl)
      cardId := trelloIdFromTitle(pr.PullRequest.Title)
      postPrLinkToTrelloCard(cardId, pr.PullRequest.HtmlUrl)
    }

    c.JSON(http.StatusOK, gin.H{"url": pr.PullRequest.HtmlUrl, "status": http.StatusOK})
  })

  router.Run(":" + port)
}

func trelloIdFromTitle(title string) (string) {
  re := regexp.MustCompile(`\[([A-Za-z0-9]{8})\]`)
  return re.FindStringSubmatch(title)[1]
}

func trelloIdFromBranch(branch string) (string) {
  re := regexp.MustCompile(`\[([A-Za-z0-9]{8})\]`)
  return re.FindStringSubmatch(branch)[1]
}

func postPrLinkToTrelloCard(cardId string, url string) {
  appKey := os.Getenv("TRELLO_TOKEN")
  token := os.Getenv("TRELLO_KEY")

  client := trello.NewClient(appKey, token)

  card, err := client.GetCard(cardId, trello.Defaults())
  if err != nil {
    log.Fatal(err)
  }

  attachment := trello.Attachment {
    Name: "PR",
    URL: url,
  }

  cardErr := card.AddURLAttachment(&attachment)
  if cardErr != nil {
    fmt.Println(cardErr)
  }
}
