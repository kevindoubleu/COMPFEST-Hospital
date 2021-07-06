# COMPFEST-Hospital
Hospital Information System - Software Engineering Academy COMPFEST Selection Task 2021

## Backend

Go with standard library

notes:
- admin panel credentials are `admin:compfesthospitaladmin`

## Frontend

[Medilab Bootstrap 5 template from BootstrapMade](https://bootstrapmade.com/medilab-free-medical-bootstrap-theme/download/), modified to meet requirements

### Sketch

Sitemap / functionalities / features

![](sketch/sketch.png)

DB

![](sketch/db.png)

#### Requirements

- [ ] authentication
  - [x] jwt

- [ ] acc roles
  - [ ] admin
    - [ ] default superuser acc
  - [ ] patient
    - [x] /register
    - [x] /login
    - [ ] /user/username
    - [x] /logout

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
- login
- jwt
- postgres db

###### References

[Implementing JWT based authentication in Golang](https://www.sohamkamani.com/golang/jwt-authentication/)

###### To fix

- more specific text field types in "patients" table
