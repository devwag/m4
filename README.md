# m4

This project is an Event Grid Web Hook written in Go that supports receiving and processing the person message type.

## Notes

** TODO [cover.html](./cover.html) contains the latest code coverage map. The Go tools have some "challenges" building this, so it's not an automated process yet.

## Developer Prerequisites

* Developer workstation will need
  * VSTS / Git access
  * Go 1.10 (I can't get Delve to work on 1.11 Windows ...)
  * Azure subscription access for SSH access (optional)

## Getting Started

* docker
  * contains docker build files for debug and release
* sendmessage
  * a simple app / shell script for sending a message to event grid
    * not part of the package
* src/m4
  * sampleapp
    * contains a sample Web Hook implementation
  * logb
    * simple log wrapper for chaining requests
  * eventgrid
    * handler that parses the event grid "envelope" and handles validation events

## Flags

* port - int - 8080 - port for web hook to listen on
* logpath - string - ./logs/ - path to the log files
  * /home/LogFiles/ is the CIFS share mounted by App Service
  * the docker files override the default and use /home/LogFiles
  * using /home/LogFiles is problematic on the Mac because of permissions on /home

## Sharing Plan

* Event Grid webhook sample for App Services / Go
* Setting up CI/CD for App Services / Go in Azure DevOps
* Debugging App Services with SSH and Go
* App Insights integration for App Services / Go

## To Do List

* TODO - PR for app analytics
* TODO - automate building infrastructure
* TODO - add Apache log handler to write Apache style web logs?
* TODO - automate end to end test
* TODO - complete and automate performance test
* TODO - complete and automate scale test
* TODO - automate cover.html
  * there's a bug but can be worked around with sed
