CREATE TABLE magic
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  id_two int NOT NULL,
  id_three int,
  bit_zero bit,
  bit_one bit NULL,
  bit_two bit NOT NULL,
  bit_three bit NULL DEFAULT 0,
  bit_four bit NULL DEFAULT 1,
  bit_five bit NOT NULL DEFAULT 0,
  bit_six bit NOT NULL DEFAULT 1,
  string_zero VARCHAR(1),
  string_one VARCHAR(1) NULL,
  string_two VARCHAR(1) NOT NULL,
  string_three VARCHAR(1) NULL DEFAULT 'a',
  string_four VARCHAR(1) NOT NULL DEFAULT 'b',
  string_five VARCHAR(1000),
  string_six VARCHAR(1000) NULL,
  string_seven VARCHAR(1000) NOT NULL,
  string_eight VARCHAR(1000) NULL DEFAULT 'abcdefgh',
  string_nine VARCHAR(1000) NOT NULL DEFAULT 'abcdefgh',
  string_ten VARCHAR(1000) NULL DEFAULT '',
  string_eleven VARCHAR(1000) NOT NULL DEFAULT '',
  big_int_zero bigint,
  big_int_one bigint NULL,
  big_int_two bigint NOT NULL,
  big_int_three bigint NULL DEFAULT 111111,
  big_int_four bigint NOT NULL DEFAULT 222222,
  big_int_five bigint NULL DEFAULT 0,
  big_int_six bigint NOT NULL DEFAULT 0,
  int_zero int,
  int_one int NULL,
  int_two int NOT NULL,
  int_three int NULL DEFAULT 333333,
  int_four int NOT NULL DEFAULT 444444,
  int_five int NULL DEFAULT 0,
  int_six int NOT NULL DEFAULT 0,
  float_zero float,
  float_one float,
  float_two float(24),
  float_three float(24),
  float_four float(24) NULL,
  float_five float(24) NOT NULL,
  float_six float(24) NULL DEFAULT 1.1,
  float_seven float(24) NOT NULL DEFAULT 1.1,
  float_eight float(24) NULL DEFAULT 0.0,
  float_nine float(24) NULL DEFAULT 0.0,
  bytea_zero binary NOT NULL,
  bytea_one binary NOT NULL,
  bytea_two binary NOT NULL,
  bytea_three binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),'a'),
  bytea_four binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),'b'),
  bytea_five binary(100) NOT NULL DEFAULT CONVERT(VARBINARY(MAX),'abcdefghabcdefghabcdefgh'),
  bytea_six binary(100) NOT NULL DEFAULT  CONVERT(VARBINARY(MAX),'hgfedcbahgfedcbahgfedcba'),
  bytea_seven binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),''),
  bytea_eight binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),''),
  time_zero timestamp NOT NULL,
  time_one date,
  time_eleven date NULL,
  time_twelve date NOT NULL,
  time_fifteen date NULL DEFAULT '19990108',
  time_sixteen date NOT NULL DEFAULT '1999-01-08'
);
GO

CREATE TABLE magicest
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  kk float NULL,
  ll float NOT NULL,
  mm tinyint NULL,
  nn tinyint NOT NULL,
  oo bit NULL,
  pp bit NOT NULL,
  qq smallint NULL,
  rr smallint NOT NULL,
  ss int NULL,
  tt int NOT NULL,
  uu bigint NULL,
  vv bigint NOT NULL,
  ww float NULL,
  xx float NOT NULL,
  yy float NULL,
  zz float NOT NULL,
  aaa double precision NULL,
  bbb double precision NOT NULL,
  ccc real NULL,
  ddd real NOT NULL,
  ggg date NULL,
  hhh date NOT NULL,
  iii datetime NULL,
  jjj datetime NOT NULL,
  kkk timestamp NOT NULL,
  mmm binary NOT NULL,
  nnn binary NOT NULL,
  ooo varbinary(100) NOT NULL,
  ppp varbinary(100) NOT NULL,
  qqq varbinary NOT NULL,
  rrr varbinary NOT NULL,
  www varbinary(max) NOT NULL,
  xxx varbinary(max) NOT NULL,
  yyy varchar(100) NULL,
  zzz varchar(100) NOT NULL,
  aaaa char NULL,
  bbbb char NOT NULL,
  cccc VARCHAR(MAX) NULL,
  dddd VARCHAR(MAX) NOT NULL,
  eeee tinyint NULL,
  ffff tinyint NOT NULL
);
GO

