# Ivan's Internal Platform API

A Go-based API for managing Ivan's internal platform and Ivan's applications... created using a variety of AI tools.

## Features


- CRUD operations for DNS records
- Domain management
- IP usage checking



### Running the API

To start the server:

```
go run main.go
```

## API Endpoints

- `GET /domains/:domain/records/:id`: Read a DNS record
- `PUT /domains/:domain/records/:id`: Update a DNS record
- `DELETE /domains/:domain/records/:id`: Delete a DNS record
- `GET /ip-usage/:ip`: Check IP usage

## Contributing

[Instructions for contributors](CONTRIBUTING.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
