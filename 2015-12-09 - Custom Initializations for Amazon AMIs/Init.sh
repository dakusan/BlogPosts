#!/bin/bash
PrimaryInstance=www.example.com

### Create a custom init file that runs on boot to take care of all the commands that need to be run
#Create the script and make sure the full path (+all other root environment variables) are set when it is ran
echo '#!/bin/bash -l' > /etc/rc.d/init.d/custom_init

#Set the script as executable
chmod +x /etc/rc.d/init.d/custom_init

#Executes it as one of the last scripts on run level 3 (Multi-user mode with networking)
ln -s ../init.d/custom_init /etc/rc.d/rc3.d/S99custom_init

### Allow login via password authentication:
perl -i -pe 's/^PasswordAuthentication.*$/PasswordAuthentication yes/mg' /etc/ssh/sshd_config
service sshd reload

### Sync SSH settings from the PrimaryInstance:
#Remove the known_hosts file, in case something on the PrimaryInstance has changed that would block ssh commands.
rm -f ~/.ssh/known_hosts

#Sync the SSH settings from the PrimaryInstance
rsync -e 'ssh -o StrictHostKeyChecking=no' -a root@"$PrimaryInstance":~/.ssh/ ~/.ssh/

#Sync required files from the PrimaryInstance. In this case, the default web root folder
rsync -at root@"$PrimaryInstance":/var/www/html/ /var/www/html/