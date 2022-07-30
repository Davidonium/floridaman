# Terraform

Terraform scripts create the necessary resources to have a working instance in amazon lightsail.

First a terraform.tfvars file is needed.

```bash
cp terraform.tfvars.dist terraform.tfvars
```
then fill the variables with the correct values.


To create aws lightsail instance with its own dns zone (this has consequences on monthly billing in aws)
```bash
# plan first to see the changes
terraform plan

# apply the changes into aws
terraform apply -auto-approve
```

This will make the instance available through ssh
```bash
terraform output private_key | base64 -D | gpg --decrypt --input - > ~/.ssh/id_floridaman
# make the permissions of the private key secure
chmod 600 ~/.ssh/id_floridaman 
terraform output pub_key > ~/.ssh/id_floridaman.pub
terraform output ssh_config >> ~/.ssh/config
```