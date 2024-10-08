<<<<<<< HEAD
# Сервис по продаже квартир
Тестовое задание Avito Bootcamp

![build workflow](https://github.com/NastyaAR/avito-test-task/actions/workflows/ci.yml/badge.svg)

## Стек

- язык сервиса Go
- база данных Postgresql
- балансировщик нагрузки nginx
- migrate для миграций
- docker & docker-compose для развертывания

## Описание

В рамках задания был реализован сервис по продаже квартир. 
Сервис предоставляет следующие возможности:

### Авторизация пользователей
- Endpoint /dummyLogin:
    - Позволяет передать желаемый тип пользователя (client, moderator).
    - В ответ возвращается токен с соответствующим уровнем доступа (обычный пользователь или модератор).
    - Токен необходимо передавать во все endpoints, требующие авторизации.

### Регистрация и авторизация пользователей по почте и паролю
- Endpoint /register:
    - Используется для регистрации нового пользователя.
    - В базе данных создается и сохраняется новый пользователь желаемого типа: обычный пользователь (client) или модератор (moderator).

- Endpoint /login:
    - У созданного пользователя появляется токен после успешной авторизации по почте и паролю.
    - Возвращается токен для пользователя с соответствующим уровнем доступа.

### Создание дома
- Endpoint /house/create:
    - Только модератор имеет возможность создать дом.
    - При успешном запросе возвращается полная информация о созданном доме.

### Создание квартиры
- Endpoint /flat/create:
    - Квартиру может создать любой пользователь.
    - При успешном запросе возвращается полная информация о квартире.
    - Объявление получает статус модерации created.
    - Обновляется дата последнего добавления жилья для дома, в котором была создана новая квартира.

### Модерация квартиры
- Статусы модерации квартиры:
    - Возможные статусы: created, approved, declined, on moderation.

- Endpoint /flat/update:
    - Только модератор может изменить статус модерации квартиры.
    - При успешном запросе возвращается полная информация об обновленной квартире.

### Получение списка квартир по номеру дома
- Endpoint /house/{id}:
    - Обычный пользователь и модератор могут получить список квартир по номеру дома.
    - Обычный пользователь видит только квартиры со статусом модерации approved, а модератор — жильё с любым статусом модерации.

### Подписка на уведомления
- Endpoint /house/{id}/subscribe:
    - Обычный пользователь может подписаться на уведомления о новых квартирах в доме по его номеру.

## Описание архитектуры и среды развертывания

Для развертывания используются Docker и Docker-Compose. Было необходимо обеспечить следующие условия: Всего квартир (до 20kk), RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%.

Чтобы выдерживать такую нагрузку было поднято три реплики сервиса и nginx в качестве балансировщика нагрузки (доступ к сервисам осуществляется через порт, на котором разввернут nginx, nginx распределяет запросы по репликам сервиса).

### Схема среды развертывания

![services](./img/services.png)

При запуске контейнеров при помощи docker-compose up сервис migrate накатывает стартовую миграцию на основную базу данных postgres-db.
Затем поднимается первый backend (в это же время создается образ для сервиса приложения). Остальные две реплики приложения поднимаются в зависимости от успеха/неуспеха запуска первой реплики.
Это сделано, чтобы не пересоздавать образ с приложением много раз.

nginx поднимается после запуска всех контейнеров приложения.

postgres-test - тестовая база данных.

### Инструкция по запуску

 - git clone https://github.com/NastyaAR/avito-test-task
 - cd avito-test-task
 - export POSTGRES_USER=creator POSTGRES_PASSWORD=password POSTGRES_DB=postgres-db KEY=any_key
 - docker-compose up -d

Тесты можно запустить с помощью команды: go test ./tests

Запросы к приложению выполняются по адресу: localhost:80

### Схема базы данных

База данных содержит сущности: Дом, Квартира, Пользователь. Для реализации отправки уведомлений 
при помощи паттерна Transactional Outbox есть сущность Уведомление.

Для реализации подписки создана дополнительная таблица Subscribers, которая реализует связь многие-ко-многим.

![bd](./img/bd.png)

### Особенности реализации

- Для соблюдения SLI успешности ответа — 99.99% сделан retryAdapter, он осуществляет доп запросы в бд, если данные не удалось получить после первого запроса (также там стоят задержки)
- Для RPS — 1k и SLI времени ответа — 50 мс - поднят nginx в качестве балансировщика нагрузки. Также создан индекс в базе данных для ускорения получения данных при фильтрации и распараллелена операция считывания данных
- Для надежной отправки уведомлений был использован паттерн Transactional Outbox и создана горутина, которая на протяжении работы приложения ходит в таблицу с уведомлениями и отсылает их
- Реализована пользовательская авторизация с помощью jwt-токенов. В токене шифруется id пользователя и роль для соблюдения ролевой модели
- Разработана модель прав доступа
- Пароли в базе данных хранятся в зашифрованном виде
- Создан логгер (вывод в файл и в консоль структурированных логов)
- Код для отправки уведомлений гибко встроен в сервис

### Аутентификация и авторизация

Аутентификация реализована на основе jwt-токена, который через некоторое время становится невалидным.
Разработан middleware, который проверяет токен в заголовке HTTP-запроса.

Авторизация реализована на основе ролей пользователей, который зашифрованы в токене доступа. В зависимости от роли в access middleware разрешается или запрещается доступ к тем или иным ресурсам.

### Отправка писем при подписке на дом

Для надежной at-least-once доставки письма адресату был использован паттерн Transactional Outbox.
При появлении новой квартиры в доме, на который подписан клиент, срабатывает триггер базы данных. Уведомление с ссылками на квартиру и адрес
пользователя добавляется в таблицу бд.

Создана горутина по типу демона. Она периодически ходит в бд и осуществляет отправку писем адресатам.
Контролируется каналом.

### CI

С помощью github actions настроен ci со сборкой, запуском и тестированием сервиса.
Создан бейдж со статусом билда.

### Вопросы, проблемы и решения

1) квартира должна кому-то принадлежать, чтобы в будущем пользователь смог посмотреть конкретно свои
квартиры

