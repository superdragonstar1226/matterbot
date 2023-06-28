--
-- create databases and DB users for local purposes
--
CREATE DATABASE mattermost;
CREATE USER mattermost WITH PASSWORD 'secret';
GRANT ALL ON DATABASE mattermost TO mattermost;
