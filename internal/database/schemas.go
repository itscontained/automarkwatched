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
    machine_identifier VARCHAR(80) PRIMARY KEY NOT NULL,
	name VARCHAR(80) NOT NULL,
	scheme VARCHAR(80) NOT NULL,
	host VARCHAR(120) NOT NULL,
	port INTEGER NOT NULL,
	local_addresses VARCHAR(120),
	address VARCHAR(120) NOT NULL,
    owner_id INTEGER NOT NULL,
	version VARCHAR(120) NOT NULL,
	enabled BOOLEAN NOT NULL,
	FOREIGN KEY (owner_id) REFERENCES users(id)
);`
	userServerSchema = `
CREATE TABLE IF NOT EXISTS user_servers (
	user_id INTEGER NOT NULL,
    server_machine_identifier VARCHAR(80) NOT NULL,
	access_token VARCHAR(80) NOT NULL UNIQUE,
	enabled BOOLEAN NOT NULL,
	UNIQUE(user_id,server_machine_identifier)
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier)
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
	server_machine_identifier VARCHAR(80) NOT NULL,
	UNIQUE(key,server_machine_identifier),
	FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier)
);`
	userLibrarySchema = `
CREATE TABLE IF NOT EXISTS user_libraries (
	user_id INTEGER NOT NULL,
	server_machine_identifier VARCHAR(80) NOT NULL,
	library_uuid VARCHAR(80) NOT NULL,
	enabled BOOLEAN NOT NULL,
	UNIQUE(user_id,server_machine_identifier,library_uuid),
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier),
	FOREIGN KEY (library_uuid) REFERENCES libraries(uuid)
);`
	seriesSchema = `
CREATE TABLE IF NOT EXISTS series (
	rating_key INTEGER PRIMARY KEY NOT NULL,
	title VARCHAR(120) NOT NULL,
	year INTEGER NOT NULL,
    studio VARCHAR(120) NOT NULL,
    content_rating VARCHAR(80) NOT NULL,
    server_machine_identifier VARCHAR(80) NOT NULL,
    library_uuid VARCHAR(80) NOT NULL,
	enabled BOOLEAN NOT NULL,
	UNIQUE(rating_key,server_machine_identifier,library_uuid),
    FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier),
    FOREIGN KEY (library_uuid) REFERENCES libraries(uuid)
);`
	userSeriesSchema = `
CREATE TABLE IF NOT EXISTS user_series (
	user_id INTEGER NOT NULL,
	server_machine_identifier VARCHAR(80) NOT NULL,
	library_uuid VARCHAR(80) NOT NULL,
	series_rating_key INTEGER NOT NULL,
	scrobble BOOLEAN NOT NULL DEFAULT 0,
	enabled BOOLEAN NOT NULL,
	UNIQUE(user_id,server_machine_identifier,library_uuid,series_rating_key),
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY (server_machine_identifier) REFERENCES servers(machine_identifier),
	FOREIGN KEY (library_uuid) REFERENCES libraries(uuid),
	FOREIGN KEY(series_rating_key) REFERENCES series(rating_key)
);`
)
