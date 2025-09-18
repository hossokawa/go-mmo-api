CREATE TABLE IF NOT EXISTS inventory (
  player_id INT NOT NULL REFERENCES player(id),
  item_id UUID NOT NULL REFERENCES item(id),
  PRIMARY KEY (player_id, item_id)
);

CREATE INDEX ON inventory(player_id);

INSERT INTO inventory (player_id, item_id) VALUES
(1, '42c6294c-56de-49d2-be2e-055b2a2151a6'),
(2, '2e9e9593-c5ec-4554-9e15-131aa0b63127'),
(3, '0faffedc-2047-4616-9357-51d22fe80ff7'),
(4, '42c6294c-56de-49d2-be2e-055b2a2151a6'),
(5, '0faffedc-2047-4616-9357-51d22fe80ff7');
