# Gator CLI

Welcome to the Gator CLI! This tool helps you aggregate and manage RSS feeds directly from your terminal using a PostgreSQL database.

## Prerequisites

Before running this program, you must have the following installed on your machine:

1.  **Go (Golang):**
    You need Go installed to compile and install the CLI tool.
    *   [Download and Install Go](https://go.dev/doc/install)
    *   Verify installation by running: `go version`

2.  **PostgreSQL:**
    You need a Postgres database running to store users and feed data.
    *   [Download and Install PostgreSQL](https://www.postgresql.org/download/)
    *   Make sure your Postgres service is running in the background.
    *   You will need a connection string (e.g., `postgres://username:password@localhost:5432/gator?sslmode=disable`).

## Installation

To install the `gator` CLI tool globally on your system, use the `go install` command from the root of the project directory (or using the remote repository path):

```bash
# Run this from inside the project directory
go install .

# OR run this using the remote path (replace with your actual repo path)
go install github.com/yourusername/gator@latest