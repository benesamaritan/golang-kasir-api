# Project Overview

This is a Go-based web API designed for a kasir (cashier) system. It implements a clean architecture pattern, organizing code into distinct layers: handlers for API interaction, services for business logic, repositories for data access, and models for data structures. The API provides comprehensive endpoints for managing products, categories, customer transactions (checkout), and generating reports.

For data persistence, the application connects to a PostgreSQL database. Configuration management is handled efficiently using the `viper` library, allowing settings to be sourced from environment variables or a `.env` file.

The project leverages `devbox` for environment management and containerization, ensuring a consistent development and deployment experience. It also incorporates basic load testing using `k6` to verify the application's health and responsiveness.

# Building and Running

This project uses `devbox` to manage its development environment and execution.

## Prerequisites

*   **Install Devbox:** If you don't have `devbox` installed, run:
    ```sh
    curl -fsSL https://get.jetpack.io/devbox | bash
    ```

## Commands

1.  **Enter Devbox Shell:**
    To enter the project's isolated development environment, use:
    ```sh
    devbox shell
    ```

2.  **Run the Application:**
    This command builds the Go application and starts it in the background on port `8000`.
    ```sh
    devbox run kasir
    ```
    You should see output similar to: `Kasir application started in background on port 8000.`

3.  **Run Load Tests:**
    This executes a basic load test using `k6` against the `/health` endpoint of the running API.
    ```sh
    devbox run test-load
    ```

# Development Conventions

*   **Language:** Go
*   **Dependencies:** Managed using `go.mod`. Key dependencies include `github.com/lib/pq` for PostgreSQL database interaction and `github.com/spf13/viper` for robust configuration management.
*   **Configuration:** The application loads its configuration (e.g., `PORT`, `DB_CONN`, `API_KEY`) from environment variables or a `.env` file using `viper`.
*   **Architecture:** The codebase adheres to a layered architecture, separating concerns into handlers (API layer), services (business logic), repositories (data access layer), and models (data structures).
*   **API Design:** The API follows RESTful principles, utilizing standard HTTP methods (GET, POST, PUT, DELETE) and communicating via JSON payloads for requests and responses.
*   **Middleware:** Custom middleware is implemented for essential functionalities such as API key authentication, Cross-Origin Resource Sharing (CORS), and request logging, enhancing security and observability.
*   **Containerization:** The project is set up for containerization using `devbox` and a `Dockerfile`, facilitating consistent environments across development and potential deployment stages.
*   **Testing:** A basic load testing setup is included using `k6` to perform health checks on the API. Comprehensive unit or integration tests for business logic are not explicitly defined in the provided `README.md` or `devbox.json` (the `test` script is currently a placeholder).
