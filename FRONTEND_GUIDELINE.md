# Panduan Implementasi Frontend: Go + HTMX + Templ + Tailwind

Dokumen ini adalah panduan teknis bagi AI Agent untuk mengimplementasikan dashboard frontend sesuai dengan blok "Frontend Implementation" pada `TODO.md`.

## 🛠 Tech Stack Utama
- **Go 1.22+**: Core backend.
- **HTMX**: Untuk interaksi AJAX tanpa JavaScript kompleks. [Ref: htmx.org](https://htmx.org/)
- **Templ**: Type-safe HTML components dalam Go. [Ref: templ.guide](https://templ.guide/)
- **Tailwind CSS**: Styling via CDN atau JIT build.
- **Alpine.js**: (Opsional) Untuk interaksi UI sisi klien yang sangat ringan (seperti modal/dropdown).

---

## 🏗 Arsitektur Folder Baru
```text
.
├── cmd/
│   └── web/            # Entry point untuk web server jika ingin dipisah dari API
├── internal/
│   └── ui/             # Logika spesifik UI
│       ├── components/ # Templ components (button, card, navbar)
│       └── pages/      # Templ pages (dashboard, products, reports)
├── static/             # File statis (css, js, images)
└── main.go             # Registrasi route web handlers
```

---

## 📝 Instruksi Implementasi Per Bagian

### 1. Setup Infrastruktur (Templ & Static Files)
- **Tooling**: Gunakan `templ generate` untuk compile file `.templ`.
- **Static Handlers**:
  ```go
  fs := http.FileServer(http.Dir("./static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  ```
- **Base Layout**: Buat `layout.templ` yang membungkus konten dengan Bootstrap/Tailwind dan script HTMX.

### 2. Product Management (Dynamic UI)
- **Live Search**: Gunakan atribut HTMX berikut pada input search:
  ```html
  <input type="search" name="search" 
         hx-get="/dashboard/produk/list" 
         hx-trigger="keyup delay:500ms, search" 
         hx-target="#product-table-body" 
         placeholder="Cari produk...">
  ```
- **Inline Delete**: Gunakan `hx-delete` dengan `hx-confirm` untuk proteksi. Kembalikan response kosong atau snippet kecil untuk menghapus baris dari DOM.

### 3. Dashboard & Reporting
- **Partial Updates**: Gunakan `hx-get` pada tab atau interval (`hx-trigger="every 30s"`) untuk memperbarui data laporan "Today's Sales" secara otomatis tanpa reload halaman.

---

## 🚦 Standar Koding untuk Agent
1. **Component-Based**: Setiap bagian UI (misal: satu baris tabel produk) harus menjadi fungsi `templ` sendiri agar bisa dirender ulang secara parsial oleh HTMX.
2. **RESTful Integration**: Web handlers harus memanggil `Service` yang sama dengan API handlers untuk menjaga konsistensi data.
3. **Graceful Degradation**: Pastikan link tetap memiliki atribut `href` dasar meskipun menggunakan HTMX.
4. **Loading States**: Selalu tambahkan indikator loading menggunakan class `htmx-indicator`.

## 🔗 Referensi Penting
- **Templ Documentation**: [https://templ.guide](https://templ.guide) - Gunakan komponen `.templ` untuk type-safety.
- **HTMX Examples**: [https://htmx.org/examples/](https://htmx.org/examples/) - Gunakan pola "Active Search" dan "Infinite Scroll" jika perlu.
- **Go standard library `http.ServeMux`**: Gunakan fitur `r.PathValue("id")` untuk routing bersih.
