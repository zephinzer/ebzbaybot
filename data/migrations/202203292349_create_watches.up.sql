CREATE TABLE IF NOT EXISTS watches(
   chat_id INT PRIMARY KEY,
   collection VARCHAR (42) UNIQUE NOT NULL,
   last_updated TIMESTAMP NOT NULL DEFAULT NOW()
);
