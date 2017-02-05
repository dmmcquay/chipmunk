package chipmunk

const createdb = `
CREATE TABLE IF NOT EXISTS
users (
    id uuid PRIMARY KEY,
    email varchar(64) UNIQUE,
    admin boolean DEFAULT false
);

CREATE TABLE IF NOT EXISTS
categories (
    id SERIAL PRIMARY KEY,
    name varchar(1024) UNIQUE
);

CREATE TABLE IF NOT EXISTS
tranx (
    id SERIAL PRIMARY KEY,
	cost numeric DEFAULT 0,
    store varchar(64) UNIQUE,
    info varchar(1024) UNIQUE,
    user uuid
);
`

const primeCategories = `
INSERT INTO categories (name) VALUES
    ('Derek'),
    ('Colleen'),
    ('Groceries'),
    ('Restaurant'),
    ('Misc')
;
`
