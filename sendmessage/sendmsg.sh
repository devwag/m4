#!/bin/sh
# requires Azure cli and access to the subscription
# run using source or . or it won't work because of the return


if [[ -z ${m4key} ]] || [[ -z ${m4endpoint} ]]
then
        if [[ -z $1  ]]
        then
                echo "usage: $0 subscriptionName"
                exit 1
        fi

        az account set -s $1
        m4endpoint=$(az eventgrid topic show --resource-group 7-11-Dev --name minus4dev --query endpoint --output tsv)
        m4key=$(az eventgrid topic key list --resource-group 7-11-Dev --name minus4dev  --output tsv --query key1)
fi

m4body='[{"id": "'"$RANDOM"'","topic":"","subject":"person","eventType":"person","eventTime":"2018-09-08T07:16:46Z","data":{"firstName":"'${USER}'","lastName":"Doe"}}]'

curl -X POST -H "aeg-sas-key: $m4key" -d "$m4body" $m4endpoint
