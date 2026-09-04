<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Custom_Initializations_for_Amazon_AMIs">Custom Initializations for Amazon AMIs</a><br>
Date: 2015-12-09<br>
Tags: Amazon, AMIs, AWS, CentOS, Cloud, EC2, Examples, Servers, SSH<br><br>
</div>

<p>I was recently hired to move a client's site from our primary server in Houston to the Amazon cloud, as it was about to take a big hit in traffic. The normal setup for this kind of job is pretty straightforward. Move the database over to RDS, set up an AMI of an EC2 instance, a load balancer, and ec2 auto scaling. However, there were a couple of problems I needed to solve this time around for the instances launched via the auto scalar that I had not really needed to do before. This includes syncing the SSH settings and current codebase from the primary instance, as opposed to recreating AMIs every time there was a change. So, long story short, here are the problems and solutions that need to be added before the AMI image is created.</p><br>

<p>This all assumes you are running as root. Most of these commands should work on any Linux distribution that Amazon has default AMIs for, but some of these may only work in the Amazon and CentOS AMIs.</p><br>

### Pre-setup
<ul>
	<li>
		Your first instance that you are creating the AMI from should be a permanent instance. This is important for 2 reasons.
		<ol>
			<li>When changing configurations for the auto scalar, if and when your instances are terminated and recreated, this instance will always be available on the load balancer, so there is no downtime.</li>
			<li>This instance can act as a central repository for other instances to sync from.</li>
		</ol>
		So make sure this instance has an elastic IP assigned to it. From here on out, we will refer to this instance as <b><i>PrimaryInstance</i></b> (you can set this physically in the host file, or change it in all scripts to however you want to refer to your elastic IP [most likely through a DNS domain]).
	</li><li>
		Create your ssh private key for the instances: (For all prompts, use default settings)
		<br><code>ssh-keygen -t rsa -b 4096 -C "your_email@example.com"</code>
	</li><li>
		Make sure your current ssh authorized_keys contains your new ssh private key:
		<br><code>cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys</code>
	</li><li>
		Make sure your ssh known_hosts includes your primary instance, so all future ssh calls to it are automatically accept it as a known host:
		<br><code>ssh <b><i>PrimaryInstance</i></b> -o StrictHostKeyChecking=no</code><br>
		You do not have to finish the login process. This just makes sure our primary instance will be recognized by other instances.
	</li><li>
		Turn on <b><i>PermitRootLogin</i></b> in /etc/ssh/sshd_config and reload the sshd config <code>service sshd reload</code><br>
		I just recommend this because it makes life way, way easier. The scripts below assume that you did this.
	</li>
</ul><br>

### Create a custom init file that runs on boot to take care of all the commands that need to be run
```bash
#Create the script and make sure the full path (+all other root environment variables) are set when it is ran
echo '#!/bin/bash -l' > /etc/rc.d/init.d/custom_init

#Set the script as executable
chmod +x /etc/rc.d/init.d/custom_init

#Executes it as one of the last scripts on run level 3 (Multi-user mode with networking)
ln -s ../init.d/custom_init /etc/rc.d/rc3.d/S99custom_init
```
All of the below commands in this post will go into this script.<br><br>

### Allow login via password authentication
```bash
perl -i -pe 's/^PasswordAuthentication.*$/PasswordAuthentication yes/mg' /etc/ssh/sshd_config
service sshd reload
```
Notes:<br>
You may not want to do this. It was just required by my client in this case.<br>
This is required in the startup script because Amazon likes to mess with the sshd_config (and authorized_keys) in new instances it boots.<br><br>

### Sync SSH settings from the <b><i>PrimaryInstance</i></b>
<pre><code>#Remove the known_hosts file, in case something on the <b><i>PrimaryInstance</i></b> has changed that would block ssh commands.
rm -f ~/.ssh/known_hosts

#Sync the SSH settings from the <b><i>PrimaryInstance</i></b>
rsync -e 'ssh -o StrictHostKeyChecking=no' -a root@<b><i>PrimaryInstance</i></b>:~/.ssh/ ~/.ssh/</code></pre>

### Sync required files from the <b><i>PrimaryInstance</i></b>. In this case, the default web root folder
<pre><code>rsync -at root@<b><i>PrimaryInstance</i></b>:/var/www/html/ /var/www/html/</code></pre>

<p>That's it for the things that need to be configured/added to the instance. From there, create your AMI and launch config, and create/modify your launch group and load balancer.</p><br>

<p>Also, as a very important note about your load balancer, make sure if you are mirroring its IP on another domain to use a CNAME record, and not the IP in an A record, as the load balancer IP is subject to change.</p>

Continued at <a href="../2015-12-10 - Syncing Amazon EC2 Instances">Syncing Amazon EC2 Instances</a>