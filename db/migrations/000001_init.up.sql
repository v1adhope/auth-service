create table if not exists auth_whitelist(
  id uuid,
  token varchar(255),

  constraint auth_whitelist_id primary key (id)
);
