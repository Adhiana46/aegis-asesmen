CREATE TABLE IF NOT EXISTS "organizations" (
    id CHAR(26) NOT NULL,
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by CHAR(26) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by CHAR(26) NOT NULL,
    PRIMARY KEY (id)
);