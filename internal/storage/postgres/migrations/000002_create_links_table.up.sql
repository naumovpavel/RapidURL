CREATE TABLE IF NOT EXISTS links(
    id serial PRIMARY KEY,
    alias varchar(50) not null,
    url text not null
);
create index if not exists alias_id on links(alias);