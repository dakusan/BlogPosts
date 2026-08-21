<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Slackware_64_Linux_Install_on_UEFI">Slackware 64 Linux Install on UEFI</a><br>
Date: 2019-08-11<br>
Tags: Boot Loader, Install, Linux, Operating Systems, Quick Fixes, Slackware, UEFI<br><br>

<p>I recently tried to install Slackware 4.2 64-bit (Linux) onto a new mini PC I just bought. The new PC only supports UEFI so I had major issues getting the darn setup on the install cd to actually run. I never DID actually get the install cd to boot properly on the system, so I used an alternative. While the slack install usb key was in, I also added and loaded up an ubuntu live cd usb key. The following is what I used to run the slackware setup in Ubuntu.</p>

```bash
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
```