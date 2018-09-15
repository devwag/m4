<#
 .SYNOPSIS
    Deploys a template to Azure

 .DESCRIPTION
    Deploys an Azure Resource Manager template

 .PARAMETER subscriptionId
    The subscription id where the template will be deployed.

 .PARAMETER resourceGroupName
    The resource group where the template will be deployed. Can be the name of an existing or a new resource group.

 .PARAMETER resourceGroupLocation
    Optional, a resource group location. If specified, will try to create a new resource group in this location. If not specified, assumes resource group is existing.

 .PARAMETER appName
    A unique name to be used as the resource and DNS prefix. 
#>

param(
 [string]
 $subscriptionId,

 [Parameter(Mandatory=$True)]
 [string]
 $resourceGroupName,

 [string]
 $resourceGroupLocation,

 [Parameter(Mandatory=$True)]
 [string]
 $appName,

 [string]
 $templateFilePath = "template.json"
)

$ErrorActionPreference = "Stop";

if($subscriptionId -eq '') {
    $activeSub = ((Get-AzureRmContext).Subscription)
    if(($result = Read-Host "subscriptionID [$activeSub]") -eq '') {
        $subscriptionId = $activeSub
    } else {
        $subscriptionId = $result
    }
}

#******************************************************************************
# Helper functions
#******************************************************************************
<#
.SYNOPSIS
    Registers RPs
#>
Function RegisterRP {
    Param(
        [string]$ResourceProviderNamespace
    )

    Write-Host "Registering resource provider '$ResourceProviderNamespace'";
    Register-AzureRmResourceProvider -ProviderNamespace $ResourceProviderNamespace;
}

function Login
{
    $needLogin = $true
    Try 
    {
        $content = Get-AzureRmContext
        if ($content) 
        {
            $needLogin = ([string]::IsNullOrEmpty($content.Account))
        } 
    } 
    Catch 
    {
        if ($_ -like "*Login-AzureRmAccount to login*") 
        {
            $needLogin = $true
        } 
        else 
        {
            throw
        }
    }

    if ($needLogin)
    {
        Login-AzureRmAccount
    }
}

#******************************************************************************
# Script body
# Execution begins here
#******************************************************************************
$ErrorActionPreference = "Stop"

# sign in
Write-Host "Logging in...";
Login;

# select subscription
Write-Host "Selecting subscription '$subscriptionId'";
Select-AzureRmSubscription -SubscriptionID $subscriptionId;

#Create or check for existing resource group
$resourceGroup = Get-AzureRmResourceGroup -Name $resourceGroupName -ErrorAction SilentlyContinue
if(!$resourceGroup)
{
    Write-Host "Resource group '$resourceGroupName' does not exist. To create a new resource group, please enter a location.";
    if(!$resourceGroupLocation) {
        $resourceGroupLocation = Read-Host "resourceGroupLocation";
    }
    Write-Host "Creating resource group '$resourceGroupName' in location '$resourceGroupLocation'";
    New-AzureRmResourceGroup -Name $resourceGroupName -Location $resourceGroupLocation
}
else{
    $resourceGroupLocation = (Get-AzureRmResourceGroup -Name $resourceGroupName).Location
    Write-Host "Using existing resource group '$resourceGroupName'";
}

# Register RPs
$resourceProviders = @("microsoft.eventgrid","microsoft.web");
if($resourceProviders.length) {
    Write-Host "Registering resource providers"
    foreach($resourceProvider in $resourceProviders) {
        #RegisterRP($resourceProvider);
    }
}

# Start the deployment
Write-Host "Deploying...";
$parameters = @{}
$parameters.Add("location", "$resourceGroupLocation")
$parameters.Add("dns_prefix", "$appName")
try {
    New-AzureRmResourceGroupDeployment -ResourceGroupName $resourceGroupName -TemplateFile $templateFilePath -TemplateParameterObject $parameters;
}
catch {
    # "esvalidation" indicates the error triggered because the new event grid webhook subscription could not be verified.
    # That is expected and normal because the app service is not yet up and running and can't yet respond to the validation
    # request. Thus, there is no real error. 
    if (($_.Exception.Message).Contains("esvalidation")) {
        Write-Host "Deployment succeeded. The app service typically takes a few minutes to finish initializing before it can accept requests."
    } else {
        Write-Error "$_.Exception.Message"
    }
}
