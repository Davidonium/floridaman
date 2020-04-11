provider "aws" {
  region = local.region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

resource "aws_lightsail_static_ip_attachment" "app" {
  static_ip_name = aws_lightsail_static_ip.app.name
  instance_name  = aws_lightsail_instance.app.name
}

resource "aws_lightsail_static_ip" "app" {
  name = local.appname
}

resource "aws_lightsail_instance" "app" {
  name              = "${local.appname}.${var.domain}"
  availability_zone = "${local.region}b"
  blueprint_id      = "ubuntu_18_04"
  bundle_id         = "micro_2_0"
  key_pair_name     = aws_lightsail_key_pair.app.name
    provisioner "local-exec" {
    command = "aws lightsail put-instance-public-ports --instance-name=${local.appname}.${var.domain} --port-infos fromPort=443,toPort=443,protocol=https"
  }
}

data "aws_route53_zone" "zone" {
  name = "${var.domain}."
}

resource "aws_route53_record" "app" {
  zone_id = data.aws_route53_zone.zone.zone_id
  name    = "${local.appname}.${data.aws_route53_zone.zone.name}"
  type    = "A"
  ttl     = "86400"
  records = [aws_lightsail_static_ip.app.ip_address]
}

resource "aws_lightsail_key_pair" "app" {
  name    = var.domain
  pgp_key = var.pgp_key
}