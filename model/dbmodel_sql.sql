CREATE TABLE user (
  id VARCHAR(16) PRIMARY KEY,
  user_name varchar(64) UNIQUE,
  auth BLOB
);

CREATE TABLE device (
  id VARCHAR(16) PRIMARY KEY,
  user VARCHAR(255),
  info BLOB
);

CREATE TABLE problem (
  id VARCHAR(16) PRIMARY KEY,
  title VARCHAR(255)
);

CREATE TABLE problem_source (
  id VARCHAR(16) PRIMARY KEY,
  data BLOB
);

CREATE TABLE problem_judger (
  id VARCHAR(16) PRIMARY KEY,
  problem VARCHAR(255),
  user VARCHAR(255),
  type VARCHAR(255),
  data BLOB
);

CREATE TABLE problem_statement (
  id VARCHAR(16) PRIMARY KEY,
  problem VARCHAR(255),
  user VARCHAR(255),
  data BLOB
);

CREATE TABLE submission (
  id VARCHAR(16) PRIMARY KEY,
  problem_judger VARCHAR(255),
  user VARCHAR(255),
  data BLOB
);


