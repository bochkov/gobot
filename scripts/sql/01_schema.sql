drop table if exists props;
create table props
(
    id    int primary key generated by default as identity (sequence name props_seq start with 1),
    key   varchar(20)  not null unique,
    value varchar(255) not null
);

drop table if exists code;
drop table if exists region;

create table region
(
    id   int primary key generated by default as identity (sequence name region_seq start with 1),
    name varchar(255)
);

create table code
(
    id        int primary key generated by default as identity (sequence name code_seq start with 1),
    val       varchar(3),
    region_id int,
    foreign key (region_id) references region (id)
);

drop table if exists currency_item;
drop table if exists currency_rate cascade;
drop table if exists currency_rate_record cascade;

create table currency_item
(
    id            varchar(10) primary key,
    name          varchar(100),
    eng_name      varchar(100),
    nominal       int,
    parent_code   varchar(10),
    iso_num_code  int,
    iso_char_code varchar(3)
);

create table currency_rate
(
    id         bigint primary key generated by default as identity (sequence name currency_rate_seq start with 1),
    date       date unique,
    fetch_time timestamp not null default now(),
    name       varchar(100)
);

create table currency_rate_record
(
    id         bigint primary key generated by default as identity (sequence name currency_rate_record_seq start with 1),
    curr_id    varchar(10),
    rate_id    bigint,
    num_code   varchar(3),
    char_code  varchar(3),
    nominal    int,
    name       varchar(100),
    rate_value numeric(10, 4),
    foreign key (rate_id) references currency_rate (id)
);