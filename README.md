# Floridaman

This project is an HTTP API with an endpoint that provides a random floridaman heading each time and an integration with slack to add a slash command to slack to provide the same functionality.

Also includes a script that retrieves all the floridaman headings from /r/floridaman in Reddit and dumps them in redis.

The project uses Redis for storage and [go-reddit](https://github.com/vartanbeno/go-reddit) for querying the reddit API.

## Compile

```bash
make build
```

## Run

Before running the application on your machine, make sure the `.env` file exists. If not create it and fill it.
```bash
cp .env.dist .env
```

API server:
```bash
go run ./cmd/floridaman serve
```

Read reddit script:

```bash
go run ./cmd/floridaman readreddit
```

## API Spec

`GET /health`

Application Monitoring

`200 OK`
```json
{
  "status": "UP",
  "redis": "UP"
}
```

`GET /random`

Returns a random entry in each request

`200 OK`
```json
{
  "title": "Florida Man stops a robbery",
  "link": "https://i.imgur.com/2fwv8Iz.jpg",
  "source": "reddit"
}
```

`POST /slack/random`

This endpoint handles slack slash commands that one must setup following https://api.slack.com/interactivity/slash-commands#app_command_handling

`SLACK_SIGNING_SECRET` environment variable must be present with a valid slack signing secret provided in their UI.


It validates the slack request as documented in https://api.slack.com/docs/verifying-requests-from-slack

It responds with the most basic response.

`200 OK`
```json
{
  "response_type": "in_channel",
  "text": "Florida Man stops a robbery (https://i.imgur.com/2fwv8Iz.jpg)"
}
```

`400 Bad Request`

when the slack request could not be validated
```json
{
  "message": "Invalid slack request"
}
```

---

Generic Errors:

`500 Internal Server Error`
```json
{
  "message": "Internal server error"
}
```

### Infrastructure

Terraform and Ansible are used to provide a very simple infrastructure of a single machine in amazon lightsail with redis installed and an nginx for proxying requests to the go app. The code is under `infra`.

See the [ansible README](./infra/ansible/README.md) and the [terraform README](./infra/terraform/README.md) for more details on how to use them.


### Disclamer
I did this project for fun as a result of a joke with some friends at work and me wanting to learn terraform and ansible. It features 0 tests and bad practices, I don't work like this though :D.

I documented some things for the sole purpose that if I want to rekindle this project, I can remember what I did.

Hopefully I haven't leaked any aws or reddit credentials, if I have please notify me through an issue. 
 
### Help
Random commands that I usually forget.
```bash
# get logs
ssh floridaman "sudo journalctl -u floridaman.service"

# check for enabled timers and last execution
ssh floridaman "sudo systemctl list-timers"
```