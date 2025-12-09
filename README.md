# abc-user-service
Demo User Service for ABC

## Code Structure

The project structure is based on the [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

The project follows a clean architecture approach, separating concerns into distinct layers. Here's a breakdown of the main directories:

-   `cmd/server`: Main application entry point.
-   `internal`: Contains all the core application and business logic. The structure of this directory is inspired by Uncle Bob's [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and Domain-Driven Design principles.
    -   `application/service`: Implements the use cases of the application.
    -   `domain`: Core domain models, repository interfaces, and events. This layer is independent of any other layer.
    -   `infrastructure`: Implements the interfaces defined in the domain layer (e.g., database repositories, message brokers, HTTP handlers).
-   `pkg`: Shared libraries and utilities.
-   `api`: API contracts and specifications (e.g., OpenAPI, gRPC).
-   `build`: Docker configurations for different environments.
-   `scripts`: Helper scripts for development and operational tasks.
-   `mocks`: Generated mocks for testing purposes.

## Development

### Prerequisites

-   Go (version 1.21 or higher)
-   Docker and Docker Compose
-   `make`

### Running Locally

The easiest way to get the development environment up and running is by using the provided `docker-compose` setup.

1.  **Start the environment:**
    This command will build the Go application, create the Docker containers for the service and its dependencies (like the database), and run them.

    ```bash
    make docker-compose-up
    ```

2.  **Stopping the environment:**
    To stop and remove the containers, run:

    ```bash
    make docker-compose-down
    ```

Alternatively, you can run the service directly on your host machine using `make run-local`. This requires you to manually set up the database and other dependencies.

## API Documentation

The service provides OpenAPI documentation via Swagger UI. Once the application is running, you can access the documentation by navigating to:

[/swagger/index.html](/swagger/index.html)

This will display the interactive API documentation, allowing you to explore the available endpoints and their details.

## Makefile Commands

The `Makefile` provides several commands to streamline development:

-   `make binary`: Compiles the Go application into a binary located in the `bin/` directory.
-   `make docker-compose-up`: Starts the local development environment using Docker Compose. It also builds the application.
-   `make docker-compose-down`: Stops the local development environment.
-   `make run-local`: Runs the application directly on the host machine (requires dependencies to be running).
-   `make test`: Runs the unit and integration tests and generates a coverage report.
-   `make lint`: Lints the codebase using `golangci-lint`.
-   `make lint-install`: Installs the `golangci-lint` tool.
-   `make gomod`: Tidies the `go.mod` file, adding missing and removing unused modules.
-   `make mocks`: Generates mocks for the repository and other interfaces.
-   `make update-mocks`: Updates mockery and then generates mocks.
-   `make test-mocks`: A convenient command that updates mocks and then runs tests.
-   `make clean`: Removes build artifacts from the `bin/` directory.
-   `make install-go-test-coverage`: Installs the `go-test-coverage` tool.
-   `make install-swag`: Installs the `swag` tool for Swagger documentation.
-   `make swagger`: Generates Swagger documentation.
