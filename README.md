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

More random stuff
-----------------
FirefoxOS webservice
====================

- this pad aims to centralize the discussion and design of implementing web services for FirefoxOS.
- Written in Go (easy to build a crosstoolchain for arm)

Configuration
 - User/Password
 - ListenPort

Problems.
 - some OS block connecting to local ports
 - How to start the program at the begining?
 - as long as its listening in localhost all apps have access to it, requires auth
   - http auth proposed
 - installation: requires adb push
 - how to start/stop a service?
 - language prefered
   - probably i would prefer JS
   - but it is easier to crosscompile in Go
 - the service must survive even after an OTA upgrade
   - store daemon in webapps folder?
   - looks like /system is not updated after the OTA updates
   - init.rc is replaced
 - which services aim to provide?
   - mainly native interaction with system
     - execute programs and get the output
     - get a shell!
     - show /proc files
     - launch instances of r2
     - start of stop services
       - ssh
       - openvpn
       - syncthing
       - btsync
       - run Tor
       - adb over wifi
     - list processes
     - kill a process
     - reboot
     - enable/disable ssh
