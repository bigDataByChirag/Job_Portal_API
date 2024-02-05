CREATE TABLE companies (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  address TEXT NOT NULL,
  userId SERIAL,
  FOREIGN KEY (userId) REFERENCES users (id)
);