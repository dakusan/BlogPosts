<?
/*Safe include/require functions that make sure UTF8 BOM is not output
Use like: eval(safe_INCLUDETYPE($INCLUDE_FILE_NAME));
where INCLUDETYPE is one of the following: include, require, include_once, require_once
An eval statement is used to maintain current scope
*/

//The different include type functions
function safe_include($FileName)	{ return real_safe_include($FileName, 'include'); }
function safe_require($FileName)	{ return real_safe_include($FileName, 'require'); }
function safe_include_once($FileName)	{ return real_safe_include($FileName, 'include_once'); }
function safe_require_once($FileName)	{ return real_safe_include($FileName, 'require_once'); }

//Start the processing and return the eval statement
function real_safe_include($FileName, $IncludeType)
{
	ob_start();
	return "$IncludeType('".strtr($FileName, Array("\\"=>"\\\\", "'", "\\'"))."'); safe_output_handler();";
}

//Do the actual processing and return the include data
function safe_output_handler()
{
	$Output=ob_get_clean();
	while(substr($Output, 0, 3)=='?') //Remove all instances of UTF8 BOM at the beginning of the output
		$Output=substr($Output, 3);
	print $Output;
}
?>