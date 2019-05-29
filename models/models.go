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
				websocket TEXT UNIQUE,
				beingserved bool
				);`
	_, err = db.Exec(statement)
	check(err)

	// Anteroom table.
	statement = `CREATE TABLE IF NOT EXISTS anteroom (
		sessioncookie TEXT UNIQUE,
		username TEXT UNIQUE,
		age TEXT UNIQUE,
		gender TEXT UNIQUE,
		issues TEXT UNIQUE
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
    
        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);
        
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;`
	_, err = db.Exec(statement)
	check(err)

	statement = `CREATE TRIGGER products_notify_event
	AFTER INSERT OR UPDATE OR DELETE ON rooms
		FOR EACH ROW EXECUTE PROCEDURE notify_event();`
	_, err = db.Exec(statement)
	check(err)
}
