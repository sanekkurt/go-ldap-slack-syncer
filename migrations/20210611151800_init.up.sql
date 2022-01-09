create table users
(
    slack_id       varchar(100)              not null,
    mail           varchar(100)              not null,
    action         enum ('disable','enable') not null,
    date_of_action date                      not null,
    unique index slack_ldap_relationship (slack_id, mail)
) DEFAULT CHARSET = utf8;