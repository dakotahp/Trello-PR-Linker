package main

import (
  "crypto/hmac"
  "crypto/sha1"
  "crypto/subtle"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "regexp"

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

    if pr.Action == "opened" || pr.Action == "edited" {
      fmt.Println("Operating on PR:", pr.PullRequest.HtmlUrl)

      titleId := trelloIdFromTitle(pr.PullRequest.Title)
      branchId := trelloIdFromBranch(pr.PullRequest.Head.Ref)

      if titleId != "" {
        postPrLinkToTrelloCard(titleId, pr.PullRequest.HtmlUrl)
      } else if branchId != "" {
        postPrLinkToTrelloCard(branchId, pr.PullRequest.HtmlUrl)
      } else {
        fmt.Println("Skipping, ticket ID not found in title or branch", pr.PullRequest.Title, pr.PullRequest.Head.Ref)
      }
    } else {
      fmt.Println("Skipping, action not relevant", pr.PullRequest.HtmlUrl, pr.Action)
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
  matches := re.FindStringSubmatch(title)

  if len(matches) > 0 {
    return matches[1]
  }

  return ""
}

func trelloIdFromBranch(branch string) (string) {
  re := regexp.MustCompile(`([A-Za-z0-9]{8})(--|\/)`)
  matches := re.FindStringSubmatch(branch)

  if len(matches) > 0 {
    return matches[1]
  }

  return ""
}

func postPrLinkToTrelloCard(cardId string, url string) {
  key := os.Getenv("TRELLO_KEY")
  token := os.Getenv("TRELLO_TOKEN")

  client := trello.NewClient(key, token)

  fmt.Println("Getting card from Trello:", cardId)
  card, err := client.GetCard(cardId, trello.Defaults())
  if err != nil {
    log.Fatal(err)
  }

  if !prAlreadyAttached(card, url) {
    attachment := trello.Attachment {
      Name: "PR",
      URL: url,
    }
    fmt.Println(prAlreadyAttached(card, url))
    fmt.Println("Attaching URL:", url)
    cardErr := card.AddURLAttachment(&attachment)
    if cardErr != nil {
      fmt.Println(cardErr)
    }
  }
}

func prAlreadyAttached(card *trello.Card, url string) (bool) {
  fmt.Println("Attachments", card.Attachments)
  fmt.Println("Card", card)
  for i := 0; i < len(card.Attachments); i++ {
    fmt.Println(card.Attachments[i])
    fmt.Println(card.Attachments[i].URL)
    fmt.Println(url)
    if card.Attachments[i].URL == url {
      return true
    }
  }

  return false
}
