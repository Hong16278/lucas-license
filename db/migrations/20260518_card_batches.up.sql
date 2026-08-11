-- Card batches and activation-started terms for external auto-delivery platforms.

ALTER TABLE plans
    ADD COLUMN IF NOT EXISTS duration_days INTEGER NOT NULL DEFAULT 0
    CHECK (duration_days >= 0 AND duration_days <= 36500);

CREATE TABLE IF NOT EXISTS license_batches (
    id         TEXT PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES products(id),
    plan_id    TEXT NOT NULL REFERENCES plans(id),
    name       TEXT NOT NULL,
    channel    TEXT NOT NULL DEFAULT '',
    quantity   INTEGER NOT NULL CHECK (quantity > 0 AND quantity <= 5000),
    created_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_license_batches_product_created
    ON license_batches(product_id, created_at DESC);

ALTER TABLE licenses
    ADD COLUMN IF NOT EXISTS batch_id TEXT REFERENCES license_batches(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_licenses_batch ON licenses(batch_id);
