# About

This is my attempt to build a chat SAAS.


# Todos

- Delete empty rooms after an hour.
- Rebuild RoomRegistry in case program crashes.
- Prioritise database over in-memory storage for persistence.
  - esp. for RoomRegistry.
- Delete all members and traces of organisation upon deletion of organisation by manager.
- Save messages so they can be reloaded if connection breaks for whatever reason.
- Delete saved messages when cleaning up rooms.



# Learning sources


## Listen and notify on row with postgres

https://godoc.org/github.com/lib/pq/example/listen

http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/

https://dzone.com/articles/notify-events-from-postgresql-to-external-listener

## Listen and notify on column with postgres

https://tapoueh.org/blog/2018/07/postgresql-listen/notify/

## Access struct in map without copying

https://stackoverflow.com/questions/17438253/access-struct-in-map-without-copying

## How to have iframes in tabs

https://howto.caspio.com/tech-tips-and-articles/advanced-customizations/create-embeddable-tabbed-interface/