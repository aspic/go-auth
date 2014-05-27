CREATE TABLE Identity (
        id SERIAL PRIMARY KEY NOT NULL,
        username VARCHAR UNIQUE,
        pw_hash VARCHAR,
        email VARCHAR,
        salt VARCHAR,
        created TIMESTAMP DEFAULT now(),
        requested_key VARCHAR,
        requested_at TIMESTAMP
);
CREATE TABLE Realm (
        id SERIAL PRIMARY KEY NOT NULL,
        name VARCHAR NOT NULL,
        key VARCHAR NOT NULL
);
CREATE TABLE InRealm (
        id integer NOT NULL REFERENCES identity(id),
        realm integer NOT NULL REFERENCES realm(id),
        PRIMARY KEY(id, realm)
);
