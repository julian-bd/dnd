SET @abilityId = (SELECT id FROM ability WHERE ability.name='CHA' LIMIT 1);
INSERT INTO starting_ability_bonus (playable_race_id, ability_id, amount)
VALUES (4, @abilityId, 2)
