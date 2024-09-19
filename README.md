# Ivan's Internal Platform API

A Go-based cli and API for managing Ivan's internal platform and Ivan's applications... created using a variety of AI tools.

## Features

- CRUD operations for DNS records
- Domain management
- IP usage checking
- Providers:
  - Google Cloud
  - Cloudflare
  - AWS
- Docker container management
- Virtual Machine (VM) operations
- Proxmox cluster management

## CLI

There are the main commands:

- `i2 vms`: Manage virtual machines
- `i2 dns`: Manage DNS records
- `i2 apps`: Manage applications
- `i2 containers`: Manage containers
- `i2 cp`: Copy files to and from containers and VMs
- `i2 config`: config i2

These are the commands in the backlog:

- `i2 logs`: Manage logs
- `i2 backups`: Manage backups
- `i2 certs`: Manage certificates
- `i2 ssh`: Manage SSH keys and connections


## API Endpoints

- `GET /dns/:zone/records/:id`: Read a DNS record
- `PUT /dns/:zone/records/:id`: Update a DNS record
- `DELETE /dns/:zone/records/:id`: Delete a DNS record
- `GET /dns/ip/:ip`: Returns the domains using an IP
- // `POST /auth/login`: User login
- // `POST /auth/logout`: User logout
- `GET /apps`: List all applications
- `POST /apps`: Create a new application
- `GET /apps/:id`: Get application details
- `PUT /apps/:id`: Update application settings
- `DELETE /apps/:id`: Delete an application
- `POST /apps/:id/deploy`: Deploy an application
- `GET /monitoring`: Get system monitoring data
- `GET /logs`: Retrieve application logs
- `POST /backups`: Create a new backup
- `GET /backups`: List all backups
- `POST /certificates`: Request a new SSL/TLS certificate
- `GET /certificates`: List all SSL/TLS certificates

## Contributing

[Instructions for contributors](CONTRIBUTING.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
