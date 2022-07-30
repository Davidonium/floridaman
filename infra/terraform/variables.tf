variable "domain" {
  description = "base domain host the instance will be assigned, must not include subdomains. e.g. davidonium.com"
}

variable "pgp_key" {
  description = "pgp key used to generate the public and private rsa keys to be able to ssh into the instance"
}