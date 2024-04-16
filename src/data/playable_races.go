package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type PlayableRace struct {
	ID                         int                         `json:"id"`
	Name                       string                      `json:"name"`
	Speed                      int                         `json:"speed"`
	AbilityBonuses             []AbilityBonus              `json:"ability_bonuses"`
	StartingLanguages          []string                    `json:"starting_languages"`
	StartingProficiencies      []string                    `json:"starting_proficiencies"`
	StartingProficiencyOptions []StartingProficiencyOption `json:"starting_proficiency_options"`
	Traits                     []string                    `json:"traits"`
	SubRaces                   []string                    `json:"sub_races"`
}

type AbilityBonus struct {
	Ability string `json:"ability"`
	Bonus   int    `json:"bonus"`
}

type trait struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StartingProficiencyOption struct {
	Count   int      `json:"count"`
	Options []string `json:"options"`
}

func PlayableRaceNames() ([]string, error) {
	var names []string
	query := `SELECT name FROM playable_race`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pr PlayableRace
		if err := rows.Scan(&pr.Name); err != nil {
			return nil, err
		}
		names = append(names, pr.Name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

func hydrateRace(playableRace PlayableRace) (PlayableRace, error) {
	startingLanguages, err := startingLanguages(playableRace.ID)
	if err != nil {
		return playableRace, err
	}
	playableRace.StartingLanguages = startingLanguages

	startingProficiencies, err := startingProficiencies(playableRace.ID)
	if err != nil {
		return playableRace, err
	}
	playableRace.StartingProficiencies = startingProficiencies

	startingAbilityBonuses, err := startingAbilityBonuses(playableRace.ID)
	if err != nil {
		return playableRace, err
	}
	playableRace.AbilityBonuses = startingAbilityBonuses

	startingTraits, err := startingTraits(playableRace.ID)
	if err != nil {
		return playableRace, err
	}
	playableRace.Traits = startingTraits

	subRaces, err := subRaces(playableRace.ID)
	if err != nil {
		return playableRace, err
	}
	playableRace.SubRaces = subRaces

	startingProficiencyOptions, err := startingProficiencyOptions(playableRace.ID)
	if err != nil {
		return playableRace, err
	}
	playableRace.StartingProficiencyOptions = startingProficiencyOptions

	return HydrateSubRace(playableRace)
}

func PlayableRaceById(id int) (PlayableRace, error) {
	var playableRace PlayableRace

	row := db.QueryRow("SELECT id, name, speed FROM playable_race WHERE id = ?", id)
	if err := row.Scan(&playableRace.ID, &playableRace.Name, &playableRace.Speed); err != nil {
		if err == sql.ErrNoRows {
			return playableRace, err
		}
		return playableRace, err
	}

    return hydrateRace(playableRace)
}

func PlayableRaceByName(name string) (PlayableRace, error) {
	var playableRace PlayableRace

	row := db.QueryRow("SELECT id, name, speed FROM playable_race WHERE name = ?", name)
	if err := row.Scan(&playableRace.ID, &playableRace.Name, &playableRace.Speed); err != nil {
		if err == sql.ErrNoRows {
			return playableRace, err
		}
		return playableRace, err
	}

    return hydrateRace(playableRace)
}

func HydrateSubRace(subRace PlayableRace) (PlayableRace, error) {
    row := db.QueryRow(`
        SELECT main_race_id AS id 
        FROM sub_race
        WHERE sub_race_id = ?
        LIMIT 1
        `, 
        subRace.ID)

    var mainRaceId int
    if err := row.Scan(&mainRaceId); err != nil {
        if err == sql.ErrNoRows {
            return subRace, nil
        }
        return subRace, err
    }

    mainRace, err := PlayableRaceById(mainRaceId)
    if err != nil {
        return subRace, err
    }

    if subRace.Speed == 0 {
        subRace.Speed = mainRace.Speed
    }

    subRace.AbilityBonuses = append(subRace.AbilityBonuses, mainRace.AbilityBonuses...)
    subRace.StartingLanguages = append(subRace.StartingLanguages, mainRace.StartingLanguages...)
    subRace.StartingProficiencies = append(subRace.StartingProficiencies, mainRace.StartingProficiencies...)
    subRace.StartingProficiencyOptions = append(subRace.StartingProficiencyOptions, mainRace.StartingProficiencyOptions...)
    subRace.Traits = append(subRace.Traits, mainRace.Traits...)

    return subRace, nil
}

func InsertTrait(trait string) error {
    _, err := db.Exec("INSERT INTO trait (name) VALUES (?)", trait)
    if err != nil {
        return err
    }
    return nil
}

func InsertProficiency(trait string) error {
    _, err := db.Exec("INSERT INTO proficiency (name) VALUES (?)", trait)
    if err != nil {
        return err
    }
    return nil
}

func InsertLanguage(trait string) error {
    _, err := db.Exec("INSERT INTO language (name) VALUES (?)", trait)
    if err != nil {
        return err
    }
    return nil
}

func InsertPlayableRace(playableRace PlayableRace) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO playable_race (name, speed) VALUES (?,?)",
		playableRace.Name,
		playableRace.Speed)
	if err != nil {
		return 0, err
	}

	playableRaceId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if playableRace.AbilityBonuses != nil {
		for _, bonus := range playableRace.AbilityBonuses {
			var abilityId int64
			row := tx.QueryRow(`SELECT id FROM ability WHERE ability.name=? LIMIT 1`, bonus.Ability)
			if err := row.Scan(&abilityId); err != nil {
				return 0, err
			}
			_, err := tx.Exec(`
                    INSERT INTO starting_ability_bonus (playable_race_id, ability_id, amount)
                    VALUES (?, ?, ?)
                `,
				playableRaceId,
				abilityId,
				bonus.Bonus,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	if playableRace.StartingLanguages != nil {
		for _, language := range playableRace.StartingLanguages {
			var languageId int64
			row := tx.QueryRow(`SELECT id FROM language WHERE language.name=? LIMIT 1`, language)
			if err := row.Scan(&languageId); err != nil {
				if err != sql.ErrNoRows {
					return 0, err
				}
				result, err = tx.Exec(`INSERT INTO language (name) VALUES (?)`, language)
				if err != nil {
					return 0, err
				}
				languageId, err = result.LastInsertId()
				if err != nil {
					return 0, err
				}
			}
			_, err := tx.Exec(`
                    INSERT INTO starting_language (playable_race_id, language_id)
                    VALUES (?, ?)
                `,
				playableRaceId,
				languageId,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	if playableRace.StartingProficiencies != nil {
		for _, proficiency := range playableRace.StartingProficiencies {
			var proficiencyId int64
			row := tx.QueryRow(`SELECT id FROM proficiency WHERE proficiency.name=? LIMIT 1`, proficiency)
			if err := row.Scan(&proficiencyId); err != nil {
				if err != sql.ErrNoRows {
					return 0, err
				}
				result, err = tx.Exec(`INSERT INTO proficiency (name) VALUES (?)`, proficiency)
				if err != nil {
					return 0, err
				}
				proficiencyId, err = result.LastInsertId()
				if err != nil {
					return 0, err
				}
			}
			_, err := tx.Exec(`
                    INSERT INTO starting_proficiency (playable_race_id, proficiency_id)
                    VALUES (?, ?)
                `,
				playableRaceId,
				proficiencyId,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	if playableRace.StartingProficiencyOptions != nil {
		for _, o := range playableRace.StartingProficiencyOptions {
			groupId := uuid.New().String()
			for _, opt := range o.Options {
				var proficiencyId int64
				row := tx.QueryRow(`SELECT id FROM proficiency WHERE proficiency.name=?`, opt)
				if err := row.Scan(&proficiencyId); err != nil {
					return 0, err
				}
				_, err := tx.Exec(`
                    INSERT INTO starting_proficiency_option (group_id, proficiency_id, playable_race_id, count)
                    VALUES (UUID_TO_BIN(?), ?, ?, ?)
                    `,
					groupId,
					proficiencyId,
					playableRaceId,
					o.Count,
				)
				if err != nil {
					return 0, err
				}
			}
		}
	}

	if playableRace.Traits != nil {
		for _, trait := range playableRace.Traits {
			var traitId int64
			row := tx.QueryRow(`SELECT id FROM trait WHERE trait.name=? LIMIT 1`, trait)
			if err := row.Scan(&traitId); err != nil {
				if err != sql.ErrNoRows {
					return 0, err
				}
				result, err = tx.Exec(`INSERT INTO trait (name) VALUES (?)`, trait)
				if err != nil {
					return 0, err
				}
				traitId, err = result.LastInsertId()
				if err != nil {
					return 0, err
				}
			}
			_, err := tx.Exec(`
                    INSERT INTO starting_trait (playable_race_id, trait_id)
                    VALUES (?, ?)
                `,
				playableRaceId,
				traitId,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	if playableRace.SubRaces != nil {
		for _, subRace := range playableRace.SubRaces {
			var subRaceId int64
			row := tx.QueryRow(`SELECT id FROM playable_race WHERE playable_race.name=?  LIMIT 1`, subRace)
			if err := row.Scan(&subRaceId); err != nil {
				return 0, err
			}
			_, err := tx.Exec(`
                    INSERT INTO sub_race (sub_race_id, main_race_id)
                    VALUES (?, ?)
                `,
				subRaceId,
				playableRaceId,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return playableRaceId, nil
}

func startingLanguages(id int) ([]string, error) {
	query := `
        SELECT language.name AS name
        FROM starting_language
        JOIN language
        ON language.id = starting_language.language_id
        WHERE starting_language.playable_race_id = ?`
	var startingLanguages []string
	languageRows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer languageRows.Close()
	for languageRows.Next() {
		var language string
		if err := languageRows.Scan(&language); err != nil {
			return nil, err
		}
		startingLanguages = append(startingLanguages, language)
	}
	if err := languageRows.Err(); err != nil {
		return nil, err
	}
	return startingLanguages, nil
}

func subRaces(id int) ([]string, error) {
	query := `
        SELECT playable_race.name
        FROM sub_race
        JOIN playable_race
        ON playable_race.id = sub_race.sub_race_id
        WHERE sub_race.main_race_id = ?`
	var startingLanguages []string
	subRaceRows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer subRaceRows.Close()
	for subRaceRows.Next() {
		var language string
		if err := subRaceRows.Scan(&language); err != nil {
			return nil, err
		}
		startingLanguages = append(startingLanguages, language)
	}
	if err := subRaceRows.Err(); err != nil {
		return nil, err
	}
	return startingLanguages, nil
}

func startingProficiencies(id int) ([]string, error) {
	query := `
        SELECT proficiency.name AS name
        FROM starting_proficiency
        JOIN proficiency
        ON proficiency.id = starting_proficiency.proficiency_id
        WHERE starting_proficiency.playable_race_id = ?`
	var startingProficiencies []string
	proficiencyRows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer proficiencyRows.Close()
	for proficiencyRows.Next() {
		var proficiency string
		if err := proficiencyRows.Scan(&proficiency); err != nil {
			return nil, err
		}
		startingProficiencies = append(startingProficiencies, proficiency)
	}
	if err := proficiencyRows.Err(); err != nil {
		return nil, err
	}
	return startingProficiencies, nil
}

type startingProficiencyOptionRow struct {
	GroupId string
	Name    string
	Count   int
}

/*
DROP TABLE IF EXISTS ability;
CREATE TABLE ability (

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
*/
func abilityBonuses(id int) ([]AbilityBonus, error) {
	query := `
        SELECT ability.name, bonus.amount
        FROM ability
        JOIN starting_ability_bonus AS bonus
        ON ability.id = bonus.ability_id
        WHERE bonus.player_race_id = ?
    `
	var bs []AbilityBonus
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b AbilityBonus
		if err := rows.Scan(b); err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}

	return bs, nil
}

func startingProficiencyOptions(id int) ([]StartingProficiencyOption, error) {
	query := ` 
        SELECT starting_proficiency_option.group_id, proficiency.name, starting_proficiency_option.count
        FROM starting_proficiency_option
        JOIN proficiency
        ON starting_proficiency_option.proficiency_id = proficiency.id
        WHERE starting_proficiency_option.playable_race_id = ?
    `
	var m map[string]StartingProficiencyOption
	m = make(map[string]StartingProficiencyOption)

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row startingProficiencyOptionRow
		if err := rows.Scan(&row.GroupId, &row.Name, &row.Count); err != nil {
			return nil, err
		}
		r := m[row.GroupId]
		r.Count = row.Count
		r.Options = append(r.Options, row.Name)
		m[row.GroupId] = r

	}
	defer rows.Close()

	var startingProficiencyOptions []StartingProficiencyOption
	for _, val := range m {
		startingProficiencyOptions = append(startingProficiencyOptions, val)
	}
	return startingProficiencyOptions, nil
}

func startingTraits(id int) ([]string, error) {
	query := `
        SELECT trait.name AS name
        FROM starting_trait
        JOIN trait
        ON trait.id = starting_trait.trait_id
        WHERE starting_trait.playable_race_id = ?`
	var startingTraits []string
	traitRows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer traitRows.Close()
	for traitRows.Next() {
		var trait string
		if err := traitRows.Scan(&trait); err != nil {
			return nil, err
		}
		startingTraits = append(startingTraits, trait)
	}
	if err := traitRows.Err(); err != nil {
		return nil, err
	}
	return startingTraits, nil
}

func startingAbilityBonuses(id int) ([]AbilityBonus, error) {
	query := `
        SELECT name, amount
        FROM starting_ability_bonus
        JOIN ability
        ON ability.id = starting_ability_bonus.ability_id
        WHERE starting_ability_bonus.playable_race_id = ?`
	var abilityBonuses []AbilityBonus
	abilityBonusRows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer abilityBonusRows.Close()
	for abilityBonusRows.Next() {
		var abilityBonus AbilityBonus
		if err := abilityBonusRows.Scan(&abilityBonus.Ability, &abilityBonus.Bonus); err != nil {
			return nil, err
		}
		abilityBonuses = append(abilityBonuses, abilityBonus)
	}
	if err := abilityBonusRows.Err(); err != nil {
		return nil, err
	}
	return abilityBonuses, nil
}
