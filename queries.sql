-- find user
select uuid, name from users where email = 'user1@mail.ru';

-- get user lists
select * from users_lists u join lists l on u.list_id = l.id where u.user_uuid = (select uuid from users where email = 'user1@mail.ru' limit 1);

-- get user lists and items
with list as (
    select list_id from users_lists
    where user_uuid in (select uuid from users where email = 'user1@mail.ru')
)
select i.*, l.title from lists l join items i on l.id = i.list_id
where l.id in (select list_id from list);


-- get user lists and items V2
with l as (
    select * from lists l
    join users_lists u on u.list_id = l.id 
    where u.user_uuid = (select uuid from users where email = 'user1@mail.ru' limit 1)
)
select i.*, l.title from l join items i on l.id = i.list_id;