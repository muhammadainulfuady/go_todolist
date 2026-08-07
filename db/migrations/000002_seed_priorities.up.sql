-- =====================================================================
-- Todo List REST API - Migration 000002: seed_priorities
-- Up: mengisi 4 master data prioritas (Eisenhower Matrix)
-- =====================================================================

INSERT INTO priorities (id_priorities, name, description) VALUES
    (1, 'Penting & Mendesak',          'Kerjakan sekarang'),
    (2, 'Penting tapi Tidak Mendesak', 'Jadwalkan waktu khusus'),
    (3, 'Tidak Penting tapi Mendesak', 'Delegasikan jika memungkinkan'),
    (4, 'Tidak Penting & Tidak Mendesak', 'Hapus atau tunda dulu');