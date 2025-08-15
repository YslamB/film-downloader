drop table if exists episodes;
drop table if exists seasons;
drop table if exists films;


create table films (
    "id" serial primary key,
    "name" varchar(100) not null,
    "status" integer not null default 0, -- 0 created, 1- uploaded, 3- error, 
    "created_at" timestamp default now() not null
);

create table seasons (
    "id" serial primary key,
    "film_id" integer not null,
    "name" varchar(100) not null,
    "created_at" timestamp default now() not null,    
    CONSTRAINT seasons_film_id_fk
        FOREIGN KEY (film_id)
            REFERENCES films(id)
                ON DELETE CASCADE
                on update CASCADE
);

create table episodes (
    "id" serial primary key,
    "season_id" integer not null,
    "name" varchar(100) not null,
    "status" integer not null default 0, -- 0 created, 1- uploaded, 3- error, 
    "created_at" timestamp default now() not null,    
    CONSTRAINT seasons_season_id_fk
        FOREIGN KEY (season_id)
            REFERENCES seasons(id)
); 