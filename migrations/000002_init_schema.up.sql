create index users_email_idx on users using btree (email text_pattern_ops);
