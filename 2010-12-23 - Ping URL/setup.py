#Overwrite py2exe version to include Author in file version information
import py2exe.resources.VersionInfo
OldVersionInfo=py2exe.resources.VersionInfo.Version
class OverwriteVersionInfo(OldVersionInfo):
    def __init__(self, version, **NameParameters):
        OldVersionInfo.__init__(self, version, **NameParameters)
        self.strings.insert(0, ("Author", "Dakusan"))
py2exe.resources.VersionInfo.Version=OverwriteVersionInfo

#Prepare for distributing
from distutils.core import setup
import py2exe

#Version/Resource information
class PingURLTarget:
    def __init__(self, **kw):
        #Version information
        self.__dict__.update(kw)
        self.version="1, 0, 0, 0"
        self.product_version="1, 0, 0, 0"
        self.author="Dakusan"
        self.comments="http://www.castledragmire.com/Posts/Ping_URL"
        self.company_name = "Hyrulean Productions"
        self.description = "Periodically pings a requested URL"
        self.copyright = "(Dakusan) Copyright \xa9 2010"
        self.dest_base = "PingURL" #Output Filename
        self.name = "Ping URL"
        #Import file information
        self.script = "PingURL.py"
        self.icon_resources = [(10, "HyruleanProductions.ico")]

#Compilation information
setup(
    windows=[PingURLTarget()], #Resource information
    data_files=[("", [ #Put data files in top level directory
        "PingURL.cfg", "Readme.txt", #Readme, Configuration
        r"WinStuff\Microsoft.VC90.CRT.manifest", r"WinStuff\msvcr90.dll" #Required manifest and dlls for Python v2.6
    ])],
    options={ "py2exe": {
        "excludes": [ #Files to exclude
            "pyreadline", "locale", "doctest", "pdb", "unittest", "difflib", "inspect", "optparse", "pickle", "calendar"], #Unused base python modules
        "includes": ["encodings.utf_16", "encodings.utf_8"], #utf-16 and utf-8 are the only encoding that needs to be added
        "compressed": True, #Compress library.zip
        "ascii": True, #Only include needed encodings
        "dll_excludes": ["unicodedata", "w9xpopen.exe"], #Unused dlls
        "bundle_files": 2, #Bundle all files into library.zip
        "optimize": 2 #Optimize byte codes
    } }
)