CREATE TABLE IF NOT EXISTS watches_channel(
   chat_id VARCHAR (128) NOT NULL,
   collection_id VARCHAR (64) NOT NULL,
   last_updated TIMESTAMP NOT NULL,
   UNIQUE (chat_id, collection_id)
);
