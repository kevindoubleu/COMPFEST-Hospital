# checklist

- [x] authentication
  - [x] jwt

- [x] acc roles
  - [x] admin
    - [x] default superuser acc
  - [x] patient
    - [x] /register
    - [x] /login
    - [x] /profile
    - [x] /logout

- [x] admin funcs
  - [x] doctor appointment crud
    - [x] c
    - [x] r
      - [x] r registrants per appointment
    - [x] u
    - [x] d

- [x] patient funcs
  - [x] see list of appointments (all?)
  - [x] apply for appointment
    - [x] check if fully booked registrant (how many?)
  - [x] cancel an appointment

###### Todo

- ~~change "patients" table into "users" table with an additional "role" field~~
  - ~~this adds the abilty to manually add more admin role accounts~~
- ~~change "users" primary key to "username" instead of unnecessary field "id"~~
- more specific text field types in "users" table

## timeline

july 5
- start proj
- sketch sitemap / functionalities
- object definitions
- in memory db
- frontend
- register

july 6
- logout
- login
- jwt
- sketch + implement postgres db
- admin frontend + read

july 7
- admin create, update, delete
- patient read, cancel, apply
- update db schema

at this point, minimum requirements are satisfied

july 8
- user profile frontend + read, update (not a req but why not)
- toasts
- heroku
- documentation

at this point, the project is basically done \
gonna add more stuff for extra points, will keep continuing

july 9
- admin read on patient data
- admin create patient account from `/administration/patients`
- admin update on patient data
- admin update on patient password
- user profile delete
- admins can cancel patient's applied appointment in appointment listing
