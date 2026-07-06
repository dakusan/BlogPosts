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