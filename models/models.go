package models

import (
	"database/sql"
	"log"

	"github.com/alexalexyang/hayat/config"
	_ "github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DBSetup() {
	db, err := sql.Open(config.Driver, config.DBconfig)
	check(err)
	defer db.Close()

	// REMEMBER TO MAKE TOKEN UNIQUE.
	// Rooms table. roomid is for clientlist. token is to identify customer.
	statement := `CREATE TABLE IF NOT EXISTS rooms (
				timestamptz TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				roomid TEXT UNIQUE,
				organisation TEXT,
				sessioncookie TEXT UNIQUE,
				beingserved bool
				);`
	_, err = db.Exec(statement)
	check(err)

	// Anteroom table.
	statement = `CREATE TABLE IF NOT EXISTS clientprofiles (
		sessioncookie TEXT UNIQUE,
		username TEXT,
		age TEXT,
		gender TEXT,
		issues TEXT
		);`
	_, err = db.Exec(statement)
	check(err)

	// Customer table.
	statement = `CREATE TABLE IF NOT EXISTS consultants (
		firstname TEXT,
		lastname TEXT,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password TEXT UNIQUE,
		organisation TEXT,
		sessioncookie TEXT UNIQUE
		);`
	_, err = db.Exec(statement)
	check(err)

	// Organisation table.
	statement = `CREATE TABLE IF NOT EXISTS organisations (
		id SERIAL UNIQUE,
		orgname TEXT,
		phone TEXT UNIQUE,
		email TEXT UNIQUE,
		managername TEXT UNIQUE,
		password TEXT UNIQUE,
		organisation TEXT
		);`
	_, err = db.Exec(statement)
	check(err)

	statement = `CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
    
    BEGIN
    
	PERFORM (
		with payload(roomid, beingserved) as
		(
		  select NEW.roomid,
				 NEW.beingserved
		)
		select pg_notify('events', row_to_json(payload)::text)
		  from payload
	 );
        
        RETURN NULL; 
    END;
    
	$$ LANGUAGE plpgsql;`

	_, err = db.Exec(statement)
	check(err)

	statement = `DROP TRIGGER IF EXISTS products_notify_event ON rooms;
				CREATE TRIGGER products_notify_event
				AFTER INSERT OR UPDATE OR DELETE ON rooms
				FOR EACH ROW EXECUTE PROCEDURE notify_event();`
	_, err = db.Exec(statement)
	check(err)
}
