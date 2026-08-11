DROP INDEX IF EXISTS idx_licenses_batch;
ALTER TABLE licenses DROP COLUMN IF EXISTS batch_id;
DROP TABLE IF EXISTS license_batches;
ALTER TABLE plans DROP COLUMN IF EXISTS duration_days;
