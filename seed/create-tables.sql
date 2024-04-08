DROP TABLE IF EXISTS playable_race;
CREATE TABLE playable_race (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    speed       INT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

DROP TABLE IF EXISTS ability;
CREATE TABLE ability (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

DROP TABLE IF EXISTS proficiency;
CREATE TABLE proficiency (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

DROP TABLE IF EXISTS trait;
CREATE TABLE trait (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

DROP TABLE IF EXISTS language;
CREATE TABLE language (
    id          INT AUTO_INCREMENT NOT NULL,
    name        VARCHAR(128) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

DROP TABLE IF EXISTS starting_ability_bonus;
CREATE TABLE starting_ability_bonus (
    playable_race_id INT NOT NULL,
    ability_id INT NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (ability_id) REFERENCES ability(id),
    CONSTRAINT id UNIQUE (playable_race_id, ability_id)
);


DROP TABLE IF EXISTS starting_proficiency;
CREATE TABLE starting_proficiency (
    playable_race_id INT NOT NULL,
    proficiency_id INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (proficiency_id) REFERENCES proficiency(id),
    CONSTRAINT id UNIQUE (playable_race_id, proficiency_id)
);

DROP TABLE IF EXISTS starting_language;
CREATE TABLE starting_language (
    playable_race_id INT NOT NULL,
    language_id INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (language_id) REFERENCES language(id),
    CONSTRAINT id UNIQUE (playable_race_id, language_id)
);

DROP TABLE IF EXISTS starting_trait;
CREATE TABLE starting_trait (
    playable_race_id INT NOT NULL,
    trait_id INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (trait_id) REFERENCES trait(id),
    CONSTRAINT id UNIQUE (playable_race_id, trait_id)
);

DROP TABLE IF EXISTS sub_race;
CREATE TABLE sub_race (
    sub_race_id INT PRIMARY KEY,
    main_race_id INT NOT NULL,
    FOREIGN KEY (sub_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (main_race_id) REFERENCES playable_race(id)
);

DROP TABLE IF EXISTS starting_proficiency_option;
CREATE TABLE starting_proficiency_option (
    id INT NOT NULL,
    proficiency_id INT NOT NULL,
    playable_race_id INT NOT NULL,
    count INT NOT NULL,
    FOREIGN KEY (playable_race_id) REFERENCES playable_race(id),
    FOREIGN KEY (proficiency_id) REFERENCES proficiency(id)
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

INSERT INTO proficiency
    (id, name)
VALUES
    (1, "Battleaxes"),
    (2, "Brewer's Supplies"),
    (3, "Handaxes"),
    (4, "Light hammers"),
    (5, "Mason's Tools"),
    (6, "Smith's Tools"),
    (7, "Warhammers");

INSERT INTO language
    (id, name)
VALUES
    (1, "Draconic"),
    (2, "Dwarvish"),
    (3, "Common");

INSERT INTO trait
    (id, name)
VALUES
    (1, "Breath Weapon"),
    (3, "Damage Resistance"),
    (4, "Darkvision"),
    (5, "Draconic Ancestry"),
    (6, "Dwarven Combat Training"),
    (7, "Dwarven Resilience"),
    (8, "Stonecunning"),
    (9, "Tool Proficiency");

INSERT INTO playable_race
    (id, name, speed)
VALUES
    (1, "Dragonborn", 30), 
    (2, "Dwarf", 25), 
    (3, "Hill Dwarf", 25);

INSERT INTO starting_ability_bonus
    (playable_race_id, ability_id, amount)
VALUES
    (1, 5, 2),
    (1, 1, 1),
    (2, 2, 2);

INSERT INTO starting_proficiency
    (playable_race_id , proficiency_id)
VALUES
    (2, 1),
    (2, 2),
    (2, 3),
    (2, 4),
    (2, 5),
    (2, 6),
    (2, 7);

INSERT INTO starting_language
    (playable_race_id, language_id)
VALUES
    (1, 1),
    (1, 3),
    (2, 3),
    (2, 2);

INSERT INTO starting_trait
    (playable_race_id , trait_id)
VALUES
    (1, 1),
    (1, 3),
    (1, 5),
    (2, 4),
    (2, 6),
    (2, 7),
    (2, 8),
    (2, 9);

INSERT INTO sub_race
    (sub_race_id, main_race_id)
VALUES
    (3, 2);

INSERT INTO starting_proficiency_option
    (id, proficiency_id, playable_race_id, count)
VALUES
    (1, 6, 2, 1),
    (1, 2, 2, 1),
    (1, 5, 2, 1);
