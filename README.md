# COMPFEST-Hospital
Hospital Information System - Software Engineering Academy COMPFEST Selection Task 2021

## Backend

Go with standard library

notes:
- admin panel credentials are `admin:compfesthospitaladmin`

## Frontend

[Medilab Bootstrap 5 template from BootstrapMade](https://bootstrapmade.com/medilab-free-medical-bootstrap-theme/download/)

### Sketch

![](sketch/sketch.png)

#### requirements

- [ ] authentication
  - [ ] jwt

- [ ] acc roles
  - [ ] admin
    - [ ] default superuser acc
  - [ ] patient
    - [x] /register
    - [ ] /login
    - [ ] /user/username
    - [ ] /logout

- [ ] admin funcs
  - [ ] doctor appointment
    - [ ] c
    - [ ] r
    - [ ] u
    - [ ] d

- [ ] patient funcs
  - [ ] see list of appointments (all?)
  - [ ] apply for appointment
    - [ ] check if fully booked registrant (how many?)
  - [ ] cancel an appointment

##### timeline

july 5
- start proj
- sketch sitemap / functionalities
- object definitions
- in memory db
- frontend
- register

july 6
- logout
