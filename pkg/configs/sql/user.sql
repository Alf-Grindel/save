create database if not exists save;

use save;

create table if not exists users
(
    id bigint auto_increment primary key,
    account varchar(256) not null,
    password varchar(512) not null,
    name varchar(256) null,
    avatar varchar(1024) null,
    profile varchar(512) null,
    tags varchar(1024) null,
    create_time datetime default current_timestamp not null,
    update_time datetime default current_timestamp not null on update current_timestamp,
    is_delete tinyint default 0 not null
) comment 'user' collate = utf8mb4_unicode_ci;

alter table users
add column user_account varchar(256) generated always as (if(is_delete = 0, account, null)) STORED ;

alter table users
add unique key uk_user_account(user_account);