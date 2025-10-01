-- public.likes definition

-- Drop table

-- DROP TABLE public.likes;

CREATE TABLE public.likes (
	user_id int4 NOT NULL,
	post_id int4 NOT NULL,
	CONSTRAINT likes_pkey PRIMARY KEY (user_id, post_id)
);


-- public.likes foreign keys

ALTER TABLE public.likes ADD CONSTRAINT likes_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;
ALTER TABLE public.likes ADD CONSTRAINT likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;