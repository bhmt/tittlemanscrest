create table if not exists container (
    id bigserial primary key,
    name text unique not null,
    classification text not null default 'goodware'
);

create table if not exists component (
    id bigserial primary key,
    score int not null default 5
);

create table  if not exists association (
    id bigserial primary key,
    container_id  bigint references container(id) on delete cascade,
    component_id  bigint references component(id) on delete cascade,
    constraint uq_mapping unique (container_id, component_id)
);
