# GolangApp API with SQLite

This API, built with Go using the Echo framework, GORM for database management, and SQLite, provides user and group management features, including user authentication. This guide will walk you through setting up and running the API using Docker.

## Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your machine
- [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine

## Getting Started

Follow these steps to set up and run the GolangApp API in a Docker container.

### 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/osgioia/challenge-Iskaypet.git
cd challenge-Iskaypet
```

### 2. Build and Run the Containers
This command will build the Go application into a Docker image (which may take 1–2 minutes, depending on your network speed) and start the application container.

```bash
docker-compose up --build
```

### 3. Verify the Setup
Once the containers are up and running, you can verify the API by accessing the following endpoints in Postman. Use Basic Auth for authorization with Username: `spadmin` and Password: `admin`.

1. Get all users: `http://localhost:8080/api/v1/users`
2. Get all groups: `http://localhost:8080/api/v1/groups`
3. Get all clients: `http://localhost:8080/api/v1/clients`

### 4. API Endpoints

Here are some of the available API endpoints:

| Method | Endpoint                                  | Description                                                                           |
|--------|-------------------------------------------|---------------------------------------------------------------------------------------|
| GET    | /api/v1/users/:id                         | Fetch a single user by ID                                                             |
| GET    | /api/v1/users                             | Fetch all users                                                                       |
| POST   | /api/v1/users                             | Create a new user                                                                     |
| PUT    | /api/v1/users/:id                         | Update a user by ID                                                                   |
| DELETE | /api/v1/users/:id                         | Delete a user by ID                                                                   |
| PUT    | /api/v1/users/:id/enable                  | Enable a user                                                                         |
| PUT    | /api/v1/users/:id/disable                 | Disable a user                                                                        |
| PUT    | /api/v1/users/:id/reset_password          | Reset a user's password                                                               |
| GET    | /api/v1/groups/:id                        | Fetch a single group by ID                                                            |
| GET    | /api/v1/groups                            | Fetch all groups                                                                      |
| POST   | /api/v1/groups                            | Create a new group                                                                    |
| POST   | /api/v1/users/:id/groups/:group_id        | Assign a group to a user                                                              |
| DELETE | /api/v1/users/:id/groups/:group_id        | Remove a group from a user                                                            |
| DELETE | /api/v1/groups/:group_id                  | Delete a group                                                                        |
| GET    | /api/v1/clients/:id                       | Fetch a single client by ID                                                           |
| GET    | /api/v1/clients/kpi                       | Fetch client KPIs                                                                     |
| GET    | /api/v1/clients                           | Fetch all clients                                                                     |
| POST   | /api/v1/clients                           | Create a new client                                                                   |
| PUT    | /api/v1/clients/:id                       | Update a client by ID                                                                 |
| DELETE | /api/v1/clients/:id                       | Delete a client by ID                                                                 |

### 5. Stopping the Containers
To stop the running containers, press `Ctrl+C` in the terminal where Docker Compose is running. You can also use the following command to stop and remove the containers:

```bash
docker-compose down
```

### 6. Additional Commands
To rebuild the images without using the cache:
```bash
docker-compose build --no-cache
```

To view the logs:
```bash
docker-compose logs
```

To start the containers in the background (detached mode):
```bash
docker-compose up -d
```

### 7. Deploying to AWS Lambda
If you want to deploy the app as an AWS Lambda function, copy the `lambda/main.go` file to the root directory and run the following script:

```bash
sh lambda/script.sh
```

The script will guide you through the deployment process.

### 8. Customizing Database Initialization
If you want to customize the initial database schema or seed data, you can modify the code in `config.go`, where the SQLite database is initialized and migrated automatically.

### 9. Troubleshooting
- If you encounter issues with the container not starting, ensure Docker and Docker Compose are installed correctly, and check for error messages in the terminal.
- Ensure port 8080 (for the API) is not in use by other applications.

## Author

This project was created by [Oz](https://www.linkedin.com/in/osvaldogioia/).