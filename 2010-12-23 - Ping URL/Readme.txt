Ping URL v1.0.0

Description:
	Periodically pings a requested URL.

Updates/Recent Information:
	See http://www.castledragmire.com/Posts/Ping_URL for up to date information on this project.

Configuration Options (See PingURL.cfg): (This file only supports standard ASCII)
	Line #1: The URL to ping
	Line #2: The millisecond interval between pings (Default=600000=10 minutes)

Dependencies: (When running from source)
	This application uses python and associated libraries.

To Run From Source:
	./PingURL.py

To Compile: (py2exe required)
	python setup.py py2exe