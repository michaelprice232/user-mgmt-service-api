CREATE TABLE IF NOT EXISTS users (
      user_id serial PRIMARY KEY,
      full_name  VARCHAR ( 50 ) NOT NULL,
      email VARCHAR ( 255 ) NOT NULL
);