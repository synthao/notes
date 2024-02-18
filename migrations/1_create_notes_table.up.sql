create table notes
(
    id         int auto_increment primary key,
    name       varchar(255)                         null,
    text       varchar(1000)                        null,
    created_at datetime default current_timestamp() null,
    constraint notes_pk2
        unique (id)
);