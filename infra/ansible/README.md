# Ansible playbooks

Copy versioned vars file. Then modify to match your infrastructure
```bash
cp vars.local.yml.dist vars.local.yml
```

Machine provisioning:
```bash
ansible-playbook floridaman.yml
```

Application deploy:
```bash
ansible-playbook floridaman-deploy.yml
```