-- =====================================================================
-- Todo List REST API - Migration 000001: create_priorities_table (up)
-- Membuat tabel master prioritas (Matriks Eisenhower, skala 1-4)
-- =====================================================================

CREATE TABLE IF NOT EXISTS priorities (
    id_priorities   INT PRIMARY KEY,
    name            VARCHAR(50) NOT NULL,
    description     TEXT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;