create table owner
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  name varchar(255) not null
);
GO

create table cats
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  name varchar(255) not null,
  owner_id int
);
GO

ALTER TABLE cats ADD CONSTRAINT cats_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES owner(id);
GO

create table toys
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  name varchar(255) not null
);
GO

create table cat_toys
(
  cat_id int not null references cats (id),
  toy_id int not null references toys (id),
  primary key (cat_id, toy_id)
);
GO

create table dog_toys
(
  dog_id int not null,
  toy_id int not null,
  primary key (dog_id, toy_id)
);
GO

create table dragon_toys
(
  dragon_id varchar(100),
  toy_id varchar(100),
  primary key (dragon_id, toy_id)
);
GO

create table spider_toys
(
  spider_id varchar(100) primary key,
  name varchar(100)
);
GO

create table pals
(
  pal varchar(100) primary key,
  name varchar(100)
);
GO

create table friend
(
  friend varchar(100) primary key,
  name varchar(100)
);
GO

create table bro
(
  bros varchar(100) primary key,
  name varchar(100)
);
GO

create table enemies
(
  enemies varchar(100) primary key,
  name varchar(100)
);
GO

create table chocolate
(
  dog varchar(100) primary key
);
GO

create table waffles
(
  cat varchar(100) primary key
);
GO

create table tigers
(
  id binary primary key,
  name binary NOT NULL
);
GO

create table elephants
(
  id binary primary key,
  name binary not null,
  tiger_id binary NOT NULL unique
);
GO

ALTER TABLE elephants ADD CONSTRAINT elephants_tiger_id_fkey FOREIGN KEY (tiger_id) REFERENCES tigers(id);
GO

create table wolves
(
  id binary primary key,
  name binary not null,
  tiger_id binary not null unique
);
GO

ALTER TABLE wolves ADD CONSTRAINT wolves_tiger_id_fkey FOREIGN KEY (tiger_id) REFERENCES tigers(id);
GO

create table ants
(
  id binary primary key,
  name binary not null,
  tiger_id binary not null
);
GO

ALTER TABLE ants ADD CONSTRAINT ants_tiger_id_fkey FOREIGN KEY (tiger_id) REFERENCES tigers(id);
GO

create table worms
(
  id binary primary key,
  name binary not null,
  tiger_id binary NOT NULL
);
GO

ALTER TABLE worms ADD CONSTRAINT worms_tiger_id_fkey FOREIGN KEY (tiger_id) REFERENCES tigers(id);
GO

create table byte_pilots
(
  id binary primary key not null,
  name varchar(255)
);
GO

create table byte_airports
(
  id binary primary key not null,
  name varchar(255)
);
GO

create table byte_languages
(
  id binary primary key not null,
  name varchar(255)
);
GO

create table byte_jets
(
  id binary primary key not null,
  name varchar(255),
  byte_pilot_id binary unique NOT NULL,
  byte_airport_id binary NOT NULL
);
GO

ALTER TABLE byte_jets ADD CONSTRAINT byte_jets_byte_pilot_id_fkey FOREIGN KEY (byte_pilot_id) REFERENCES byte_pilots(id);
GO
ALTER TABLE byte_jets ADD CONSTRAINT byte_jets_byte_airport_id_fkey FOREIGN KEY (byte_airport_id) REFERENCES byte_airports(id);
GO

create table byte_pilot_languages
(
  byte_pilot_id binary not null,
  byte_language_id binary not null
);
GO

ALTER TABLE byte_pilot_languages ADD CONSTRAINT byte_pilot_languages_pkey PRIMARY KEY (byte_pilot_id,byte_language_id);
GO

