# Ansible playbooks

Copy versioned vars file. Then modify to match your infrastructure
```bash
cp floridaman_vars.yml.dist floridaman_vars.yml
```

Machine provisioning:
```bash
ansible-playbook -i hosts floridaman.yml
```

Application deploy:
```bash
ansible-playbook -i hosts floridaman-deploy.yml
```