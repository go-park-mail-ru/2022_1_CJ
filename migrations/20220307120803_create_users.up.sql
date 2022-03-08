CREATE TABLE users
(
    user_id            bigserial not null primary key,

    first_name         varchar   not null,
    last_name          varchar   not null,
    path_to_img        varchar   not null,

    email              varchar   not null unique,
    phone              varchar   not null unique,

    password           varchar   not null,
    encrypted_password varchar   not null
);

CREATE TABLE post
(
    post_id   bigserial not null primary key,
    post_text varchar,
    img_path  varchar,
    like_num  bigserial,
    repost    bigserial
);

CREATE TABLE relation_post_user
(
    user_id int foreign key references users(user_id),
    post_id int foreign key references post_id(post_id),
    primary key (user_id, post_id)
);


