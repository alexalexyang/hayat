# About

This is my attempt to build a chat hotline SAAS.

# Use cases

To generalise use cases more, I refer to service side persons such as social workers, counsellors, and customer service officers as **consultants**. Their clients remain **clients**.

## Primary use case

This is primarily built as a chat hotline for counsellors/social workers and their clients.

## Extended use cases

- Customer service
- Psychological therapy?

# How to use

Clients enter via the /anteroom page. They wait for a consultant to enter.

## Clients

Clients who need to talk arrive at the /anteroom page. They fill in their details and click submit.

Upon the submission of the form, two main things happen:

- Form details are sent to /clientprofile page.
- Client is brought to /chatclient page to wait for a consultant to connect.

## Consultants

### Registration

A consultant must first register as an organisation at /register/org. The first consultant to do so will be the manager of this organisation.

The manager can then invite other consultants to join at /invite. Here, the manager can choose to make other consultants managers or staff. Managers get to delete the whole organisation, other managers, and staff.

### Dashboard page

After registration, log in at /login. Once logged in, the consultant is brought to the /dashboard page. You'll see a completely blank page. This is normal.

Whenever clients start a chat (from /anteroom), their username will appear on this page. Clicking it will open a tab with their profile inside an embedded window.

Click "Chat now." after you're done reading their profile. Their chatroom will open in the same embedded window. If you have more than one client, they will be in different tabs.


# Please note

This is a young project being built by one person with a lot of other pressing matters to deal with so there are flaws. In particular:

- Because this is a testing phase, I will regularly delete the site and its database, which means you will have to register again.
- This is not a stable product yet. It will probably take anywhere between now and 2021 to become truly usable and secure unless I get an injection of $$$.


<!-- # Installation

The program is in two parts:

## The /anteroom page

You'll probably need help for this one. You'll need to embed anteroom.html in an iframe or something anywhere on your own website. Replace "secret_key_here" in `token` with the organisation name you registered with.

This allows you to maintain your own branding (that is, after I figure out how to let you use your own CSS without introducing security risks).

## The rest of it

There's no installation for the rest of the program. Register yourself as manager of your organisation at /register/org.

Once you register and log in, you can invite other members of your staff to join at /invite. Set whether or not the invitee should be Manager or Staff. An email with a one-time-only link to a registration page will be emailed to the invitee.

Managers get to delete the entire organisation, including all staff. So, be careful. -->


# Todos

## Priority
- Add separate div so consultant can continue to refer to client's profile.
- Make "Chat now" link on /clientprofile into a button.
- Add line to tell client to wait for a consultant when they first enter.
- Have /dashboard listen on database table messages so that we can change tab name CSS on new messages.
- Fix postgres notification so it properly garbage collects dead notification instances.
- Write errors to log.
- Enforce client-side password best practices.

## UX
- Show only the first 200 lines, and then load all lines if user requests for it?

## ExplicitAuth
- Consider session cookies rather than cookies, esp. using the Gorilla package.
- Delete all members and traces of organisation upon deletion of organisation by manager.
- Double-check placeholders. `?` vs. `$1`.
- Consider JWT rather than bcrypt.
- Need separate port for domain name as opposed to localhost port?
- Delete all cookies on logout.
- Extend lifetime of all cookies, esp. consultant ones.
- Add timestamp to ExplicitAuth sessions table.
  - Delete consultant session after X hours of inactivity.
- Timestamp to invite table as well.
  - Delete invite if not used within a week.
- Delete all rooms, messages, sessions if not used after an hour from emptySince.

## External demo site
- Build small external demo site.
- JAMstack CMS?
- Host on Heroku free tier.

## Misc.
- Allow users who don't have websites to use my subdomain.
- Add slice to displayTemplate helper function for unlimited templates.


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

## Show different things on page depending on type of user.

https://www.calhoun.io/intro-to-templates-p3-functions/