package vanilla

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	row := cs.db.QueryRowContext(ctx, `SELECT id, actor_id, name FROM characters WHERE id = $1`, id)

	var c Character
	err := row.Scan(&c.ID, &c.ActorID, &c.Name)
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
	row := cs.db.QueryRowContext(ctx, `INSERT INTO characters (actor_id, name) VALUES ($1, $2) RETURNING id`, c.ActorID, c.Name)
	return row.Scan(&c.ID)
}

func (cs *CharacterStore) update(ctx context.Context, c *Character) error {
	res, err := cs.db.ExecContext(ctx, `UPDATE characters SET actor_id = $1, name = $2 WHERE id = $3`, c.ActorID, c.Name, c.ID)
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
	res, err := cs.db.ExecContext(ctx, `DELETE FROM characters WHERE id = $1`, id)
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
	var args []interface{}
	query := "SELECT c.id, c.actor_id, c.name FROM characters c"
	joins := []string{}
	where := []string{}

	if filters != nil {
		if filters.ActorID != 0 {
			where = append(where, "c.actor_id = ?")
			args = append(args, filters.ActorID)
		} else if filters.ActorName != "" {
			joins = append(joins, "JOIN actors a ON a.id = c.actor_id")
			where = append(where, "LOWER(a.name) LIKE ?")
			args = append(args, "%"+strings.ToLower(filters.ActorName)+"%")
		}

		if filters.Name != "" {
			where = append(where, "LOWER(c.name) LIKE ?")
			args = append(args, "%"+strings.ToLower(filters.Name)+"%")
		}

		if filters.SceneNumber != 0 {
			joins = append(joins, "JOIN scene_characters sc ON sc.character_id = c.id")
			where = append(where, "sc.scene_id = ?")
			args = append(args, filters.SceneNumber)
		}
	}

	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	rows, err := cs.db.QueryContext(ctx, query, args...)
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
