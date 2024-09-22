CREATE TABLE sessions ( 
    token char(43) primary key,
    data blob not null,
    expiry timestamp(6) not null
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
