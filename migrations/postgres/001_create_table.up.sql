CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
                                    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    name VARCHAR(150) NOT NULL,
                                    email VARCHAR(350) NOT NULL,
                                    password VARCHAR(60) NOT NULL,
                                    path_image VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS list (
                                    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    name VARCHAR(150) NOT NULL,
                                    created_by uuid,
                                    FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE ,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS todo (
                                    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    description VARCHAR(400) NOT NULL,
                                    status boolean not null default false,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    created_in_list uuid,
                                    FOREIGN KEY(created_in_list) REFERENCES list(id) ON DELETE CASCADE
);