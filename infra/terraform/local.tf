terraform {
  backend "local" {}
}

locals {
  appname  = "floridaman"
  region   = "eu-west-3"
}