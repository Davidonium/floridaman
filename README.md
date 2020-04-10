# Floridaman

This project is an HTTP API with a single endpoint that provides a random floridaman heading each time and an integration to add a slash command to slack to provide the same functionality.

Also includes a script to retrieve all the floridaman headings from /r/floridaman in Reddit.

The project uses Redis for storage and [graw](https://github.com/turnage/graw) for querying reddit.

Terraform and Ansible are used to provide a very simple infrastructure of a single machine in amazon lightsail with redis installed and an nginx for proxying requests to the go app. The code is under `infra`.

* Terraform version used: `v0.12.18`
* Ansible version used: `2.9.2`

## Compile

```bash
make build
```

## Run

Api:
```bash
make api
```

Read reddit script:

```bash
make readreddit
```

### Disclamer
I did this project for fun as a result of a joke with some friends at work and me wanting to learn terraform and ansible. It features absolutely 0 tests and bad practices, I don't work like this though :D.

I documented some things for the sole purpose that if I want to rekindle this project, I can remember what I did.

Hopefully I haven't leaked any aws or reddit credentials, if I have please notify me through an issue. 
 