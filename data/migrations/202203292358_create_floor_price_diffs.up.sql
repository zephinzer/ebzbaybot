CREATE TABLE IF NOT EXISTS floor_price_diffs(
   collection_id VARCHAR(64) UNIQUE,
   previous_price TEXT NOT NULL,
   current_price TEXT NOT NULL,
   last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
   CONSTRAINT fk_collection_id
      FOREIGN KEY(collection_id)
      REFERENCES collections(id)
);
