<?php
if(preg_match('/^(?:win|cygwin)/i', PHP_OS))
{
	$CWD=getcwd();
	require_once('Scripts/PHPExcel.phar');
	require_once('PHPExcel/IOFactory.php');
	chdir($CWD);
}
else
	require_once('phar://Scripts/PHPExcel.phar/PHPExcel/IOFactory.php');