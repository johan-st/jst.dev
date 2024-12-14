# jst.dev

jst.dev is my personal website and playground. I am in the middle of rebuilding it from scratch in go. I am staying close to the standard ilbrary and will take the oppotunity to explore concepts capabilities with very few dependencies. I will in several places happily "reinvent the wheel". My oppinion is that to understand and internalize a technology or concept you need to build it yourself.

As I enjoy the slightly lower levels of abstraction and dev-ops, this project is somewhere I will do the heavy lifting myself. I have the feeling a lot of services are wildly over-complicated for most use cases and thus I will try to focus on the simpleset possible setup that gets the job done.

This whole project is ment to be on a lower level of abstraction than my day job.

I intend to incorporate my image service from a previous project and build a cdn for the static assets of the app. I will also deploy the service as distrubuted application on fly.io. Hopefully I will be able to create a nice setup for future projects.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Development

```bash
# watch wont work on wsl
templ generate --watch 

# run server
go run cmd/web/main.go
```

## Architecture

- A bespoke web-server
  - templ for templating and server side rendering
- fly.io for deployment
- Tigris for persistent storage of larger files (images, etc)
- Turso for persistent of data (and logging?)
- Sentry is on the table for error tracking
- (if ever I need analytics, PostHog is on the table)

The app is built around a bespoke (go) web-server. This ensures a very flexible setup with very few limitations and little overhead.

The frontend resources are embedded in the binary and rendered on the server.

### persistant storage

#### truso.tech

For persistence between deployments we use truso.tech. It is a database service built around sqlite (libsql to be precise). It provides a great dev experience and is fast and easy to manage with a great cli and dashboard for managing the database and its performance. They provide a feature for embedding local replicas for read operations that are easily synced with the remote master.

(Turso will happily create thousands of databases for you. An intressting pattern they enable is to have a database per user. This ensures isolation and performance. The provide tools for syncing the schema of several databases.)

#### tigris

Tigris is a s3 compatible blob storage service. It is a great way to store files and other large data. What makes them stand out is that they are build on fly.io and will sync the data between any region the data is used from. This will enable us to easily build a cdn for the static assets of the app when the time comes.

### fly.io

Fly.io is a platform for deploying and scaling web apps.

## TODO

### blog

- create a separate db for blog posts
- embed a replica of the db in the app
- consider prerendering the markdown to html
  - i.e. have a db-column with the rendered html
  - postpone indefinately. Probablt not necceary and will create two sources of truth for the blog posts
- Set up a way to import posts from obsidian
- Set up a way to manage assets (images, etc)
- Set up a way to manage the blog posts
  - create, update, delete
  - draft / published
  - schedule for publishing
  - tags
  - ...

### development

- set up live reload for frontend development.

### deployment

- consider dagger if the setup starts to get complex

###
