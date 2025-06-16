/*所有測試/開發環境會重建用*/
CREATE TABLE IF NOT EXISTS members
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT        NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

