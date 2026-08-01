<?php
//Get old domains
$BlockedFile='/var/named/blacklisted.conf';
$CurrentZones=Array();
foreach(explode("\n", file_get_contents($BlockedFile)) as $Line)
        if(preg_match('/^zone "([\w\._-]+)"/', $Line, $Results))
                $CurrentZones[]=$Results[1];

//List domains
if(isset($_REQUEST['List']))
        return print implode('<br>', $CurrentZones);

//Get new domains
if(!isset($_REQUEST['Domains']))
        return print 'Missing Domains';
$Domains=$_REQUEST['Domains'];
if(!preg_match('/^[\w\._-]+(,[\w\._-]+)*$/uD', $Domains))
        return print 'Invalid domains string';
$Domains=explode(',', $Domains);

//Remove domains
if(isset($_REQUEST['Remove']))
{
        $CurrentZones=array_flip($CurrentZones);
        foreach($Domains as $Domain)
                unset($CurrentZones[$Domain]);
        $FinalDomainList=array_keys($CurrentZones);
}
else //Combine domains
        $FinalDomainList=array_unique(array_merge($Domains, $CurrentZones));

//Output to the file
$FinalDomainData=Array();
foreach($FinalDomainList as $Domain)
        $FinalDomainData[]=
                "zone \"$Domain\" { type master; file \"blacklisted.db\"; };";
file_put_contents($BlockedFile, implode("\n", $FinalDomainData));

//Reload named
print `sudo /var/www/html/AddRules/restart_named`;
?>