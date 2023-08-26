CREATE TABLE IF NOT EXISTS links(
    id serial PRIMARY KEY,
    alias varchar(50) not null,
    url text not null,
    user_id int not null references users(id)
);
create index if not exists alias_id on links(alias);