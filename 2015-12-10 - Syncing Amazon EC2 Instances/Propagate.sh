#The first command line of the script contains the master’s IP, so it does not sync with itself.
export LocalIP="Your_Master_IP_Here"; #FILL THIS IN

#Get the IPs of all slave instances
export NewIPs=`aws ec2 describe-instances | grep '"PrivateIpAddress"' | perl -i -pe 's/(^.*?: "|",?\s*?$)//gm' | sort -u | grep -v $LocalIP`

#Loop over all slave instances
for i in $NewIPs; do
        echo "Syncing to: $i";
        #Run an rsync from the master to the slave
        rsync -aP -e 'ssh -o StrictHostKeyChecking=no' /var/www/html/ root@$i:/var/www/html/;
done