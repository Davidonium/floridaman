variable "ssh_public_key_path" {
  default = "~/.ssh/id_rsa.pub"
}

variable "domain" {
  description = "base domain host the instance will be assigned, must not include sudomains"
}

variable "aws_access_key" {}

variable "aws_secret_key" {}

variable "pgp_key" {}