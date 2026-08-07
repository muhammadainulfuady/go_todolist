-- =====================================================================
-- Todo List REST API - Migration 000002: seed_priorities
-- Down: menghapus data seed prioritas
-- =====================================================================

DELETE FROM priorities WHERE id_priorities IN (1, 2, 3, 4);