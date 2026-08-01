<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Syncing_Amazon_EC2_Instances">Syncing Amazon EC2 Instances</a><br>
Date: 2015-12-10<br>
Tags: Amazon, AWS, Cloud, EC2, Examples, Servers, SSH<br><br>

<p>In continuation of <a href="../2015-12-09 - Custom Initializations for Amazon AMIs">yesterday’s post, in which I showed how to create Amazon AMIs to keep your newly created EC2 instances up to date</a>, today I will cover syncing already-live instances from the master to slaves. All of the below takes place on the master instance, and assumes all other instances are part of the slave group. You may have to use extra filters on the below “aws” command to only pull IPs from a certain group of instances.</p>

<p>Here is a simple bash script (hereby referred to as “Propagate.sh”) which syncs /var/www/html/ to all of your slave instances. It uses the “aws” command line interface provided by Amazon, which comes default with the Amazon Linux starter AMI.</p>
</div>

```bash
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
```

<p>You may also want to add “-o UserKnownHostsFile=/dev/null” to the SSH command (directly after “-o StrictHostKeyChecking=no”), as a second EC2 instance may end up having the same IP as a previously terminated instance. Another solution to that problem is syncing the “/etc/ssh/ssh_host_rsa_key*” from the master when an instance initializes, so all instances keep the same SSH fingerprint.</p><br>

To let other people manually execute this script, you can create a PHP file with the following in it. (Change /var/www/ in all below examples to where you place your Propagate.sh)
```php
<? print nl2br(htmlentities(shell_exec('sudo /var/www/Propagate.sh 2<&1'))); ?>
```

If your Propagate.sh needs to be ran as root, which it may if your PHP environment is not run as the user root (usually “apache”), then you need to make sure it CAN run as root without intervention. To do this, add the following to the /etc/sudoers file
```
apache  ALL=(ALL)       NOPASSWD: /usr/bin/whoami, /var/www/Propagate.sh
```
Change the user from “apache” to the user which PHP runs as (when running through apache).<br>
I included “whoami” as a valid sudoer application for testing purposes.<br>
Also, in the sudoers file, if “Defaults requiretty” is turned on, you will need to comment it/turn it off.<br><br>

<p>While I did not mention it in yesterday's post, I thought I should at least mention it here. There are other ways to keep file systems in sync with each other. This is just a good use case for when you want to keep all instances as separate independent entities. Another solution to many of the previously mentioned problems is using Amazon's new <a href="https://aws.amazon.com/efs/">EFS</a>, which is currently still in preview mode.</p>