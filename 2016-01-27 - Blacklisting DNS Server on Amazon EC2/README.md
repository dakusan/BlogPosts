<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Blacklisting_DNS_Server_on_Amazon_EC2">Blacklisting DNS Server on Amazon EC2</a><br>
Date: 2016-01-27<br>
Tags: Amazon, Blacklist, DNS, EC2, Examples, Linux, Security, Servers<br><br>

<p>Amazon EC2 is a great resource for cheap virtual servers to do simple things, like DNS or (low bandwidth) VPNs. I had the need this morning to set up a DNS server for a company which needed to blacklist a list of domains. The simplest way to do this is by editing all the computers’ hostfiles, but that method leaves a lot to be desired. Namely, blocking entire domains (as opposed to single subdomains), and deploying changes. Centralizing in a single place makes the job instant, immediate, and in the end, faster.</p>
<p>The following are the steps I used to set this up on an EC2 server. All command line instructions are followed by a single command you can run to execute the step. There is a <a href="#full-serverside-script">full script below, at the end of the post</a>, containing all steps from when you first login to SSH ("Login to root") to the end.</p><br>

<p>I am not going to go into the details of setting up an <a href="https://console.aws.amazon.com/ec2/">EC2 instance</a>, as that information can be found elsewhere. I will also be skipping over some of the more obvious steps. Just create a default EC2 instance with the “Amazon Linux AMI”, and I will list all the changes that need to be made beyond that.</p>
<ul>
	<li><b>Creating the instance</b><ul>
		<li>For the first year, for the instance type, you might as well use a t2.micro, as it is free. After that, a t2.nano (which is a new lower level) currently at $56.94/year ($0.0065/Hour), should be fine.</li>
		<li>After you select your instance type, click “Review and Launch” to launch the instance with all of the defaults.</li>
		<li>After the confirmation screen, it will ask you to create a key pair. You can see other tutorials about this and how it enables you to log into your instance.</li>
	</ul></li>
	<li><b>Edit the security group</b><ul>
		<li>Next, you need to edit the security group for your instance to allow incoming connections.</li>
		<li>Go to “Instances” under the “Instances” group on the left menu, and click your instance.</li>
		<li>In the bottom of the window, in the “Descriptions” tab, click the link next to “Security Groups”, which will bring you to the proper group in the security groups tab.</li>
		<li>Right click it and “Edit inbound Rules”.</li>
		<li>Make sure it has the following rules with Source=Anywhere: ALL ICMP [For pinging], SSH, HTTP, DNS (UDP), DNS (TCP)</li>
	</ul></li>
	<li><b>Assign a permanent IP to your instance</b><ul>
		<li>To do this, click the “Elastic IPs” under “Network &amp; Security” in the left menu.</li>
		<li>Click “Allocate New Address”.</li>
		<li>After creating it, right click the new address, then “Associate Address”, and assign it to your new instance.</li>
		<li>You should probably set this IP up as an A record somewhere. I will refer to this IP as dns.yourdomain.com from now on.</li>
	</ul></li>
	<li><b>Login to root</b><ul>
		<li>SSH into your instance as the ec2-user via “ssh ec2-user@dns.yourdomain.com”. If in windows, you could also use <a href="http://www.putty.org/">putty</a>.</li>
		<li>Sudo into root via <code>sudo su</code>.</li>
	</ul></li>
	<li><b>Allow root login</b><ul>
		<li>At this point, I recommend setting it up so you can directly root into the server. <b><i>Warning</i></b>: some people consider this a security risk.</li>
		<li>Copy your key pair(s) to the root user via <code>cat /home/ec2-user/.ssh/authorized_keys &gt; /root/.ssh/authorized_keys</code></li>
		<li>Set SSHD to permit root logins by changing the PermitRootLogin variable to “yes” in /etc/ssh/sshd_config. A quick command to do this is <code>perl -pi -e 's/^\s*#?\s*PermitRootLogin.*$/PermitRootLogin yes/igm' /etc/ssh/sshd_config</code>, and then reload the SSHD config with <code>service sshd reload</code>. Make sure to attempt to directly log into SSH as root before exiting your current session to make sure you haven’t locked yourself out.</li>
	</ul></li>
	<li><b>Install apache (the web server), bind/named (the DNS server), and PHP (a scripting language)</b><ul>
		<li><code>yum -y install bind httpd php</code></li>
	</ul></li>
	<li><b>Start and set services to run at boot</b><ul>
		<li><code>service httpd start; service named start; chkconfig httpd on; chkconfig named on;</code></li>
	</ul></li>
	<li><b>Set the DNS server to be usable by other computers</b><ul>
		<li>Edit /etc/named.conf and change the 2 following lines to have the value “any”: “listen-on port 53” and “allow-query”</li>
		<li><code>perl -pi -e 's/^(\s*(?:listen-on port 53|allow-query)\s*{).*$/$1 any; };/igm' /etc/named.conf; service named reload;</code></li>
	</ul></li>
	<li><b>Point the DNS server to the blacklist files</b><ul>
		<li>This is done by adding <code>include "/var/named/blacklisted.conf";</code> to /etc/named.conf</li>
		<li><code>echo -ne '\ninclude "/var/named/blacklisted.conf";' &gt;&gt; /etc/named.conf</code></li>
	</ul></li>
	<li><b>Create the blacklist domain list file</b><ul>
		<li><code>touch /var/named/blacklisted.conf</code></li>
	</ul></li>
	<li><b>Create the blacklist zone file</b><ul>
		<li>Put the following into /var/named/blacklisted.db . Make sure to change dns.yourdomain.com to your domain (or otherwise, “localhost”), and 1.1.1.1 to dns.yourdomain.com’s (your server’s) IP address. Make sure to keep all periods intact.<br>
