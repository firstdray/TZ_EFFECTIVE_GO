CREATE TABLE people
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    surname     VARCHAR(255) NOT NULL,
    patronymic  VARCHAR(255),
    age         INTEGER,
    gender      VARCHAR(20)  NOT NULL,
    national VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_surname ON people(surname);