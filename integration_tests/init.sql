-- SCHEMA START
create type semester as enum ('FALL', 'SPRING', 'SUMMER');

create type subscription_tier as enum ('BRONZE', 'SILVER', 'GOLD', 'PLATINUM');

create table sponsors
(
    id          serial,
    name        varchar           not null,
    tier        subscription_tier not null,
    since       date              not null,
    description varchar,
    website     varchar,
    logo_url    varchar,
    constraint sponsors_pk
        primary key (id, name)
);

create unique index sponsors_id_uindex
    on sponsors (id);

create unique index sponsors_name_uindex
    on sponsors (name);

create table terms
(
    id       serial
        constraint terms_pk
            primary key,
    year     integer  not null,
    semester semester not null
);

create unique index terms_id_uindex
    on terms (id);

create table hackathons
(
    id         serial
        constraint hackathons_pk
            primary key,
    term_id    serial
        constraint hackathons_terms_id_fk
            references terms,
    start_date timestamp not null,
    end_date   timestamp not null
);

create unique index hackathons_id_uindex
    on hackathons (id);

create unique index hackathons_term_id_uindex
    on hackathons (term_id);

create table pronouns
(
    id         serial
        constraint pronouns_pk
            primary key,
    subjective varchar not null,
    objective  varchar not null
);

create unique index pronouns_id_uindex
    on pronouns (id);

create table users
(
    id                  serial
        constraint users_pk
            primary key,
    email               varchar not null,
    phone_number        varchar,
    last_name           varchar not null,
    age                 integer,
        pronoun_id          integer,
    first_name          varchar not null,
    role                varchar not null,
    oauth_uid           varchar not null
        constraint users_oauth_uid_unique
            unique,
    oauth_provider      varchar not null,
    years_of_experience double precision,
    shirt_size          varchar not null,
    race                character varying[],
    gender              varchar
);

create unique index users_email_uindex
    on users (email);

create unique index users_phone_number_uindex
    on users (phone_number);

create table hackathon_sponsors
(
    hackathon_id integer not null
        constraint hackathon_sponsors_hackathons_null_fk
            references hackathons,
    sponsor_id   integer not null
        constraint hackathon_sponsors_sponsors_null_fk
            references sponsors (id)
);

create table events
(
    id           serial
        constraint events_pk
            primary key,
    hackathon_id integer   not null
        constraint events_hackathons_id_fk
            references hackathons,
    location     varchar   not null,
    start_date   timestamp not null,
    end_date     timestamp not null,
    name         varchar   not null,
    description  varchar   not null
);

create table hackathon_applications
(
    id                        serial
        constraint hackathon_applications_pk
            primary key,
    user_id                   integer                 not null
        constraint hackathon_applications_users_id_fk
            references users,
    hackathon_id              integer                 not null
        constraint hackathon_applications_hackathons_id_fk
            references hackathons,
    why_attend                character varying[]     not null,
    what_do_you_want_to_learn character varying[]     not null,
    share_info_with_sponsors  boolean                 not null,
    application_status        varchar                 not null,
    created_time              timestamp default now() not null,
    status_change_time        timestamp
);

create table mailing_addresses
(
    user_id       integer             not null
        constraint mailing_addresses_pk
            primary key
        constraint mailing_addresses_users_id_fk
            references users,
    country       varchar             not null,
    state         varchar             not null,
    city          varchar             not null,
    postal_code   varchar             not null,
    address_lines character varying[] not null
);

create table mlh_terms
(
    user_id         integer not null
        constraint mlh_terms_pk
            primary key
        constraint mlh_terms_users_null_fk
            references users,
    send_messages   boolean not null,
    share_info      boolean not null,
    code_of_conduct boolean not null
);

create table education_info
(
    user_id         integer   not null
        constraint education_info_pk
            primary key
        constraint education_info_users_id_fk
            references users,
    name            varchar   not null,
    major           varchar   not null,
    graduation_date timestamp not null,
    level           varchar
);

create table event_attendance
(
    event_id integer                 not null
        constraint event_attendance_events_id_fk
            references events,
    user_id  integer                 not null
        constraint event_attendance_users_id_fk
            references users,
    time     timestamp default now() not null,
    constraint event_attendance_pk
        primary key (event_id, user_id)
);

create table meals
(
    hackathon_id integer             not null
        constraint meals_hackathons_null_fk
            references hackathons,
    user_id      integer             not null
        constraint meals_users_null_fk
            references users,
    meals        character varying[] not null,
    constraint meals_pk
        primary key (hackathon_id, user_id)
);

create table hackathon_checkin
(
    hackathon_id integer   not null
        constraint hackathon_checkin_hackathons_id_fk
            references hackathons,
    user_id      integer   not null
        constraint hackathon_checkin_users_id_fk
            references users,
    time         timestamp not null,
    constraint hackathon_checkin_pk
        primary key (hackathon_id, user_id)
);

create table api_keys
(
    user_id integer   not null
        constraint api_keys_pk
            primary key
        constraint api_keys_users_id_fk
            references users,
    key     varchar   not null,
    created timestamp not null
);

create unique index api_keys_key_uindex
    on api_keys (key);

-- SCHEMA END

-- INTEGRATION TEST DATA START

INSERT INTO pronouns (subjective, objective)
VALUES ('he', 'him'); -- ID = 1

INSERT INTO users (email, phone_number, last_name, age, pronoun_id, first_name, role, oauth_uid,
                   oauth_provider, years_of_experience, shirt_size, race, gender)
VALUES ('joe.bob@example.com'::varchar, '100-200-3000'::varchar, 'Bob'::varchar, 22::integer, 1::integer,
        'Joe'::varchar, 'NORMAL'::varchar, '1'::varchar, 'GITHUB'::varchar, 3.5::double precision, 'L'::varchar,
        ARRAY ['CAUCASIAN'], 'MALE'::varchar);
-- ID = 1

INSERT INTO mlh_terms (user_id, send_messages, share_info, code_of_conduct)
VALUES (1, true, true, true);

INSERT INTO mailing_addresses (user_id, country, state, city, postal_code, address_lines)
VALUES (1, 'United States', 'Florida', 'Orlando', '32765', ARRAY ['1000 Abc Rd', 'APT 69']);

INSERT INTO users (email, phone_number, last_name, age, pronoun_id, first_name, role, oauth_uid,
                   oauth_provider, years_of_experience, shirt_size, race, gender)
VALUES ('joe.biron@example.com'::varchar, '123-456-7890'::varchar, 'Biron'::varchar, 21::integer, 1::integer,
        'Joe'::varchar, 'NORMAL'::varchar, '4'::varchar, 'GITHUB'::varchar, 3.5::double precision, 'L'::varchar,
        ARRAY ['AFRICAN_AMERICAN'], 'MALE'::varchar);
-- ID = 2

INSERT INTO api_keys (user_id, key, created)
VALUES (2, '1234567890abc', '2022-11-09')
-- ID = 1


INSERT INTO users (first_name, last_name, email, phone_number, age, pronoun_id, oauth_uid, oauth_provider, role,years_of_experience, shirt_size, race, gender, race) 
    VALUES ("dough"::varchar, "boy"::varchar, "doughboy@gmail.com"::varchar, "4071234567"::varchar, 16::integer, 1::integer, 
    "12velofabo12"::varchar, "GITHUB"::varchar, "ADMIN"::varchar, "2.5"::double precision, "Sm"::varchar, 
    ARRAY ['AFRICAN_AMERICAN'], "MALE"::varchar,);
--ID=3

-- INTEGRATION TEST DATA END