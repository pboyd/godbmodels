grail.db: schema.sql seed.sql
	sqlite3 grail.db < schema.sql && sqlite3 grail.db < seed.sql
