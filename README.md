# Ivan's Internal Platform API

A Go-based API for managing Ivan's internal platform and Ivan's applications... created using a variety of AI tools.

## Features


- CRUD operations for DNS records
- Domain management
- IP usage checking
- Providers:
  - Google Cloud
  - Cloudflare
  - AWS




### Running the API

To start the server:

```
go run main.go
```

## API Endpoints

- `GET /dns/:zone/records/:id`: Read a DNS record
- `PUT /dns/:zone/records/:id`: Update a DNS record
- `DELETE /dns/:zone/records/:id`: Delete a DNS record
- `GET /dns/ip/:ip`: Returns the domains using an IP

## Contributing

[Instructions for contributors](CONTRIBUTING.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
