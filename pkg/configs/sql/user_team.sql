use save;

create table if not exists user_team
(
    id bigint auto_increment primary key,
    user_id bigint,
    team_id bigint,
    join_time datetime not null comment 'join team time',
    create_time datetime default current_timestamp not null comment 'create team time',
    update_time datetime default current_timestamp not null on update current_timestamp comment 'update team info time',
    is_delete   tinyint  default 0                 not null comment '0 - exist 1 - delete'
) comment 'user - team relationship';