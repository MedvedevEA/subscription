CREATE TABLE IF NOT EXISTS public.services(
    service_id bigserial  NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT service_pk PRIMARY KEY (service_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS services_name
    ON public.services USING btree
    (name COLLATE pg_catalog."default" ASC NULLS LAST);
CREATE TABLE IF NOT EXISTS public.subscriptions(
    subscription_id bigserial NOT NULL,
    service_id bigint NOT NULL,
    price integer NOT NULL,
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    stop_date date,
    CONSTRAINT subscriptions_pk PRIMARY KEY (subscription_id),
    CONSTRAINT services_fk FOREIGN KEY (service_id)
        REFERENCES public.services (service_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT start_date_stop_date CHECK (start_date <= stop_date)
);
CREATE INDEX IF NOT EXISTS subscriptions_service_id
    ON public.subscriptions USING btree
    (service_id ASC NULLS LAST);