# TODO Microservice

A simple todo backend. Built to practice go and have a backend from creating todo applications when testing out new front-end frameworks.

The app is hard-coded to use port 8080.

Build the image:
```sh
docker build -f ./DOCKERFILE  -t behlers/todo-app .
```

Run the image:
```sh
docker run -p 8080:8080 behlers/todo-app
```