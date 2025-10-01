-- public.follows definition

-- Drop table

-- DROP TABLE public.follows;

CREATE TABLE public.follows (
	follower int4 NOT NULL,
	"following" int4 NOT NULL,
	CONSTRAINT follows_pkey PRIMARY KEY (follower, following)
);


-- public.follows foreign keys

ALTER TABLE public.follows ADD CONSTRAINT fk_follower FOREIGN KEY (follower) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE public.follows ADD CONSTRAINT fk_following FOREIGN KEY ("following") REFERENCES public.users(id) ON DELETE CASCADE;