ALTER TABLE byte_pilot_languages ADD CONSTRAINT byte_pilot_languages_byte_pilot_id_fkey FOREIGN KEY (byte_pilot_id) REFERENCES byte_pilots(id);
GO
ALTER TABLE byte_pilot_languages ADD CONSTRAINT byte_pilot_languages_byte_language_id_fkey FOREIGN KEY (byte_language_id) REFERENCES byte_languages(id);
GO

create table cars
(
  id integer not null,
  name VARCHAR(MAX),
  primary key (id)
);
GO

create table car_cars
(
  car_id integer not null,
  awesome_car_id integer not null,
  relation VARCHAR(MAX) not null,
  primary key (car_id, awesome_car_id)
);
GO

ALTER TABLE car_cars ADD CONSTRAINT car_id_fkey FOREIGN KEY (car_id) REFERENCES cars(id);
GO
ALTER TABLE car_cars ADD CONSTRAINT awesome_car_id_fkey FOREIGN KEY (awesome_car_id) REFERENCES cars(id);
GO

create table trucks
(
  id integer not null,
  parent_id integer,
  name VARCHAR(MAX),
  primary key (id)
);
GO

ALTER TABLE trucks ADD CONSTRAINT parent_id_fkey FOREIGN KEY (parent_id) REFERENCES trucks(id);
GO

CREATE TABLE race
(
  id integer PRIMARY KEY NOT NULL,
  race_date datetime,
  track VARCHAR(MAX)
);
GO

CREATE TABLE race_results
(
  id integer PRIMARY KEY NOT NULL,
  race_id integer,
  name VARCHAR(MAX)
);
GO

ALTER TABLE race_results ADD CONSTRAINT race_id_fkey FOREIGN KEY (race_id) REFERENCES race(id);
GO

CREATE TABLE race_result_scratchings
(
  id integer PRIMARY KEY NOT NULL,
  results_id integer NOT NULL,
  name VARCHAR(MAX) NOT NULL
);
GO

ALTER TABLE race_result_scratchings ADD CONSTRAINT results_id_fkey FOREIGN KEY (results_id) REFERENCES race_results(id);
GO

CREATE TABLE pilots
(
  id integer NOT NULL,
  name VARCHAR(MAX) NOT NULL
);
GO

ALTER TABLE pilots ADD CONSTRAINT pilot_pkey PRIMARY KEY (id);
GO

CREATE TABLE jets
(
  id integer NOT NULL,
  pilot_id integer NOT NULL,
  age integer NOT NULL,
  name VARCHAR(MAX) NOT NULL,
  color VARCHAR(MAX) NOT NULL
);
GO

ALTER TABLE jets ADD CONSTRAINT jet_pkey PRIMARY KEY (id);
GO
ALTER TABLE jets ADD CONSTRAINT pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
GO

CREATE TABLE languages
(
  id integer NOT NULL,
  language VARCHAR(MAX) NOT NULL
);
GO

ALTER TABLE languages ADD CONSTRAINT language_pkey PRIMARY KEY (id);
GO

-- Join table
CREATE TABLE pilot_languages
(
  pilot_id integer NOT NULL,
  language_id integer NOT NULL,
  uniqueid uniqueidentifier NOT NULL,
);
GO

-- Composite primary key
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_pkey PRIMARY KEY (pilot_id, language_id);
GO
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
GO
ALTER TABLE pilot_languages ADD CONSTRAINT languages_fkey FOREIGN KEY (language_id) REFERENCES languages(id);
GO

CREATE TABLE powers_of_two
(
  vid int NOT NULL IDENTITY(1,1),
  name varchar(255) NOT NULL DEFAULT '',
  machine_name varchar(255) NOT NULL,
  description VARCHAR(MAX),
  hierarchy tinyint NOT NULL DEFAULT '0',
  module varchar(255) NOT NULL DEFAULT '',
  weight int NOT NULL DEFAULT '0',
  PRIMARY KEY (vid),
  CONSTRAINT machine_name UNIQUE(machine_name)
);
GO
