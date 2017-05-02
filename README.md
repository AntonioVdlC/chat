# Simple Go Chat

Just a simple chat written in Go!
It is mobile-first (-only?) and can be deployed fairly easily on Heroku, that way you can have total control on how your chat messages are handled.
For now, the only auth provider is Facebook.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

To get this little chat running, you'll need to get Go on your machine!
Your best bet is to look into the official documentation to get set: https://golang.org/doc/install

> Note: if you are on macOS, have `brew` and want to have a more personalized install you can look into this: http://stackoverflow.com/a/40129734

### Installing

Once you have Go, running the server is as simple as

```
go run *.go
```

and open your browser to http://localhost:8000/

> NB: make sure you set the environment variables `SESSION_SECRET`, `FACEBOOK_KEY` and `FACEBOOK_SECRET`!

## Running the tests

We still have to work on those ... ahem :grimacing:

## Deployment

I ain't no DevOps, so these instructions will help you set up a production instance of this simple chat on Heroku!
> If you have set this up on something else and want to share your knowledge, please feel free to open a PR! :)

### Prerequisites

Regardless of your prefered method of deployment, you will need to tweak a bit the source code to get it running!
- First of, you may want to clone this repository.
- Once that's done, please update your hostname (`HOST`) in the `.env` file.
- Finally, you can customize the title of the app in the `locales/` folder, and you may want to modify it too in the `manifest.json`. For a different set of colors, please take a look at the files inside `public/styles/`! Feel free to change the `favicon.ico` too, and the files inside the `public/icons/` folder.

Now you're all set to deploy!

### Heroku

Create a new Heroku app and provision a PostgreSQL database.
In the dashboard, connect your GitHub repo to the app and Deploy ... there you go! :tada:

## Contributing

No contributing guidelines for now ... just fork this repo, start working on an issue in a branch, push and open a PR! :)

## License

This project is licensed under the MIT License.

## Acknowledgments

* [Ed Zynda III](https://github.com/ezynda3) for his [chat example in Go](https://github.com/ezynda3/go-chat) and subsequent [article](https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets)
* [Gorilla](https://github.com/gorilla) contributors for their libs and, most importantly, for their [chat example in Go](https://github.com/gorilla/websocket/tree/master/examples/chat)!
* [Billie Thompson](https://github.com/PurpleBooth) for this [README.md template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2) ... it's awesome!
