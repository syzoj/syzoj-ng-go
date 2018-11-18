CREATE TABLE users (
    id bytea PRIMARY KEY NOT NULL,
    user_name varchar(64) NOT NULL,
    auth_info jsonb NOT NULL,
    can_login bool NOT NULL,
    git_password varchar,
    user_profile_info jsonb NOT NULL DEFAULT '{}'
);
CREATE UNIQUE INDEX users_user_name ON users (user_name);

CREATE TABLE groups (
    id bytea PRIMARY KEY NOT NULL,
    group_name varchar(64) NOT NULL,
    policy_info jsonb NOT NULL
);
CREATE UNIQUE INDEX groups_group_name ON groups (group_name);
CREATE TABLE group_users (
    group_id bytea NOT NULL REFERENCES groups(id),
    user_id bytea NOT NULL REFERENCES users(id),
    role_info jsonb NOT NULL,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE teams (
    id bytea PRIMARY KEY NOT NULL,
    group_id bytea NOT NULL REFERENCES groups(id),
    name varchar(64) NOT NULL
);
CREATE UNIQUE INDEX team_name ON teams (group_id, name);
CREATE TABLE team_users (
    team_id bytea NOT NULL REFERENCES teams(id),
    user_id bytea NOT NULL REFERENCES users(id),
    PRIMARY KEY (team_id, user_id)
);

CREATE TABLE problemsets (
    id bytea PRIMARY KEY NOT NULL,
    name varchar(64) NOT NULL,
    group_id bytea REFERENCES users(id)
);
CREATE TABLE problemset_teams (
    problemset_id bytea NOT NULL REFERENCES problemsets(id),
    team_id bytea NOT NULL REFERENCES teams(id),
    role int NOT NULL,
    PRIMARY KEY (problemset_id, team_id)
);
CREATE TABLE problemset_users (
    problemset_id bytea NOT NULL REFERENCES problemsets(id),
    user_id bytea NOT NULL REFERENCES users(id),
    role int NOT NULL,
    PRIMARY KEY (problemset_id, user_id)
);

CREATE TABLE problems (
    id bytea PRIMARY KEY NOT NULL,
    problemset_id bytea NOT NULL REFERENCES problemsets(id),
    type int NOT NULL,
    content jsonb NOT NULL,
    push_token varchar
);

CREATE TABLE submissions (
    id bytea PRIMARY KEY NOT NULL,
    problem_id bytea NOT NULL REFERENCES problems(id),
    status int NOT NULL,
    score float DEFAULT NULL,
    content jsonb,
    result jsonb
);

CREATE TABLE git_repos (
    id bytea PRIMARY KEY NOT NULL
);
