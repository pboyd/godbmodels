package mapper

import (
	"context"
	"testing"

	"github.com/pboyd/godbmodels/common"
	"github.com/stretchr/testify/assert"
)

func TestCharacters(t *testing.T) {
	assert := assert.New(t)

	cs := NewCharacterStore(common.TestDB(t))
	c := &Character{
		Name:    "Sir Not-Appearing-in-this-Film",
		ActorID: 1,
	}

	// Create
	err := cs.Store(context.Background(), c)
	if !assert.NoError(err) {
		return
	}
	assert.NotZero(c.ID)

	// Read
	c2, err := cs.Get(context.Background(), c.ID)
	if !assert.NoError(err) {
		return
	}
	assert.Equal(c, c2)

	// Update
	c.Name = "Sir Maybe-Appearing-in-this-Film"
	c.ActorID = 2
	err = cs.Store(context.Background(), c)
	if !assert.NoError(err) {
		return
	}
	c2, err = cs.Get(context.Background(), c.ID)
	if !assert.NoError(err) {
		return
	}
	assert.Equal(c, c2)

	// Delete
	err = cs.Delete(context.Background(), c.ID)
	if !assert.NoError(err) {
		return
	}
	c2, err = cs.Get(context.Background(), c.ID)
	if assert.NoError(err) {
		assert.Nil(c2)
	}

	// Delete again
	err = cs.Delete(context.Background(), c.ID)
	assert.Equal(ErrNotFound, err)
}

func TestListCharacters(t *testing.T) {
	cases := map[string]struct {
		filters       *CharacterFilters
		expected      int
		expectedNames []string
	}{
		"All": {
			filters:  &CharacterFilters{},
			expected: 81,
		},
		"Eric Idle": {
			filters: &CharacterFilters{
				ActorID: 3,
			},
			expected: 8,
			expectedNames: []string{
				"Brother Maynard",
				"Concorde",
				"Dead Collector",
				"First Swamp Castle Guard",
				"Knight of Camelot",
				"Peasant 1",
				"Roger the Shrubber",
				"Sir Robin the Not-Quite-So-Brave-as-Sir Launcelot",
			},
		},
		"Sandy": {
			filters: &CharacterFilters{
				ActorName: "Sandy",
			},
			expected: 6,
			expectedNames: []string{
				"Girl in Castle Anthrax #3",
				"Knight in Battle",
				"Knight of Ni",
				"Monk",
				"Musician at Wedding",
				"Villager at Witch Burning",
			},
		},
		"Brother Maynard": {
			filters: &CharacterFilters{
				Name: "Brother Maynard",
			},
			expected: 2,
			expectedNames: []string{
				"Brother Maynard",
				"Brother Maynard's Brother",
			},
		},
		"The violence inherent in the system": {
			filters: &CharacterFilters{
				SceneNumber: 3,
			},
			expected: 4,
			expectedNames: []string{
				"Dennis",
				"Dennis's Mother",
				"King Arthur",
				"Patsy",
			},
		},
	}

	assert := assert.New(t)
	cs := NewCharacterStore(common.TestDB(t))

	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			ctx := context.Background()
			characters, err := cs.List(ctx, c.filters)
			if !assert.NoError(err) {
				return
			}
			assert.Len(characters, c.expected)
			if c.expectedNames != nil {
				names := make([]string, 0, len(characters))
				for _, c := range characters {
					names = append(names, c.Name)
				}
				assert.ElementsMatch(c.expectedNames, names)
			}
		})
	}
}
