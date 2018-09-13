# m4

This project is an Event Grid Web Hook written in Go that supports receiving and processing the person message type.

## Developer Prerequisites

* Developer workstation will need
  * Git
  * Go 1.11
  * optional
    * Azure DevOps access to import / edit CI/CD pipelines
    * Azure subscription access for Kudu / SSH access to debug build
    * Azure subscription access for Kudu access to release build log files

## Repo Organization

* cicd
  * these are the import files for the CI/CD pipelines
    * from Azure DevOps / Pipelines / Build, click "new" and choose import a pipeline then specify either the m4-debug-pipeline.json file or the m4-release-pipeline.json file
  * Note that if you import into a different project, you will have to setup credentials for Docker hub and perhaps github

* docker
  * contains docker build files for debug and release

* sendmessage
  * a simple app / shell script for sending a message to event grid
    * not part of the package

* src/m4
  * samplewebhook
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

* TODO - move these into VSTS as work items / bugs
* TODO - Pete - Add app analytics
* TODO - Jeff - automate building infrastructure
* TODO - add Apache log handler to write Apache style web logs?
* TODO - automate end to end test
* TODO - complete and automate performance test
* TODO - complete and automate scale test
* TODO - Jeff suggested defaulting logs to ./ and checking to see if /home/LogFiles exists and, if so, using that folder
* TODO - verify that the request came from a trusted source
  * one thought is to put a token in the webhook URL and pull that token from Key Vault
  * Jeff is looking into other options
* TODO - bug - Jeff pointed out that the dockerfiles don't have an ARG parameter, so you can't override the flags
* TODO - do we want to pull the packages from github.com/bartr/m4?
  * for debugging, I like being able to use the local packages
  * using the github link, can we pull packages from a different branch than master?
