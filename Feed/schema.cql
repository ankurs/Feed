CREATE KEYSPACE IF NOT EXISTS user WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1} AND durable_writes = true;
CREATE TABLE IF NOT EXISTS user.users (  ID uuid  PRIMARY KEY, Email text, FirstName text, LastName text, UserName text, Password text, Salt text) WITH compaction = {'class': 'LeveledCompactionStrategy'} AND default_time_to_live = 0;
CREATE INDEX IF NOT EXISTS ON user.users (email);
CREATE INDEX IF NOT EXISTS ON user.users (username);


CREATE KEYSPACE IF NOT EXISTS follow WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1} AND durable_writes = true;
CREATE TABLE IF NOT EXISTS follow.follower (  user uuid  , follower uuid, PRIMARY KEY (user, follower)) WITH CLUSTERING ORDER BY (follower DESC) AND compaction = {'class': 'LeveledCompactionStrategy'} AND default_time_to_live = 0;
CREATE TABLE IF NOT EXISTS follow.following (  user uuid  , following uuid, PRIMARY KEY (user, following)) WITH CLUSTERING ORDER BY (following DESC) AND compaction = {'class': 'LeveledCompactionStrategy'} AND default_time_to_live = 0;
