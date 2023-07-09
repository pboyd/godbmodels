package builder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

// ErrNotFound is returned when updating or deleting a character that does not
// exist in the database.
var ErrNotFound = errors.New("not found")

// Character is one character from the database.
type Character struct {
	ID      int64
	ActorID int64
	Name    string
}

// CharacterStore loads and updates characters in the database.
type CharacterStore struct {
	db *sql.DB
}

// NewCharacterStore creates a new CharacterStore.
func NewCharacterStore(db *sql.DB) *CharacterStore {
	return &CharacterStore{db: db}
}

// Get loads a character from the database by ID.
//
// If no character is found, Get returns a nil Character and no error.
func (cs *CharacterStore) Get(ctx context.Context, id int64) (*Character, error) {
	var c Character
	err := squirrel.
		Select("id", "actor_id", "name").
		From("characters").
		Where("id = ?", id).
		RunWith(cs.db).
		QueryRowContext(ctx).
		Scan(&c.ID, &c.ActorID, &c.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &c, err
}

// Store saves a character to the database. If the character has an ID, it will
// be updated. Otherwise, it will be inserted and the ID will be set.
//
// If the character has an ID and it does not exist in the database, Store
// returns ErrNotFound.
func (cs *CharacterStore) Store(ctx context.Context, c *Character) error {
	if c.ID == 0 {
		return cs.insert(ctx, c)
	}

	return cs.update(ctx, c)
}

func (cs *CharacterStore) insert(ctx context.Context, c *Character) error {
	return squirrel.
		Insert("characters").
		Columns("actor_id", "name").
		Values(c.ActorID, c.Name).
		Suffix("RETURNING id").
		RunWith(cs.db).
		QueryRowContext(ctx).
		Scan(&c.ID)
}

func (cs *CharacterStore) update(ctx context.Context, c *Character) error {
	res, err := squirrel.
		Update("characters").
		Set("actor_id", c.ActorID).
		Set("name", c.Name).
		Where("id = ?", c.ID).
		RunWith(cs.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("update character: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete removes a character from the database.
//
// If the character does not exist in the database, Delete returns ErrNotFound.
func (cs *CharacterStore) Delete(ctx context.Context, id int64) error {
	res, err := squirrel.
		Delete("characters").
		Where("id = ?", id).
		RunWith(cs.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("delete character: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
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

// List searches for characters in the database.
//
// If filters is nil, all characters are returned. Otherwise, the results are
// filtered by the criteria in filters.
func (cs *CharacterStore) List(ctx context.Context, filters *CharacterFilters) ([]*Character, error) {
	q := squirrel.
		Select("c.id", "c.actor_id", "c.name").
		From("characters c").
		RunWith(cs.db)

	if filters != nil {
		if filters.ActorID != 0 {
			q = q.Where("actor_id = ?", filters.ActorID)
		} else if filters.ActorName != "" {
			q = q.
				Join("actors a ON a.id = c.actor_id").
				Where("LOWER(a.name) LIKE ?", "%"+strings.ToLower(filters.ActorName)+"%")
		}

		if filters.Name != "" {
			q = q.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filters.Name)+"%")
		}

		if filters.SceneNumber != 0 {
			q = q.
				Join("scene_characters sc ON sc.character_id = c.id").
				Where("sc.scene_id = ?", filters.SceneNumber)
		}
	}

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("list characters: %w", err)
	}
	defer rows.Close()

	var characters []*Character
	for rows.Next() {
		var c Character
		err := rows.Scan(&c.ID, &c.ActorID, &c.Name)
		if err != nil {
			return nil, fmt.Errorf("list characters: %w", err)
		}

		characters = append(characters, &c)
	}

	return characters, nil
}
