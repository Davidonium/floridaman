terraform {
  backend "local" {}
}

locals {
  appname  = "floridaman"
  domain   = "davidonium.com"
  region   = "eu-west-3"
}