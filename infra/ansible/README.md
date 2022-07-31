# Ansible playbooks


## Installation

In order to run ansible, this guide should be followed:
- https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#installing-and-upgrading-ansible

Copy versioned vars file. Then modify to match your infrastructure
```bash
cp vars.local.yml.dist vars.local.yml
```

Machine provisioning:
```bash
ansible-playbook floridaman.yml
```

After provisioning the machine, it is required to fill the `{{ app_dir }}/environment.conf` file with the needed secrets of the application.
It should contain the syntax for declaring environment variables, same as `.env.dist` in the root of the repository.
