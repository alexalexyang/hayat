# About

This is my attempt to build a chat SAAS.

This is primarily built as a chat hotline for counsellors and their clients. To generalise the use case a bit more, I refer to counsellors as **consultants**. Their clients remain **clients**.

The following sections are mainly for consultants.

# How to use

## Installation

The program is in two parts:

### The /anteroom page

You'll probably need help for this one. You'll need to embed anteroom.html in an iframe or something anywhere on your own website. Replace "secret_key_here" in `token` with the organisation name you registered with.

This allows you to maintain your own branding (that is, after I figure out how to let you use your own CSS without introducing security risks).

### The rest of it

There's no installation for the rest of the program. Register yourself as manager of your organisation at /register/org.

Once you register and log in, you can invite other members of your staff to join at /invite. Set whether or not the invitee should be Manager or Staff. An email with a one-time-only link to a registration page will be emailed to the invitee.

Managers get to delete the entire organisation, including all staff. So, be careful.

## Clientlist

Once you log in, you'll be brought to /clientlist. Whenever clients start a chat, their username will appear on this page. Clicking it will open their profile inside an embedded window.

Click "Chat now." after you're done reading their profile. Their chatroom will open in the same embedded window.


# Please note

This is a young project being built by one person with a lot of other pressing matters to deal with so there are flaws. In particular:

- Please set a strong and unique password for yourself because I haven't enforced password best practices yet.
- If a client disconnects from a chat before a consultant enters, the chat will still show on the clientlist until a consultant clicks into it. So, consultants may sometimes see an empty room.
- I didn't think to add a separate section on the page for the client profile. I'll figure this out later.
- If you accidentally reload your page, all chats are lost. Another thing for me to think about.
- This is not a stable product yet. It will probably take anywhere between now and 2021 to become truly usable and secure unless I get an injection of $$$.


# Todos

## UX
- Add navbar/links to essential pages.
- Add button to remove chat tab and close websocket.
- Add separate div so consultant can continue to refer to client's profile.

## Clean-up
- Delete saved messages when cleaning up rooms.
- Add timestamp to ExplicitAuth sessions table.
  - Delete consultant session after X hours of inactivity.

## Persistence
- Prioritise database over in-memory storage for persistence.
- Reload all messages and chats in case user accidentally reloads page.
- Save messages so they can be reloaded if connection breaks for whatever reason.

## Security
- Consider session cookies rather than cookies, esp. using the Gorilla package.
- Delete all members and traces of organisation upon deletion of organisation by manager.
- Double-check placeholders. `?` vs. `$1`.
- Consider JWT rather than bcrypt.
- Enforce client-side password best practices.

## Misc.
- Allow users who don't have websites to use my subdomain.


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