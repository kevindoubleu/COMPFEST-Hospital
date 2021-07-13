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

july 10
- basic middlewares
- ajax for patient cancel
- ajax for patient apply
- modal for ajax reply messages
- ajax for admin kicking patients from appointments

july 11
- add middlewares
- appointments can have images
  - patients can view appointment images

july 12
- admin create appointment supports image uploads
- admin crud appointment images
  - admin can (read) view appointment images (uses same endpoint as patients)
  - admin can (update) add images on existing appointment
  - admin can delete images (per appointment)
- appointment image in patient's "my appointment" section

july 13
- appointments can have comments
- patients and admins can comment on any appointment
  - uses JSON in both req + resp!

july 14
- admins can delete comments (per appointment)

planned features
- patients can upload a profile picture
  - has a default generic picture
  - shows up in comments
- ~~admins can assign patients to appointments~~ no time
- ~~admins can cancel patient's appointment from `/administration/patients`~~ no time
- ~~create new role under admin: moderator, for each appointment~~ violates requirement 1: "There are two types of account roles: Administrator and Patient."
