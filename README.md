# Trello-PR-Linker

An app that listens to GitHub PR webhooks and parses out Trello ticket IDs so it can attach the PR to them. It was written in Go and uses the gin-tonic framework for routing and exception handling. Hosting should be on Heroku due to it being created from a starter boilerplate repo for it.

## How to use

### Create a Webhooks on GitHub
1. Go to a repository settings (`https://github.com/username/repo-name/settings`).
2. Go to webhook settings (`https://github.com/username/repo-name/settings/hooks`).
3. Click "Add webhook".
4. Enter Payload URL of `https://your-host-name.herokuapp.com/webhook`.
5. Choose content type of `application/json`.
6. [Generate a secret token](https://developer.github.com/webhooks/securing/) in your terminal with `ruby -rsecurerandom -e 'puts SecureRandom.hex(20)'`. (*Save this token for later*)
7. Enter the output from the terminal into the Secret field.
8. Under "Which events would you like to trigger this webhook?", select "Let me select individual events."
9. Check off "Pull requests" (and probably _uncheck_ "Pushes")
10. Submit the form and you're done with this step.

### Save Webhook Secret as ENV VAR
Save the secret you created for the webhook and set it as an environment variable, either in `.env` locally or on the server as `heroku config:set SECRET=123`.

### Get Trello API Key
1. [Generate a set of credentials](https://trello.com/app-key)
2. Set these as `TRELLO_KEY` and `TRELLO_TOKEN` as environment variables.

You should be all setup and the app should run.

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

[x] Secure request with hash
[ ] Parse out ticket ID from branch

## Attribution

This application started from the heroku Go boilerplate article [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go).
