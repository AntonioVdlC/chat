# Simple Go Chat

Just a simple chat written in Go!
It is mobile-first (-only?) and can be deployed fairly easily on Heroku, that way you can have total control on how your chat messages are handled.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

#### Go
To get this little chat running, you'll need to get Go on your machine!
Your best bet is to look into the official documentation to get set: https://golang.org/doc/install

> Note: if you are on macOS, have `brew` and want to have a more personalized install you can look into this: http://stackoverflow.com/a/40129734

#### PostgreSQL

You will also need to install PostgreSQL on your machine!
Once that's done, you'll need to create a `chat_dev` database as well as a `dev` user with `dev` as password that has all the rights on the `chat_dev` database! All the tables will be created when lauching the app, so no need to worry about that!

#### NodeJS

To be able to build the frontend assets, you will need to have an up-to-date version of NodeJS on your machine! We recommand also installing `yarn`.

Once NodeJS and `yarn` are installed, simply run `yarn install` to install all the dependencies!

### Installing

For the app to run in dev, you will need to start PostgreSQL (and make sure it listens on port 5432, which is the default).

Also, you will need to copy `.env.example`, rename the file `.env` and complete all the fields except `DATABASE_URL` and `PORT`.

Finally, to build the assets and watch them, you will need to run `npm run dev`.
> Every change on CSS or JS files will trigger a build, and a simple refresh of the browser will pull up the new file versions!

Once that is up and running, starting the server is as simple as

```
ENV=dev go run *.go
```

and open your browser to http://localhost:8000/

## Running the tests

We still have to work on those ... ahem :grimacing:

## Deployment

I ain't no DevOps, so these instructions will help you set up a production instance of this simple chat on Heroku!
> If you have set this up on something else and want to share your knowledge, please feel free to open a PR! :)

### Prerequisites

Regardless of your prefered method of deployment, you will need to tweak a bit the source code to get it running!
- First of, you may want to clone this repository.
- You will also need to create Facebook Login access tokens, as well as a Twitter API key and secret!
- Once that's done, please copy the `.env.example` file into `.env` and update all the fields.
- Finally, you can customize the title of the app in the `locales/` folder, and you may want to modify it too in the `manifest.json`. For a different set of colors, please take a look at the files inside `source/styles/`! Feel free to change the `favicon.ico` too, and the files inside the `public/icons/` folder.

> NB: Make sure you precompile the assets by running `npm run build` before deploying!

Now you're all set to deploy!

### Heroku

Create a new Go Heroku app and provision a PostgreSQL database.
In the dashboard, connect your GitHub repo to the app and Deploy ... there you go! :tada:

> NB: As the `.env` file is ignored by git, you may want to copy the environment variables into Heroku directly!

## Contributing

No contributing guidelines for now ... just fork this repo, start working on an issue in a branch, push and open a PR! :)

## License

This project is licensed under the MIT License.

## Acknowledgments

* [Ed Zynda III](https://github.com/ezynda3) for his [chat example in Go](https://github.com/ezynda3/go-chat) and subsequent [article](https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets)
* [Gorilla](https://github.com/gorilla) contributors for their libs and, most importantly, for their [chat example in Go](https://github.com/gorilla/websocket/tree/master/examples/chat)!
* [Billie Thompson](https://github.com/PurpleBooth) for this [README.md template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2) ... it's awesome!
