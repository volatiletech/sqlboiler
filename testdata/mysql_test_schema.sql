CREATE TABLE event_one (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  day    enum('monday','tuesday','wednesday')
);

CREATE TABLE event_two (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  face    enum('happy','sad','bitter')
);

CREATE TABLE event_three (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  face   enum('happy','sad','bitter'),
  mood   enum('happy','sad','bitter'),
  day    enum('monday','tuesday','wednesday')
);

CREATE TABLE magic (
  id         int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  id_two     int NOT NULL,
  id_three   int,
  bool_zero  bool,
  bool_one   bool NULL,
  bool_two   bool NOT NULL,
  bool_three bool NULL DEFAULT FALSE,
  bool_four  bool NULL DEFAULT TRUE,
  bool_five  bool NOT NULL DEFAULT FALSE,
  bool_six   bool NOT NULL DEFAULT TRUE,
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
  float_zero  float,
  float_one   float,
  float_two   float(2,1),
  float_three float(2,1),
  float_four  float(2,1) NULL,
  float_five  float(2,1) NOT NULL,
  float_six   float(2,1) NULL DEFAULT 1.1,
  float_seven float(2,1) NOT NULL DEFAULT 1.1,
  float_eight float(2,1) NULL DEFAULT 0.0,
  float_nine  float(2,1) NULL DEFAULT 0.0,
  bytea_zero  binary,
  bytea_one   binary NULL,
  bytea_two   binary NOT NULL,
  bytea_three binary NOT NULL DEFAULT 'a',
  bytea_four  binary NULL DEFAULT 'b',
  bytea_five  binary(100) NOT NULL DEFAULT 'abcdefghabcdefghabcdefgh',
  bytea_six   binary(100) NULL DEFAULT 'hgfedcbahgfedcbahgfedcba',
  bytea_seven binary NOT NULL DEFAULT '',
  bytea_eight binary NOT NULL DEFAULT '',
  time_zero   timestamp,
  time_one    date,
  time_two    timestamp NULL DEFAULT NULL,
  time_three  timestamp NULL,
  time_five   timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  time_nine   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  time_eleven date NULL,
  time_twelve date NOT NULL,
  time_fifteen date NULL DEFAULT '19990108',
  time_sixteen date NOT NULL DEFAULT '1999-01-08'
);

CREATE TABLE magicest (
  id   int primary key not null auto_increment,
  aa   json NULL,
  bb   json NOT NULL,
  kk   double precision NULL,
  ll   double precision NOT NULL,
  mm   tinyint NULL,
  nn   tinyint NOT NULL,
  oo   tinyint(1) NULL,
  pp   tinyint(1) NOT NULL,
  qq   smallint NULL,
  rr   smallint NOT NULL,
  ss   mediumint NULL,
  tt   mediumint NOT NULL,
  uu   bigint NULL,
  vv   bigint NOT NULL,
  ww   float NULL,
  xx   float NOT NULL,
  yy   double NULL,
  zz   double NOT NULL,
  aaa  double precision NULL,
  bbb  double precision NOT NULL,
  ccc  real NULL,
  ddd  real NOT NULL,
  eee  boolean NULL,
  fff  boolean NOT NULL,
  ggg  date NULL,
  hhh  date NOT NULL,
  iii  datetime NULL,
  jjj  datetime NOT NULL,
  kkk  timestamp NULL,
  lll  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  mmm  binary NULL,
  nnn  binary NOT NULL,
  ooo  varbinary(100) NULL,
  ppp  varbinary(100) NOT NULL,
  qqq  tinyblob NULL,
  rrr  tinyblob NOT NULL,
  sss  blob NULL,
  ttt  blob NOT NULL,
  uuu  mediumblob NULL,
  vvv  mediumblob NOT NULL,
  www  longblob NULL,
  xxx  longblob NOT NULL,
  yyy  varchar(100) NULL,
  zzz  varchar(100) NOT NULL,
  aaaa char NULL,
  bbbb char NOT NULL,
  cccc text NULL,
  dddd text NOT NULL,
  eeee tinyint(2) NULL,
  ffff tinyint(2) NOT NULL
);

create table owner (
  id    int primary key not null auto_increment,
  name  varchar(255) not null
);

create table cats (
  id    int primary key not null auto_increment,
  name     varchar(255) not null,
  owner_id int references owner (id)
);

create table toys (
  id    int primary key not null auto_increment,
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
  dragon_id varchar(100),
  toy_id    varchar(100),
  primary key (dragon_id, toy_id)
);

create table spider_toys (
  spider_id varchar(100) primary key,
  name      varchar(100)
);

create table pals (
  pal varchar(100) primary key,
  name varchar(100)
);

create table friend (
  friend varchar(100) primary key,
  name varchar(100)
);

create table bro (
  bros varchar(100) primary key,
  name varchar(100)
);

create table enemies (
  enemies varchar(100) primary key,
  name varchar(100)
);

create table chocolate (
  dog varchar(100) primary key
);

create table waffles (
  cat varchar(100) primary key
);

create table tigers (
  id    binary primary key,
  name  binary null
);

create table elephants (
  id        binary primary key,
  name      binary not null,
  tiger_id  binary null unique,
  foreign key (tiger_id) references tigers (id)
);

create table wolves (
  id        binary primary key,
  name      binary not null,
  tiger_id  binary not null unique,
  foreign key (tiger_id) references tigers (id)
);

create table ants (
  id        binary primary key,
  name      binary not null,
  tiger_id  binary not null,
  foreign key (tiger_id) references tigers (id)
);

create table worms (
  id        binary primary key,
  name      binary not null,
  tiger_id  binary null,
  foreign key (tiger_id) references tigers (id)
);

create table byte_pilots (
  id   binary primary key not null,
  name varchar(255)
);

create table byte_airports (
  id   binary primary key not null,
  name varchar(255)
);

create table byte_languages (
  id   binary primary key not null,
  name varchar(255)
);

create table byte_jets (
  id              binary primary key not null,
  name            varchar(255),
  byte_pilot_id   binary unique,
  byte_airport_id binary,

  foreign key (byte_pilot_id) references byte_pilots (id),
  foreign key (byte_airport_id) references byte_airports (id)
);

create table byte_pilot_languages (
  byte_pilot_id    binary not null,
  byte_language_id binary not null,

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
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
ALTER TABLE pilot_languages ADD CONSTRAINT languages_fkey FOREIGN KEY (language_id) REFERENCES languages(id);
