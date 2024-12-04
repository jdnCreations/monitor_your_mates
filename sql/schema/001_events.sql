-- +goose Up
CREATE TABLE events (
  id INT,
  message TEXT,
  created_at TIMESTAMP,
  severity TEXT,
  PRIMARY KEY (id, created_at)
);

-- +goose Down
DROP TABLE if exists events; 