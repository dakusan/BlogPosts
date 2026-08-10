<?php
//Please note, for the recursion test to go beyond 256, xdebug.max_nesting_level needs to be raised.
//You also may need to update your memory_limit depending on the number of iterations

//Output the start memory
print 'Start: '.memory_get_usage()."B<br><br>Below test results are in MB<br>";

//Our 1MB string
global $OneMB, $NumIterations;
$OneMB=str_repeat('x', 1024*1024);
$NumIterations=500;

//Run the tests
$ConcatTest=RunTest('ConcatTest');
$ImplodeTest=RunTest('ImplodeTest');
$RecurseTest=RunTest('RecurseTest');

//Output the results in a table
OutputResults(
  Array('ConcatTest', 'ImplodeTest', 'RecurseTest'),
  Array($ConcatTest, $ImplodeTest, $RecurseTest)
);

//Start a test run by initializing the array that will hold the results and manipulating those results after the test is complete
function RunTest($TestName)
{
  $CurrentTestNums=Array();
  $TestStartMem=memory_get_usage();
  $StartTime=microtime(true);
  RunTestReal($TestName, $CurrentTestNums, $StrLen);
  $CurrentTestNums[]=memory_get_usage();

  //Subtract $TestStartMem from all other numbers
  foreach($CurrentTestNums as &$Num)
    $Num-=$TestStartMem;
  unset($Num);

  $CurrentTestNums[]=$StrLen;
  $CurrentTestNums[]=microtime(true)-$StartTime;

  return $CurrentTestNums;
}

//Initialize the test and store the memory allocated at the end of the test, with the result
function RunTestReal($TestName, &$CurrentTestNums, &$StrLen)
{
  $R=$TestName($CurrentTestNums);
  $CurrentTestNums[]=memory_get_usage();
  $StrLen=strlen($R);
}

//Concatenate 1MB string over and over onto a single string
function ConcatTest(&$CurrentTestNums)
{
  global $OneMB, $NumIterations;
  $Result='';
  for($i=0;$i<$NumIterations;$i++)
  {
    $Result.=$OneMB;
    $CurrentTestNums[]=memory_get_usage();
  }
  return $Result;
}

//Create an array of 1MB strings and then join w/ an implode
function ImplodeTest(&$CurrentTestNums)
{
  global $OneMB, $NumIterations;
  $Result=Array();
  for($i=0;$i<$NumIterations;$i++)
  {
    $Result[]=$OneMB;
    $CurrentTestNums[]=memory_get_usage();
  }
  return implode('', $Result);
}

//Recursively add strings onto each other
function RecurseTest(&$CurrentTestNums, $TestNum=0)
{
  Global $OneMB, $NumIterations;
  if($TestNum==$NumIterations)
    return '';

  $NewStr=RecurseTest($CurrentTestNums, $TestNum+1).$OneMB;
  $CurrentTestNums[]=memory_get_usage();
  return $NewStr;
}

//Output the results in a table
function OutputResults($TestNames, $TestResults)
{
  global $NumIterations;
  print '<table border=1 cellspacing=0 cellpadding=2><tr><th>Test Name</th><th>'.implode('</th><th>', $TestNames).'</th></tr>';
  $FinalNames=Array('Final Result', 'Clean');
  for($i=0;$i<$NumIterations+2;$i++)
  {
    $TestName=($i<$NumIterations ? $i : $FinalNames[$i-$NumIterations]);
    print "<tr><th>$TestName</th>";
    foreach($TestResults as $TR)
      printf('<td>%07.4f</td>', $TR[$i]/1024/1024);
    print '</tr>';
  }

  //Other result numbers
  print '<tr><th>Final String Size</th>';
  foreach($TestResults as $TR)
    printf('<td>%d</td>', $TR[$NumIterations+2]);
  print '</tr><tr><th>Runtime</th>';
    foreach($TestResults as $TR)
      printf('<td>%s</td>', $TR[$NumIterations+3]);
  print '</tr></table>';
}
?>