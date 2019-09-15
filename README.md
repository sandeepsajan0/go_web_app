# go_web_app

Go web app for beginners. Good ref. to start with go in web.

## Features implemented:

1. Created a web app and bind with port 8080.

2. Simple listening of application on url: `/hello` and `/goodbye`.

3. Implemented comment feature for both get and post requests using redis on url: `localhost:8080/get_comments`.

4. Able to store session on url: `/login` and can test session username on url: `/login/session`.

5. Implemented the login and register feature. Session will be stored only for 10 min(you can change the maxAge). Comments feature can be used only after login.  

6. Structurize the whole code in modules models, session, routes, middleware. Add static/css file.

## Installation:

Clone the repository and run `go run main.go` inside the repo folder.

You also need to start redis server on its default port `:6379`.

* Feel free to add more features and make more easy for beginners.