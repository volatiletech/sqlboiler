CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE workday AS ENUM('monday', 'tuesday', 'wednesday', 'thursday', 'friday');
CREATE TYPE faceyface AS ENUM('angry', 'hungry', 'bitter');

CREATE TABLE event_one (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  day    workday NOT NULL
);

CREATE TABLE event_two (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  day    workday NOT NULL
);

CREATE TABLE event_three (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  day    workday NOT NULL,
  face   faceyface NOT NULL,
  thing  workday NULL,
  stuff  faceyface NULL
);

CREATE TABLE facey (
  id serial PRIMARY KEY NOT NULL,
  name VARCHAR(255),
  face faceyface NOT NULL
);

CREATE TABLE magic (
  id       serial PRIMARY KEY NOT NULL,
  id_two   serial NOT NULL,
  id_three serial,

  bool_zero   bool,
  bool_one    bool NULL,
  bool_two    bool NOT NULL,
  bool_three  bool NULL DEFAULT FALSE,
  bool_four   bool NULL DEFAULT TRUE,
  bool_five   bool NOT NULL DEFAULT FALSE,
  bool_six    bool NOT NULL DEFAULT TRUE,

  string_zero   VARCHAR(1),
  string_one    VARCHAR(1) NULL,
  string_two    VARCHAR(1) NOT NULL,
  string_three  VARCHAR(1) NULL DEFAULT 'a',
  string_four   VARCHAR(1) NOT NULL DEFAULT 'b',
  string_five   VARCHAR(1000),
  string_six    VARCHAR(1000) NULL,
  string_seven  VARCHAR(1000) NOT NULL,
  string_eight  VARCHAR(1000) NULL DEFAULT 'abcdefgh',
  string_nine   VARCHAR(1000) NOT NULL DEFAULT 'abcdefgh',
  string_ten    VARCHAR(1000) NULL DEFAULT '',
  string_eleven VARCHAR(1000) NOT NULL DEFAULT '',

  nonbyte_zero   CHAR(1),
  nonbyte_one    CHAR(1) NULL,
  nonbyte_two    CHAR(1) NOT NULL,
  nonbyte_three  CHAR(1) NULL DEFAULT 'a',
  nonbyte_four   CHAR(1) NOT NULL DEFAULT 'b',
  nonbyte_five   CHAR(1000),
  nonbyte_six    CHAR(1000) NULL,
  nonbyte_seven  CHAR(1000) NOT NULL,
  nonbyte_eight  CHAR(1000) NULL DEFAULT 'a',
  nonbyte_nine   CHAR(1000) NOT NULL DEFAULT 'b',

  byte_zero   "char",
  byte_one    "char" NULL,
  byte_two    "char" NULL DEFAULT 'a',
  byte_three  "char" NOT NULL,
  byte_four   "char" NOT NULL DEFAULT 'b',

  big_int_zero  bigint,
  big_int_one   bigint NULL,
  big_int_two   bigint NOT NULL,
  big_int_three bigint NULL DEFAULT 111111,
  big_int_four  bigint NOT NULL DEFAULT 222222,
  big_int_five  bigint NULL DEFAULT 0,
  big_int_six   bigint NOT NULL DEFAULT 0,

  int_zero  int,
  int_one   int NULL,
  int_two   int NOT NULL,
  int_three int NULL DEFAULT 333333,
  int_four  int NOT NULL DEFAULT 444444,
  int_five  int NULL DEFAULT 0,
  int_six   int NOT NULL DEFAULT 0,

  float_zero  decimal,
  float_one   numeric,
  float_two   numeric(2,1),
  float_three numeric(2,1),
  float_four  numeric(2,1) NULL,
  float_five  numeric(2,1) NOT NULL,
  float_six   numeric(2,1) NULL DEFAULT 1.1,
  float_seven numeric(2,1) NOT NULL DEFAULT 1.1,
  float_eight numeric(2,1) NULL DEFAULT 0.0,
  float_nine  numeric(2,1) NULL DEFAULT 0.0,

  bytea_zero  bytea,
  bytea_one   bytea NULL,
  bytea_two   bytea NOT NULL,
  bytea_three bytea NOT NULL DEFAULT 'a',
  bytea_four  bytea NULL DEFAULT 'b',
  bytea_five  bytea NOT NULL DEFAULT 'abcdefghabcdefghabcdefgh',
  bytea_six   bytea NULL DEFAULT 'hgfedcbahgfedcbahgfedcba',
  bytea_seven bytea NOT NULL DEFAULT '',
  bytea_eight bytea NOT NULL DEFAULT '',

  time_zero      timestamp,
  time_one       date,
  time_two       timestamp NULL DEFAULT NULL,
  time_three     timestamp NULL,
  time_four      timestamp NOT NULL,
  time_five      timestamp NULL DEFAULT '1999-01-08 04:05:06.789',
  time_six       timestamp NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
  time_seven     timestamp NULL DEFAULT 'January 8 04:05:06 1999 PST',
  time_eight     timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789',
  time_nine      timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
  time_ten       timestamp NOT NULL DEFAULT 'January 8 04:05:06 1999 PST',
  time_eleven    date NULL,
  time_twelve    date NOT NULL,
  time_thirteen  date NULL DEFAULT '1999-01-08',
  time_fourteen  date NULL DEFAULT 'January 8, 1999',
  time_fifteen   date NULL DEFAULT '19990108',
  time_sixteen   date NOT NULL DEFAULT '1999-01-08',
  time_seventeen date NOT NULL DEFAULT 'January 8, 1999',
  time_eighteen  date NOT NULL DEFAULT '19990108',

  uuid_zero  uuid,
  uuid_one   uuid NULL,
  uuid_two   uuid NULL DEFAULT NULL,
  uuid_three uuid NOT NULL,
  uuid_four  uuid NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
  uuid_five  uuid NOT NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',

  strange_one   integer DEFAULT '5'::integer,
  strange_two   varchar(1000) DEFAULT 5::varchar,
  strange_three timestamp without time zone default (now() at time zone 'utc'),
  strange_four  timestamp with time zone default (now() at time zone 'utc'),
  strange_five  interval NOT NULL DEFAULT '21 days',
  strange_six   interval NULL DEFAULT '23 hours',

  aa  json NULL,
  bb  json NOT NULL,
  cc  jsonb NULL,
  dd  jsonb NOT NULL,
  ee  box NULL,
  ff  box NOT NULL,
  gg  cidr NULL,
  hh  cidr NOT NULL,
  ii  circle NULL,
  jj  circle NOT NULL,
  kk  double precision NULL,
  ll  double precision NOT NULL,
  mm  inet NULL,
  nn  inet NOT NULL,
  oo  line NULL,
  pp  line NOT NULL,
  qq  lseg NULL,
  rr  lseg NOT NULL,
  ss  macaddr NULL,
  tt  macaddr NOT NULL,
  uu  money NULL,
  vv  money NOT NULL,
  ww  path NULL,
  xx  path NOT NULL,
  yy  pg_lsn NULL,
  zz  pg_lsn NOT NULL,
  aaa point NULL,
  bbb point NOT NULL,
  ccc polygon NULL,
  ddd polygon NOT NULL,
  eee tsquery NULL,
  fff tsquery NOT NULL,
  ggg tsvector NULL,
  hhh tsvector NOT NULL,
  iii txid_snapshot NULL,
  jjj txid_snapshot NOT NULL,
  kkk xml NULL,
  lll xml NOT NULL,
  mmm citext NULL,
  nnn citext NOT NULL
);

