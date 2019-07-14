package main

import (
  "log"
  "net/http"
  "os"

  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
)

type PullRequest struct {
  Action string `json:"action"`
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
    var pr PullRequest
    c.BindJSON(&pr)
    // buf, e := ioutil.ReadAll(c.Request.Body)
    //
    // if e != nil {
    //   return e
    // }

    // ioutil.NopCloser(bytes.NewReader(buf))
    // return json.Unmarshal(buf, dest)

    c.JSON(http.StatusOK, gin.H{"message": pr.Action, "status": http.StatusOK})
  })

  router.Run(":" + port)
}
