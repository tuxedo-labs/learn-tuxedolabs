# API Server

## Overview
This directory contains the source code for the API server of the Tuxedo Labs project. The server is responsible for handling client requests, processing data, and communicating with the database.

## Features
- **RESTful API**: Provides a set of RESTful endpoints for various operations.
- **Authentication**: Implements JWT-based authentication.
- **Database Integration**: Connects to a PostgreSQL database to store and retrieve data.
- **Logging**: Uses structured logging for better traceability and debugging.
- **Error Handling**: Comprehensive error handling mechanisms.

## Project Structure
```
server/
├── cmd/               # Entry points to start the server
├── config/            # Configuration files
├── controllers/       # API endpoint controllers
├── middleware/        # Middleware functions
├── models/            # Database models
├── routes/            # API route definitions
├── services/          # Business logic and services
├── utils/             # Utility functions
├── main.go            # Main entry point of the server
└── README.md          # This file
```

## Getting Started

### Prerequisites
- Go 1.16 or later
- PostgreSQL database

### Installation
1. **Clone the repository**
    ```bash
    git clone https://github.com/tuxedo-labs/learn-tuxedolabs.git
    cd learn-tuxedolabs/server
    ```

2. **Install dependencies**
    ```bash
    go mod download
    ```

### Configuration
Create a `.env` file in the `server` directory and add the following environment variables:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret
```

### Running the Server
Start the server using the following command:
```bash
go run main.go
```

## Usage
The API server provides various endpoints that can be accessed via HTTP requests. Below are some example endpoints:

- `GET /api/v1/users`: Retrieve a list of users.
- `POST /api/v1/auth/login`: Authenticate a user and receive a JWT.
- `GET /api/v1/profile`: Get the profile of the authenticated user.
