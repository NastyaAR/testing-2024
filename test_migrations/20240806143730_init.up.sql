create type user_role as enum ('moderator', 'client');
create type flat_status as enum ('created', 'approved', 'declined', 'on moderation');
create type flat_update_msg_status as enum ('send', 'no send');

create table users (
                       user_id uuid primary key ,
                       mail varchar(50) not null,
                       password text not null,
                       role user_role not null
);

create table houses (
                        house_id serial primary key,
                        address text not null,
                        construct_year int,
                        developer text,
                        create_house_date timestamp(0) without time zone,
                        update_flat_date timestamp(0) without time zone
);

create table flats (
                       flat_id int not null,
                       house_id int references houses(house_id),
                       user_id uuid references users(user_id),
                       price int not null,
                       rooms int not null,
                       status flat_status not null,
                       moderator_id uuid references users(user_id),
                       primary key (flat_id, house_id)
);

create index houses_id_on_flats
    on flats (house_id);

create table subscribers (
                             user_id uuid references users(user_id),
                             house_id int references houses(house_id)
);

create table new_flats_outbox (
                                  id serial primary key,
                                  flat_id int,
                                  house_id int,
                                  mail text not null ,
                                  status flat_update_msg_status not null,
                                  foreign key (flat_id, house_id) references flats(flat_id, house_id)
);

create or replace function insert_flat_to_outbox()
    returns trigger as $$
declare
    subscriber_mail text;
    subscriber_mails text[];
begin
    select array_agg(u.mail)
    into subscriber_mails
    from subscribers s
             join users u on u.user_id = s.user_id
    where s.house_id = new.house_id;

    if subscriber_mails is null then
        return new;
    end if;

    foreach subscriber_mail in array subscriber_mails
        loop
            insert into new_flats_outbox(flat_id, house_id, mail, status)
            values (new.flat_id, new.house_id, subscriber_mail, 'no send');
        end loop;

    return new;
end;
$$ language plpgsql;

CREATE TRIGGER flat_create_trigger
    AFTER INSERT ON flats
    FOR EACH ROW
EXECUTE FUNCTION insert_flat_to_outbox();

create or replace  function check_exists_subscriber()
    returns trigger as $$
declare
    usr uuid;
begin
    select subscribers.user_id into usr from subscribers
    where subscribers.user_id=NEW.user_id and subscribers.house_id=NEW.house_id;

    if usr is not null then
        return NULL;
    end if;

    return NEW;
end;
$$ language plpgsql;

CREATE TRIGGER insert_subscribe_trigger
    BEFORE INSERT ON subscribers
    FOR EACH ROW
EXECUTE FUNCTION check_exists_subscriber();

create type flat as (flat_id int, house_id int, user_id uuid,
    price int,
    rooms int,
    status flat_status,
    moderator_id uuid);

create or replace  function update_status(new_status flat_status, new_flat_id int, new_house_id int, new_moderator_id uuid)
    returns flat as $$
declare
    mod_id uuid;
    f flat;
begin
    if new_status = 'on moderation' then
        select flats.moderator_id into mod_id from flats
        where flats.flat_id=new_flat_id and flats.house_id=new_house_id;

        if mod_id != new_moderator_id then
            raise exception 'flat already on moderation';
        end if;
    end if;

        update flats set status=new_status, moderator_id=new_moderator_id
            where flats.flat_id=new_flat_id and flats.house_id=new_house_id
            returning flat_id, house_id, user_id, price, rooms, status, moderator_id
        into f;
    return f;
end;
$$ language plpgsql;


insert into houses (address, construct_year, developer, create_house_date, update_flat_date) VALUES
('ул. Спортивная, д. 1', 2021, 'OOO Строй', now(), now()),
('ул. Спортивная, д. 2', 2020, 'ЗАО Строительство', now(), now()),
('ул. Спортивная, д. 3', 2019, 'ИП Строитель', now(), now()),
('ул. Спортивная, д. 4', 2022, 'ЗАО Новострой', now(), now()),
('ул. Спортивная, д. 5', 2021, 'OOO Строй', now(), now()),
('ул. Спортивная, д. 6', 2023, 'Компания Реал', now(), now()),
('ул. Спортивная, д. 7', 2020, 'ООО БК Строй', now(), now()),
('ул. Спортивная, д. 8', 2021, 'ЗАО Строительство', now(), now()),
('ул. Спортивная, д. 9', 2022, 'ООО Комфорт', now(), now()),
('ул. Спортивная, д. 10', 2018, 'ИП СтройПроект', now(), now());

insert into users(user_id, mail, password, role)
values ('019126ee-2b7d-758e-bb22-fe2e45b2db22', 'test@mail.ru', 'password', 'client');

insert into users(user_id, mail, password, role)
values ('019126ee-2b7d-758e-bb22-fe2e45b2db23', 'test@mail.ru', 'password', 'moderator');

INSERT INTO users (user_id, mail, password, role) VALUES
('019126ee-2b7d-758e-bb22-fe2e45b2db24', 'user1@mail.ru', 'password1', 'client'),
('019126ee-2b7d-758e-bb22-fe2e45b2db30', 'user2@mail.ru', 'password2', 'moderator'),
('019126ee-2b7d-758e-bb22-fe2e45b2db25', 'user3@mail.ru', 'password3', 'client'),
('019126ee-2b7d-758e-bb22-fe2e45b2db31', 'user4@mail.ru', 'password4', 'client'),
('019126ee-2b7d-758e-bb22-fe2e45b2db33', 'user5@mail.ru', 'password5', 'moderator'),
('019126ee-2b7d-758e-bb22-fe2e45b2db26', 'user6@mail.ru', 'password6', 'client'),
('019126ee-2b7d-758e-bb22-fe2e45b2db27', 'user7@mail.ru', 'password7', 'client'),
('019126ee-2b7d-758e-bb22-fe2e45b2db32', 'user8@mail.ru', 'password8', 'moderator'),
('019126ee-2b7d-758e-bb22-fe2e45b2db28', 'user9@mail.ru', 'password9', 'client'),
('019126ee-2b7d-758e-bb22-fe2e45b2db29', 'user10@mail.ru', 'password10', 'client');

insert into flats(flat_id, house_id, user_id, price, rooms, status)
values (10, 1, '019126ee-2b7d-758e-bb22-fe2e45b2db22', 100, 2, 'created');

INSERT INTO flats (flat_id, house_id, user_id, price, rooms, status) VALUES
(1, 1, '019126ee-2b7d-758e-bb22-fe2e45b2db22', 100, 2, 'created'),
(2, 1, '019126ee-2b7d-758e-bb22-fe2e45b2db24', 150, 3, 'approved'),
(3, 2, '019126ee-2b7d-758e-bb22-fe2e45b2db24', 200, 2, 'declined'),
(4, 2, '019126ee-2b7d-758e-bb22-fe2e45b2db25', 250, 4, 'on moderation'),
(5, 3, '019126ee-2b7d-758e-bb22-fe2e45b2db26', 300, 1, 'created'),
(6, 3, '019126ee-2b7d-758e-bb22-fe2e45b2db27', 350, 2, 'approved'),
(7, 4, '019126ee-2b7d-758e-bb22-fe2e45b2db28', 400, 3, 'declined'),
(8, 4, '019126ee-2b7d-758e-bb22-fe2e45b2db29', 450, 4, 'on moderation'),
(9, 5, '019126ee-2b7d-758e-bb22-fe2e45b2db29', 500, 2, 'created'),
(10, 5, '019126ee-2b7d-758e-bb22-fe2e45b2db29', 550, 3, 'approved');