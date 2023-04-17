# **`jst.dev`**

This is an attempt at consolidating a somewhat ambitious set of projects into a single domain. For now it is a place to jot down my thoughts, ideaa and dreams. I will update this document as I progress.

My goal is to help myself focus on a synergistic set of projects and gain experience with devops and infrastructure.

documentation of services, portfolio, infrastructure administration, and other projects on the domain jst.dev.

## HTTP services

| service                             | protocol | path             | description                    |              by               |                   online                   |
| :---------------------------------- | -------: | :--------------- | :----------------------------- | :---------------------------: | :----------------------------------------: |
| [Portfolio](#portfolio)             |     http | `jst.dev`        | portfolio website              |              me               |                    yes                     |
| [Image Server](#image-server)       |     http | `img.jst.dev`    | image server                   |              me               |                     no                     |
| [Meal Calculator](#meal-calculator) |     http | `meal.jst.dev`   | a meal calculator              |              me               | yes ([here](https://strandersson.se/meal)) |
| [Status](#status)                   |     http | `status.jst.dev` | a dashboard for the server     |              me               |                     no                     |
| [API](#api)                         |     http | `api.jst.dev`    | consolidated backend           |              me               |                     no                     |
| [Wishlist](#wishlist)               |      ssh | `jst.dev`        | a catalog of ssh endpoints     | [charm.sh](https://charm.sh/) |                     no                     |
| [Git Server](#git-server)           |      ssh | `git.jst.dev`    | git server                     | [charm.sh](https://charm.sh/) |                     no                     |
| [Status (ssh)](#status-ssh)         |      ssh | `status.jst.dev` | ssh a dashboard for the server |              me               |                     no                     |
| [gRPC](#gRPC)                       |     http | `grpc.jst.dev`   | a gRPC server                  |              me               |                     no                     |

## Infrastructure

- For now, I plan to run all services on a single server. I will likely use a small DigitalOcean droplet.

### Alternatives

- I may build a custom go web server to handle all the http requests.
- I may use kubernetes to run the services. I have some experience with kubernetes but it would be a good learning experience since it was some time since I last worked with it.
- I may use digital oceans App Platform to run the services.
- I may look into Terraform to manage the infrastructure.

# Portfolio

_source code: [github](https://github.com/johan-st/portfolio)_

A react site buildt to list and describe some of my projects

# Wishlist

_source code: [github](https://github.com/charmbracelet/wishlist)_

A private catalog of ssh endpoints. Wishlist can be used to store and share ssh endpoints. It also gates access to the endpoints. Access is granted when the user is added to the config file. 

# Git Server

_source code: [github](https://github.com/charmbracelet/soft-serve)_

A git server accessible via ssh and git. ssh to `git.jst.dev` to browse the repositories.

# Image Server

_source code: [github](https://github.com/johan-st/go-image-server)_

- A simple image server written in go. It is used to host images for any website I host. It stores a master copy of the image and then generates and caches the sizes requested.

# Meal Calculator

not on jst.dev yet but it is online [here](https://strandersson.se/meal).

buildt with Elm.

# Status

- monitoring uptime and other metrics.

# Status (ssh)

- ssh to `status.jst.dev` to get a dashboard of my services.
- monitoring uptime and other metrics.

# API

- planed consolidated backend for all web services
- and maybe the ssh-applications

# gRPC

- This is a far future maybe. Might be a good way to learn gRPC.
- Considering using this behind the scenes for the API and ssh applications.
