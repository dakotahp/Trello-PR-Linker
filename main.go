package main

import (
  "log"
  "net/http"
  "os"
  "fmt"

  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
)

type Payload struct {
  Action string `json:"action"`
  PullRequest
}

type PullRequest struct {
  Url string `json:"url"`
  Title string `json:"title"`
}

func main() {
  port := os.Getenv("PORT")

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
    c.BindJSON(&pr)
    // buf, e := ioutil.ReadAll(c.Request.Body)
    //
    // if e != nil {
    //   return e
    // }

    // ioutil.NopCloser(bytes.NewReader(buf))
    // return json.Unmarshal(buf, dest)
    fmt.Println("logging output:", pr.Action)
    fmt.Println(pr)
    fmt.Println(pr.PullRequest.title)
    c.JSON(http.StatusOK, gin.H{"message": pr.Action, "status": http.StatusOK})
  })

  router.Run(":" + port)
}
