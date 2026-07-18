<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Installing_PSLib_for_PHP">Installing PSLib for PHP</a><br>
Date: 2013-08-16<br>
Tags: Examples, Installation, Libraries, PHP, PSLib<br><br>

I was setting up a brand new server yesterday running Cent OS 6.4 and had the need to install PSLib (PostScript) for PHP.<br>
The initial setup commands all worked as expected:
</div>

```bash
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
```

At this point, the make failed with
```
/src/ps-1.3.6/ps.c:58: error: expected ‘=’, ‘,’, ‘;’, ‘asm’ or ‘__attribute__’ before ‘ps_functions’
```

After a little bit of browsing the code and a really lucky first guess, I found that changing the following fixed the problem:<br>
In <b>ps.c</b> change: `function_entry ps_functions[] = {` to `zend_function_entry ps_functions[] = {`

Then to finish the install, just run:
```bash
make
make install
```

and you’re done!