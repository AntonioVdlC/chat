# Efrei Int Chat

Just a simple chat for past, present and future members of Efrei International!

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

To get this baby running, you'll need to get Go on your machine!
Your best bet is to look into the oficial documentation to get set: https://golang.org/doc/install

> Note: if you are on macOS, have `brew` and want to have a more personalized install you can look into this: http://stackoverflow.com/a/40129734

### Installing

Once you have Go, running the server is as simple as

```
cd ./src
go get github.com/gorilla/websocket
go run main.go
```

and open your browser to http://localhost:8000/

## Running the tests

We still have to work on those ... ahem :grimacing:

## Deployment

This one is also up for grabs, but this is probably gonna end up in a dyno on Heroku or somethings!
We can also use the domain http://efrei.international to serve the chat (either as a subdomain or a path)!

## Contributing

No contributing guidelines for now ... just start working on an issue in a branch, push and open a PR! :)

## License

This project is licensed under the MIT License.

## Acknowledgments

* [Ed Zynda III](https://github.com/ezynda3) for his [chat example in Go](https://github.com/ezynda3/go-chat) and subsequent [article](https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets)
* [Billie Thompson](https://github.com/PurpleBooth) for this [README.md template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2) ... it's awesome!
