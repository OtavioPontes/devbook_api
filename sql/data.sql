insert into users (name, nick, email, password)
values
("User1", "user1", "user1@gmail.com","$2a$10$sqQYpR87SvNfkzqsspVzleSFXs1hAuLrKlXnYFHF41hbJMAo.WNCO"),
("User2", "user2", "user2@gmail.com","$2a$10$sqQYpR87SvNfkzqsspVzleSFXs1hAuLrKlXnYFHF41hbJMAo.WNCO"),
("User3", "user3", "user3@gmail.com","$2a$10$sqQYpR87SvNfkzqsspVzleSFXs1hAuLrKlXnYFHF41hbJMAo.WNCO"),
("User4", "user4", "user4@gmail.com","$2a$10$sqQYpR87SvNfkzqsspVzleSFXs1hAuLrKlXnYFHF41hbJMAo.WNCO")

insert into followers (user_id, follower_id)
values
(1, 2),
(3, 1),
(1, 3)


insert into posts (title, content, author_id)
values
("Post 1 User 2", "Post 1 User 2",2),
("Post 2 User 2", "Post 2 User 2",2),
("Post 1 User 3", "Post 1 User 3",3),
("Post 1 User 4", "Post 1 User 4",4);