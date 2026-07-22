
CREATE TABLE IF NOT EXISTS migrations (
    idx integer,
    enabled TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (idx)
);

INSERT into migrations (idx)
values
(0),
(1),
(2),
(3),
(4),
(5),
(6),
(7),
(8),
(9),
(10),
(11),
(12),
(13),
(14),
(15),
(16)
