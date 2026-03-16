CREATE TABLE IF NOT EXISTS todo (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    description VARCHAR(400) NOT NULL,
    status boolean not null default false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_in_list uuid,
    FOREIGN KEY(created_in_list) REFERENCES list(id)
);