
# go-getting-started

A barebones Go app, which can easily be deployed to Heroku.

This application supports the [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go) article - check it out.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ git clone https://github.com/heroku/go-getting-started.git
$ cd go-getting-started
$ go build -o bin/card-linker -v .
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/go-getting-started
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```


## To Do

* Secure request with hash
