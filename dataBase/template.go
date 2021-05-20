package dataBase

var (
	loggerTableCreate = ` create table if not exists logger (
	    tm 		timestamp with time zone NOT NULL,
		login 	text 	NOT NULL,
		key   	text NOT NULL,
    	txt 	text 	NOT NULL);`
	phoneTableCreate = `	create table if not exists phones (
		login 		text NOT NULL primary key,
		phone 		json		);`
	areaTableCreate = `create table if not exists area (
		area 		integer NOT NULL primary key,
		namearea 	text  NOT NULL	);`
	crossesTableCreate = ` create table if not exists crosses (
		key   	text NOT NULL primary key,
    	crossT   json 			);`
)
