#User defined variables
VARIABLES_SET=0; #Set this to 1 to allow the script to run
YOUR_DOMAIN="localhost";
YOUR_IP="1.1.1.1";
BLOCKED_ERROR_MESSAGE="Domain is blocked";
ADDRULES_USERNAME="YourUserName";
ADDRULES_PASSWORD="YourPassword";

#Confirm script is ready to run
if [ $VARIABLES_SET != 1 ]; then
    echo 'Variables need to be set in the script';
    exit 1;
fi
if [ `whoami` != 'root' ]; then
    echo 'Must be root to run script. When running the script, add "sudo" before it to' \
        'run as root';
    exit 1;
fi

#Allow root login
cat /home/ec2-user/.ssh/authorized_keys > /root/.ssh/authorized_keys;
perl -pi -e 's/^\s*#?\s*PermitRootLogin.*$/PermitRootLogin yes/igm' /etc/ssh/sshd_config;
service sshd reload;

#Install services
yum -y install bind httpd php;
chkconfig httpd on;
chkconfig named on;
service httpd start;
service named start;

#Set the DNS server to be usable by other computers
perl -pi -e 's/^(\s*(?:listen-on port 53|allow-query)\s*{).*$/$1 any; };/igm' \
    /etc/named.conf;
service named reload;

#Create/link the blacklist files
echo -ne '\ninclude "/var/named/blacklisted.conf";' >> /etc/named.conf;
touch /var/named/blacklisted.conf;

#Create the blacklist zone file
echo -ne "\$TTL 14400
@       IN SOA $YOUR_DOMAIN. $YOUR_DOMAIN ( 2003052800  86400  300  604800  3600 )
@       IN      NS   $YOUR_DOMAIN.
@       IN      A    $YOUR_IP
*       IN      A    $YOUR_IP" > /var/named/blacklisted.db;

#Fix the permissions on the blacklist files
chgrp named /var/named/blacklisted.*;
chmod 660 /var/named/blacklisted.*;

#Set the server’s domain resolution name servers
perl -pi -e 's/^(?!;)/;/gm' /etc/resolv.conf;
echo -ne '\nnameserver localhost' >> /etc/resolv.conf;

#Run a test
echo 'zone "example.com" { type master; file "blacklisted.db"; };' >> \
    /var/named/blacklisted.conf;
service named reload;
FOUND_IP=`dig -t A example.com | grep -ioP "^example\.com\..*?"'in\s+a\s+[\d\.:]+' | \
     grep -oP '[\d\.:]+$'`;
if [ "$YOUR_IP" == "$FOUND_IP" ]
then
  echo 'Success: Example domain matches your given IP' > /dev/stderr;
else
  echo 'Warning: Example domain does not match your given IP' > /dev/stderr;
fi

#Have the server return a message when a blacklisted domain is accessed
echo "$BLOCKED_ERROR_MESSAGE" > /var/www/html/index.html;
perl -0777 -pi -e 's~(<Directory "/var/www/html">.*?\n\s*AllowOverride).*?\n~$1 All~s' \
     /etc/httpd/conf/httpd.conf;
echo -n 'RewriteEngine on
RewriteCond %{REQUEST_URI} !index.html
RewriteCond %{REQUEST_URI} !AddRules/
RewriteRule ^(.*)$ /index.html [L]' > /var/www/html/.htaccess;
service httpd graceful;

#Create a script that allows apache to refresh the name server’s settings
mkdir /var/www/html/AddRules;
echo '/sbin/service named reload' > /var/www/html/AddRules/restart_named;
chmod 755 /var/www/html/AddRules/restart_named;

echo 'apache ALL=(root) NOPASSWD:/var/www/html/AddRules/restart_named
Defaults!/var/www/html/AddRules/restart_named !requiretty' >> /etc/sudoers;

#Create a script that allows the user to add, remove, and list the blacklisted domains
echo -n $'<?php
//Get old domains
$BlockedFile=\'/var/named/blacklisted.conf\';
$CurrentZones=Array();
foreach(explode("\\n", file_get_contents($BlockedFile)) as $Line)
        if(preg_match(\'/^zone "([\\w\\._-]+)"/\', $Line, $Results))
                $CurrentZones[]=$Results[1];

//List domains
if(isset($_REQUEST[\'List\']))
        return print implode(\'<br>\', $CurrentZones);

//Get new domains
if(!isset($_REQUEST[\'Domains\']))
        return print \'Missing Domains\';
$Domains=$_REQUEST[\'Domains\'];
if(!preg_match(\'/^[\\w\\._-]+(,[\\w\\._-]+)*$/uD\', $Domains))
        return print \'Invalid domains string\';
$Domains=explode(\',\', $Domains);

//Remove domains
if(isset($_REQUEST[\'Remove\']))
{
        $CurrentZones=array_flip($CurrentZones);
        foreach($Domains as $Domain)
                unset($CurrentZones[$Domain]);
        $FinalDomainList=array_keys($CurrentZones);
}
else //Combine domains
        $FinalDomainList=array_unique(array_merge($Domains, $CurrentZones));

//Output to the file
$FinalDomainData=Array();
foreach($FinalDomainList as $Domain)
    $FinalDomainData[]="zone \\"$Domain\\" { type master; file \\"blacklisted.db\\"; };";
file_put_contents($BlockedFile, implode("\\n", $FinalDomainData));

//Reload named
print `sudo /var/www/html/AddRules/restart_named`;
?>' > /var/www/html/AddRules/index.php;

usermod -a -G named apache;
service httpd graceful;

#Password protect the domain update script
echo -n 'AuthType Basic
AuthName "Admins Only"
AuthUserFile "/var/www/html/AddRules/.htpasswd"
require valid-user' > /var/www/html/AddRules/.htaccess;

htpasswd -bc /var/www/html/AddRules/.htpasswd "$ADDRULES_USERNAME" "$ADDRULES_PASSWORD";

echo 'Script complete';