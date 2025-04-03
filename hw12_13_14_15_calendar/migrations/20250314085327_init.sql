-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS public.sid
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE public.sid OWNER TO postgres;

CREATE TABLE IF NOT EXISTS public.tuser
(
    id bigint NOT NULL,
    name character varying(250) COLLATE pg_catalog."default" NOT NULL,
    address character varying(250) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT tuser_pkey PRIMARY KEY (id)
) 
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tuser OWNER to postgres;
GRANT ALL ON TABLE public.tuser TO postgres;
COMMENT ON TABLE public.tuser IS 'Пользователь';

CREATE TABLE IF NOT EXISTS public.tevent
(
    id bigint NOT NULL,
    name character varying(25) COLLATE pg_catalog."default" NOT NULL,
    start_date_time timestamp without time zone NOT NULL,
    end_date_time timestamp without time zone,
    user_id bigint NOT NULL,
    description character varying(500) COLLATE pg_catalog."default" NOT NULL,
    remind_for integer,
    CONSTRAINT tevent_pkey PRIMARY KEY (id),
    CONSTRAINT fk_event_user FOREIGN KEY (user_id)
        REFERENCES public.tuser (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)
TABLESPACE pg_default;
ALTER TABLE IF EXISTS public.tevent OWNER to postgres;
GRANT ALL ON TABLE public.tevent TO postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.tevent;
DROP TABLE IF EXISTS public.tuser;
DROP SEQUENCE IF EXISTS public.sid;
-- +goose StatementEnd
