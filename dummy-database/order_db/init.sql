-- Hapus tabel jika sudah ada (untuk pengembangan/pengujian)
DROP TABLE IF EXISTS orders CASCADE;

-- 1. Buat Tabel 'orders'
CREATE TABLE orders (
      order_id SERIAL PRIMARY KEY,
      customer_id INTEGER NOT NULL,
      order_date DATE NOT NULL DEFAULT CURRENT_DATE,
      total_amount NUMERIC(10, 2) NOT NULL
);

-- 2. Tambahkan 5 Data Awal
INSERT INTO orders (customer_id, order_date, total_amount) VALUES
(101, '2025-10-20', 150.50),
(102, '2025-10-21', 45.99),
(101, '2025-10-25', 300.75),
(103, '2025-11-01', 12.00),
(104, '2025-11-05', 520.00);