create table owner (
  id    serial primary key not null,
  name  varchar(255) not null
);

create table cats (
  id       serial primary key not null,
  name     varchar(255) not null,
  owner_id int references owner (id)
);

create table toys (
  id    serial primary key not null,
  name  varchar(255) not null
);

create table cat_toys (
  cat_id int not null references cats (id),
  toy_id int not null references toys (id),
  primary key (cat_id, toy_id)
);

create table dog_toys (
  dog_id int not null,
  toy_id int not null,
  primary key (dog_id, toy_id)
);

create table dragon_toys (
  dragon_id uuid,
  toy_id    uuid,
  primary key (dragon_id, toy_id)
);

create table spider_toys (
  spider_id uuid,
  name      character varying,
  primary key (spider_id)
);

create table pals (
  pal character varying,
  name character varying,
  primary key (pal)
);

create table friend (
  friend character varying,
  name character varying,
  primary key (friend)
);

create table bro (
  bros character varying,
  name character varying,
  primary key (bros)
);

create table enemies (
  enemies character varying,
  name character varying,
  primary key (enemies)
);

create table chocolate (
  dog varchar(100) primary key
);

create table waffles (
  cat varchar(100) primary key
);

create table fun_arrays (
  id serial,
  fun_one integer[] null,
  fun_two integer[] not null,
  fun_three boolean[] null,
  fun_four boolean[] not null,
  fun_five varchar[] null,
  fun_six varchar[] not null,
  fun_seven decimal[] null,
  fun_eight decimal[] not null,
  fun_nine bytea[] null,
  fun_ten bytea[] not null,
  fun_eleven jsonb[] null,
  fun_twelve jsonb[] not null,
  fun_thirteen json[] null,
  fun_fourteen json[] not null,
  primary key (id)
);