<pre><code>$TTL 14400
@       IN SOA <kbd>dns.yourdomain.com</kbd>. <kbd>dns.yourdomain.com</kbd> ( 2003052800  86400  300  604800  3600 )
@       IN      NS   <kbd>dns.yourdomain.com</kbd>.
@       IN      A    <kbd>1.1.1.1</kbd>
*       IN      A    <kbd>1.1.1.1</kbd></code></pre>
		</li>
		<li>The first 2 lines tell the server the domains belong to it. The 3rd line sets the base blacklisted domain to your server’s IP. The 4th line sets all subdomains of the blacklisted domain to your server’s IP.</li>
		<li>This can be done via (Update the first line with your values)<br>
<pre><code>YOURDOMAIN="<kbd>dns.yourdomain.com</kbd>"; YOURIP="<kbd>1.1.1.1</kbd>";
echo -ne "\$TTL 14400\n@       IN SOA $YOURDOMAIN. $YOURDOMAIN ( 2003052800  86400  300  604800  3600 )\n@       IN      NS   $YOURDOMAIN.\n@       IN      A    $YOURIP\n*       IN      A    $YOURIP" > /var/named/blacklisted.db;</code></pre>
		</li>
	</ul></li>
	<li><b>Fix the permissions on the blacklist files</b><ul>
		<li><code>chgrp named /var/named/blacklisted.*; chmod 660 /var/named/blacklisted.*;</code></li>
	</ul></li>
	<li><b>Set the server’s domain resolution name servers</b><ul>
		<li>The server always needs to look at itself before other DNS servers. To do this, comment out everything in /etc/resolv.conf and add to it “nameserver localhost”. This is not the best solution. I’ll find something better later.</li>
		<li><code>perl -pi -e 's/^(?!;)/;/gm' /etc/resolv.conf; echo -ne '\nnameserver localhost' &gt;&gt; /etc/resolv.conf</code></li>
	</ul></li>
	<li><b>Run a test</b><ul>
		<li>At this point, it’s a good idea to make sure the DNS server is working as intended. So first, we’ll add an example domain to the DNS server.</li>
		<li>Add the following to /var/named/blacklisted.conf and restart named to get the server going with example.com: <code>zone "example.com" { type master; file "blacklisted.db"; };</code>”</li>
		<li><code>echo 'zone "example.com" { type master; file "blacklisted.db"; };' &gt;&gt; /var/named/blacklisted.conf; service named reload;</code></li>
		<li>Ping “test.example.com” and make sure it’s IP is your server’s IP</li>
		<li>Set your computer’s DNS to your server’s IP in your computer’s network settings, ping “test.example.com” from your computer, and make sure the returned IP is your server’s IP. If it works, you can restore your computer’s DNS settings.</li>
	</ul></li>
	<li><b>Have the server return a message when a blacklisted domain is accessed</b><ul>
			<li>Add your message to /var/www/html</li>
			<li><code>echo '<b>Domain is blocked</b>' &gt; /var/www/html/index.html</code></li>
			<li>Set all URL paths to show the message by adding the following to the /var/www/html/.htaccess file<br>
<pre><code>RewriteEngine on
RewriteCond %{REQUEST_URI} !index.html
RewriteCond %{REQUEST_URI} !AddRules/
RewriteRule ^(.*)$ /index.html [L]</code></pre></li>
			<li><code>echo -ne 'RewriteEngine on\nRewriteCond %{REQUEST_URI} !index.html\nRewriteCond %{REQUEST_URI} !AddRules/\nRewriteRule ^(.*)$ /index.html [L]' &gt; /var/www/html/.htaccess</code></li>
			<li>Turn on AllowOverride in the /etc/httpd/conf/httpd.conf for the document directory (/var/www/html/) via <code> perl -0777 -pi -e 's~(&lt;Directory "/var/www/html"&gt;.*?\n\s*AllowOverride).*?\n~$1 All~s' /etc/httpd/conf/httpd.conf</code></li>
			<li>Start the server via <code>service httpd graceful</code></li>
	</ul></li>
	<li><b>Create a script that allows apache to refresh the name server’s settings</b><ul>
		<li>Create a script at /var/www/html/AddRules/restart_named with <code>/sbin/service named reload</code> and set it to executable</li>
		<li><code>mkdir /var/www/html/AddRules; echo '/sbin/service named reload' &gt; /var/www/html/AddRules/restart_named; chmod 755 /var/www/html/AddRules/restart_named</code></li>
		<li>Allow the user to run the script as root by adding to /etc/sudoers <code>apache ALL=(root) NOPASSWD: /var/www/html/AddRules/restart_named</code> and <code>Defaults!/var/www/html/AddRules/restart_named !requiretty</code> </li>
		<li><code>echo -e 'apache ALL=(root) NOPASSWD:/var/www/html/AddRules/restart_named\nDefaults!/var/www/html/AddRules/restart_named !requiretty' &gt;&gt; /etc/sudoers</code></li>
	</ul></li>
	<li><b>Create a script that allows the user to add, remove, and list the blacklisted domains</b><ul>
		<li>Add the following to /var/www/html/AddRules/index.php (one line command not given. You can use “nano” to create it)
