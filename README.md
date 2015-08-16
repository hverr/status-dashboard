# Status Dashboard
[![Build Status](https://travis-ci.org/hverr/status-dashboard.svg?branch=master)](https://travis-ci.org/hverr/status-dashboard)

**Status Dashboard** is a highly configurable **AngularJS** application backed by a **Go** API to show information about **multiple servers** in a dasboard style way.

 - [Features](#features)
 - [Development](#development)
 - [Screenshot](#screenshot)

## Features

 - **Efficient**
   - Machines only send information when someone has opened the dashboard.
   - No database or permanent file storage needed.
   - Only update requested information.
 - **Dynamic**
   - Add and remove widgets.
   - Reorder and resize widgets.
   - Add and remove columns and rows.
   - Can refresh as often as every second.
 - **Persistent**
   - Bookmark a layout in your browser.

## Development
The following dependencies are needed to build and run the application.

  - Node and NPM
  - Bower
  - Go with a working `GOPATH`

Get the project: `go get https://github.com/hverr/status-dashboard`

### Server

Build and run the server

```sh
cd $GOPATH/src/github.com/hverr/status-dashboard/server
npm install
npm install -g bower
bower install
cd $GOPATH/src/github.com/hverr/status-dashboard/server/main
go run *.go -c dev_config.json
```

Point your browser to [http://localhost:8050/](http://localhost:8050)

### Client

Build and run the client(s)

```sh
cd $GOPATH/src/github.com/hverr/status-dashboard/client/main
go run *.go -c dev_config.json
```

## Screenshot
![Screenshot](screenshot.png)
