create table users (
	id int primary key not null
);

create table sponsors (
	id int primary key not null
);

create table videos (
	id int primary key not null,

	user_id int not null,
	sponsor_id int unique,

	foreign key (user_id) references users (id),
	foreign key (sponsor_id) references sponsors (id)
);

create table tags (
	id int primary key not null
);

create table video_tags (
	video_id int not null,
	tag_id   int not null,

	primary key (video_id, tag_id),
	foreign key (video_id) references videos (id),
	foreign key (tag_id) references tags (id)
);

create table type_monsters (
	id int primary key not null,

	id_two     int not null,
	id_three   int,
	bool_zero  bool,
	bool_one   bool null,
	bool_two   bool not null,
	bool_three bool null default false,
	bool_four  bool null default true,
	bool_five  bool not null default false,
	bool_six   bool not null default true,

	string_zero   varchar(1),
	string_one    varchar(1) null,
	string_two    varchar(1) not null,
	string_three  varchar(1) null default 'a',
	string_four   varchar(1) not null default 'b',
	string_five   varchar(1000),
	string_six    varchar(1000) null,
	string_seven  varchar(1000) not null,
	string_eight  varchar(1000) null default 'abcdefgh',
	string_nine   varchar(1000) not null default 'abcdefgh',
	string_ten    varchar(1000) null default '',
	string_eleven varchar(1000) not null default '',

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
	float_four  float(2,1) null,
	float_five  float(2,1) not null,
	float_six   float(2,1) null default 1.1,
	float_seven float(2,1) not null default 1.1,
	float_eight float(2,1) null default 0.0,
	float_nine  float(2,1) null default 0.0,
	bytea_zero  binary,
	bytea_one   binary null,
	bytea_two   binary not null,
	bytea_three binary not null default 'a',
	bytea_four  binary null default 'b',
	bytea_five  binary(100) not null default 'abcdefghabcdefghabcdefgh',
	bytea_six   binary(100) null default 'hgfedcbahgfedcbahgfedcba',
	bytea_seven binary not null default '',
	bytea_eight binary not null default '',
	time_zero   timestamp,
	time_one    date,
	time_two    timestamp null default null,
	time_three  timestamp null,
	time_five   timestamp null default current_timestamp,
	time_nine   timestamp not null default current_timestamp,
	time_eleven date null,
	time_twelve date not null,
	time_fifteen date null default '19990108',
	time_sixteen date not null default '1999-01-08',

	json_null  json null,
	json_nnull json not null,

	tinyint_null    tinyint null,
	tinyint_nnull   tinyint not null,
	tinyint1_null   tinyint(1) null,
	tinyint1_nnull  tinyint(1) not null,
	tinyint2_null   tinyint(2) null,
	tinyint2_nnull  tinyint(2) not null,
	smallint_null   smallint null,
	smallint_nnull  smallint not null,
	mediumint_null  mediumint null,
	mediumint_nnull mediumint not null,
	bigint_null     bigint null,
	bigint_nnull    bigint not null,

	float_null       float null,
	float_nnull      float not null,
	double_null      double null,
	double_nnull     double not null,
	doubleprec_null  double precision null,
	doubleprec_nnull double precision not null,

	real_null  real null,
	real_nnull real not null,

	boolean_null  boolean null,
	boolean_nnull boolean not null,

	date_null  date null,
	date_nnull date not null,

	datetime_null  datetime null,
	datetime_nnull datetime not null,

	timestamp_null  timestamp null,
	timestamp_nnull timestamp not null default current_timestamp,

	binary_null      binary null,
	binary_nnull     binary not null,
	varbinary_null   varbinary(100) null,
	varbinary_nnull  varbinary(100) not null,
	tinyblob_null    tinyblob null,
	tinyblob_nnull   tinyblob not null,
	blob_null        blob null,
	blob_nnull       blob not null,
	mediumblob_null  mediumblob null,
	mediumblob_nnull mediumblob not null,
	longblob_null    longblob null,
	longblob_nnull   longblob not null,

	varchar_null  varchar(100) null,
	varchar_nnull varchar(100) not null,
	char_null     char null,
	char_nnull    char not null,
	text_null     text null,
	text_nnull    text not null
);

-- all table defintions will not cause sqlite autoincrement primary key without rowid tables to be generated
create table autoinctest (
	id INTEGER PRIMARY KEY
);

-- additional fields should not be marked as auto generated, when the AUTOINCREMENT keyword is present
create table autoinckeywordtest (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	b INTEGER
);

-- An INTEGER primary key column is an alias for the table rowid only if it is
-- the only primary key column for the table. i.e. composite primary keys do
-- not exhibit the rowid alias behaviour.
create table compositeprimarykeytest (
	a INTEGER,
	b INTEGER,
	PRIMARY KEY (a, b)
);

create view user_videos as 
select u.id user_id, v.id video_id, v.sponsor_id sponsor_id
from users u
inner join videos v on v.user_id = u.id;

CREATE TABLE has_generated_columns (
   a INTEGER PRIMARY KEY,
   b INT,
   c TEXT,
   d INT GENERATED ALWAYS AS (a*abs(b)) VIRTUAL,
   e TEXT GENERATED ALWAYS AS (substr(c,b,b+1)) STORED
);
