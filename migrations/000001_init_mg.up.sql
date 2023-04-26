-- public.manga_updates definition

-- Drop table

-- DROP TABLE public.manga_updates;

CREATE TABLE public.manga_updates (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" varchar NOT NULL,
	raw_url varchar NOT NULL,
	last_chapter int4 NULL,
	last_checked_at timestamp NULL,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);