insert into users (name, email, pass, salt) values
  ('User1', 'user1@mail.ru', 'd36f9d30acb4e2857a2818aa8420f7b7', '111'),
  ('User2', 'user2@mail.ru', 'd36f9d30acb4e2857a2818aa8420f7b7', '111'),
  ('Admin', 'admin@mail.ru', '66e1a360ee8070ba822aca90526dec47', '222');

insert into lists (user_uuid, title, description) values
  ((select uuid from users where email = 'user1@mail.ru'), 'Купить в магазине', ''),
  ((select uuid from users where email = 'user1@mail.ru'), 'Список задач', 'Что нужно сделать в ближайшее время');

insert into items (list_id, title, description, due_date) values
  (1, 'Молоко', '', now() + interval '1 day'),
  (1, 'Хлеб', '', now() + interval '1 day'),
  (1, 'Курица', '', now() + interval '1 day'),
  (1, 'Чай', '', now() + interval '1 day'),
  (2, 'сходить в магазин', '', now() + interval '2 day'),
  (2, 'получить посылку', '', now() + interval '2 day'),
  (2, 'сделать домашку', '', now() + interval '2 day'),
  (2, 'почить про CAP теорему', '', now() + interval '2 day'),
  (2, 'полить цветы', '', now() + interval '2 day');