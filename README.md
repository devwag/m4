# m4

This project is an Event Grid Web Hook written in Go that supports receiving and processing 1 custom message type: Person. 

## Notes

** TODO [cover.html](./cover.html) contains the latest code coverage map. The Go tools have some "challenges" building this, so it's not an automated process yet.

## Developer Prerequisites

* Developer workstation will need
  * VSTS / Git access
  * Go (latest)
  * SQL Server query tool (VSCode extension, SQLCMD, SSMS, etc.)
  * Docker on a Linux VM (CI/CD can also be used)
  * Azure subscription access for SSH access
    * TODO - what level access is required for SSH?

## Getting Started

* // TODO - not all checked in yet
* app - contains the Web Hook that is deployed to: minus4dev.azurewebsites.net
* docker -  contains production and dev docker build files
* end2endtest - is an end to end test app that posts messages to Event Grid and verifies they get logged in SQL Server
* SQL - is the sql DDL for creating the person table for messages

## Flags

* port - int - 8080 - port for web hook to listen on
* logpath - string - /home/LogFiles/ - local path to the log files
  * this is the CIFS share mounted by App Service)
## Sharing Plan

* Event Grid webhook sample for App Services / Go
* SSHD sample for  App Service / Go
* Multi stage build sample for App Services / Go
* App Insights integration for App Services / Go
* Securing secrets (Key Vault?) for App Services / Go

## To Do List

* TODO - setup CI / CD pipeline
* TODO - PR for app analytics
* TODO - automate building infrastructure
* TODO - add Apache log handler to write Apache style web logs?
* TODO - automate end to end test
* TODO - complete and automate performance test
* TODO - complete and automate scale test
* TODO - automate cover.html
  * there's a bug but can be worked around with sed
  * need to move to src/app to make work
  * need to modify GOPATH to make work

## Code Review Feedback

* TODO - Add -verbose flag to logging
* TODO - Log / Reject bad EG messages
* TODO - handle upcoming EG change that can send multiple items in a message
