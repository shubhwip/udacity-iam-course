# Udacity Go Auth Example

## Prerequisites
- Docker and Docker Compose
- Go

## Backend
- Backend is written in go programming language
- It has one main.go file which defines several endpoints
    - `/register` - Registers a new account
    - `/login` - Logs in a user
    - `/admin/users` - Shows all registered user - Only accesible to admin user
- We covers password based Authentication in this application
- We also cover Role Based Access Control for e.g. admin has elevated permissions and can see all registered users while a regular non admin user can't see all registered users.
- This demo application is configured with two prepopulated users
    - Admin User : `username:admin, password:admin`
    - Non Admin User : `username:user, password:user`

## Frontend
- Backend is written in react programming language
- There are different views defined for login, register and allUsers screen.

## Run
### Build Docker Image and Start Backend and Frontend Services
```
docker compose up --build
```

## Access
```
http://localhost:3000
```
### Login with Admin User
```
Username : admin
Password : admin
```
Once you're logged with admin user, you would be able to see all the registered users(http://localhost:3000/admin/users).

### Login with Non Admin User
```
Username : user
Password : user
```
Once you're logged with admin user, you would not be able to see all the registered users as you will be unauthorised(http://localhost:3000/admin/users).

## Footnotes
- Note : This is a very basic application covering fundamentals of Authentication and Authorisation.
- Future enhancements will include
    - Logout - Session Invalidation
- Data store to persist users information