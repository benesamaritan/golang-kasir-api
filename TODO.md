# TODO List for Kasir API Improvement

This document outlines current issues, security considerations, and potential enhancements for the Kasir API.

## Current Issues & Bugs (To be Fixed)
- [ ] **SQL Syntax Errors in Update Queries:** (FIXED in previous commits) Extra commas in `UPDATE` queries for products and categories.
- [ ] **Empty Best Seller in Reports:** `GetReport` currently returns an empty `BestSeller` object. Need to implement SQL logic to find the most sold product. (FIXED)
- [ ] **Inconsistent Error Response Format:** Some handlers use `http.Error` (plain text), while others use JSON. Should be standardized to JSON for frontend compatibility.
- [ ] **Product CategoryID Bug:** (FIXED in previous commits) `GetByID` was resetting `CategoryID` to `0`.
- [ ] **Standard Library Routing Limitation:** Using `http.HandleFunc("/categories/{id}", ...)` works in Go 1.22+, but `HandleCategoryByID` still manually trims the prefix. It should be updated to use `r.PathValue("id")`.

## Security & Reliability
- [ ] **Input Validation:**
    - Products: Ensure `price` and `stock` are not negative. Ensure `name` is not empty.
    - Categories: Ensure `name` is not empty.
    - Transactions: Validate `quantity` against available stock (partially implemented in repository).
- [ ] **Transaction Management:** Improve `CreateTransaction` to ensure `tx.Rollback()` is always called on error before the function returns (currently uses a `defer` that might be complex to follow).
- [ ] **Hardcoded Database String:** Ensure all database configurations are strictly from `.env`.
- [ ] **CORS Policy:** Current `Access-Control-Allow-Origin: "*"` is fine for development but should be restricted in production.
- [ ] **Authentication/Authorization for all API Endpoints:** Implement robust authentication and authorization mechanisms for all API endpoints beyond just API key for specific handlers.

## Simple Improvements (Bootcamp Level)
- [ ] **Implement Best Seller Logic:** (FIXED)
    ```sql
    SELECT p.name, SUM(td.quantity) as qty
    FROM transaction_details td
    JOIN products p ON td.product_id = p.id
    GROUP BY p.name
    ORDER BY qty DESC LIMIT 1
    ```
- [ ] **Standardize JSON Responses:** Create a helper function in `utils/response.go` to send consistent JSON responses (e.g., `{"status": "success", "data": ...}`).
- [ ] **Add "Active" status to Categories:** Similar to products, categories could have an `active` status to "soft-delete" them.
- [ ] **Refactor Routing:** Use `r.PathValue("id")` for all ID-based endpoints to simplify code.
- [ ] **Add Recovery Middleware:** Add a middleware to catch panics and return a 500 error instead of crashing the server.
- [ ] **Error Handling Improvement:** Implement more granular and user-friendly error messages across the application.

## Future Enhancements
- [ ] **Add Transaction History Handler (GET /api/transaction):** Implement an endpoint to retrieve a list of past transactions.
- [ ] **Pagination:** Add `limit` and `offset` to `GetAll` endpoints and for the new transaction history handler.
- [ ] **Add Customer Name and Contact to Checkout:** Extend the database schema and logic flow to include customer name and contact information during the checkout process.
- [ ] **API Documentation:** Create a simple Swagger/OpenAPI spec.
- [ ] **Unit Tests & Integration Tests:** Develop comprehensive unit tests for business logic (Service layer) and integration tests for API endpoints and database interactions.
- [ ] **Logging Enhancement:** Implement structured logging with different levels (e.g., INFO, DEBUG, ERROR) to aid in debugging and monitoring.
- [ ] **Product Image Support:** Add functionality to store and retrieve product images.
- [ ] **User Management:** Implement a basic user management system to handle different roles (e.g., cashier, manager) and access levels.
