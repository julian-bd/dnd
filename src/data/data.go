package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db *sql.DB

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

func PlayableRaceByName(name string) (PlayableRace, error) {
	var pr PlayableRace

	row := db.QueryRow("SELECT id, name, speed FROM playable_race WHERE name = ?", name)
	if err := row.Scan(&pr.ID, &pr.Name, &pr.Speed); err != nil {
		if err == sql.ErrNoRows {
			return pr, fmt.Errorf("no such race: %v", name)
		}
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}

	startingLanguages, err := startingLanguages(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.StartingLanguages = startingLanguages

	startingProficiencies, err := startingProficiencies(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.StartingProficiencies = startingProficiencies

	startingAbilityBonuses, err := startingAbilityBonuses(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.AbilityBonuses = startingAbilityBonuses

	startingTraits, err := startingTraits(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Traits = startingTraits

	subRaces, err := subRaces(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.SubRaces = subRaces

	startingProficiencyOptions, err := startingProficiencyOptions(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.StartingProficiencyOptions = startingProficiencyOptions

	return pr, nil
}

func InsertPlayableRace(playableRace PlayableRace) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error creating ctx")
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
		return 0, fmt.Errorf("insertion error (2)")
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

func startingProficiencyOptions(id int) ([]startingProficiencyOption, error) {
	query := ` 
        SELECT starting_proficiency_option.group_id, proficiency.name, starting_proficiency_option.count
        FROM starting_proficiency_option
        JOIN proficiency
        ON starting_proficiency_option.proficiency_id = proficiency.id
        WHERE starting_proficiency_option.playable_race_id = ?
    `
	var m map[string]startingProficiencyOption
	m = make(map[string]startingProficiencyOption)

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

	var startingProficiencyOptions []startingProficiencyOption
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

func startingAbilityBonuses(id int) ([]abilityBonus, error) {
	query := `
        SELECT name, amount
        FROM starting_ability_bonus
        JOIN ability
        ON ability.id = starting_ability_bonus.ability_id
        WHERE starting_ability_bonus.playable_race_id = ?`
	var abilityBonuses []abilityBonus
	abilityBonusRows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer abilityBonusRows.Close()
	for abilityBonusRows.Next() {
		var abilityBonus abilityBonus
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

func InitDB() error {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "dnd",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return db.Ping()
}
