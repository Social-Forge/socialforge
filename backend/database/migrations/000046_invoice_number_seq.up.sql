-- Global sequence for human-friendly, cross-tenant-unique invoice numbers.
-- (Per-tenant MAX+1 would collide under the global UNIQUE constraint.)
CREATE SEQUENCE IF NOT EXISTS invoice_number_seq START 1000;
