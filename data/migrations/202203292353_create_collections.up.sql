CREATE TABLE IF NOT EXISTS collections(
   id VARCHAR(64) PRIMARY KEY,
   label VARCHAR(128) NOT NULL,
   aliases TEXT,
   average_sale_price TEXT NOT NULL,
   floor_price TEXT NOT NULL,
   number_active TEXT NOT NULL,
   number_of_sales TEXT NOT NULL,
   total_royalties TEXT NOT NULL,
   total_volume TEXT NOT NULL,
   last_updated TIMESTAMP NOT NULL DEFAULT NOW()
);