create table tigers (
  id    bytea primary key,
  name  bytea null
);

create table elephants (
  id        bytea primary key,
  name      bytea not null,
  tiger_id  bytea null unique,
  foreign key (tiger_id) references tigers (id)
);

create table wolves (
  id        bytea primary key,
  name      bytea not null,
  tiger_id  bytea not null unique,
  foreign key (tiger_id) references tigers (id)
);

create table ants (
  id        bytea primary key,
  name      bytea not null,
  tiger_id  bytea not null,
  foreign key (tiger_id) references tigers (id)
);

create table worms (
  id        bytea primary key,
  name      bytea not null,
  tiger_id  bytea null,
  foreign key (tiger_id) references tigers (id)
);

create table addresses (
  id bytea primary key,
  name bytea null
);

create table houses (
  id bytea primary key,
  name bytea not null,
  address_id bytea not null unique,
  foreign key (address_id) references addresses (id)
);

create table byte_pilots (
  id   bytea primary key not null,
  name character varying
);

create table byte_airports (
  id   bytea primary key not null,
  name character varying
);

create table byte_languages (
  id   bytea primary key not null,
  name character varying
);

create table byte_jets (
  id              bytea primary key not null,
  name            character varying,
  byte_pilot_id   bytea unique,
  byte_airport_id bytea,

  foreign key (byte_pilot_id) references byte_pilots (id),
  foreign key (byte_airport_id) references byte_airports (id)
);

create table byte_pilot_languages (
  byte_pilot_id    bytea not null,
  byte_language_id bytea not null,

  primary key (byte_pilot_id, byte_language_id),
  foreign key (byte_pilot_id) references byte_pilots (id),
  foreign key (byte_language_id) references byte_languages (id)
);

create table cars (
  id integer not null,
  name text,
  primary key (id)
);

create table car_cars (
  car_id integer not null,
  awesome_car_id integer not null,
  relation text not null,
  primary key (car_id, awesome_car_id),
  foreign key (car_id) references cars(id),
  foreign key (awesome_car_id) references cars(id)
);

create table trucks (
  id integer not null,
  parent_id integer,
  name text,
  primary key (id),
  foreign key (parent_id) references trucks(id)
);

CREATE TABLE race (
    id integer PRIMARY KEY NOT NULL,
    race_date timestamp,
    track text
);

CREATE TABLE race_results (
    id integer PRIMARY KEY NOT NULL,
    race_id integer,
    name text, 
    foreign key (race_id) references race(id)
);

CREATE TABLE race_result_scratchings (
    id integer PRIMARY KEY NOT NULL,
    results_id integer NOT NULL,
    name text NOT NULL,
    foreign key (results_id) references race_results(id)
);

CREATE TABLE pilots (
  id integer NOT NULL,
  name text NOT NULL
);

ALTER TABLE pilots ADD CONSTRAINT pilot_pkey PRIMARY KEY (id);

CREATE TABLE jets (
  id integer NOT NULL,
  pilot_id integer NOT NULL,
  age integer NOT NULL,
  name text NOT NULL,
  color text NOT NULL
);

ALTER TABLE jets ADD CONSTRAINT jet_pkey PRIMARY KEY (id);
-- The following fkey remains poorly named to avoid regressions related to psql naming
ALTER TABLE jets ADD CONSTRAINT pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);

CREATE TABLE languages (
  id integer NOT NULL,
  language text NOT NULL
);

ALTER TABLE languages ADD CONSTRAINT language_pkey PRIMARY KEY (id);

-- Join table
CREATE TABLE pilot_languages (
  pilot_id integer NOT NULL,
  language_id integer NOT NULL
);

-- Composite primary key
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_pkey PRIMARY KEY (pilot_id, language_id);
-- The following fkey remains poorly named to avoid regressions related to psql naming
ALTER TABLE pilot_languages ADD CONSTRAINT pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
ALTER TABLE pilot_languages ADD CONSTRAINT languages_fkey FOREIGN KEY (language_id) REFERENCES languages(id);
