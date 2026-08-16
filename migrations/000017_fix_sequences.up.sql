-- Migration Up: 000017_fix_sequences
--
-- Seed migration 000002 memasukkan unit dengan id EKSPLISIT (1..8) tanpa
-- memajukan sequence SERIAL masing-masing.  Pada DB yang baru di-clone,
-- urutan nextval() masih kembali ke nilai lama, sehingga INSERT baru gagal:
--
--   ERROR: duplicate key value violates unique constraint "units_db_pkey"
--   (SQLSTATE 23505)
--
-- Perbaikan universal: loncengkan SELURUH sequence SERIAL ke MAX(id) tiap
-- tabel (atau ke 1 bila tabel kosong), supaya setiap nextval() menghasilkan
-- id bebas.  Idempotent — aman dijalankan berulang kali.
--
-- Nama sequence dibaca langsung dari kolom column_default (bentuk
-- nextval('...seq'::regclass)) lewat regexp_match, sehingga tidak bergantung
-- pada resolusi overload pg_get_serial_sequence() yang tidak konsisten di
-- environment ini.

DO $$
DECLARE
    rec     RECORD;
    seqname TEXT;
    col     TEXT;
    tbl     TEXT;
    mx      BIGINT;
BEGIN
    FOR rec IN
        SELECT c.table_name::text     AS t,
               c.column_name::text    AS c,
               c.column_default::text AS def
        FROM information_schema.columns c
        WHERE c.table_schema = 'public'
          AND c.column_default LIKE 'nextval%'
          AND c.table_name <> 'schema_migrations'
        ORDER BY c.table_name, c.column_name
    LOOP
        seqname := NULL;
        SELECT (regexp_match(rec.def, '^nextval\(''(.*)''::regclass\)$'))[1]
            INTO seqname;

        IF seqname IS NULL THEN
            RAISE LOG '[seq-fix] skip %.% (pola nextval tak dikenal)', rec.t, rec.c;
            CONTINUE;
        END IF;

        col := rec.c;
        tbl := rec.t;

        EXECUTE format('SELECT MAX(%I) FROM %I.%I', col, 'public', tbl)
            INTO mx;

        IF mx IS NULL THEN
            EXECUTE format('SELECT setval(%L::regclass, 1, false)', seqname);
            RAISE LOG '[seq-fix] reset %.% -> 1 (tabel kosong)', tbl, col;
        ELSE
            EXECUTE format('SELECT setval(%L::regclass, %s, true)', seqname, mx);
            RAISE LOG '[seq-fix] align %.% -> %s', tbl, col, mx;
        END IF;
    END LOOP;
END $$;