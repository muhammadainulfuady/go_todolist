-- =====================================================================
-- Todo List REST API - Migration 000001: create_tables
-- Down: menghapus tabel (urutan terbalik dari FK dependency)
-- =====================================================================

DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS priorities;
DROP TABLE IF EXISTS users;