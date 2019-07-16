package main

import (
  "log"
  "net/http"
  "os"
  "fmt"
  "regexp"
  "crypto/hmac"
  "crypto/sha1"
  "crypto/subtle"
  "encoding/hex"
  "io/ioutil"
  "encoding/json"

  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "github.com/adlio/trello"
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
  secret := os.Getenv("SECRET_TOKEN")

  if port == "" {
    log.Fatal("$PORT must be set")
  }

  router := gin.New()
  router.Use(gin.Logger())
  router.LoadHTMLGlob("templates/*.tmpl.html")

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl.html", nil)
  })

  router.POST("/webhook", func(c *gin.Context) {
    var pr Payload
    var reqBody []byte
    reqBody, _ = ioutil.ReadAll(c.Request.Body)
    json.Unmarshal(reqBody, &pr)

    if !verifySignature(secret, string(reqBody), c.Request.Header.Get("X-Hub-Signature")) {
      log.Fatal("Signatures didn't match")
    }

    if pr.Action == "opened" {
      fmt.Println("Operating on PR:", pr.PullRequest.HtmlUrl)
      cardId := trelloIdFromTitle(pr.PullRequest.Title)
      postPrLinkToTrelloCard(cardId, pr.PullRequest.HtmlUrl)
    } else {
      fmt.Println("Skipping due to action:", pr.Action)
      fmt.Println("Payload:", pr)
    }

    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
  })

  router.Run(":" + port)
}

func generateSignature(secretToken, payloadBody string) string {
	mac := hmac.New(sha1.New, []byte(secretToken))
	mac.Write([]byte(payloadBody))
	expectedMAC := mac.Sum(nil)
	return "sha1=" + hex.EncodeToString(expectedMAC)
}

func verifySignature(secretToken, payloadBody string, signatureToCompareWith string) bool {
	signature := generateSignature(secretToken, payloadBody)
	return subtle.ConstantTimeCompare([]byte(signature), []byte(signatureToCompareWith)) == 1
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

  fmt.Println("Getting card from Trello:", cardId)
  card, err := client.GetCard(cardId, trello.Defaults())
  if err != nil {
    log.Fatal(err)
  }

  attachment := trello.Attachment {
    Name: "PR",
    URL: url,
  }

  fmt.Println("Attaching URL:", url)
  cardErr := card.AddURLAttachment(&attachment)
  if cardErr != nil {
    fmt.Println(cardErr)
  }
}
