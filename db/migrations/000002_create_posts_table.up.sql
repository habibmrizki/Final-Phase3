-- public.posts definition

-- Drop table

-- DROP TABLE public.posts;

CREATE TABLE public.posts (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	user_id int4 NOT NULL,
	"content" text NULL,
	image text NULL,
	created_at timestamp NULL,
	CONSTRAINT posts_pkey PRIMARY KEY (id)
);


-- public.posts foreign keys

ALTER TABLE public.posts ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;