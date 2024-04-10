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
		return nil, fmt.Errorf("error getting playable races (1)")
	}
	defer rows.Close()
	for rows.Next() {
		var pr PlayableRace
		if err := rows.Scan(&pr.Name); err != nil {
			return nil, fmt.Errorf("error getting playable races (2)")
		}
		names = append(names, pr.Name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting playable races (3)")
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

	starting_languages, err := get_starting_languages(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Starting_Languages = starting_languages

	starting_proficiencies, err := get_starting_proficiencies(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Starting_Proficiencies = starting_proficiencies

	starting_ability_bonuses, err := get_starting_ability_bonuses(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Ability_Bonuses = starting_ability_bonuses

	starting_traits, err := get_starting_traits(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Traits = starting_traits

	sub_races, err := get_sub_races(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Sub_Races = sub_races

	starting_proficiency_options, err := get_starting_proficiency_options(pr.ID)
	if err != nil {
		return pr, fmt.Errorf("query error %v: %d", name, err)
	}
	pr.Starting_Proficiency_Options = starting_proficiency_options

	return pr, nil
}

func InsertPlayableRace(playable_race PlayableRace) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error creating ctx")
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO playable_race (name, speed) VALUES (?,?)",
		playable_race.Name,
		playable_race.Speed)
	if err != nil {
		return 0, err
	}

	playableRaceId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertion error (2)")
	}

	if playable_race.Ability_Bonuses != nil {
		for _, bonus := range playable_race.Ability_Bonuses {
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

	if playable_race.Starting_Languages != nil {
		for _, language := range playable_race.Starting_Languages {
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

	if playable_race.Starting_Proficiencies != nil {
		for _, proficiency := range playable_race.Starting_Proficiencies {
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

	if playable_race.Starting_Proficiency_Options != nil {
		for _, o := range playable_race.Starting_Proficiency_Options {
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

	if playable_race.Traits != nil {
		for _, trait := range playable_race.Traits {
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

	if playable_race.Sub_Races != nil {
		for _, subRace := range playable_race.Sub_Races {
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

func get_starting_languages(id int) ([]string, error) {
	query := `
        SELECT language.name AS name
        FROM starting_language
        JOIN language
        ON language.id = starting_language.language_id
        WHERE starting_language.playable_race_id = ?`
	var starting_languages []string
	language_rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("language query error (1)")
	}
	defer language_rows.Close()
	for language_rows.Next() {
		var language string
		if err := language_rows.Scan(&language); err != nil {
			return nil, fmt.Errorf("language query error (2)")
		}
		starting_languages = append(starting_languages, language)
	}
	if err := language_rows.Err(); err != nil {
		return nil, fmt.Errorf("language query error (3)")
	}
	return starting_languages, nil
}

func get_sub_races(id int) ([]string, error) {
	query := `
        SELECT playable_race.name
        FROM sub_race
        JOIN playable_race
        ON playable_race.id = sub_race.sub_race_id
        WHERE sub_race.main_race_id = ?`
	var starting_languages []string
	sub_race_rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("language query error (1)")
	}
	defer sub_race_rows.Close()
	for sub_race_rows.Next() {
		var language string
		if err := sub_race_rows.Scan(&language); err != nil {
			return nil, fmt.Errorf("language query error (2)")
		}
		starting_languages = append(starting_languages, language)
	}
	if err := sub_race_rows.Err(); err != nil {
		return nil, fmt.Errorf("language query error (3)")
	}
	return starting_languages, nil
}

func get_starting_proficiencies(id int) ([]string, error) {
	query := `
        SELECT proficiency.name AS name
        FROM starting_proficiency
        JOIN proficiency
        ON proficiency.id = starting_proficiency.proficiency_id
        WHERE starting_proficiency.playable_race_id = ?`
	var starting_proficiencies []string
	proficiency_rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("proficiency query error (1)")
	}
	defer proficiency_rows.Close()
	for proficiency_rows.Next() {
		var proficiency string
		if err := proficiency_rows.Scan(&proficiency); err != nil {
			return nil, fmt.Errorf("proficiency query error (2)")
		}
		starting_proficiencies = append(starting_proficiencies, proficiency)
	}
	if err := proficiency_rows.Err(); err != nil {
		return nil, fmt.Errorf("proficiency query error (3)")
	}
	return starting_proficiencies, nil
}

type starting_proficiency_option_row struct {
	GroupId string
	Name    string
	Count   int
}

func get_starting_hroficiency_options(id int) ([]starting_proficiency_options, error) {
	query := ` 
        SELECT starting_proficiency_option.group_id, proficiency.name, starting_proficiency_option.count
        FROM starting_proficiency_option
        JOIN proficiency
        ON starting_proficiency_option.proficiency_id = proficiency.id
        WHERE starting_proficiency_option.playable_race_id = ?
    `
	var m map[string]starting_proficiency_options
	m = make(map[string]starting_proficiency_options)

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("proficiency option query error (1)")
	}
	for rows.Next() {
		var row starting_proficiency_option_row
		if err := rows.Scan(&row.GroupId, &row.Name, &row.Count); err != nil {
			return nil, fmt.Errorf("proficiency option query error (2)")
		}
		r := m[row.GroupId]
		r.Count = row.Count
		r.Options = append(r.Options, row.Name)
		m[row.GroupId] = r

	}
	defer rows.Close()

	var starting_proficiency_options []starting_proficiency_options
	for _, val := range m {
		starting_proficiency_options = append(starting_proficiency_options, val)
	}
	return starting_proficiency_options, nil
}

func get_starting_traits(id int) ([]string, error) {
	query := `
        SELECT trait.name AS name
        FROM starting_trait
        JOIN trait
        ON trait.id = starting_trait.trait_id
        WHERE starting_trait.playable_race_id = ?`
	var starting_traits []string
	trait_rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("trait query error (1)")
	}
	defer trait_rows.Close()
	for trait_rows.Next() {
		var trait string
		if err := trait_rows.Scan(&trait); err != nil {
			return nil, fmt.Errorf("trait query error (2)")
		}
		starting_traits = append(starting_traits, trait)
	}
	if err := trait_rows.Err(); err != nil {
		return nil, fmt.Errorf("trait query error (3)")
	}
	return starting_traits, nil
}

func get_starting_ability_bonuses(id int) ([]ability_bonus, error) {
	query := `
        SELECT name, amount
        FROM starting_ability_bonus
        JOIN ability
        ON ability.id = starting_ability_bonus.ability_id
        WHERE starting_ability_bonus.playable_race_id = ?`
	var ability_bonuses []ability_bonus
	ability_bonus_rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("ability_bonus query error (1)")
	}
	defer ability_bonus_rows.Close()
	for ability_bonus_rows.Next() {
		var ability_bonus ability_bonus
		if err := ability_bonus_rows.Scan(&ability_bonus.Ability, &ability_bonus.Bonus); err != nil {
			return nil, fmt.Errorf("ability_bonus query error (2)")
		}
		ability_bonuses = append(ability_bonuses, ability_bonus)
	}
	if err := ability_bonus_rows.Err(); err != nil {
		return nil, fmt.Errorf("ability_bonus query error (3)")
	}
	return ability_bonuses, nil
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
