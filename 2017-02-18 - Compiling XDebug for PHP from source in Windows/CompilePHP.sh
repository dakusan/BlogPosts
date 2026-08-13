cd\php\src
buildconf --add-modules-dir=..\src-extensions
:: You may want some different flags for the below command
configure "--enable-snapshot-build" "--enable-debug-pack" "--without-pdo-mssql" "--enable-com-dotnet=shared" "--with-mcrypt=static" "--without-analyzer"
:: At this point we are going to actually build php to get the phpize.bat
:: You may receive some warnings at the end that say something to the extent of
:: “'-9' is not recognized as an internal or external command”. You can ignore these
nmake snap
:: This next command will “install” php to “c:\php”.
nmake install