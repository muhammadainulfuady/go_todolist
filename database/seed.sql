-- =====================================================================
-- Todo List REST API - Seed Master Data Priorities (Eisenhower Matrix)
-- =====================================================================

INSERT INTO priorities (id_priorities, name, description) VALUES
    (1, 'Penting & Mendesak',          'Kerjakan sekarang'),
    (2, 'Penting tapi Tidak Mendesak', 'Jadwalkan waktu khusus'),
    (3, 'Tidak Penting tapi Mendesak', 'Delegasikan jika memungkinkan'),
    (4, 'Tidak Penting & Tidak Mendesak', 'Hapus atau tunda dulu')
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    description = VALUES(description);