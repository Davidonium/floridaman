output "ssh_config" {
  value = <<EOF

Host ${local.app_name}
	HostName		${aws_route53_record.app.name}
	User			${aws_lightsail_instance.app.username}
    IdentityFile    ~/.ssh/id_${local.app_name}
EOF
}

output "pub_key" {
  value = aws_lightsail_key_pair.app.public_key
}

output "private_key" {
  value     = aws_lightsail_key_pair.app.encrypted_private_key
  sensitive = true
}