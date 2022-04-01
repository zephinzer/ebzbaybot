BEGIN;

ALTER TABLE floor_price_diffs ADD COLUMN listing_id VARCHAR(32) NULL;
ALTER TABLE floor_price_diffs ADD COLUMN image_url TEXT NULL;
ALTER TABLE floor_price_diffs ADD COLUMN edition VARCHAR(32) NULL;
ALTER TABLE floor_price_diffs ADD COLUMN score VARCHAR(32) NULL;

COMMIT;
