cd ..\src-extensions\xdebug
..\..\SDK\phpize
configure --with-xdebug
nmake php_xdebug.dll
copy Release_TS\php_xdebug.dll ..\..\ext