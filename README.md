# notes-backend

This is the backend for the `notes` project. It is a [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) HTTP server written in [Go](https://en.wikipedia.org/wiki/Go_(programming_language)) which allows you to keep notes. These notes are saved in a file containing an [SQLite](https://en.wikipedia.org/wiki/SQLite) database (which you can open yourself with the proper tools).

## Installation

Before you can get started you need to [install Go](https://golang.org/doc/install). After that you can install the backend by opening a command prompt and entering:

```
go get -u -v github.com/ojz/notes-backend
```

## Usage

Now that you have installed the backend, you can start it using the following command:

```
notes-backend
```

You can see all the options that are supported by typing:

```
notes-backend -help
```

## API

When you have the server running you can interact with it by sending it HTTP requests. You could, for example, use the [fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) in JavaScript for this. Take a look at the `build.sh` file to see what kinds of requests you can send.

Note that if you type `source build.sh` in a bash terminal, you will be able to send HTTP requests through the command line using the `get`, `post`, etc... functions defined in `build.sh`.