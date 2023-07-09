CREATE TABLE actors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    actor_id INTEGER NOT NULL,
    FOREIGN KEY (actor_id) REFERENCES actors (id)
);

CREATE TABLE scenes (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE scene_characters (
    scene_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    PRIMARY KEY (scene_id, character_id),
    FOREIGN KEY (scene_id) REFERENCES scenes (id),
    FOREIGN KEY (character_id) REFERENCES characters (id)
);

CREATE TABLE quotes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    scene_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    FOREIGN KEY (character_id) REFERENCES characters (id),
    FOREIGN KEY (scene_id) REFERENCES scenes (id)
);
