#!/bin/bash
#-----------------------------------Variables----------------------------------
#Current root drive
CurPart="nvme0n1p3" #The current drive partition in /dev. This example uses nvme disk #0 partition #3
CurCryptVol="vgmint" #What you named your LVM under LUKS
CurCryptRoot="root" #What you named your root partition under the LVM
CurCryptRootSubVol="/" #The path of the subvolume that is used as the root partition. For example, I use “@”
CurCryptSwap="swap_1" #What you named your swap partition under the LVM
CurCryptExtra="extra" #What you WANT to name your extra partition under the LVM
CurCryptExtraSize="100M" #How big you want your extra partition that will hold your key file
CurKeyPath="" #The path to a key file that will unlock both drives. If left blank then one will be created

#New drive
NewDrive="nvme1n1" #The new drive in /dev. This example uses nvme disk #1
NewPart="nvme1n1p3" #The new partition in /dev. You will be creating this with the cfdisk. This example uses nvme disk#1 partition#3
NewCryptName="raid1_crypt" #What we’ll name the root LUKS partition (no LVM)

#Other variables you do not need to set
CurMount="/mnt/primary"
ExtraMountPath="$CurMount/mnt/extra"
BtrfsReleasePath="kdave/btrfs-progs"
BtrfsReleaseFile="btrfs.box.static"
DriveKeyName="drivekey"

echo "---------------------------------Update BTRFS---------------------------------"
echo "Make sure you are using the latest btrfs-progs"
cd "$(dirname "$(which btrfs)")"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$BtrfsReleasePath/releases/latest" | grep tag_name | cut -d \" -f4)
wget "https://github.com/$BtrfsReleasePath/releases/download/$LATEST_RELEASE/$BtrfsReleaseFile"
chmod +x "$BtrfsReleaseFile"

echo "Link all btrfs programs to btrfs.box.static. Rename old files as .old.FILENAME"
if ! [ -L ./btrfs ]; then
  for v in $(\ls btrfs*); do
    if [ "$v" != "$BtrfsReleaseFile" ]; then
      mv "$v" ".old.$v"
      ln -s "$BtrfsReleaseFile" "$v"
    fi
  done
fi

echo "--------------------------Current drive and key setup-------------------------"
echo "Mount the current root partition"
cryptsetup luksOpen "/dev/$CurPart" "$CurCryptVol"
vgchange -ay "$CurCryptVol"
mkdir -p "$CurMount"
mount -o "subvol=$CurCryptRootSubVol" "/dev/$CurCryptVol/$CurCryptRoot" "$CurMount"

echo "If the extra volume has not been created, then resize the swap and create it"
if ! [ -e "/dev/$CurCryptVol/$CurCryptExtra" ]; then
  lvremove -y "/dev/$CurCryptVol/$CurCryptSwap"

  lvcreate -n "$CurCryptExtra" -L "$CurCryptExtraSize" "$CurCryptVol"
  mkfs.ext4 "/dev/$CurCryptVol/$CurCryptExtra"

  lvcreate -n "$CurCryptSwap" -l 100%FREE "$CurCryptVol"
  mkswap "/dev/$CurCryptVol/$CurCryptSwap"
fi

echo "Make sure the key file exists, if it does not, either copy it (if given in $CurKeyPath) or create it"
mkdir -p "$ExtraMountPath"
mount "/dev/$CurCryptVol/$CurCryptExtra" "$ExtraMountPath"
if ! [ -e "$ExtraMountPath/$DriveKeyName" ]; then
  if [ "$CurKeyPath" != "" ]; then
    if ! [ -e "$CurKeyPath" ]; then
      echo "Not found: $CurKeyPath"
      exit 1
    fi
    cp "$CurKeyPath" "$ExtraMountPath/$DriveKeyName"
  else
    openssl rand -out "$ExtraMountPath/$DriveKeyName" 512
  fi
  chmod 400 "$ExtraMountPath/$DriveKeyName"
  chown root:root "$ExtraMountPath/$DriveKeyName"
