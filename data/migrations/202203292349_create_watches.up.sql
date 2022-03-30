CREATE TABLE IF NOT EXISTS watches(
   chat_id INT,
   collection_id VARCHAR (64) NOT NULL,
   last_updated TIMESTAMP NOT NULL,
   UNIQUE (chat_id, collection_id)
);
