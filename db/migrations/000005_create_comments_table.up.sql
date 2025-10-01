-- public."comments" definition

-- Drop table

-- DROP TABLE public."comments";

CREATE TABLE public."comments" (
	id int4 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1 NO CYCLE) NOT NULL,
	user_id int4 NOT NULL,
	post_id int4 NOT NULL,
	"content" text NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT comments_pkey PRIMARY KEY (id)
);


-- public."comments" foreign keys

ALTER TABLE public."comments" ADD CONSTRAINT comments_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;
ALTER TABLE public."comments" ADD CONSTRAINT comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;