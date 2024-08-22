create table if not exists auth_whitelist(
  token varchar(255),

  constraint auth_whitelist_token primary key (token)
);
