CREATE TABLE magic (
  id        serial PRIMARY KEY NOT NULL,
  id_two    serial NOT NULL,
  id_three  serial,

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

  time_zero       timestamp,
  time_one        date,
  time_two        timestamp NULL DEFAULT NULL,
  time_three      timestamp NULL,
  time_four       timestamp NOT NULL,
  time_five       timestamp NULL DEFAULT '1999-01-08 04:05:06.789',
  time_six        timestamp NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
  time_seven      timestamp NULL DEFAULT 'January 8 04:05:06 1999 PST',
  time_eight      timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789',
  time_nine       timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
  time_ten        timestamp NOT NULL DEFAULT 'January 8 04:05:06 1999 PST',
  time_eleven     date NULL,
  time_twelve     date NOT NULL,
  time_thirteen   date NULL DEFAULT '1999-01-08',
  time_fourteen   date NULL DEFAULT 'January 8, 1999',
  time_fifteen    date NULL DEFAULT '19990108',
  time_sixteen    date NOT NULL DEFAULT '1999-01-08',
  time_seventeen  date NOT NULL DEFAULT 'January 8, 1999',
  time_eighteen   date NOT NULL DEFAULT '19990108',

  uuid_zero   uuid,
  uuid_one    uuid NULL,
  uuid_two    uuid NULL DEFAULT NULL,
  uuid_three  uuid NOT NULL,
  uuid_four   uuid NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
  uuid_five   uuid NOT NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',

  strange_one   integer DEFAULT '5'::integer,
  strange_two   varchar(1000) DEFAULT 5::varchar,
  strange_three timestamp without time zone default (now() at time zone 'utc'),
  strange_four  timestamp with time zone default (now() at time zone 'utc'),
  strange_five  interval NOT NULL DEFAULT '21 days',
  strange_six   interval NULL DEFAULT '23 hours'
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
  toy_id uuid,
  primary key (dragon_id, toy_id)
);

create table spider_toys (
  spider_id uuid,
  name character varying,
  primary key (spider_id)
);

/*
  Test:
  * Variations of capitalization
  * Single value columns
  * Primary key as only value
*/
create table pals (
  pal character varying,
  primary key (pal)
);

create table friend (
  friend character varying,
  primary key (friend)
);

create table bro (
  bros character varying,
  primary key (bros)
);

create table enemies (
  enemies character varying,
  primary key (enemies)
);
