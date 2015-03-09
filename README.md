LocalWebService
===============

The LWS is a system management daemon for FirefoxOS, but it can be
used in any other platform.

It's a web server that runs as root and accepts clients using http
auth only in localhost. Once permission is granted, the rest of
pages are provided to manage the services.

--pancake

Services For LWS
----------------
The services that run under lws are mainly native programs that need
to run always in background and offers a way to communicate with them
via HTTP API. Basic status information is provided (pid running,
stopped, logs, ..)

* Tor
* OpenVPN
* WebServer
* Syncthing
* WebShell
* Pebble Link
* TcpDump
* GPG
* Run NodeJS apps ?
* Launch r2 sessions
* Interact with programs
  - run command
  - rest api for stdin/stdout/kill/eof/...

Security
--------
- Downgrade permissions / sandbox / change user when running a service
- The webserver must be run without root perms
- Login/Password information must be stored in hashed form
- The hashed password must be salted with a random value on each build
