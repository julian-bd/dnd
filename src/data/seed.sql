SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS ability;
DROP TABLE IF EXISTS language;
DROP TABLE IF EXISTS playable_race;
DROP TABLE IF EXISTS proficiency;
DROP TABLE IF EXISTS starting_ability_bonus;
DROP TABLE IF EXISTS starting_language;
DROP TABLE IF EXISTS starting_proficiency;
DROP TABLE IF EXISTS starting_proficiency_option;
DROP TABLE IF EXISTS starting_trait;
DROP TABLE IF EXISTS sub_race;
DROP TABLE IF EXISTS trait;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE playable_race (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    speed       INT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE TABLE ability (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE TABLE proficiency (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE TABLE trait (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(2048) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE TABLE language (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);


CREATE TABLE starting_ability_bonus (
    playable_race_id INT NOT NULL,
    ability_id INT NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (ability_id) REFERENCES ability(id),
    CONSTRAINT id UNIQUE (playable_race_id, ability_id)
);



CREATE TABLE starting_proficiency (
    playable_race_id INT NOT NULL,
    proficiency_id INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (proficiency_id) REFERENCES proficiency(id),
    CONSTRAINT id UNIQUE (playable_race_id, proficiency_id)
);


CREATE TABLE starting_language (
    playable_race_id INT NOT NULL,
    language_id INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (language_id) REFERENCES language(id),
    CONSTRAINT id UNIQUE (playable_race_id, language_id)
);


CREATE TABLE starting_trait (
    playable_race_id INT NOT NULL,
    trait_id INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (trait_id) REFERENCES trait(id),
    CONSTRAINT id UNIQUE (playable_race_id, trait_id)
);

CREATE TABLE sub_race (
    sub_race_id INT PRIMARY KEY,
    main_race_id INT NOT NULL,
    FOREIGN KEY (sub_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (main_race_id) REFERENCES playable_race(id)
);

CREATE TABLE starting_proficiency_option (
    id          INT AUTO_INCREMENT NOT NULL,
    group_id    binary(16) NOT NULL,
    proficiency_id INT NOT NULL,
    playable_race_id INT NOT NULL,
    count INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (proficiency_id) REFERENCES proficiency(id),
    PRIMARY KEY (id)
);

INSERT INTO ability
    (id, name)
VALUES
    (1, "CHA"),
    (2, "CON"),
    (3, "DEX"),
    (4, "INT"),
    (5, "STR"),
    (6, "WIS");
