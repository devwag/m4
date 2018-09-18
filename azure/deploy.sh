#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# -e: immediately exit if any command has a non-zero exit status
# -o: prevents errors in a pipeline from being masked
# IFS new value is less likely to cause confusing bugs when looping arrays or arguments (e.g. $@)

#******************************************************************************
# Helper functions
#******************************************************************************
usage() { echo "Usage: $0 -i <subscriptionId> -g <resourceGroupName> -n <appName> -l <resourceGroupLocation>" 1>&2; exit 1; }

validatedRead() {
	prompt=$1
	regex=$2
	error=$3

	userInput=""
	while [[ ! $userInput =~ $regex ]]; do
		if [[ (-n $userInput) ]]; then
			printf "'%s' is not valid. %s\n" $userInput $error
		fi
		printf $prompt
    read userInput
	done
}

readSubscriptionId () {
	currentSub="$(az account show -o tsv | cut -f2)"
	subNames="$(az account list -o tsv | cut -f4)"
	subIds="$(az account list -o tsv | cut -f2)"
	
	while ([[ -z "$subscriptionId" ]]); do
		printf "Enter your subscription ID [%s]: " $currentSub
		read userInput

		if [[ (-z "$userInput") && (-n "$currentSub")]]; then
			userInput=$currentSub
		fi

		set +e
		nameExists="$(echo $subNames | grep $userInput)"
		idExists="$(echo $subIds | grep $userInput)"
	
		if [[ (-z "$nameExists") && (-z "$idExists") ]]; then
			printf "'${userInput}' is not a valid subscription name or ID.\n"
		else
			subscriptionId=$userInput
		fi
	done
}

readLocation() {
	if [[ -z "$resourceGroupLocation" ]]; then
		locations="$(az account list-locations --output tsv | cut -f5 | tr '\n' ', ' | sed "s/,/, /g")"
		printf "\n%s\n" "${locations%??}"

		declare locationExists
		while ([[ -z $resourceGroupLocation ]]); do
			validatedRead "\nEnter resource group location: " "^[a-zA-Z0-9]+$" "Only letters & numbers are allowed."
			locationExists="$(echo $locations | grep $userInput)"
			if [[ -z $locationExists ]]; then
				printf "'${userInput}' is not a valid location.\n"
			else
				resourceGroupLocation=$userInput
			fi
		done
	fi
}

testAppService() {
	(
		set +e
		httpStatus="$(curl -s -w "%{http_code}" -o /dev/null -X POST -H "application/json" -m 120 --data-binary @validation.json $webhookUrl)"
		exitCode=$?
		if [[ httpStatus -eq 200 ]]; then
			printf "App service running.\n"
		elif [[ exitCode -eq 28 ]]; then
			printf "Warning: App service did not respond after 120 seconds. It may stil be starting.\n"
		else
			printf "App service not responding correctly.\n"
		fi
	)
}

#******************************************************************************
# Main script
#******************************************************************************
declare subscriptionId=""
declare resourceGroupName=""
declare appName=""
declare resourceGroupLocation=""

# Initialize parameters specified from command line
while getopts ":i:g:n:l:" arg; do
	case "${arg}" in
		i)
			subscriptionId=${OPTARG}
			;;
		g)
			resourceGroupName=${OPTARG}
			;;
		n)
			appName=${OPTARG}
			;;
		l)
			resourceGroupLocation=${OPTARG}
			;;
	esac
done
shift $((OPTIND-1))

(
	set +e
	#login to azure using your credentials
	az account show &> /dev/null

	if [ $? != 0 ];
	then
		echo "Azure login required..."
		az login -o table
	else
		az account list -o table
	fi
)

#Prompt for parameters if some required parameters are missing
if [[ -z "$subscriptionId" ]]; then
	echo
	readSubscriptionId
fi
printf "Using subscriptijef %s.\n" $subscriptionId

if [[ -z "$appName" ]]; then
	validatedRead "\nEnter a unique application name: " "^[a-zA-Z0-9]+$" "Only letters & numbers are allowed." # app name used as DNS prefix so no underscores allowed
  appName=$userInput
fi

if [[ -z "$resourceGroupName" ]]; then
	validatedRead "\nEnter a resource group name: " "^[a-zA-Z0-9_]+$" "Only letters, numbers and underscores are allowed."
	resourceGroupName=$userInput
fi

#set the default subscription id
az account set --subscription $subscriptionId

set +e

#Check for existing RG
az group show --name $resourceGroupName &> /dev/null
if [ $? != 0 ]; then
	echo "To create a new resource group, please enter an Azure location:"
	readLocation

	(set -ex;	az group create --name $resourceGroupName --location $resourceGroupLocation)
else
	resourceGroupLocation="$(az group show -n $resourceGroupName -o tsv | cut -f2)"
	printf "Using existing resource group...\n"
fi

siteName=$appName
printf -v servicePlanName "%s-serviceplan" $appName
printf -v topicName "%s-topic" $appName
printf -v webhookName "%s-person-webhook" $appName
printf -v webhookUrl "https://%s.azurewebsites.net/person" $appName

set -e

printf "\nDeploying App Service Plan...\n"
(set -x; az appservice plan create --resource-group $resourceGroupName --name $servicePlanName --sku B1 --is-linux)
printf "\nDeploying App Service...\n"
(set -x; az webapp create --resource-group $resourceGroupName --plan $servicePlanName --name $servicePlanName --name $siteName --deployment-container-image-name bartr/m4debug)
printf "\nEnabling continuous deployment (CD)...\n"
(set -x; az webapp deployment container config --resource-group $resourceGroupName --name $siteName --enable-cd true)
printf "\nEnabling storage for log files...\n"
(
	set -x
	az webapp config appsettings set --resource-group $resourceGroupName --name $siteName --settings WEBSITES_ENABLE_APP_SERVICE_STORAGE=true
	az webapp restart --resource-group $resourceGroupName --name $siteName
)
printf "\nVerifying the app service is running...\n"
testAppService
printf "\nDeploying Event Grid topic...\n"
(set -x; az eventgrid topic create --resource-group $resourceGroupName --name $topicName	--location $resourceGroupLocation)
printf "\nDeploying Event Grid subscription (webhook)...\n"
(set -x; az eventgrid event-subscription create --resource-group $resourceGroupName --name $webhookName --topic-name $topicName --endpoint $webhookUrl)

printf "\nDeployment Complete.\n\n"