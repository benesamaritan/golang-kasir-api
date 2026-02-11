# TODO List for Kasir API Improvement

This document outlines current issues, security considerations, and potential enhancements for the Kasir API.

## 🛠️ Missing API Route Handlers (Recommended)
- [ ] **Transaction Management:**
    - [ ] `GET /api/transactions`: List all transaction history (currently only `POST /api/checkout` exists).
    - [ ] `GET /api/transactions/{id}`: Get detailed information for a specific transaction including its items.
- [ ] **Advanced Product & Category Features:**
    - [ ] `PATCH /api/product/{id}/active`: Dedicated endpoint to toggle product active status without sending full body.
    - [ ] `GET /api/categories/{id}/products`: List all products belonging to a specific category.
- [ ] **Enhanced Reporting:**
    - [ ] `GET /api/report/range?from=YYYY-MM-DD&to=YYYY-MM-DD`: Get sales report for a specific date range.
    - [ ] `GET /api/report/top-categories`: Identify which categories are generating the most revenue.

## 🚀 Frontend Implementation (HTMX + Go Templates)
- [ ] **Infrastructure Setup:**
    - [ ] Create `templates/` directory for HTML files.
    - [ ] Setup `static/` directory for CSS/JS assets (optional, can use CDNs).
    - [ ] Configure `main.go` to serve static files and parse templates.
- [ ] **Core Layout:**
    - [ ] Create `layout.html` with Bootstrap 5 and HTMX script.
    - [ ] Implement a responsive Navbar.
- [ ] **Product Management Dashboard:**
    - [ ] Create `/dashboard/produk` route.
    - [ ] Implement dynamic product list with "Delete" and "Toggle Active" buttons using HTMX (`hx-delete`, `hx-patch`).
    - [ ] Add "Search" functionality with live filtering (`hx-trigger="keyup delay:500ms"`).
    - [ ] Create a Modal or side-panel for "Add Product" using HTMX.
- [ ] **Category Management:**
    - [ ] Create `/dashboard/categories` route.
    - [ ] Implement category list and creation form.
- [ ] **Cashier/Checkout Interface:**
    - [ ] Build a simple "Point of Sale" UI where products can be added to a virtual cart.
    - [ ] Implement "Process Checkout" button that calls `/api/checkout` via HTMX and displays a success/receipt message.
- [ ] **Reporting Dashboard:**
    - [ ] Visual dashboard for "Today's Report" using simple HTML tables or charts.

## Current Issues & Bugs (To be Fixed)
- [ ] **Inconsistent Error Response Format:** Some handlers use `http.Error` (plain text), while others use JSON. Should be standardized to JSON for frontend compatibility.
- [ ] **Standard Library Routing Limitation:** Using `http.HandleFunc("/categories/{id}", ...)` works in Go 1.22+, but `HandleCategoryByID` still manually trims the prefix. It should be updated to use `r.PathValue("id")`.
- [ ] **Dynamic API Slash Redirect:** Ensure the fix for `/api/` vs `/api` is robust and handled at the middleware or router level.

## Security & Reliability
- [ ] **Input Validation:**
    - Products: Ensure `price` and `stock` are not negative. Ensure `name` is not empty.
    - Categories: Ensure `name` is not empty.
    - Transactions: Validate `quantity` against available stock (partially implemented in repository).
- [ ] **Transaction Management:** Improve `CreateTransaction` to ensure `tx.Rollback()` is always called on error before the function returns (currently uses a `defer` that might be complex to follow).
- [ ] **CORS Policy:** Current `Access-Control-Allow-Origin: "*"` is fine for development but should be restricted in production.
- [ ] **Authentication/Authorization:** Implement session-based or JWT-based auth for the Dashboard routes.

## Simple Improvements (Bootcamp Level)
- [ ] **Standardize JSON Responses:** Create a helper function in `utils/response.go` to send consistent JSON responses.
- [ ] **Add "Active" status to Categories:** Similar to products, categories could have an `active` status.
- [ ] **Refactor Routing:** Use `r.PathValue("id")` for all ID-based endpoints to simplify code.
- [ ] **Add Recovery Middleware:** Add a middleware to catch panics and return a 500 error.

## Future Enhancements
- [ ] **Add Transaction History Handler (GET /api/transaction).**
- [ ] **Pagination:** Add `limit` and `offset` to `GetAll` endpoints.
- [ ] **Customer Info:** Extend checkout to include customer name and contact.
- [ ] **API Documentation:** Create a Swagger/OpenAPI spec.
- [ ] **Unit Tests & Integration Tests.**
- [ ] **Product Image Support.**
- [ ] **User Management / Roles.**