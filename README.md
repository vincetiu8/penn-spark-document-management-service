# Penn SPARK Document Management System

This is a document management system for Penn SPARK. I've built this to meet the
technical challenge the team has set as part of the Penn SPARK Red application process.

## Development Stack

Frontend: React, Redux, Axios, Material UI  
Backend: Golang, Mux, JWT, GORM
Database: SQLite

I'm aware that Penn SPARK uses the MERN stack, and I've also done some work with Express, MongoDB and Node.
However, given the time constants of this project, I preferred to develop the backend in Golang, which I'm more
comfortable with. I am very much willing to learn more about the MERN stack and improve my skills in it if I get to join
the team!

I based some code off a previous project, but had to fix lots of errors and add some new features. I also had to set up 
the hosting service on DigitalOcean, which is accessible at `http://206.189.185.232/`.

## Running the project

1. Clone the repository
2. Set up the backend
    1. Install golang
    2. `cd server`
    3. `go install`
    4. `go run main.go`
3. Set up the frontend
    1. Install nodejs
    2. `cd client`
    3. `npm install`
    4. `npm start`
4. You should be able to access the website on `localhost:3000`.
5. You can access the API on `localhost:8080`, if you want to make direct calls.
6. Alternatively, check out the project at this link! `http://206.189.185.232/`.

## Usage Instructions
- The default username is "admin" and the default password is "password".
- Once you log in, you need to give the admin role access to the filesystem.
  - This can be done by navigating to the "Access Roles" page and adding a new role with publisher access.
  - This will allow you to upload, publish, and delete documents. Try playing around with the filesystem!
- All documents are shared with all users in the system.
  - You can create new users in the "Users" page, and attach new access roles to them to let them into the filesystem.
- If you have any questions, please reach out to me at `vincetiu@seas.upenn.edu`!
