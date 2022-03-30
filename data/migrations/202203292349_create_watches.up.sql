CREATE TABLE IF NOT EXISTS watches(
   chat_id INT PRIMARY KEY,
   collection_id VARCHAR (64) NOT NULL,
   last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
   UNIQUE (chat_id, collection_id)
);