fi

echo "Make sure the key file works on the current drive"
if cryptsetup --test-passphrase luksOpen --key-file "$ExtraMountPath/$DriveKeyName" "/dev/$CurPart" test; then
  echo "Keyfile successfully opened the LUKS partition."
  #cryptsetup luksClose test #This doesn’t seem to be needed
else
  echo "Adding keyfile to the LUKS partition"
  cryptsetup luksAddKey "/dev/$CurPart" "$ExtraMountPath/$DriveKeyName"
fi

echo "--------------------------------New drive setup-------------------------------"
echo "Use cfdisk to set the new disk as GPT and add partitions."
echo "Make sure to mark the partition you want to use for the raid disk as type “Linux Filesystem”."
echo "Also make it the same size as /dev/$CurPart to avoid errors"
cfdisk "/dev/$NewDrive"

echo "Encrypt the new partition"
cryptsetup luksFormat "/dev/$NewPart"

echo "Open the encrypted partition"
cryptsetup luksOpen "/dev/$NewPart" "$NewCryptName"

echo "Add the key to the partition"
cryptsetup luksAddKey "/dev/$NewPart" "$ExtraMountPath/$DriveKeyName"

echo "Add the new partition to the root btrfs file system"
btrfs device add "/dev/mapper/$NewCryptName" "$CurMount"

echo "Convert to RAID1"
btrfs balance start -dconvert=raid1 -mconvert=raid1 "$CurMount"

echo "Confirm both disks are in use"
btrfs filesystem usage "$CurMount"

echo "--------------------Booting script to load encrypted drives-------------------"
echo "Get the UUID of the second btrfs volume"
Drive2_UUID=$(lsblk -o UUID -d "/dev/$NewPart" | tail -n1)

echo "Create a script to open your second luks volumes before mounting the partition"
echo "Note: In some scenarios this may need to go into “scripts/local-premount” instead of “scripts/local-bottom”"
cat <<EOF > "$CurMount/etc/initramfs-tools/scripts/local-bottom/unlock_drive2"
#!/bin/sh
PREREQ=""

prereqs()
{
    echo "\$PREREQ"
}

case "\$1" in
    prereqs)
        prereqs
        exit 0
        ;;
esac

. /scripts/functions
cryptroot-unlock
vgchange -ay "$CurCryptVol"
mkdir -p /mnt/keyfile
mount "/dev/$CurCryptVol/$CurCryptExtra" /mnt/keyfile
cryptsetup luksOpen /dev/disk/by-uuid/$Drive2_UUID "$NewCryptName" "--key-file=/mnt/keyfile/$DriveKeyName"
umount /mnt/keyfile
rmdir /mnt/keyfile

mount -t btrfs -o "subvol=$CurCryptRootSubVol" "/dev/$CurCryptVol/$CurCryptRoot" /root

#If you are weird like me and /usr is stored elsewhere, here is where you would need to mount it.
#It cannot be done through your fstab in this setup.
#mount --bind /root/sub/sys/usr /root/usr

mount --bind /dev /root/dev
mount --bind /proc /root/proc
mount --bind /sys /root/sys
EOF

chmod 755 "$CurMount/etc/initramfs-tools/scripts/local-bottom/unlock_drive2"

echo "--------------------Setup booting from the root file system-------------------"
echo "Prepare a chroot environment"
for i in dev dev/pts proc sys run tmp; do
  mount -o bind /$i "$CurMount/$i"
done

echo "Run commands in the chroot environment to update initramfs and grub"
chroot "$CurMount" <<EOF
echo "Mount the other partitions (specifically for “boot” and “boot/efi”)"
mount -a

echo "Update initramfs and grub"
update-initramfs -u -k all
update-grub
EOF

echo "-----------------------------------Finish up----------------------------------"
echo "Reboot and pray"
reboot