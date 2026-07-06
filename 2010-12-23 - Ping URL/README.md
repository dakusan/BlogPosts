<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Ping_URL">Ping URL</a><br>
Date: 2010-12-23<br>
SubTitle: Cause Python is quick to program in and can make executables<br>
Tags: Examples, Python<br><br>

<p>The following is a Python script that automatically pings a requested web address at a given interval. It was made as a quick favor for a friend. Here is the downloadable <a href="PingURL_v1.0_Binary.zip">Windows binary</a>.</p>
The Configuration Options File (PingURL.cfg) contains 2 lines:
<ol>
	<li>The URL to ping (The default pings the <a href="http://www.castledragmire.com/Projects/Other_Web_Scripts#GetIP">GetIP</a> script)</li>
	<li>The millisecond interval between pings (Default=600000=10 minutes)</li>
</ol><br>

</div>

```python
#!python
from sys import stderr
from urllib import urlopen
from time import localtime, strftime, sleep

#Main function
def Main():
    #Open the settings file
    SettingFileName='PingURL.cfg';
    Settings=[] #Blank out settings file variable in case it can't be opened
    try:
        Settings=open(SettingFileName, 'r').readlines()
    except IOError as (errno, strerror):
        stderr.write('Cannot open {0} configuration file: I/O error({1}): {2}\n'.format(SettingFileName, errno, strerror))
        return

    #Confirm valid settings were passed
    if(len(Settings)<2):
        stderr.write('Not enough settings found in settings file\n')
        return
    try:
        IntervalTime=int(Settings[1])
    except:
        stderr.write('Invalid interval time\n')
        return
        
    #Ping the URL indefinitely
    while(True):
        try:
            URLText=urlopen(Settings[0]).read()
        except:
            URLText='READ FAILED'
        print 'URL Pinged At {0}: {1}'.format(strftime('%Y-%m-%d %H:%M:%S', localtime()), URLText)
        sleep(IntervalTime/1000)

#Run the program
Main()
```