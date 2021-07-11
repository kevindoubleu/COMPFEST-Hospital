# COMPFEST-Hospital
Hospital Information System - Software Engineering Academy COMPFEST Selection Task 2021

  - [Backend](#backend)
  - [Frontend](#frontend)
  - [Assumptions (Notes for COMPFEST)](#assumptions-notes-for-compfest)
  - [Documentation](#documentation)

## Backend

Go, PostgreSQL, Heroku ([compfesthospital.herokuapp.com](https://compfesthospital.herokuapp.com/))

Dependencies:
- [github.com/dgrijalva/jwt-go](https://pkg.go.dev/github.com/dgrijalva/jwt-go) \
  for JWT-based session authentication
- [github.com/lib/pq](https://pkg.go.dev/github.com/lib/pq) \
  postgresql driver
- [github.com/urfave/negroni](https://pkg.go.dev/github.com/urfave/negroni) \
  middleware for basic request logging
- [golang.org/x/crypto](https://pkg.go.dev/crypto) \
  bcrypt for encrypting passwords to store in the database
- standard library

notes:
- default admin credentials are username:`admin` password:`compfesthospitaladmin`
- sample patient credential username:`aboots` password:`andi`

## Frontend

Go ([html/template](https://pkg.go.dev/html/template))

[Bootstrap 5](https://getbootstrap.com/)

HTML, CSS, JS

## Assumptions (Notes for COMPFEST)

Some unclear requirements are given, and we made these assumptions to make sure we deliver the end product as best as possible while still being on time.

1. Authentication
   - There are two types of account roles \
     It is not explicitly specified how this is implemented, so we went with a boolean field in the "users" table named "admin"
2. Administrator
   - "Administrator can see a list of patients that are registered in each appointment" \
     Reading is the only necessary requirement, no editing
   - Admins have no privileges to edit patient data
3. Patient
   - "Patients can see a list of appointments" \
     Appointment data to be displayed is not specified, we assumed this to be public data; doctor name, description, current registrants count, and max registrants. Also it is stated that only patients can see this list, so users need to be registered as a patient to use this functionality
   - "Patients cannot apply for an appointment with a fully booked registrant" \
     It is not specified the details of what an appointment with a fully booked registrant is, so we decided that each appointment will have it's own maximum registrants count

## Documentation

This project has a [detailed documentation here](docs/spec.md), made conforming to [the requirements document](docs/requirements.pdf).

___

###### Mistakes made

- must've started with a design pattern in mind, for example: MVC
- must've designed RESTful APIs / microservice style APIs (?) for easier AJAX stuff, instead of using POST forms
- must've used a more general template where I understand the deps, or make 1 myself, couldn't use toasts with this template for some reason
- must've used sessions (instead of jwt, but that's a requirement in this project)
- must've gone with a more "struct oriented" approach, with methods instead of regular functions, bcs i didn't do this, i had to edit filenames to (barely) control go's compilation order

would love to get another project like this where I could correct my mistakes

###### Questions

Along making this project, I come accross some questions I hope to ask in the academy

- how to deal with not liking frontend? (i wanna be a backend sofeng, not fullstack)
- exactly how expensive is a DB query?
- should everything be AJAX?
  - how much AJAX is too much AJAX?
- where are you supposed to write js code? especially AJAX related stuff, especially within templates (in this project, i wrote them in the templates just under the html element they're responsible for)
- is it possible to do AJAX / microservice APIs without js?
- how do u document the API responses of every AJAX endpoint like the JSON structure for example?

###### Footnotes

I hope the frontend part is acceptable, I suck at it, but I really want to get in the academy though so I packed this with as many backend features as I can think of. Also I only knew about this event 1 week prior to the registration deadline :')

That aside though, I genuinely really enjoy making this webapp, using minimal frameworks, with a new language I just learned, facing all it's challenges.

I hope to learn more of this in the academy, I really want to be able to do more of this, and to make it my job, and learn even more in the field!
