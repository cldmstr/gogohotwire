# GOGO Hotwire

This is an example project to show one way to use the [Hotwire HTML-over-the-wire](https://hotwired.dev) frameworks
Turbo and Stimulus in a Golang service.

It attempts to cover all the basic and a few of the less basic use-cases you might use the frameworks
for. Although Hotwire is very straight-forward and mostly backend language agnostic, there are one or two
instances where it is not quite obvious how to use the framework on the server side. There are many good examples
for Ruby and other languages. This project attempts to give you a template to work from for Golang.
It is not meant as a substitute, but as an addtion to the excellent Hotwire documentation and
will hopefully offer enough bread-crumbs to get started or get past some not quite obvious issues encountered when
working with the framework.

## Run

```shell
go run ./cmd/gogohotwire
``` 

Service will be available on `http://localhost:8088`

## Service

The service is a silly (and very minimal) take on a race management application. Most of the functionality is hopefully obvious.
There are two "domains": races and drivers.

- Create a race by typing in a name and pressing enter or the + button.
- Click on a race title to see the race details.
- Start a race from the race details page.
- Click on drivers to open a sidebar with the list of driver information.

## Hotwire Turbo and Stimulus Features

Here are some quick links to where in the HTML templates the different Hotwire features are used:

- [Eager-Loading Frame](internal/app/views/application.html) line 24
- [Lazy-Loading Frame](internal/app/views/application.html) line 29
- [Targeting Navigation Into or Out of a Frame](internal/app/views/application.html) line 15
- [Turbo Stream Prepend](races/views/add_race.stream.html)
- [Turbo Stream Update](races/views/race_details.stream.html)
- [Form Within Frame](races/views/list.html)
- [Simple Stimulus Controller](internal/app/views/application.html) line 19 
  [[javascript controller]](assets/js/src/controllers/toggle_drivers_controller.ts)
- [Server Sent Events Turbo Stream](races/views/race_running.partial.html)
  [[Start function]](races/service_model.go) [[javascript controller]](assets/js/src/controllers/race_update_controller.ts)

## Templating HTML

As opposed to a RESTful API which generally returns JSON or XML, sending HTML over the wire requires templating HTML
snippets or whole pages to return the necessary response. Hotwire comes with very good server side integration with a
Ruby framework created by the same team that created Hotwire. Golang is very well suited to this task with the built-in
templating, but does not quite come with the ease of use afforded by the extra work put into Ruby.

The solution used here is to add a [custom renderer](internal/template/template.go) to the echo HTML handler which
handles the templating and some special cases like adding `Content-Type` headers for Turbo streams, etc.

Special mention to the `yield` [function](internal/template/tools.go) that allows the rendering of a template inside a
template. This is built in to Ruby, but can be added to the go templates by implementing a custom templating function.

## Development

- Install [Task](https://taskfile.dev)
- Run `task dev` to start service and reload on file changes

## Special Thanks

Awesome images by and (c) Erick Zelaya