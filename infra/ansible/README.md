# Ansible playbooks

first copy the variables file 
```bash
cp floridaman_vars.yml.dist floridaman_vars.yml
```

file the file with correct values.

to provision the machines
```bash
ansible-playbook -i hosts floridaman.yml
```

to deploy the go executable
```bash
ansible-playbook -i hosts floridaman-deploy.yml
```