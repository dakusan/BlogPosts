#!/bin/bash
cd /src #A good directory to install stuff in. You may need to create it first

#Install intltool (required for pslib)
yum install intltool

#Install pslib
wget http://sourceforge.net/projects/pslib/files/latest/download?source=files
tar -xzvf pslib-*.tar.gz
cd pslib-*
./configure
make
make install
cd ..

#Install pslib wrapper for php
pecl download ps
tar -xzvf ps-*.tgz
cd ps-*
phpize
./configure
make

echo 'In ps.c change: `function_entry ps_functions[] = {` to `zend_function_entry ps_functions[] = {`'
read

make
make install