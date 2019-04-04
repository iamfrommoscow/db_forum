DROP TABLE IF EXISTS users CASCADE;
 DROP TABLE IF EXISTS forums CASCADE;
 DROP TABLE IF EXISTS threads CASCADE;
 DROP TABLE IF EXISTS posts CASCADE;
 DROP TABLE IF EXISTS votes CASCADE;

 DROP TABLE IF EXISTS errors CASCADE;
 CREATE EXTENSION IF NOT EXISTS CITEXT;
 CREATE TABLE users(
   nickname CITEXT UNIQUE NOT NULL PRIMARY KEY ,
   email CITEXT UNIQUE NOT NULL,
   fullname TEXT NOT NULL ,
   about TEXT
 );


 CREATE TABLE forums (
   posts INTEGER DEFAULT 0 NOT NULL ,
   slug CITEXT UNIQUE NOT NULL,
   threads INTEGER DEFAULT 0 NOT NULL ,
   title TEXT,
   "user"  CITEXT  NOT NULL REFERENCES users(nickname)
 );


 CREATE TABLE threads (
   author  CITEXT  NOT NULL REFERENCES users(nickname),
   created timestamptz(3) DEFAULT now() NOT NULL ,
   forum  CITEXT  NOT NULL REFERENCES forums(slug),
   id integer NOT NULL PRIMARY KEY ,
   message TEXT NOT NULL ,
   slug TEXT,
   title TEXT,
   votes INTEGER DEFAULT 0
 );


 CREATE TABLE posts (
   author  CITEXT  NOT NULL REFERENCES users(nickname),
   created timestamptz(3) DEFAULT now() NOT NULL ,
   forum  CITEXT  NOT NULL REFERENCES forums(slug),
   id INTEGER PRIMARY KEY ,
   isEdited BOOLEAN DEFAULT FALSE,
   message TEXT NOT NULL ,
   parent INTEGER DEFAULT 0,
   thread  INTEGER  NOT NULL REFERENCES threads(id)
 );


 CREATE TABLE votes (
   thread  INTEGER,
   nickname  CITEXT  NOT NULL REFERENCES users(nickname),
   voice INTEGER
 );
