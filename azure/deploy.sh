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
	declare allowEmpty

	prompt=$1
	regex=$2
	error=$3
	if [[ $# -eq 4 ]]; then
		allowEmpty=$4
	fi

	userInput=""
	while [[ ! $userInput =~ $regex ]]; do
		if [[ (-n $userInput) ]]; then
			printf "'%s' is not valid. %s\n" $userInput $error
		fi
		printf $prompt
    read userInput
		if [[ $allowEmpty && (-z "$userInput") ]]; then
			return
		fi
	done
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
		az login
	fi
)

#Prompt for parameters if some required parameters are missing
if [[ -z "$subscriptionId" ]]; then
	printf "\nAvailable subscriptions:\n"
	
	az account list -o table
	echo
	
	currentSub="$(az account show -o tsv | cut -f2)"
	printf -v prompt "Enter your subscription ID [%s]: " $currentSub
	validatedRead $prompt "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}" "Subscription ID must be a GUID." true 

	subscriptionId=$userInput
	if [[ (-z $subscriptionId) && (-n $currentSub)]]; then
		subscriptionId=$currentSub
	else
		[[ "${subscriptionId:?}" ]]
	fi
fi
printf "Using SubscriptionId %s.\n" $subscriptionId

if [[ -z "$appName" ]]; then
	validatedRead "\nEnter a unique application name: " "^[a-zA-Z0-9]+$" "Only letters & numbers are allowed." # app name used as DNS prefix so no underscores allowed
  appName=$userInput
fi

if [[ -z "$resourceGroupName" ]]; then
	validatedRead "\nEnter a resource group name: " "^[a-zA-Z0-9_]+$" "Only letters, numbers and underscores are allowed."
	resourceGroupName=$userInput
fi

if [ -z "$subscriptionId" ] || [ -z "$resourceGroupName" ] || [ -z "$appName" ]; then
	printf "\nEither one of subscriptionId, resourceGroupName, appName is empty\n"
	usage
fi

#set the default subscription id
az account set --subscription $subscriptionId

set +e

#Check for existing RG
az group show --name $resourceGroupName 1> /dev/null
if [ $? != 0 ]; then
	echo "To create a new resource group, please enter an Azure location:"

	if [[ -z "$resourceGroupLocation" ]]; then
		locations="$(az account list-locations --output tsv | cut -f5 | tr '\n' ', ' | sed "s/,/, /g")"
		printf "\n%s\n" "${locations%??}"
		validatedRead "\nEnter resource group location: " "^[a-zA-Z0-9]+$" "Only letters & numbers are allowed."
		resourceGroupLocation=$userInput
	fi

	set -e
	(
		set -x
		az group create --name $resourceGroupName --location $resourceGroupLocation 1> /dev/null
	)
else
	resourceGroupLocation="$(az group show -n $resourceGroupName -o tsv | cut -f2)"

	printf "Using existing resource group...\n"
fi

siteName=$appName
printf -v servicePlanName "%s-serviceplan" $appName
printf -v topicName "%s-topic" $appName
printf -v webhookName "%s-person-webhook" $appName
printf -v webhookUrl "https://%s.azurewebsites.net/person" $appName

(
	set -e

	printf "\nDeploying App Service Plan...\n"
	(set -x; az appservice plan create --name $servicePlanName --resource-group $resourceGroupName --sku B1 --is-linux)
	printf "\nDeploying App Service...\n"
	(set -x; az webapp create --plan $servicePlanName --resource-group $resourceGroupName --name $siteName --deployment-container-image-name bartr/m4)
	printf "\nEnabling storage for log files...\n"
	(set -x; az webapp config appsettings set --resource-group $resourceGroupName --name $siteName --settings WEBSITES_ENABLE_APP_SERVICE_STORAGE=true)
	printf "\nConfiguring App Service TCP ports...\n"
	(set -x; az webapp config appsettings set --resource-group $resourceGroupName --name $siteName --settings WEBSITES_PORT=8080)
	printf "\nDeploying Event Grid topic...\n"
	(set -x; az eventgrid topic create --name $topicName --resource-group $resourceGroupName	--location $resourceGroupLocation)
	printf "\nDeploying Event Grid subscription (webhook)...\n"
	(set -x; az eventgrid event-subscription create --name $webhookName --topic-name $topicName --resource-group $resourceGroupName --endpoint $webhookUrl)
)

printf "\nDeployment Complete.\n\n"