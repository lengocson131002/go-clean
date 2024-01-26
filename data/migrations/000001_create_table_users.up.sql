CREATE TABLE users (
    id varchar(100) NOT NULL,
    name       varchar(100) NOT NULL,
    password   varchar(100) NOT NULL,
    token      varchar(100) NULL,
    created_at bigint       NOT NULL,
    updated_at bigint       NOT NULL,
    PRIMARY KEY (id)
) 