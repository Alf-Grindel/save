use save;

create table if not exists teams
(
    id          bigint auto_increment primary key comment 'id',
    team_name   varchar(256)                       not null comment 'team name',
    description varchar(1024)                      not null comment 'team description',
    max_num     int      default 1                 not null comment 'max allow join team member number',
    expire_time datetime                           null comment 'team expire time',
    user_id     bigint                             not null comment 'leader id',
    status      enum('public','private','encrypted')      default 'public'                 not null comment ' public, private, encrypted',
    password    varchar(512)                       null comment 'password if status is 2',
    create_time datetime default current_timestamp not null comment 'create team time',
    update_time datetime default current_timestamp not null on update current_timestamp comment 'update team info time',
    is_delete   tinyint  default 0                 not null comment '0 - exist 1 - delete'
) comment 'team';
