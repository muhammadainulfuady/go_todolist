-- =====================================================================
-- Todo List REST API - Migration 000001: create_tables
-- Up: membuat tabel users, priorities, todos
-- =====================================================================

CREATE TABLE users (
    id_users        INT AUTO_INCREMENT PRIMARY KEY,
    nama            VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    foto_profile    VARCHAR(255) NULL,
    otp_code        VARCHAR(6)   NULL,
    otp_expires_at  DATETIME     NULL,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE priorities (
    id_priorities   INT PRIMARY KEY,
    name            VARCHAR(50) NOT NULL,
    description     TEXT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE todos (
    id_todos        INT AUTO_INCREMENT PRIMARY KEY,
    id_users        INT NOT NULL,
    id_priorities   INT NOT NULL,
    title           VARCHAR(150) NOT NULL,
    slug            VARCHAR(200) NOT NULL UNIQUE,
    description     TEXT NULL,
    image           VARCHAR(255) NULL,
    is_completed    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_todos_users      FOREIGN KEY (id_users)      REFERENCES users (id_users)      ON DELETE CASCADE,
    CONSTRAINT fk_todos_priorities FOREIGN KEY (id_priorities) REFERENCES priorities (id_priorities)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;