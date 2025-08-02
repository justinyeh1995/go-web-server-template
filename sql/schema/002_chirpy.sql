-- +goose Up
CREATE TABLE chirpies (
    id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    body TEXT NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT chirpies_pk PRIMARY KEY (id),
    CONSTRAINT chirpies_users_fk 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE chirpies;
