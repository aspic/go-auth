CREATE TABLE Identity (
        id SERIAL PRIMARY KEY NOT NULL,
        username VARCHAR,
        password VARCHAR,
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
        identity REFERENCES Identity(id) NOT NULL,
        realm REFERENCES Realm(id) NOT NULL,
        PRIMARY KEY(identity, realm)
);
