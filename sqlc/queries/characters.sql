-- name: GetCharacter :one
-- GetCharacter loads a character from the database by ID.
SELECT * FROM characters WHERE id = ?;

-- name: insertCharacter :one
-- insertCharacter creates a new character record.
INSERT INTO characters (actor_id, name) VALUES (?, ?) RETURNING id;

-- name: updateCharacter :exec
-- updateCharacter updates a character's information.
UPDATE characters SET actor_id = ?, name = ? WHERE id = ?;

-- name: DeleteCharacter :exec
-- DeleteCharacter removes a character from the database.
DELETE FROM characters WHERE id = ?;

-- name: listAllCharacters :many
-- listAllCharacters returns all characters.
SELECT * FROM characters;

-- name: listCharactersByActor :many
-- listCharactersByActor returns all characters played a given actor.
SELECT * FROM characters WHERE actor_id = ?;

-- name: listCharactersByActorName :many
-- listCharactersByActorName returns all characters played by an actor with a
-- name matching the given name.
SELECT c.* FROM characters c JOIN actors a ON c.actor_id = a.id WHERE LOWER(a.name) LIKE '%' || LOWER(?) || '%';

-- name: listCharactersByName :many
-- listCharactersByName returns all characters with a name matching the given
-- name.
SELECT * FROM characters WHERE LOWER(name) LIKE '%' || LOWER(?) || '%';

-- name: listCharactersByScene :many
-- listCharactersByScene returns all characters in a given scene.
SELECT c.* FROM characters c JOIN scene_characters sc ON c.id = sc.character_id WHERE sc.scene_id = ?;
