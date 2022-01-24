SET QUOTED_IDENTIFIER ON;

-- Don't forget to maintain order here, foreign keys!
drop table if exists video_tags;
drop table if exists tags;
drop table if exists videos;
drop table if exists sponsors;
drop table if exists users;
drop table if exists type_monsters;
drop view if exists user_videos;

-- Note that if we don't explicitly name foreign keys then MS SQL will
-- generate a name that includes a random set of 8 hex digits at the end,
-- meaning that the driver result varies every time it's run.

create table users (
	id int identity (1,1) primary key not null
);

create table sponsors (
	id int identity (1,1) primary key not null
);

create table videos (
	id int identity (1,1) primary key not null,
	
	user_id int not null,
	sponsor_id int unique,

	constraint FK_videos_users foreign key (user_id) references users (id),
	constraint FK_videos_sponsors foreign key (sponsor_id) references sponsors (id)
);

create table tags (
	id int identity (1,1) primary key not null
);

create table video_tags (
	video_id int not null,
	tag_id int not null,

	primary key (video_id, tag_id),
	constraint FK_video_tags_videos foreign key (video_id) references videos (id),
	constraint FK_video_tags_tags foreign key (tag_id) references tags (id)
);

create table type_monsters (
	id int identity (1,1) primary key not null,

	id_two int not null,
	id_three int,
	bit_zero bit,
	bit_one bit null,
	bit_two bit not null,
	bit_three bit null default 0,
	bit_four bit null default 1,
	bit_five bit not null default 0,
	bit_six bit not null default 1,
	string_zero varchar(1),
	string_one varchar(1) null,
	string_two varchar(1) not null,
	string_three varchar(1) null default 'a',
	string_four varchar(1) not null default 'b',
	string_five varchar(1000),
	string_six varchar(1000) null,
	string_seven varchar(1000) not null,
	string_eight varchar(1000) null default 'abcdefgh',
	string_nine varchar(1000) not null default 'abcdefgh',
	string_ten varchar(1000) null default '',
	string_eleven varchar(1000) not null default '',
	big_int_zero bigint,
	big_int_one bigint null,
	big_int_two bigint not null,
	big_int_three bigint null default 111111,
	big_int_four bigint not null default 222222,
	big_int_five bigint null default 0,
	big_int_six bigint not null default 0,
	int_zero int,
	int_one int null,
	int_two int not null,
	int_three int null default 333333,
	int_four int not null default 444444,
	int_five int null default 0,
	int_six int not null default 0,
	float_zero float,
	float_one float,
	float_two float(24),
	float_three float(24),
	float_four float(24) null,
	float_five float(24) not null,
	float_six float(24) null default 1.1,
	float_seven float(24) not null default 1.1,
	float_eight float(24) null default 0.0,
	float_nine float(24) null default 0.0,
	bytea_zero binary not null,
	bytea_one binary not null,
	bytea_two binary not null,
	bytea_three binary not null default convert(varbinary(max),'a'),
	bytea_four binary not null default convert(varbinary(max),'b'),
	bytea_five binary(100) not null default convert(varbinary(max),'abcdefghabcdefghabcdefgh'),
	bytea_six binary(100) not null default  convert(varbinary(max),'hgfedcbahgfedcbahgfedcba'),
	bytea_seven binary not null default convert(varbinary(max),''),
	bytea_eight binary not null default convert(varbinary(max),''),
	time_zero timestamp not null,
	time_one date,
	time_eleven date null,
	time_twelve date not null,
	time_fifteen date null default '19990108',
	time_sixteen date not null default '1999-01-08',

	bit_null  bit null,
	bit_nnull bit not null,

	tinyint_null   tinyint null,
	tinyint_nnull  tinyint not null,
	smallint_null  smallint null,
	smallint_nnull smallint not null,
	int_null       int null,
	int_nnull      int not null,
	bigint_null    bigint null,
	bigint_nnull   bigint not null,

	float_null       float null,
	float_nnull      float not null,
	doubleprec_null  double precision null,
	doubleprec_nnull double precision not null,
	real_null        real null,
	real_nnull       real not null,

	date_null       date null,
	date_nnull      date not null,
	datetime_null   datetime null,
	datetime_nnull  datetime not null,

	binary_null  binary null,
	binary_nnull binary not null,

	varbinary_null     varbinary null,
	varbinary_nnull    varbinary not null,
	varbinary100_null  varbinary(100) null,
	varbinary100_nnull varbinary(100) not null,
	varbinarymax_null  varbinary(max) null,
	varbinarymax_nnull varbinary(max) not null,

	char_null        char null,
	char_nnull       char not null,
	varchar_null     varchar(max) null,
	varchar_nnull    varchar(max) not null,
	varchar100_null  varchar(100) null,
	varchar100_nnull varchar(100) not null,

	uniqueidentifier_null uniqueidentifier null,
	uniqueidentifier_nnull uniqueidentifier not null,
	datetimeoffset_null datetimeoffset null,
	datetimeoffset_nnull datetimeoffset not null,

    generated_persisted AS bigint_nnull * bigint_null PERSISTED,
    generated_virtual AS smallint_nnull * smallint_null
);

GO

create view user_videos as 
select u.id user_id, v.id video_id, v.sponsor_id sponsor_id
from users u
inner join videos v on v.user_id = u.id;
