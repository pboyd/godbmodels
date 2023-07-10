package sqlc

import (
	"context"
)

// StoreCharacter saves a character to the database. If the character has an
// ID, it will be updated. Otherwise, it will be inserted and the ID will be
// set.
func (q *Queries) StoreCharacter(ctx context.Context, c *Character) error {
	if c.ID == 0 {
		id, err := q.insertCharacter(ctx, insertCharacterParams{
			ActorID: c.ActorID,
			Name:    c.Name,
		})
		if err != nil {
			return err
		}

		c.ID = id
		return nil
	}

	return q.updateCharacter(ctx, updateCharacterParams{
		ID:      c.ID,
		ActorID: c.ActorID,
		Name:    c.Name,
	})
}

// CharacterFilters are used to filter the results of a List query.
type CharacterFilters struct {
	// ActorID matches on the actor's ID.
	ActorID int64

	// ActorName does a case-insensitive partial match on the actor name.
	ActorName string

	// Name does a case-insensitive partial match on the character name.
	Name string

	// SceneNumber filters by the scene that the character appears in.
	SceneNumber int64
}

// ListCharacters searches for characters in the database.
//
// If filters is nil, all characters are returned. Otherwise, the results are
// filtered by the criteria in filters. Only one filter option can be used at
// a time.
func (q *Queries) ListCharacters(ctx context.Context, filters *CharacterFilters) ([]Character, error) {
	switch {
	case filters.ActorID != 0:
		return q.listCharactersByActor(ctx, filters.ActorID)
	case filters.ActorName != "":
		return q.listCharactersByActorName(ctx, filters.ActorName)
	case filters.Name != "":
		return q.listCharactersByName(ctx, filters.Name)
	case filters.SceneNumber != 0:
		return q.listCharactersByScene(ctx, filters.SceneNumber)
	default:
		return q.listAllCharacters(ctx)
	}
}