</li></ul></li></ul></div>

```php
<?php
//Get old domains
$BlockedFile='/var/named/blacklisted.conf';
$CurrentZones=Array();
foreach(explode("\n", file_get_contents($BlockedFile)) as $Line)
        if(preg_match('/^zone "([\w\._-]+)"/', $Line, $Results))
                $CurrentZones[]=$Results[1];

//List domains
if(isset($_REQUEST['List']))
        return print implode('<br>', $CurrentZones);

//Get new domains
if(!isset($_REQUEST['Domains']))
        return print 'Missing Domains';
$Domains=$_REQUEST['Domains'];
if(!preg_match('/^[\w\._-]+(,[\w\._-]+)*$/uD', $Domains))
        return print 'Invalid domains string';
$Domains=explode(',', $Domains);

//Remove domains
if(isset($_REQUEST['Remove']))
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
        $FinalDomainData[]=
                "zone \"$Domain\" { type master; file \"blacklisted.db\"; };";
file_put_contents($BlockedFile, implode("\n", $FinalDomainData));

//Reload named
print `sudo /var/www/html/AddRules/restart_named`;
?>
```
<div><ul><ul>
		<li>Add the “apache” user to the “named” group so the script can update the list of domains in /var/named/blacklisted.conf via <code>usermod -a -G named apache; service httpd graceful;</code></li>
	</ul></li>
	<li><b>Run the domain update script</b><ul>
		<li>To add a domain (separate by commas): http://<kbd>dns.yourdomain.com</kbd>/AddRules/?Domains=domain1.com,domain2.com</li>
		<li>To remove a domain (add “Remove&amp;” after the “?”): http://<kbd>dns.yourdomain.com</kbd>/AddRules/?Remove&amp;Domains=domain1.com,domain2.com</li>
		<li>To list the domains: http://<kbd>dns.yourdomain.com</kbd>/AddRules/?List</li>
	</ul></li>
	<li><b>Password protect the domain update script</b><ul>
		<li>Add to AddRules/.htaccess the following<br>
<pre><code>AuthType Basic
AuthName "Admins Only"
AuthUserFile "/var/www/html/AddRules/.htpasswd"
require valid-user</code></pre></li>
		<li><code>echo -ne 'AuthType Basic\nAuthName "Admins Only"\nAuthUserFile "/var/www/html/AddRules/.htpasswd"\nrequire valid-user' &gt; /var/www/html/AddRules/.htaccess</code></li>
		<li><b><i>Warning</i></b>: Putting the password file in an http accessible directory is a security risk. I just did this for sake of organization.</li>
		<li>Create the user+password via <code>htpasswd -bc /var/www/html/AddRules/.htpasswd <kbd>USERNAME</kbd></code> and then entering the password</li>
	</ul></li>
</ul><br><br>

<b><i>[Edit on 2016-01-30 @ noon]</i></b></div>
<p>To permanently set “localhost” as the resolver DNS, add “DNS1=localhost” to “/etc/sysconfig/network-scripts/ifcfg-eth0”. I have not yet confirmed this edit.</p>

## Security Issue
<p>Soon after setting up this DNS server, it started getting hit by a <a href="https://www.us-cert.gov/ncas/alerts/TA13-088A">DNS amplification attack</a>. As the server is being used as a client’s DNS server, turning off recursion is not available. The best solution is to limit the people who can query the name server via an access list (usually a specific subnet), but that would very often not be an option either. The solution I currently have in place, which I have not actually verified if it works, is to add a forced-forward rule which only makes external requests to the name server provided by Amazon. To do this, get the name server’s IP from /etc/resolv.conf (it should be commented from an earlier step). Then add the following to your named.conf in the “options” section.
</div>

```
	forwarders {
		DNS_SERVER_IP;
	};
	forward only;
```

<div>

<p>After I added this rule, external DNS requests stopped going through completely. To fix this, I turned “dnssec-validation” to “no” in the named.conf. Don’t forget to restart the service once you have made your changes.</p>
<b></i>[End of edit]</i></b>

## Full serverside script
<i>Make sure to run this as root (login as root or sudo it)</i><br><br>
<i>Download the script <a href="dnsblacklist_install.sh">here</a>. Make sure to chmod and sudo it when running. <code>chmod +x dnsblacklist_install.sh; sudo ./dnsblacklist_install.sh;</code></i><br><br>

```bash
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
```