CREATE TABLE `users` (
	`id` BIGINT PRIMARY KEY AUTO_INCREMENT,
	`uid` VARCHAR(255),
	`username` VARCHAR(255),
	`email` VARCHAR(255),
	`password` BLOB,
	`register_time` DATETIME,
	`problem_count` BIGINT NOT NULL DEFAULT 0,
	UNIQUE (`uid`),
	UNIQUE (`username`)
);

CREATE TABLE `problems` (
	`id` BIGINT PRIMARY KEY AUTO_INCREMENT,
	`uid` VARCHAR(255),
	`body` BLOB,
	`owner_uid` VARCHAR(255),
	UNIQUE (`uid`)
);

CREATE TABLE `testdata_meta` (
	`id` BIGINT PRIMARY KEY AUTO_INCREMENT,
	`problem_id` VARCHAR(255),
	`object_name` VARCHAR(255),
	`sha256` VARCHAR(255),
	INDEX (`problem_id`, `object_name`),
	UNIQUE (`object_name`)
);
