CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
	id uuid DEFAULT public.uuid_generate_v4(),
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE user_details (
	id uuid DEFAULT public.uuid_generate_v4(),
    user_id uuid NOT NULL,
	email varchar NOT NULL,
	"address" varchar NOT NULL,
    phone_number varchar NOT NULL,
    age int NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT user_details_pkey PRIMARY KEY (id)
);

ALTER TABLE ONLY user_details ADD CONSTRAINT fk_user_details FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL;
