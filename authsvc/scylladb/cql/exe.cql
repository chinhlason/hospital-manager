create table if not exists users(
    id uuid,
    username text,
    password text,
    fullname text,
    email text,
    phone text,
    role text,
    create_at timestamp,
    update_at timestamp,
    primary key ( (id) )
);
