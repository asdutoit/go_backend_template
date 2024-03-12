# Go Gin Backend Application

This is a backend application built with Go and the Gin framework.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.16 or later)
- Gin framework

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/your-repo-name.git
```

2. Change into the project directory

```bash
cd your-repo-name
```

3. Install the dependencies

```bash
go get
```

4. Configure environment variables:

```bash
   export AWS_REGION=<your-aws-region>
   export AWS_S3_BUCKET=<your-aws-s3-bucket>
   export AWS_ACCESS_KEY_ID=<your-aws-access-key-id>
   export AWS_SECRET_ACCESS_KEY=<your-aws-secret-access-key>
   export DB_HOST=localhost
   export DB_PORT=54320
   export DB_USER=postgres
   export DB_PASSWORD=my_password
   export DB_NAME=postgres
```

5. Start a Postgres database using Docker

```bash
docker run --name local-psql -v local_psql_data_new:/var/lib/postgresql/data -p 54320:5432 -e POSTGRES_PASSWORD=my_password -d postgres:15.3
```

6. Run the application

```bash
go run main.go
```

7. The application should now be running on `http://localhost:8080`

## Running the tests

To run the tests, run the following command:

```bash
cd routes
go test
```
