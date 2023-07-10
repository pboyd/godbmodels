package orm

// Actor represents the actor that plays a character.
type Actor struct {
	ID   int64  `gorm:"id,primary_key"`
	Name string `gorm:"name"`
}
