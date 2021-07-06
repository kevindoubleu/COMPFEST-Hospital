# COMPFEST-Hospital
Hospital Information System - Software Engineering Academy COMPFEST Selection Task 2021

## Backend

Go with standard library

notes:
- admin panel credentials are `admin:compfesthospitaladmin`

## Frontend

[Medilab Bootstrap 5 template from BootstrapMade](https://bootstrapmade.com/medilab-free-medical-bootstrap-theme/download/), modified to meet requirements

### Sketch

![](sketch/sketch.png)

#### Requirements

- [ ] authentication
  - [x] jwt

- [ ] acc roles
  - [x] admin
    - [x] default superuser acc
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

###### References

[Implementing JWT based authentication in Golang](https://www.sohamkamani.com/golang/jwt-authentication/)
