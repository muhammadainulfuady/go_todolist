-- =====================================================================
-- Todo List REST API - Migration 000002: create_users_table (up)
-- Membuat tabel pengguna (auth passwordless via email OTP)
-- =====================================================================

CREATE TABLE IF NOT EXISTS users (
    id_users        INT AUTO_INCREMENT PRIMARY KEY,
    nama            VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    foto_profile    VARCHAR(255) NULL,
    otp_code        VARCHAR(6)   NULL,
    otp_expires_at  DATETIME     NULL,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;