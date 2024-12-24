-- TABLES
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    otp TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL,
    status TEXT NOT NULL,
    image_url TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    text_content TEXT NOT NULL,
    has_viewed BOOLEAN NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- FUNCTION for TIMESTAMP
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER
LANGUAGE plpgsql AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


-- TRIGGER for TIMESTAMP
CREATE OR REPLACE TRIGGER update_timestamp BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();


-- FUNCTION to 'ADD CONSTRAINT IF NOT EXISTS' (Postgresql doesn't have this built-in)
CREATE OR REPLACE FUNCTION add_constraint_if_not_exists(table_name text, constraint_name text, constraint_sql text) RETURNS void
LANGUAGE plpgsql AS
$$
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
SELECT add_constraint_if_not_exists('notifications', 'fk_notifications_users', 'FOREIGN KEY (user_id) REFERENCES users(id);');


-- INSERT FAKE USERS (all the passwords are '12345678', but they have been hashed by bcrypt)
-- INSERT INTO users (email, username, password, user_type, user_status, image_url)
-- VALUES ('fakeadmin@email.com', 'jill78', '$2a$10$OWojPnrPLxX0TfV5NCaqEu65gSOKaWCAcupoYmekuxKq1eHC68Ulq', 'admin', 'active', '');
-- INSERT INTO users (email, username, password, user_type, user_status, image_url)
-- VALUES ('fakeuser@email.com', 'sam44', '$2a$10$OWojPnrPLxX0TfV5NCaqEu65gSOKaWCAcupoYmekuxKq1eHC68Ulq', 'regular', 'active', '');


--PREPARED STATEMENTS
-- PREPARE get_user_by_email (text) AS
--     SELECT * FROM users WHERE email = $1 LIMIT 1;




-- CREATE OR REPLACE FUNCTION get_user_by_email(email_address text)
-- RETURNS TABLE (id int, email text, username text, password_hash text, user_type text, user_status text, image_url text, created_at timestamptz, updated_at timestamptz)
-- LANGUAGE SQL AS
-- $$
--     SELECT * FROM users WHERE email = email_address;
-- $$;