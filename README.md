# trsm
my idea of time restriction based social media
Currently the project consists only of it's account management component.

Structure:
dbinteractions/dbUtilStructs.go - structures required to retrive data from database
dbinteractions/dbinteractions.go - functions executing actual queries
main/handlers.go - handlers
main/main.go - actually runs the server
main/session_manager.go - defines session
main/webosckets.go - websocket handlers
pfp_storeage - stores user's profile pictures, for the sake of simplicity
static/templates - stores templates for pages
structures - stores structures required to run the app
utils - stores utility functions

Database layout:
                    Table "public.users_pfps"
 Column  |         Type          | Collation | Nullable | Default
---------+-----------------------+-----------+----------+---------
 pfp     | character varying(16) |           |          |
 user_id | integer               |           |          |
 
                                             Table "public.users"
    Column     |            Type             | Collation | Nullable |                Default
---------------+-----------------------------+-----------+----------+----------------------------------------
 user_id       | integer                     |           | not null | nextval('users_user_id_seq'::regclass)
 user_login    | character varying(50)       |           | not null |
 user_password | character varying(64)       |           | not null |
 created_on    | timestamp without time zone |           | not null |
 last_login    | timestamp without time zone |           |          |
 tag_whitelist | character varying(500)      |           |          |
 tag_blacklist | character varying(500)      |           |          |
 admin         | boolean                     |           | not null |
 rep           | integer                     |           | not null |
