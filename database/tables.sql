CREATE TABLE manufacturer (
    make TEXT PRIMARY KEY UNIQUE NOT NULL
);

CREATE TABLE models (
    model TEXT,
    make TEXT,
    PRIMARY KEY (model, make),
    FOREIGN KEY(make) REFERENCES manufacturer(make)
);

DROP TABLE manufacturer;
DROP TABLE models;

SELECT * FROM manufacturer;
SELECT * FROM models;

DELETE FROM models;
SELECT COUNT(*) AS C FROM models;

