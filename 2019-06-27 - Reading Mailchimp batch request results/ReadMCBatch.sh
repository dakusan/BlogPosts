#!/bin/bash
#Example variables
BATCHID=abc1234567;
APIKEY=abcdefg-us11@us11.api.mailchimp.com;
APIURL=us11.api.mailchimp.com;

#Request the batch information from Mailchimp
curl --request GET --url "https://dummy:$APIKEY@$APIURL/3.0/batches/$BATCHID" 2> /dev/null | \

#Get the URL to the response
grep -oP '"response_body_url":"https:.*?"' | \
grep -oP 'https:[^"]*' | \

#Get the response
xargs wget -O - 2> /dev/null | \

#The response is a .tar.gz file with a single file in it. So get the contents of this file
tar -xzvO 2> /dev/null | \

#Pretty print the json of the full return and the “response” objects within
php -r '$Response=json_decode(file_get_contents("php://stdin"), true); foreach($Response as &$R) $R["response"]=json_decode($R["response"], true); print json_encode($Response, JSON_PRETTY_PRINT);'