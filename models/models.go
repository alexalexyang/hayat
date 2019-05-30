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
				token TEXT,
				sessioncookie TEXT UNIQUE,
				websocket TEXT UNIQUE,
				beingserved bool
				);`
	_, err = db.Exec(statement)
	check(err)

	// Anteroom table.
	statement = `CREATE TABLE IF NOT EXISTS anteroom (
		sessioncookie TEXT UNIQUE,
		username TEXT,
		age TEXT,
		gender TEXT,
		issues TEXT
		);`
	_, err = db.Exec(statement)
	check(err)

	// Customer table.
	statement = `CREATE TABLE IF NOT EXISTS customers (
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

	statement = `CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
    
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        PERFORM pg_notify('events',notification::text);
        
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