Решение: добавление user_id в Сущность квартиры

2) если модератор взял квартиру на модерацию, нужно запретить другим модераторам делать то же самое

Решение: в сущность квартиры внедрен внешний ключ, указывающий на модератора. При попытке обновить квартиру,
т е поменять ее статус срабатывает функция бд, которая проверяет, можно ли это сделать в зависимости от
текущего статуса квартиры и наличия/отсутствия модератора.

3) по условию не совсем понятно, чем именно являются номер дома и квартиры.

Решение: у дома номер - id primary key в бд, номер квартиры в бд - номер квартиры в доме.
Он не является уникальным ключом, поэтому в качестве primary key используется пара (номер дома, номер квартиры).

4) так как в квартиру был добавлен внешний ключ, указывающий на клиента, на endpount /dummyLogin
также создается запись в бд.

## Нагрузочные тесты

Провела нагрузочное тестирование с apache benchmark. При нагрузке 1000rps 0 failed запросов.
 
## Примеры работы

- Регистрация
```
curl -X POST http://localhost:80/register \
-H "Content-Type: application/json" \
-d '{
  "email": "test@gmail.com",
  "password": "password",
  "user_type": "client"
}'

{"user_id":"01913279-20a4-7455-bb10-01e4099cd1f8"}
```

- Аутентификация
```
curl -X POST http://localhost:80/login \
-H "Content-Type: application/json" \
-d '{
  "id": "01913277-e7dc-7ec4-9ede-95b5b9da2663",
  "password": "password"
}'

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkX3RpbWUiOjE3MjMxMzIzNzYsInJvbGUiOiJjbGllbnQiLCJ1c2VySUQiOiIwMTkxMzI3Ny1lN2RjLTdlYzQtOWVkZS05NWI1YjlkYTI2NjMifQ.hQ-pdQqWlK81Vy-77diN5I90E5h0FR5wOxpixhZgrCM"}
```

- Создание дома (предварительно был зарегистрирован модератор и залогинен)
```
curl -X POST http://localhost:80/house/create \
-H "Content-Type: application/json" \
-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkX3RpbWUiOjE3MjMxMzI0ODQsInJvbGUiOiJtb2RlcmF0b3IiLCJ1c2VySUQiOiIwMTkxMzI3Yi1kYjJmLTcwYmQtODY3ZC0yZWRjZDQ2YjBjZmMifQ.fVNxsiLMGgnGgWmnqsE8G-J_kepnPE1dewpirgOldds" \
-d '{
  "address": "Лесная улица, 7, Москва, 125196",
  "year": 2000,
  "developer": "Мэрия города"
}'
{"id":3,"address":"Лесная улица, 7, Москва, 125196","year":2000,"developer":"Мэрия города","created_at":"2024-08-08 14:58:59","update_at":"2024-08-08 14:58:59"}
```

- Создание квартиры
```
curl -X POST http://localhost:80/flat/create -H "Content-Type: application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkX3RpbWUiOjE3MjMxMzIzNzYsInJvbGUiOiJjbGllbnQiLCJ1c2VySUQiOiIwMTkxMzI3Ny1lN2RjLTdlYzQtOWVkZS05NWI1YjlkYTI2NjMifQ.hQ-pdQqWlK81Vy-77diN5I90E5h0FR5wOxpixhZgrCM" -d '{
  "flat_id": 10,
  "house_id": 3,
  "price": 10000,
  "rooms": 4
}'
{"id":10,"house_id":3,"price":10000,"rooms":4,"status":"created"}
```

- Обновление статуса
```
curl -X POST http://localhost:80/flat/update -H "Content-Type: application/json" -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkX3RpbWUiOjE3MjMxMzI0ODQsInJvbGUiOiJtb2RlcmF0b3IiLCJ1c2VySUQiOiIwMTkxMzI3Yi1kYjJmLTcwYmQtODY3ZC0yZWRjZDQ2YjBjZmMifQ.fVNxsiLMGgnGgWmnqsE8G-J_kepnPE1dewpirgOldds" -d '{
  "id": 10,
  "house_id": 3,
  "status": "on moderation"
}'

{"id":10,"house_id":3,"price":10000,"rooms":4,"status":"on moderation"}
```
=======
# testing-2024
Тестирование и отладка ПО
>>>>>>> main
