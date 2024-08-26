CREATE TABLE item
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    url  VARCHAR(255) NOT NULL
);

CREATE TABLE price
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id    INTEGER        NOT NULL REFERENCES item (id),
    price      DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP      NOT NULL
);

CREATE TABLE source
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id   INTEGER      NOT NULL REFERENCES item (id),
    vendor_id INTEGER      NOT NULL REFERENCES vendor (id),
    url       VARCHAR(255) NOT NULL
);

CREATE TABLE vendor
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         VARCHAR(255) NOT NULL,
    css_selector VARCHAR(255) NOT NULL
)