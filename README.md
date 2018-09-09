# m4

This project is an Event Grid Web Hook written in Go that supports receiving and processing the person message type.

## Notes

** TODO [cover.html](./cover.html) contains the latest code coverage map. The Go tools have some "challenges" building this, so it's not an automated process yet.

## Developer Prerequisites

* Developer workstation will need
  * VSTS / Git access
  * Go 1.10 (I can't get Delve to work on 1.11 Windows ...)
  * Azure subscription access for SSH access
    * TODO - what level access is required for SSH?

## Getting Started

* docker -  contains docker build files (currently for developer deployment)
* sendmessage - a simple app for sending messages to event grid
* src/m4/sampleapp - contains the Web Hook that is deployed to: m4.azurewebsites.net
* src/m4/logb - a simple log wrapper for chaining requests
* src/m4/eventgrid - handler that parses the event grid "envelope" and handles validation events

## Flags

* port - int - 8080 - port for web hook to listen on
* logpath - string - /home/LogFiles/ - path to the log files
  * /home/LogFiles/ is the CIFS share mounted by App Service

## Sharing Plan

* Event Grid webhook sample for App Services / Go
* SSHD sample for  App Service / Go
* Multi stage build sample for App Services / Go
* App Insights integration for App Services / Go
* Securing secrets (Key Vault?) for App Services / Go

## To Do List

* TODO - PR for app analytics
* TODO - automate building infrastructure
* TODO - add Apache log handler to write Apache style web logs?
* TODO - automate end to end test
* TODO - complete and automate performance test
* TODO - complete and automate scale test
* TODO - automate cover.html
  * there's a bug but can be worked around with sed
