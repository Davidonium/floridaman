terraform {
  backend "local" {}
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

locals {
  app_name = "floridaman"
  region   = "eu-west-3"
}

provider "aws" {
  region = local.region
  profile = local.app_name
}

resource "aws_lightsail_static_ip_attachment" "app" {
  static_ip_name = aws_lightsail_static_ip.app.name
  instance_name  = aws_lightsail_instance.app.name
}

resource "aws_lightsail_static_ip" "app" {
  name = local.app_name
}

resource "aws_lightsail_instance" "app" {
  name              = "${local.app_name}.${var.domain}"
  availability_zone = "${local.region}b"
  blueprint_id      = "ubuntu_22_04"
  bundle_id         = "micro_2_0"
  key_pair_name     = aws_lightsail_key_pair.app.name
}

resource "aws_lightsail_instance_public_ports" "app" {
  instance_name = aws_lightsail_instance.app.name

  port_info {
    cidrs = [
      "0.0.0.0/0"
    ]
    protocol  = "tcp"
    from_port = 443
    to_port   = 443
  }

  port_info {
    cidrs = [
      "0.0.0.0/0"
    ]
    protocol  = "tcp"
    from_port = 80
    to_port   = 80
  }

  port_info {
    cidrs = [
      "0.0.0.0/0"
    ]
    protocol  = "tcp"
    from_port = 22
    to_port   = 22
  }
}

data "aws_route53_zone" "app" {
  name = "${var.domain}."
}

resource "aws_route53_record" "app" {
  zone_id = data.aws_route53_zone.app.zone_id
  name    = "${local.app_name}.${data.aws_route53_zone.app.name}"
  type    = "A"
  ttl     = "86400"
  records = [aws_lightsail_static_ip.app.ip_address]
}

resource "aws_lightsail_key_pair" "app" {
  name    = var.domain
  pgp_key = var.pgp_key
}