--Решил использовать PostgreSQL....

--create DATABASE detector;


CREATE TABLE public.url
(
    url_id      BIGSERIAL PRIMARY KEY,
    url_path    TEXT    NOT NULL UNIQUE,
    url_dropped BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE public.history
(
    hst_id     BIGSERIAL PRIMARY KEY,
    hst_check  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    hst_status INT NOT NULL,
    hst_url_id BIGINT REFERENCES public.url
);

ALTER TABLE public.history
    ADD CONSTRAINT hst_primary_key UNIQUE (hst_check, hst_url_id, hst_status);
-----------------------------------------------------------------------------
-----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION add_status(path TEXT, status INT)
    RETURNS VOID AS
$$
DECLARE
    id BIGINT = (SELECT url_id
                 FROM public.url
                 WHERE url_path = path
                 LIMIT 1);
BEGIN
    IF id IS NULL THEN
        INSERT INTO public.url(url_path) VALUES (path) RETURNING url_id INTO id;
    END IF;

    INSERT INTO public.history(hst_status, hst_url_id) VALUES (status, id);
END
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_history(path TEXT)
    RETURNS JSONB AS
$$
DECLARE
    id BIGINT = (SELECT url_id
                 FROM public.url
                 WHERE url_path = path);
BEGIN
    IF id IS NULL THEN
        RAISE EXCEPTION 'NOT FOUND %', path;
    END IF;

    RETURN (
        SELECT json_agg(batch)
        FROM (
                 SELECT u.url_path, h.hst_status, h.hst_check
                 FROM history AS h
                          INNER JOIN url AS U ON
                              h.hst_url_id = id AND
                              h.hst_id = h.hst_url_id AND
                              NOT u.url_dropped
        ) AS batch
    );
END;
$$ LANGUAGE plpgsql;


select add_status('https://ozon.ru', 200);

select get_history('https://ya.ru');

select *
from history;

select *
from url;


alter role postgres with password 'pwd';