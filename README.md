# About

This is my attempt to build a chat SAAS.


# Todos

## Specific to app
- Add button to chat tab to close websocket and remove itself.
- Prioritise database over in-memory storage for persistence.
- Reload all messages and chats in case user accidentally reloads page.
- Save messages so they can be reloaded if connection breaks for whatever reason.
- Delete saved messages when cleaning up rooms.

## ExplicitAuth
- Delete all members and traces of organisation upon deletion of organisation by manager.
- Double-check placeholders. `?` vs. `$1`.


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