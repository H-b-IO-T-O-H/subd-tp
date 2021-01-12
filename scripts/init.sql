create extension if not exists citext;

drop trigger if exists create_post_path on posts;
drop trigger if exists update_users_on_forum_after_thread_create on threads;
drop trigger if exists update_votes on votes;
drop trigger if exists update_treads_cnt_on_forum on threads;
drop function create_post_path;
drop function add_user_on_forum_after_thread_create;
drop function update_votes;
drop function update_threads_cnt;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists users_on_forum cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists votes cascade;


-- TABLES

create unlogged table users
(
    nickname citext collate "C" not null unique primary key,
    fullname text               not null,
    about    text,
    email    citext collate "C" not null unique
);



create unlogged table forums
(
    title   text               not null,
    "user"  citext collate "C" not null references users (nickname),
    slug    citext collate "C" not null unique primary key,
    posts   int default 0,
    threads int default 0
);

create unlogged table users_on_forum
(
    slug citext collate "C" not null references forums (slug),
    nick citext collate "C" not null references users (nickname),
    unique (nick, slug)
);

create unlogged table threads
(
    id      serial primary key,
    title   text               not null,
    author  citext collate "C" not null references users (nickname),
    forum   citext collate "C" not null references forums (slug),
    message text               not null,
    votes   int                default 0,
    slug    citext collate "C" default null unique,
    created timestamptz        not null
);

create unlogged table posts
(
    id       serial primary key,
    parent   int   default 0,
    author   citext collate "C" not null references users (nickname),
    message  text               not null,
    isEdited bool  default false,
    forum    citext collate "C" not null references forums (slug),
    thread   int                not null references threads (id),
    created  timestamp          not null,
    path     int[] default array []::int[]
);

create unlogged table votes
(
    nick      citext collate "C" not null references users (nickname),
    voice     bool               not null default true,
    thread_id int                not null references threads (id),
    unique (nick, thread_id)
);

-- INDEXES


create index if not exists hash_idx_user_nickname ON users using hash (nickname);
create index if not exists hash_idx_user_email ON users using hash (email);

create index if not exists hash_idx_forum_slug ON forums using hash (slug);
create unique index if not exists idx_users_o_forum on users_on_forum (slug, nick);
cluster users_on_forum using idx_users_o_forum;

create index if not exists hash_idx_thread_slug ON threads using hash (slug);
create index if not exists hash_idx_thread_forum ON threads using hash (forum);
create unique index if not exists idx_thread_slug on threads (slug) where slug is not null;
create index if not exists idx_thread_date ON threads (created);
create index if not exists idx_thread_forum ON threads using hash (forum);
create index if not exists idx_thread_forum_created ON threads (forum, created);

create unique index if not exists idx_vote on votes (nick, thread_id);

create index if not exists idx_post_child_id on posts (id, (path[1]));
create index if not exists idx_post_thread_id_child_parent on posts (thread, id, (path[1]), parent);
create index if not exists idx_post_threads_child_id on posts (thread, path, id);
create index if not exists idx_post_child_path on posts ((path[1]));
create index if not exists idx_post_thread_id on posts (thread, id);
create index if not exists idx_post_thread ON posts (thread);




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
    if TG_OP = 'INSERT' then
        if new.voice then
            update threads set votes = votes + 1 where id = new.thread_id;
        else
            update threads set votes = votes - 1 where id = new.thread_id;
        end if;
    elseif TG_OP = 'UPDATE' then
        if new.voice then
            update threads set votes = votes + 2 where id = new.thread_id;
        else
            update threads set votes = votes - 2 where id = new.thread_id;
        end if;
    end if;
    return new;
end
$body$;


create trigger update_votes
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
