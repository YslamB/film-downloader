CREATE TABLE countries (
                           id SERIAL PRIMARY KEY,
                           name_tm VARCHAR(100) NOT NULL,
                           name_ru VARCHAR(100) NOT NULL,
                           name_en VARCHAR(100) NOT NULL,
                           belet_id int,
                           created_at TIMESTAMP DEFAULT now(),
                           updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE genres (
                        id SERIAL PRIMARY KEY,
                        name_tm VARCHAR(100) NOT NULL,
                        name_ru VARCHAR(100) NOT NULL,
                        name_en VARCHAR(100) NOT NULL,
                        belet_id int,
                        created_at TIMESTAMP DEFAULT now()

);

CREATE TABLE languages (
                           id SERIAL PRIMARY KEY,
                           name_tm VARCHAR(100) NOT NULL,
                           name_ru VARCHAR(100) NOT NULL,
                           name_en VARCHAR(100) NOT NULL,
                           belet_id int,
                           created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE  image_sizes(
                             id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                             large json,
                             medium json,
                             small json,
                             "default" json,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE persons (
                         id SERIAL PRIMARY KEY,                -- Уникальный идентификатор
                         full_name VARCHAR(100) NOT NULL,      -- Имя человека
                         bio TEXT, -- Биография
                         image_id int,
                         belet_id int,
                         created_at TIMESTAMP DEFAULT now(),
                         CONSTRAINT fk_image_id
                             FOREIGN KEY (image_id)
                                 REFERENCES image_sizes(id)
                                 ON DELETE SET NULL
);


CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name_tm VARCHAR(100) NOT NULL,
                            name_ru VARCHAR(100) NOT NULL,
                            name_en VARCHAR(100) NOT NULL,
                            belet_id int,
                            created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE studios (
                         id SERIAL PRIMARY KEY,           -- Уникальный идентификатор
                         name VARCHAR(100) NOT NULL,      -- Имя студии
                         belet_id int,
                         created_at TIMESTAMP DEFAULT now()
);


CREATE TABLE files (
                       id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                       path VARCHAR(255) NOT NULL,     -- Путь к файлу
                       type VARCHAR(50) NOT NULL,      -- Тип файла (например, "video", "image")
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Дата создания
);


CREATE TABLE movies (
                        id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                        title VARCHAR(255) NOT NULL,    -- Название фильма или сериала
                        content_type VARCHAR(50) NOT NULL, -- Тип: 'movie' (фильм) или 'series' (сериал)
                        release_year INT,               -- Год выхода
                        description TEXT,               -- Описание
                        duration INT,                   -- Продолжительность (в минутах для фильмов, в сезонах для сериалов)
                        rating FLOAT,                   -- Рейтинг
                        rating_imdb FLOAT default 0,
                        rating_kinopoisk FLOAT default 0,
                        color varchar(50) default '',
                        age_restriction int default 0,
                        status VARCHAR(50) NOT NULL DEFAULT '',
                        category_id INT,                -- Внешний ключ для категории
                        language_id INT,                -- Внешний ключ для языка
                        file_id INT,                    -- Ссылка на файл (только для фильмов)
                        vertical int,
                        vertical_without_name int,
                        horizontal_with_name int,
                        horizontal_without_name int,
                        "image_name" int,
                        belet_id int,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Внешние ключи
                        CONSTRAINT fk_category
                            FOREIGN KEY (category_id)
                                REFERENCES categories(id)
                                ON DELETE SET NULL,

                        CONSTRAINT fk_language
                            FOREIGN KEY (language_id)
                                REFERENCES languages(id)
                                ON DELETE SET NULL,

                        CONSTRAINT fk_file
                            FOREIGN KEY (file_id)
                                REFERENCES files(id)
                                ON DELETE SET NULL,
                        CONSTRAINT fk_vertical
                            FOREIGN KEY (vertical)
                                REFERENCES image_sizes(id)
                                ON DELETE SET NULL,
                        CONSTRAINT fk_vertical_without_name
                            FOREIGN KEY (vertical_without_name)
                                REFERENCES image_sizes(id)
                                ON DELETE SET NULL,
                        CONSTRAINT fk_horizontal_with_name
                            FOREIGN KEY (horizontal_with_name)
                                REFERENCES image_sizes(id)
                                ON DELETE SET NULL,
                        CONSTRAINT fk_horizontal_without_name
                            FOREIGN KEY (horizontal_without_name)
                                REFERENCES image_sizes(id)
                                ON DELETE SET NULL,
                        CONSTRAINT fk_image_name
                            FOREIGN KEY (image_name)
                                REFERENCES image_sizes(id)
                                ON DELETE SET NULL

);


CREATE TABLE seasons (
                         id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                         title VARCHAR(100),              -- Название сезона
                         movie_id INT NOT NULL,          -- Ссылка на сериал
                         number INT NOT NULL,            -- Номер сезона
                         belet_id int,
                         FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);


CREATE TABLE episodes (
                          id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                          season_id INT NOT NULL,         -- Ссылка на сезон
                          number INT NOT NULL,    -- Номер эпизода в сезоне
                          title VARCHAR(255),             -- Название эпизода
                          status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
                          file_id INT,                    -- Ссылка на файл (видео эпизода)
                          image_id int,
                          duration int default 0,
                          belet_id int,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Внешние ключи
                          CONSTRAINT fk_season
                              FOREIGN KEY (season_id)
                                  REFERENCES seasons(id)
                                  ON DELETE CASCADE,

                          CONSTRAINT fk_file
                              FOREIGN KEY (file_id)
                                  REFERENCES files(id)
                                  ON DELETE SET NULL,

                          CONSTRAINT fk_image_id
                              FOREIGN KEY (image_id)
                                  REFERENCES image_sizes(id)
                                  ON DELETE SET NULL
);




CREATE TABLE movie_countries (
                                 movie_id INT NOT NULL,          -- Внешний ключ для фильма
                                 country_id INT NOT NULL,        -- Внешний ключ для страны

    -- Создание внешних ключей
                                 CONSTRAINT fk_movie
                                     FOREIGN KEY (movie_id)
                                         REFERENCES movies(id)
                                         ON DELETE CASCADE,

                                 CONSTRAINT fk_country
                                     FOREIGN KEY (country_id)
                                         REFERENCES countries(id)
                                         ON DELETE CASCADE,

    -- Составной первичный ключ
                                 PRIMARY KEY (movie_id, country_id)
);

CREATE TABLE movie_genres (
                              movie_id INT NOT NULL,          -- Внешний ключ для фильма
                              genre_id INT NOT NULL,          -- Внешний ключ для жанра

    -- Внешние ключи
                              CONSTRAINT fk_movie
                                  FOREIGN KEY (movie_id)
                                      REFERENCES movies(id)
                                      ON DELETE CASCADE,

                              CONSTRAINT fk_genre
                                  FOREIGN KEY (genre_id)
                                      REFERENCES genres(id)
                                      ON DELETE CASCADE,

    -- Составной первичный ключ
                              PRIMARY KEY (movie_id, genre_id)
);

CREATE TABLE movie_actors (
                              movie_id INT NOT NULL,          -- Внешний ключ для фильма
                              actor_id INT NOT NULL,         -- Внешний ключ для человека

    -- Внешние ключи
                              CONSTRAINT fk_movie
                                  FOREIGN KEY (movie_id)
                                      REFERENCES movies(id)
                                      ON DELETE CASCADE,

                              CONSTRAINT fk_actor
                                  FOREIGN KEY (actor_id)
                                      REFERENCES persons(id)
                                      ON DELETE CASCADE,

    -- Составной первичный ключ
                              PRIMARY KEY (movie_id, actor_id)
);

CREATE TABLE movie_directors (
                                 movie_id INT NOT NULL,          -- Внешний ключ для фильма
                                 director_id INT NOT NULL,         -- Внешний ключ для человека

    -- Внешние ключи
                                 CONSTRAINT fk_movie
                                     FOREIGN KEY (movie_id)
                                         REFERENCES movies(id)
                                         ON DELETE CASCADE,

                                 CONSTRAINT fk_director
                                     FOREIGN KEY (director_id)
                                         REFERENCES persons(id)
                                         ON DELETE CASCADE,

    -- Составной первичный ключ
                                 PRIMARY KEY (movie_id, director_id)
);

CREATE TABLE movie_studios (
                               movie_id INT NOT NULL,          -- Внешний ключ для фильма
                               studio_id INT NOT NULL,         -- Внешний ключ для студии
-- Внешние ключи
                               CONSTRAINT fk_movie
                                   FOREIGN KEY (movie_id)
                                       REFERENCES movies(id)
                                       ON DELETE CASCADE,

                               CONSTRAINT fk_studio
                                   FOREIGN KEY (studio_id)
                                       REFERENCES studios(id)
                                       ON DELETE CASCADE,

    -- Составной первичный ключ
                               PRIMARY KEY (movie_id, studio_id)
);


CREATE TABLE trailers (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255),
                          movie_id BIGINT NOT NULL,
                          status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
                          duration int default 0,
                          image_id int,
                          file_id int,
                          created_at TIMESTAMP DEFAULT now(),
                          FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
                          CONSTRAINT fk_image_id
                              FOREIGN KEY (image_id)
                                  REFERENCES image_sizes(id)
                                  ON DELETE SET NULL,
                          CONSTRAINT fk_file
                              FOREIGN KEY (file_id)
                                  REFERENCES files(id)
                                  ON DELETE SET NULL

);


CREATE TABLE  admin_users(
                             id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                             name CHARACTER VARYING(250) NOT NULL,
                             login CHARACTER VARYING(250) not null,
                             password CHARACTER VARYING(250) not null,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             UNIQUE("login")
);

CREATE TABLE main_categories (
                                 id SERIAL PRIMARY KEY,
                                 name_tm VARCHAR(100) NOT NULL,
                                 name_ru VARCHAR(100) NOT NULL,
                                 name_en VARCHAR(100) NOT NULL,
                                 position int default 0,
                                 belet_id int,
                                 created_at TIMESTAMP DEFAULT now()
);

insert into main_categories(name_tm, name_ru, name_en) values (
                                                                  'Baş sahypa',
                                                                  'Главная страница',
                                                                  'Main page'
                                                              );


create table shelves(
                        id SERIAL PRIMARY KEY,          -- Уникальный идентификатор
                        name_tm VARCHAR(100) NOT NULL,
                        name_ru VARCHAR(100) NOT NULL,
                        name_en VARCHAR(100) NOT NULL,
                        position int,
                        description_tm text,
                        description_ru text,
                        description_en text,
                        color varchar(50) default '',
                        image_id int,
                        is_visible boolean default true,
                        category_id int,
                        "type" varchar(50),
                        created_at TIMESTAMP DEFAULT now(),
    -- Внешние ключи
                        CONSTRAINT fk_image_id
                            FOREIGN KEY (image_id)
                                REFERENCES image_sizes(id)
                                ON DELETE SET NULL,

                        CONSTRAINT fk_category
                            FOREIGN KEY (category_id)
                                REFERENCES main_categories(id)
                                ON DELETE SET NULL
);




CREATE TABLE shelf_movie (
                             id SERIAL PRIMARY KEY,          -- Уникальный числовой идентификатор
                             shelf_id INT NOT NULL,          -- Ссылка на полку
                             movie_id INT NOT NULL,          -- Ссылка на фильм
                             position INT,                   -- Позиция на полке

    -- Внешние ключи
                             CONSTRAINT fk_movie
                                 FOREIGN KEY (movie_id)
                                     REFERENCES movies(id)
                                     ON DELETE CASCADE,

                             CONSTRAINT fk_shelf
                                 FOREIGN KEY (shelf_id)
                                     REFERENCES shelves(id)
                                     ON DELETE CASCADE,

    -- Уникальное ограничение (чтобы один фильм не мог быть на одной полке дважды)
                             CONSTRAINT uk_shelf_movie UNIQUE (shelf_id, movie_id)
);
