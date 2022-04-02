BEGIN;

ALTER TABLE floor_price_diffs DROP COLUMN rank;

COMMIT;
