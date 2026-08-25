#!/bin/bash

#Shared environment variables:
KP_CUSTOM=/home/$USER/keepass-browser
KP_JSON_NAME=org.keepassxc.keepassxc_browser.json

#Chrome environment variables:
KP_FLATPAK_PACKAGE=io.github.ungoogled_software.ungoogled_chromium
KP_JSON_START=~/.config/chromium/NativeMessagingHosts
KP_JSON_END=~/.var/app/$KP_FLATPAK_PACKAGE/config/chromium/NativeMessagingHosts

#Firefox environment variables:
KP_FLATPAK_PACKAGE=org.mozilla.firefox
KP_JSON_START=~/.mozilla/native-messaging-hosts
KP_JSON_END=~/.var/app/$KP_FLATPAK_PACKAGE/.mozilla/native-messaging-hosts

#Confirm browser integration
if [[ ! -f "$KP_JSON_START" || ! -f "$KP_JSON_NAME" ]]; then
	echo "Error: Required files do not exist. Make sure you have enabled KeePassXD browser integration before continuing"
	exit 1
fi

###Set up the needed files in the sandbox
#Put KeePass proxy and needed library files in user directory
mkdir -p $KP_CUSTOM/lib
mkdir -p $KP_JSON_END #Needed for firefox
cp /usr/bin/keepassxc-proxy $KP_CUSTOM/
rsync -a /usr/lib/x86_64-linux-gnu/libicudata* /usr/lib/x86_64-linux-gnu/libicuuc* /usr/lib/x86_64-linux-gnu/libicui* /usr/lib/x86_64-linux-gnu/libdouble* /usr/lib/x86_64-linux-gnu/libsodium* /usr/lib/x86_64-linux-gnu/libQt5* $KP_CUSTOM/lib

#Copy the JSON file to the Flatpak app directory and change the executable path in the file
cp $KP_JSON_START/$KP_JSON_NAME $KP_JSON_END/
sed -i "s/\/usr\/bin\//"$(echo $KP_CUSTOM | sed 's_/_\\/_g')"\//" $KP_JSON_END/$KP_JSON_NAME

###Add permissions to the Flatpak
flatpak override --user --filesystem=$KP_CUSTOM:ro $KP_FLATPAK_PACKAGE #Only required if home directory is not shared to the Flatpak
flatpak override --user --filesystem=xdg-run/org.keepassxc.KeePassXC.BrowserServer:ro $KP_FLATPAK_PACKAGE
flatpak override --user --env=LD_LIBRARY_PATH=$(flatpak info --show-permissions $KP_FLATPAK_PACKAGE | grep -oP '(?<=LD_LIBRARY_PATH=).*')";$KP_CUSTOM/lib" $KP_FLATPAK_PACKAGE