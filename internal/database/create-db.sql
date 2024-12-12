-- TABLES
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    user_type TEXT NOT NULL,
    user_status TEXT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS somethings (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- FUNCTION for TIMESTAMP
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER
    LANGUAGE plpgsql
    AS $$
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
$$;


-- TRIGGERS for TIMESTAMP
CREATE OR REPLACE TRIGGER update_timestamp BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE OR REPLACE TRIGGER update_timestamp BEFORE UPDATE ON somethings FOR EACH ROW EXECUTE FUNCTION set_updated_at();


-- FUNCTION to ADD CONSTRAINT IF NOT EXISTS (Postgresql doesn't have this built-in yet)
CREATE OR REPLACE FUNCTION add_constraint_if_not_exists(table_name text, constraint_name text, constraint_sql text) RETURNS void
    LANGUAGE plpgsql
    AS $$
    BEGIN
        IF NOT EXISTS (SELECT c.conname
                       FROM pg_constraint AS c
                       INNER JOIN pg_class AS t ON c.conrelid = t."oid"
                       WHERE t.relname = table_name AND c.conname = constraint_name)
        THEN
            EXECUTE 'ALTER TABLE ' || table_name || ' ADD CONSTRAINT ' || constraint_name || ' ' || constraint_sql;
        END IF;
    END;
$$;


-- FOREIGN KEY CONSTRAINTS
SELECT add_constraint_if_not_exists('somethings', 'fk_somethings_users', 'FOREIGN KEY (user_id) REFERENCES users(id);');
