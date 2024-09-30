# USDT Rate Service

The USDT Rate Service is a GRPC service for getting the current USDT rate.

## Configuration

The service can be configured using environment variables or command-line flags. Command-line flags take precedence over environment variables.

### Configuration Parameters

| Parameter    | Environment Variable | Command-line Flag | Default Value | Description             |
|--------------|----------------------|-------------------|---------------|-------------------------|
| DB Host      | DB_HOST              | -db-host          | localhost     | Database host           |
| DB Port      | DB_PORT              | -db-port          | 5432          | Database port           |
| DB Username  | DB_USERNAME          | -db-username      | -             | Database username       |
| DB Password  | DB_PASSWORD          | -db-password      | -             | Database password       |
| DB Name      | DB_NAME              | -db-name          | -             | Database name           |
| API Port     | API_PORT             | -                 | 8080          | Port for HTTP API       |

### Example using environment variables:

```bash
export DB_HOST=myhost
export DB_PORT=5432
export DB_USERNAME=myuser
export DB_PASSWORD=mypassword
export DB_NAME=mydatabase
export API_PORT=8080
```

### Example using command-line flags:

```bash
go run main.go -db-host=myhost -db-port=5432 -db-username=myuser -db-password=mypassword -db-name=mydatabase
```
