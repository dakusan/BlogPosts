#!/bin/bash

#This monitors the number of emails in the queue, and sends ours admins an email when a limit (1000) is reached. It would need to be run on a schedule (via cron).
export AdminEmailList="ADMIN EMAIL LIST SEPARATED BY COMMAS HERE"
export Num=`/usr/sbin/exim -bpc`
if [ $Num -gt 1000 ]; then
        echo "Too many emails! $Num" | /usr/sbin/sendmail -v "$AdminEmailList"
        #Here might be a good place to delete emails with “undeliverable” strings within them
        #Example (See the 3rd script): exim-delete-messages-with 'A message that you sent could not be delivered'
fi

#This deletes any emails in the queue from or to a specified email address (first parameter). If the address is the recipient, the from must be "<>" (root)
exiqgrep -ir $1 -f '<>' | xargs exim -Mrm
exiqgrep -if $1 | xargs exim -Mrm

#This deletes any emails in the queue which contain a given string (first parameter)
if [ "$1" == "" ]
then
  echo 'Cannot delete with empty string'
else
  grep -lir "$1" /var/spool/exim/input/ | sed -e 's/^.*\/\([a-zA-Z0-9-]*\)-[DH]$/\1/g' | xargs /usr/sbin/exim -Mrm
fi

# Get a count of emails in the queue per sender (sender email address is supplied by sender and can be faked)
exim -bp | grep -oP '<.*?>' | sort | uniq -c | sort -n

#Get a count of emails in the queue per account (running this script can take a little while)
exim -bp | grep -Po '(?<= )[-\w]+(?= <)' | xargs -n1 exim -Mvh | grep -ioP '(?<=auth_sender ).*$' | sort | uniq -c

#Bonus: Force all non-specified accounts on Exim to use a certain IP address for sending. It would need to be run on a schedule (via cron).
export IPAddress="YOUR ADDRESS HERE"
/usr/bin/perl -i -pe 's/\*:.*/*: '$IPAddress'/g' /etc/mailips