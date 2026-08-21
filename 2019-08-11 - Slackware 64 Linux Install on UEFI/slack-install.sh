#Login as root
#sudo su

#Settings
InstallDVDName=SlackDVD #This is whatever you named your slackware usb key

#/mnt will contain the new file system for running the setup
cd /mnt

#Extract the initrd.img from the slackware dvd into /mnt
cat /media/ubuntu/$InstallDVDName/isolinux/initrd.img | gzip -d | cpio -i

#Bind special linux directories into the /mnt folder
for i in proc sys dev tmp; do mount -o bind /$i ./$i; done

#Mount the cdrom folder into /mnt/cdrom
rm cdrom
mount -o bind /media/ubuntu/$InstallDVDName/ ./cdrom

#Set /mnt as our actaul (ch)root
chroot .

#Run the slackware setup
usr/lib/setup/setup

#NOTE: When installing, your package source directory is at /cdrom/slackware64