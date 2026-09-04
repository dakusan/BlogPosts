<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Installing_KeePassXC_%26_a_Flatpak_sandboxed_browser">Installing KeePassXC &amp; a Flatpak sandboxed browser</a><br>
Date: 2024-07-13<br>
Tags: Chrome, Diagnosing, Firefox, Flatpak, KeePass, KeePassXC, Linux, Linux Migration, Password Managers, Sandbox, Security, Web browsers<br><br>
</div>

<p>Since I have been moving my primary workstation to Linux, I have been in need of finding a new password manager after <a href="../2024-07-04 - Roboform (offline mode) on Linux with Wine">Roboform failed me</a>. I decided to try out <a href="https://keepassxc.org/">KeePassXC</a> but was having trouble since my <a href="https://github.com/ungoogled-software/ungoogled-chromium">Ungoogled Chromium</a> browser was sandboxed in <a href="https://www.flatpak.org/">Flatpak</a>. The <a href="https://unix.stackexchange.com/questions/584521/how-to-setup-firefox-and-keepassxc-in-a-flatpak-and-get-the-keepassxc-browser-ad">instructions for this on Unix Stack Exchange</a> were no longer working so I had to dig into it myself.</p><br>
<p>The 2 primary problems were that:<br>
1) The browsers’ sandboxes did not have the QT libs.<br>
2) The interprocess communication pipe socket had been renamed from <b>kpxc_server</b> to <b>org.keepassxc.KeePassXC.BrowserServer</b>.
</p><br>

<p>The following are the instructions to get KeePassXC running in Flatpak versions of both Chrome and Firefox. This was tested on Linux Mint 21.3 with both Ungoogled Chromium and Firefox. You will need to change the <b>KP_FLATPAK_PACKAGE</b> if you use other versions of Chrome.</p><br>

<ol>
	<li>
		<ins>Run the relevant below environment variables in your command shell before running the commands in the following steps:</ins><br>
</li></ol>

<pre><code><ins>#Shared environment variables:</ins>
KP_CUSTOM=/home/<kbd>$USER</kbd>/keepass-browser
KP_JSON_NAME=org.keepassxc.keepassxc_browser.json

<ins>#Chrome environment variables:</ins>
KP_FLATPAK_PACKAGE=io.github.ungoogled_software.ungoogled_chromium
KP_JSON_START=~/.config/chromium/NativeMessagingHosts
KP_JSON_END=~/.var/app/<kbd>$KP_FLATPAK_PACKAGE</kbd>/config/chromium/NativeMessagingHosts

<ins>#Firefox environment variables:</ins>
KP_FLATPAK_PACKAGE=org.mozilla.firefox
KP_JSON_START=~/.mozilla/native-messaging-hosts
KP_JSON_END=~/.var/app/<kbd>$KP_FLATPAK_PACKAGE</kbd>/.mozilla/native-messaging-hosts</code></pre>

<ol start=2>
	<li>
		<ins>Install and enable the browser extension:</ins><br>
		<b>KeePassXD &gt; Tools &gt; Settings &gt; Browser Integration</b>:<br>
		<ul>
			<li>Check “Enable Browser Integration”</li>
			<li>Check “Chromium” and/or “Firefox”</li>
			<li>Download the plugin listed on this screen in your browser</li>
			<li>Click "OK"</li>
		</ul>
		Note: This creates <b>$KP_JSON_START/$KP_JSON_NAME</b><br><br>
	</li>
	<li>
		<ins>Set up the needed files in the sandbox:</ins><br>
</li></ol>

<pre><code><ins>#Put KeePass proxy and needed library files in user directory</ins>
mkdir -p <kbd>$KP_CUSTOM</kbd>/lib
mkdir -p <kbd>$KP_JSON_END</kbd> #Needed for firefox
cp /usr/bin/keepassxc-proxy <kbd>$KP_CUSTOM</kbd>/
rsync -a /usr/lib/x86_64-linux-gnu/libicudata* /usr/lib/x86_64-linux-gnu/libicuuc* /usr/lib/x86_64-linux-gnu/libicui* /usr/lib/x86_64-linux-gnu/libdouble* /usr/lib/x86_64-linux-gnu/libsodium* /usr/lib/x86_64-linux-gnu/libQt5* <kbd>$KP_CUSTOM</kbd>/lib

<ins>#Copy the JSON file to the Flatpak app directory and change the executable path in the file</ins>
cp <kbd>$KP_JSON_START</kbd>/<kbd>$KP_JSON_NAME</kbd> <kbd>$KP_JSON_END</kbd>/
sed -i "s/\/usr\/bin\//"$(echo <kbd>$KP_CUSTOM</kbd> | sed 's_/_\\/_g')"\//" <kbd>$KP_JSON_END</kbd>/<kbd>$KP_JSON_NAME</kbd></code></pre>

<ol start=4>
	<li>
		<ins>Add permissions to the Flatpak:</ins><br>
</li></ol>

<pre><code>flatpak override --user --filesystem=<kbd>$KP_CUSTOM</kbd>:ro <kbd>$KP_FLATPAK_PACKAGE</kbd> #Only required if home directory is not shared to the Flatpak
flatpak override --user --filesystem=xdg-run/org.keepassxc.KeePassXC.BrowserServer:ro <kbd>$KP_FLATPAK_PACKAGE</kbd>
flatpak override --user --env=LD_LIBRARY_PATH=$(flatpak info --show-permissions <kbd>$KP_FLATPAK_PACKAGE</kbd> | grep -oP '(?<=LD_LIBRARY_PATH=).*')";<kbd>$KP_CUSTOM</kbd>/lib" <kbd>$KP_FLATPAK_PACKAGE</kbd></code></pre>