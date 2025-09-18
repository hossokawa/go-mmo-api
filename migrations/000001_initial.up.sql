CREATE TABLE IF NOT EXISTS player (
  id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  username TEXT NOT NULL,
  class TEXT NOT NULL,
  level INT NOT NULL,
  gold INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX ON player(username);

INSERT INTO player (username, class, level, gold, created_at, updated_at) VALUES
('Sephiro', 'Warrior', 32, 4309, '2024-01-15 14:23:45Z', '2024-02-10 08:30:15Z'),
('SpellQueen', 'Mage', 47, 7390, '2023-12-05 10:15:12Z', '2024-01-22 16:45:22Z'),
('StormcallerX', 'Druid', 29, 3982, '2023-11-20 09:20:34Z', '2024-01-10 12:10:58Z'),
('ShadowMaster', 'Rogue', 50, 7521, '2023-10-08 19:50:14Z', '2023-12-18 11:35:45Z'),
('FireFurry', 'Sorcerer', 55, 8097, '2023-09-03 17:45:23Z', '2023-11-14 14:20:30Z');
