package storage

var sqlInitDatabase = `
	DROP TABLE IF EXISTS urls;

    CREATE TABLE urls (
		id INT PRIMARY KEY, 
		url VARCHAR(255) NOT NULL, 
		status varchar(50) NOT NULL DEFAULT 'NEW',
		http_code INT
	);
	
	INSERT INTO urls (id, url, status, http_code) 
	VALUES 
		(1, 'https://proxify.io', 'NEW', NULL),
		(2, 'https://reddit.com', 'NEW', NULL),
		(3, 'http://go.org', 'NEW', NULL),
        (4, 'https://I am a malformed url', 'NEW', NULL),
		(5, 'https://google.com', 'DONE', 200),
		(6, 'https://google.com', 'NEW', NULL);`

var sqlCheckUrlTableExists = `
	SELECT name FROM sqlite_master WHERE type='table' AND name='urls';`

var sqlSelectNewURLs = `
	SELECT id, url, status, http_code 
	FROM urls
	WHERE status = 'NEW'`

var sqlSelectUrlById = `
	SELECT id, url, status, http_code
	FROM urls
	WHERE id = ?`

var sqlUpdateUrlById = `
	UPDATE urls
	SET status = ?, http_code = ?
	WHERE id = ?`

var sqlUpdateStatusById = `
	UPDATE urls
	SET status = ?
	WHERE id = ?`
