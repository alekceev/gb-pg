-- find user
select uuid, name from users where email = 'user1@mail.ru';

-- get user lists
select * from lists l join users u on u.uuid = l.user_uuid where u.uuid = '2d9e428c-ca39-4b55-a2cf-98e068aebbb9';

-- get user lists and items
with list as (
    select list_id from users_lists
    where user_uuid in (select uuid from users where email = 'user1@mail.ru')
)
select i.*, l.title from lists l join items i on l.id = i.list_id
where l.user_uuid = '2d9e428c-ca39-4b55-a2cf-98e068aebbb9';
