create extension if not exists citext;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists users_on_forum cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists votes cascade;
-- TABLES

create unlogged table users
(
    nickname citext not null unique primary key,
    fullname text   not null,
    about    text,
    email    citext not null unique
);

create unlogged table forums
(
    title   text   not null,
    "user"  citext not null references users (nickname),
    slug    citext not null unique primary key,
    posts   int default 0,
    threads int default 0
);

create unlogged table users_on_forum
(
    slug citext not null references forums (slug),
    nick citext not null references users (nickname),
    unique (nick, slug)
);

create unlogged table threads
(
    id      serial primary key,
    title   text   not null,
    author  citext not null references users (nickname),
    forum   citext not null references forums (slug),
    message text   not null,
    votes   int         default 0,
    slug    citext      default null unique,
    created timestamptz default current_timestamp
);

create unlogged table posts
(
    id       serial8 primary key,
    parent   int8        default 0,
    author   citext not null references users (nickname),
    message  text   not null,
    isEdited bool        default false,
    forum    citext not null references forums (slug),
    thread   int    not null references threads (id),
    created  timestamptz default current_timestamp,
    path     int8[]      default array []::int[]
);

create unlogged table votes
(
    nick      citext not null references users (nickname),
    voice     bool   not null default true,
    thread_id int    not null references threads (id),
    unique (nick, thread_id)
);

-- TRIGGERS

create or replace function create_post_path()
    returns trigger
    language plpgsql as
$body$
begin
    if (new.parent = 0) then
        new.path = new.path || new.id;
    else
        new.path = (select path from posts where id = new.parent) || new.id;
    end if;
    return new;
end
$body$;

create trigger create_post_path
    before insert
    on posts
    for each row
execute procedure create_post_path();

---------------

create or replace function add_user_on_forum_after_thread_create()
    returns trigger
    language plpgsql as
$body$
begin
    insert into users_on_forum(slug, nick)
    values (new.forum, new.author)
    on conflict do nothing;
    return new;
end
$body$;

create trigger update_users_on_forum_after_thread_create
    after insert
    on threads
    for each row
execute procedure add_user_on_forum_after_thread_create();

---------------

create or replace function update_votes()
    returns trigger
    language plpgsql as
$body$
begin
    if (TG_OP = 'INSERT') then
        if new.voice then
            update threads
            set votes = votes + 1
            where id = new.thread_id;
        else
            update threads
            set votes = votes - 1
            where id = new.thread_id;
        end if;
    else
        if new.voice then
            update threads
            set votes = votes + 2
            where id = new.thread_id;
        else
            update threads
            set votes = votes - 2
            where id = new.thread_id;
        end if;
    end if;
    return new;
end
$body$;


create trigger update_votes_trigger
    after update or insert
    on votes
    for each row
execute procedure update_votes();

---------------

create or replace function update_threads_cnt()
    returns trigger
    language plpgsql as
$BODY$
begin
    if TG_OP = 'INSERT' then
        update forums set threads = threads + 1 where slug = new.forum;
    end if;
    return new;
end
$BODY$;

create trigger update_treads_cnt_on_forum
    after insert
    on threads
    for each row
execute procedure update_threads_cnt();
