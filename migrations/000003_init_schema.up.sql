alter table lists add column user_uuid uuid;

update lists set user_uuid = ul.user_uuid from users_lists ul where ul.list_id = id;

drop table users_lists;

alter table lists add constraint lists_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE;
