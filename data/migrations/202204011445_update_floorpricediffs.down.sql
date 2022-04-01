BEGIN;

ALTER TABLE floor_price_diffs DROP COLUMN score;
ALTER TABLE floor_price_diffs DROP COLUMN edition;
ALTER TABLE floor_price_diffs DROP COLUMN image_url;
ALTER TABLE floor_price_diffs DROP COLUMN listing_id;

COMMIT;
