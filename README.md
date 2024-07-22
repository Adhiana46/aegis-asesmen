# Aegis Asesment

## Getting Started

Follow these instructions to set up and run the project on your local machine.

### Prerequisites

Make sure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Go](https://golang.org/doc/install)

### Running the Application

1. **Set up the infrastructure:**

   Start by spinning up the necessary infrastructure using Docker Compose. Open a terminal and navigate to the project directory, then run:

   ```sh
   docker compose up -d
   ```

2. **Run the application:**

   Once the infrastructure is up and running, you can start the application. In the same terminal, execute:

   ```sh
   go run ./cmd/api
   ```
