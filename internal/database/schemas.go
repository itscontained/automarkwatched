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
    machine_identifier VARCHAR(80) NOT NULL,
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
	UNIQUE(user_id, machine_identifier)
);`
	serverTrigger = `
CREATE TRIGGER IF NOT EXISTS server_check
    BEFORE INSERT ON servers
    WHEN (NEW.user_id NOT IN (SELECT id FROM users))
BEGIN
    SELECT RAISE(FAIL, 'invalid user/server');
END;`
	librarySchema = `
CREATE TABLE IF NOT EXISTS libraries (
	uuid VARCHAR(80) NOT NULL,
	title VARCHAR(120) NOT NULL,
	type VARCHAR(80) NOT NULL,
	key INTEGER NOT NULL,
    agent VARCHAR(80) NOT NULL,
	scanner VARCHAR(80) NOT NULL,
	enabled BOOLEAN NOT NULL,
    user_id INTEGER NOT NULL,
	server_machine_identifier VARCHAR(80) NOT NULL
);`
	libraryTrigger = `
CREATE TRIGGER IF NOT EXISTS library_check
    BEFORE INSERT ON libraries
	WHEN (NEW.user_id NOT IN (SELECT id FROM users)
		  OR
          NEW.server_machine_identifier NOT IN (SELECT machine_identifier FROM servers))
BEGIN
    SELECT RAISE(FAIL, 'invalid user/server');
END;`
	seriesSchema = `
CREATE TABLE IF NOT EXISTS series (
	rating_key INTEGER NOT NULL,
	title VARCHAR(120) NOT NULL,
	year INTEGER NOT NULL,
    studio VARCHAR(120) NOT NULL,
    content_rating VARCHAR(80),
	guid VARCHAR(120) NOT NULL default unknown,
	enabled BOOLEAN NOT NULL,
    scrobble BOOLEAN NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL,
    library_uuid VARCHAR(80) NOT NULL,
    server_machine_identifier VARCHAR(80) NOT NULL
);`
	seriesTrigger = `
CREATE TRIGGER IF NOT EXISTS series_check
    BEFORE INSERT ON series
	WHEN (NEW.user_id NOT IN (SELECT id FROM users)
		  OR
          NEW.server_machine_identifier NOT IN (SELECT machine_identifier FROM servers)
          OR
		  NEW.library_uuid NOT IN (SELECT uuid FROM libraries))
BEGIN
    SELECT RAISE(FAIL, 'invalid user/server/library');
END;`
)
