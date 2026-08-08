-- =====================================================================
-- Todo List REST API - Migration 000003: create_todos_table (down)
-- Menghapus tabel tugas (drop dulu karena punya foreign key)
-- =====================================================================

DROP TABLE IF EXISTS todos;