package chipmunk

const createdb = `
CREATE TABLE IF NOT EXISTS
users (
    id SERIAL PRIMARY KEY,
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
    store varchar(64),
    info varchar(1024),
	category_id integer references categories(id) DEFAULT 0, 
	date timestamp DEFAULT CURRENT_TIMESTAMP,
	user_id integer references users(id) DEFAULT 0
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
