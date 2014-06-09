elwyn
=====

Goals:

1. reduce heartbeat rate if window not visible
2. "hide screen" when window visible for too long (AFK protection)
3. cleaner input area
4. regularly replace timestamp with relative timestamps (about 3 minutes, 2 days, etc)
5. all navigation done through nav bar
	* rooms/options are _'bootstrap breadcrumbs'_ that drop down or something
		* {Elwyn} / {organization | group} / {room | whisper}
	* people in users list are links to start PM session
		* {meta-key}+click may @mention someone in the current channel
6. IRC-style commands
7. Off-The-Record
8. all the __good__ and __thought out__ encryption I can muster
9. only keep last X messages -- scrolling loads more
10. quickselect -- potentially on holding of a meta key
	* use something like [this](http://creative-punch.net/2014/02/making-animated-radial-menu-css3-javascript/)
	* link to most recent and/or popular:
		1. people/privates
		2. rooms
		3. groups
11. search

Introduction
============

Quickjump
=========

1. [Using](#using)
2. [Building](#building)
3. [Testing and Developing](#testing-and-hacking)
4. Appendix:
	* [Notable Make Commands](#make-commands)


Using
=====

__TODO: Mongo documentation__

### From Source

First, read the [building](#building) section. It will illustrate how to:

1. [Install Dependencies](#dependencies)
2. [Build the Binary and Client](#building-application)

Then decide whether you want to:

1. [Run Locally](#locally)
2. [Run In a Virtual Machine](#testing-and-hacking)

### From Binaries

Coming soon...

### Running

Now that we've got everything built/installed let's run the thing.

#### Locally

To run the server on the machine you built/downloaded on:

	make run

This also requires a `Mongo` database be running on the system.

Assuming no changes were made to the `Vagrantfile` the application will be accessible at:

	http://192.168.20.20

#### EC2

Coming soon... done in `Vagrant`.



Building
========

Much of the building and installation is done via `make`.

There are a few prereqs.

### Prerequisites

The following are requisites only for building. They are hard dependencies for it, however.

1. The primary development language is [Go](http://golang.org/)
2. Development of the reference client aided by [node](http://nodejs.org)

All secondary libraries, utilities, and actions are installed by the `make` commands.

	make deps

Do note that the `node` requirement may be expanded or completely removed. I am open to suggestions or pull requests on this matter.

For a list of what dependencies are used:

1. Go: `/deps.json`
2. node: `/package.json`


### Steps

If all hard dependencies met, a simple `make install` will do.

For sanity's sake, I'll describe the whole process. It is two parts.

##### Dependencies

1. Install `go`:
	* Debian/Ubuntu/etc: `apt-get install golang`
	* MacOSX:
		* __[brew](http://brew.sh/)__: `brew install go`
		* __without brew__: follow the [go site](http://golang.org) instructions
2. Install `node`:
	* Debian/Ubuntu/etc: `apt-get install nodejs`
	* MacOSX:
		* __[brew](http://brew.sh/)__: `brew install node`
		* __without brew__: follow the [node site](http://nodejs.org) instructions


##### Building Application

Now, to build the binary and client, simply run:

1. `make`
2. `make install`
	* binary installation directory set by `$PREFIX`
		* __inline__: `PREFIX=/my/install/dir make install`

The default `make` target will pull secondary dependencies, compile the binary, and generate distribution client code.

Intuitively, `make install` installs the binary into a globally executable location and the client into a servable directory.




Testing and Hacking
===============

In the root directory you will find a `Vagrantfile` and an `ansible/` directory.

[Vagrant](http://vagrantup.com) provides easy creation of machine environments using a variety of providers but most natable Virtualbox and EC2.

[Ansible](http://ansible.com) is a tool that allows easy provisioning of machines generated by Vagrant or some other means. Simply point it a single or set of boxes and let her rip.

### Local

To start a local environment make sure you have [Vagrant](http://vagrantup.com) and [Virtualbox](https://www.virtualbox.org/) installed.

Once this is done, simply run:

	make local

### Remote



APPENDIX
========

## Make Commands

This list does not contain all the commands. It is meant to highlight the important ones.

### default -- happens on a bare `make`

* Does a "sum-total" build of the server binary and client.
* It __does not__ attempt to install it, however.

### deps

* Install all secondary dependencies for `go` and `node`

### binary

* Only build the server -- go portion of the codebase

### client

* Build the reference client into the ouput directory

### clean

* Remove intermediate files and built directory

### install

* Debian-focused, though should work on most *nices or BSDs
	* verified on `Ubuntu 14.04` and `MacOSX 10.9`
* Install the binary to the `$PREFIX` directory
	* if none supplied, `/usr/local/bin` is used
* Install the client to `/srv/elwyn`
	* this __will__ be made variable soon

### local

* Build a local test environment using [Vagrant](http://vagrantup.com) and [Ansible]	(http://ansible.com).

### run

* Run the server on the machine issuing the command
* Will build the system first
