DROP DATABASE IF EXISTS deepface_data;
CREATE DATABASE deepface_data;
\c deepface_data;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION pg_trgm;

CREATE OR REPLACE FUNCTION update_changetimestamp_column()
RETURNS TRIGGER AS $$
  BEGIN
   NEW.uts = now(); 
   RETURN NEW;
  END;
$$ language 'plpgsql';


--组织表

CREATE TABLE public.org_structure
(
    uts timestamp without time zone NOT NULL DEFAULT now(),
    ts bigint NOT NULL DEFAULT 0,
    org_id bigint NOT NULL,
    org_name character varying(1024) COLLATE pg_catalog."default" DEFAULT ''::character varying,
    org_level bigint NOT NULL DEFAULT 0,
    superior_org_id bigint NOT NULL DEFAULT 0,
    comment character varying COLLATE pg_catalog."default" DEFAULT ''::character varying,
    status smallint NOT NULL DEFAULT 1,
    data_org_id bigint NOT NULL DEFAULT 0,
    org_type smallint NOT NULL DEFAULT 0,
    area_code smallint NOT NULL DEFAULT 0,
    ext_data jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT org_structure_pkey PRIMARY KEY (org_id)
);
ALTER TABLE public.org_structure OWNER to postgres;

CREATE TRIGGER update_org_structure_changetimestamp BEFORE UPDATE
  ON org_structure FOR EACH ROW EXECUTE PROCEDURE 
  update_changetimestamp_column();

create index idx_org on org_structure using gin(ext_data);
INSERT INTO org_structure (org_id, org_name, org_level, superior_org_id, comment) VALUES(1, '根组织', 0, 1, 'System Default Org');
ALTER TABLE org_structure ADD CONSTRAINT superior_org_id_fk FOREIGN KEY(superior_org_id) REFERENCES org_structure (org_id) ON DELETE CASCADE;





CREATE TABLE public.statistic_day
(
    id bigint NOT NULL ,
    category character varying(1024) COLLATE pg_catalog."default",
    ts bigint DEFAULT 0,
    gas_station_code character varying(1024) COLLATE pg_catalog."default",
    date_time timestamp without time zone,
    group_data jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT statistic_day_pkey PRIMARY KEY (id)
);

CREATE TABLE public.statistic_hour
(
    id bigint NOT NULL ,
    category character varying(1024) COLLATE pg_catalog."default",
    ts bigint DEFAULT 0,
    gas_station_code character varying(1024) COLLATE pg_catalog."default",
    date_time timestamp without time zone,
    group_data jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT statistic_hour_pkey PRIMARY KEY (id)
);
create index data_index on statistic_day using gin(group_data);
create index datah_index on statistic_hour using gin(group_data);

CREATE TABLE public.stats_group_data
(
    type text COLLATE pg_catalog."default",
    count bigint,
    plate_count bigint,
    avg_data real,
    sum_data real
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.stats_group_data
    OWNER to postgres;
-- func_role
-- Table Definition ----------------------------------------------

CREATE TABLE func_role (
    func_role_id bigint PRIMARY KEY,
    func_role_name character varying(1024) NOT NULL DEFAULT ''::character varying,
    ts bigint NOT NULL DEFAULT 0,
    uts timestamp without time zone NOT NULL DEFAULT now(),
    content jsonb DEFAULT '{}'::jsonb,
    comment character varying NOT NULL DEFAULT ''::character varying,
    status smallint NOT NULL DEFAULT 1
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX func_role_pkey ON func_role(func_role_id int8_ops);

-- account
-- Table Definition ----------------------------------------------

CREATE TABLE account (
    uts timestamp without time zone NOT NULL DEFAULT now(),
    ts bigint NOT NULL DEFAULT 0,
    user_id bigint PRIMARY KEY,
    user_name character varying(1024) NOT NULL DEFAULT ''::character varying UNIQUE,
    user_passwd character varying(1024) NOT NULL DEFAULT ''::character varying,
    org_id bigint NOT NULL REFERENCES org_structure(org_id) ON DELETE CASCADE,
    func_role_id bigint NOT NULL REFERENCES func_role(func_role_id) ON DELETE CASCADE,
    security_token character varying(1024) NOT NULL DEFAULT ''::character varying UNIQUE,
    is_valid boolean NOT NULL DEFAULT true,
    real_name character varying(1024) DEFAULT ''::character varying,
    comment character varying DEFAULT ''::character varying,
    status smallint NOT NULL DEFAULT 1
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX account_pkey ON account(user_id int8_ops);
CREATE UNIQUE INDEX account_user_name_key ON account(user_name text_ops);
CREATE UNIQUE INDEX account_security_token_key ON account(security_token text_ops);

-- default records
INSERT INTO func_role(ts,func_role_id, func_role_name, comment) VALUES(round(extract(epoch FROM now())*1000), 1, '超级管理员', 'Super Admin');
INSERT INTO account(ts,user_id, user_name, user_passwd, org_id, func_role_id, is_valid, comment,security_token) VALUES(round(extract(epoch FROM now())*1000), 1, 'admin', 'admin@2013', 1, 1, true, 'Super Admin','admin');

CREATE FUNCTION count_estimate(query text) RETURNS integer AS $$
DECLARE
  rec   record;
  rows  integer;
BEGIN
  FOR rec IN EXECUTE 'EXPLAIN ' || query LOOP
    rows := substring(rec."QUERY PLAN" FROM ' rows=([[:digit:]]+)');
    EXIT WHEN rows IS NOT NULL;
  END LOOP;

  IF(rows<10000) THEN
      EXECUTE ' select count(*) from ('||query||') c ' into rows;
    END IF;
  RETURN rows;
END;
$$ LANGUAGE plpgsql VOLATILE STRICT;