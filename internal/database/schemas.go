package database

const (
	appSchema = `
CREATE TABLE IF NOT EXISTS app (
	identifier VARCHAR(80) PRIMARY KEY,
	owner_id INTEGER DEFAULT 0,
	version VARCHAR(80)
);`
	userSchema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	uuid VARCHAR(80) NOT NULL,
	username VARCHAR(80) NOT NULL,
    email VARCHAR(120) NOT NULL,
	thumb VARCHAR(120) NOT NULL,
	owner BOOLEAN NOT NULL,
	enabled BOOLEAN NOT NULL
);`
	serverSchema = `
CREATE TABLE IF NOT EXISTS servers (
    user_id INTEGER NOT NULL,
    machine_identifier VARCHAR(80) PRIMARY KEY NOT NULL,
	name VARCHAR(80) NOT NULL,
	scheme VARCHAR(80) NOT NULL,
	host VARCHAR(120) NOT NULL,
	port INTEGER NOT NULL,
	local_addresses VARCHAR(120),
	address VARCHAR(120) NOT NULL,
    owner_id INTEGER NOT NULL,
	version VARCHAR(120) NOT NULL,
    access_token VARCHAR(80) NOT NULL UNIQUE,
	enabled BOOLEAN NOT NULL,
    UNIQUE(user_id,owner_id,machine_identifier),
	FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);`
	librarySchema = `
CREATE TABLE IF NOT EXISTS libraries (
	uuid VARCHAR(80) PRIMARY KEY NOT NULL,
	title VARCHAR(120) NOT NULL,
	type VARCHAR(80) NOT NULL,
	key INTEGER NOT NULL,
    agent VARCHAR(80) NOT NULL,
	scanner VARCHAR(80) NOT NULL,
	enabled BOOLEAN NOT NULL,
    user_id INTEGER NOT NULL,
	server_machine_identifier VARCHAR(80) NOT NULL,
	UNIQUE(uuid,user_id,server_machine_identifier),
    FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier)
);`
	seriesSchema = `
CREATE TABLE IF NOT EXISTS series (
	rating_key INTEGER PRIMARY KEY NOT NULL,
	title VARCHAR(120) NOT NULL,
	year INTEGER NOT NULL,
    studio VARCHAR(120) NOT NULL,
    content_rating VARCHAR(80),
	enabled BOOLEAN NOT NULL,
    scrobble BOOLEAN NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL,
    library_uuid VARCHAR(80) NOT NULL,
    server_machine_identifier VARCHAR(80) NOT NULL,
	UNIQUE(rating_key,user_id,server_machine_identifier,library_uuid),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier),
    FOREIGN KEY (library_uuid) REFERENCES libraries(uuid)
);`
)
