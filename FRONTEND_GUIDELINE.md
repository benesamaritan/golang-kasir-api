# Comprehensive Frontend Blueprint: Go + Templ + HTMX + Tailwind

This document serves as an expert-level technical specification for an AI agent to build a polished, interactive dashboard on top of the existing Kasir API.

---

## 1. Core Architectural Strategy
The frontend will follow the **"Locality of Behavior" (LoB)** principle. We avoid heavy client-side frameworks and use Go to serve both HTML fragments and full pages.

### Security & CORS
- **Shared Middleware**: The existing `middlewares.CORS` must be applied to Web routes if the frontend is hosted on a different subdomain.
- **CSRF Protection**: For production, implement `nosurf` or similar middleware.
- **API Key Proxy**: The web server should securely inject the `API_KEY` from environment variables when making internal service calls, so the key is never exposed to the client browser.

---

## 2. Environment Configuration
The web dashboard must respect the following `.env` variables (managed via `viper`):
- `WEB_PORT`: Port for the dashboard (default: `8081` if separate, or same as `PORT`).
- `API_BASE_URL`: The URL where the backend API is reachable.
- `ENVIRONMENT`: `development` vs `production` (to toggle live-reload/minification).

---

## 3. Dynamic Routing Map (Web Handlers)

The following routes must be implemented in `main.go` or a dedicated `handlers/web_handler.go`:

| Path | Method | Purpose | HTMX Interaction |
| :--- | :--- | :--- | :--- |
| `/` | GET | Main Dashboard Shell | Full page load |
| `/dashboard/stats` | GET | Summary Cards (Revenue, Trx) | `hx-get` on timer (30s) |
| `/produk` | GET | Product Management Page | Full page load |
| `/produk/table` | GET | Product Table Rows | `hx-get` with search params |
| `/produk/add` | GET/POST | Add Product Form/Modal | `hx-swap="beforeend"` to table |
| `/checkout` | GET | POS / Cashier Interface | Full page load |
| `/checkout/cart` | POST | Add item to virtual cart | `hx-post` partial update |

---

## 4. Technical Implementation Detail

### Templ Component Structure (`internal/ui/`)
1. **`layout.templ`**:
    - Includes Tailwind CDN or link to `static/css/main.css`.
    - Includes HTMX: `<script src="https://unpkg.com/htmx.org@1.9.10"></script>`.
    - Setup `htmx-indicator` for global loading feedback.
2. **`product_row.templ`**:
    - A single `<tr>` with `id="product-{id}"`.
    - Contains `hx-delete` for soft/hard delete.
3. **`cart_view.templ`**:
    - Calculates totals on the fly using Go logic before sending HTML to the client.

### HTMX Patterns to Use
- **Active Search**: Use `hx-trigger="keyup delay:500ms, changed"` on the search input to filter products.
- **OutOfBand Updates (OOB)**: Use `hx-swap-oob="true"` to update the total revenue card while processing a transaction.
- **Validation**: Use `hx-post` on input blur to validate stock availability before the user hits "Checkout".

---

## 5. Build Pipeline for AI Agent
To generate the UI, the agent must:
1.  **Install Templ**: `go install github.com/a-h/templ/cmd/templ@latest`.
2.  **Generate Go code from templates**: Run `templ generate`.
3.  **Static Embedding**: Use `//go:embed static/*` to bundle assets into the single binary for easy deployment to Zeabur/Docker.

## 🔗 References
- **Templ Guide**: [https://templ.guide](https://templ.guide)
- **HTMX + Go Patterns**: [https://htmx.org/essays/template-fragments/](https://htmx.org/essays/template-fragments/)
- **Alpine.js for Modals**: [https://alpinejs.dev/](https://alpinejs.dev/)