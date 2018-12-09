CREATE KEYSPACE IF NOT EXISTS user WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1} AND durable_writes = true;
CREATE TABLE IF NOT EXISTS user.users (  ID uuid  PRIMARY KEY, Email text, FirstName text, LastName text, UserName text, Password text, Salt text) WITH compaction = {'class': 'LeveledCompactionStrategy'} AND default_time_to_live = 0;
CREATE INDEX IF NOT EXISTS ON user.users (email);
CREATE INDEX IF NOT EXISTS ON user.users (username);
