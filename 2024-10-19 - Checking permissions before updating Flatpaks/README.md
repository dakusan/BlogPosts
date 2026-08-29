<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Checking_permissions_before_updating_Flatpaks">Checking permissions before updating Flatpaks</a><br>
Date: 2024-10-19<br>
Tags: AI, ChatGPT, Command Line Utility, Examples, Flatpak, Linux Migration, Permissions, PHP, Security, Technology<br><br>

<p>I wanted to make a script to check if updated flatpak app permissions have changed before updating them. I decided to use <b>ChatGPTs</b> newest model, <b>01-preview</b> for this task and holy crap am I impressed. I had to give it 14 commands to get to the final product and it seems to work great. That number could have been reduced quite a bit had I done things a bit differently.</p>
<p>I still had to find the problems, look up reference docs to correct it, and even debug in an IDE a little. But just telling it where the problems were, it got there in the end, and its user interaction output is way better than what I was planning on doing.</p>
</div>

> [!NOTE]
> To see the rest of this post that talked about working with the AI, see the <a href="https://www.castledragmire.com/Posts/Checking_permissions_before_updating_Flatpaks">original post here</a>.

# My final updated copy of the script
```php
<?php
// -- This script checks for Flatpak updates and reports permission changes. --

// Function to parse permissions from Flatpak metadata
function parsePermissions($content) {
    $permissions = [];
    $lines = is_array($content) ? $content : explode("\n", $content);
    $currentSection = '';
    $skipSections = ['Application']; // Sections to skip

    foreach ($lines as $line) {
        $line = trim($line);

        if (empty($line)) {
            continue;
        }

        // Check for section headers
        if (preg_match('/^\[(.*)\]$/', $line, $matches)) {
            $currentSection = $matches[1];

            // Skip the [Application] section
            if (in_array($currentSection, $skipSections)) {
                $currentSection = '';
                continue;
            }

            $permissions[$currentSection] = [];
        } elseif ($currentSection !== '') {
            // Only process lines within non-skipped sections
            $parts = explode('=', $line, 2);
            if (count($parts) == 2) {
                $key = $parts[0];
                $values = explode(';', trim($parts[1], ';'));
                $permissions[$currentSection][$key] = $values;
            } else {
                // Handle keys without '=' (e.g., single permissions)
                $permissions[$currentSection][$line] = [];
            }
        }
    }

    return $permissions;
}

// Function to compare permissions
function comparePermissions($current, $new) {
    $differences = [];

    // Get all sections
    $sections = array_unique(array_merge(array_keys($current), array_keys($new)));

    foreach ($sections as $section) {
        $currentSection = isset($current[$section]) ? $current[$section] : [];
        $newSection = isset($new[$section]) ? $new[$section] : [];

        // Get all keys in this section
        $keys = array_unique(array_merge(array_keys($currentSection), array_keys($newSection)));

        foreach ($keys as $key) {
            $currentValues = isset($currentSection[$key]) ? $currentSection[$key] : [];
            $newValues = isset($newSection[$key]) ? $newSection[$key] : [];

            // Compare values
            $added = array_diff($newValues, $currentValues);
            $removed = array_diff($currentValues, $newValues);

            if (!empty($added) || !empty($removed)) {
                $differences[$section][$key] = [
                    'added' => $added,
                    'removed' => $removed,
                ];
            }
        }
    }

    return $differences;
}

// Function to output a line with appId (if given), colored red (IsError=true) or green (IsError=false)
function outputColor($appId, $str, $isError=true, $NoNewline=false)
{
    // Determine if coloring should be used
    static $hasColors=null;
    if(!isset($hasColors)) {
        $hasColors=stream_isatty(STDOUT) && (!trim(`command -v tput`) || intval(`tput colors`)>=16);
    }

    echo
        (!$hasColors ? '' : ($isError ? "\e[31m" : "\e[32m")).
        ($appId ? "$appId: " : '').
        $str.
        (!$hasColors ? '' : "\e[0m").
        ($NoNewline ? '' : "\n");
}

// Function to display permission differences
function displayDifferences($appId, $differences) {
    outputColor($appId, 'Permission changes:');
    foreach ($differences as $section => $keys) {
        outputColor(null, "  [$section]");
        foreach ($keys as $key => $changes) {
            if (!empty($changes['added'])) {
                outputColor(null, "    + $key = " . implode(';', $changes['added']));
            }
            if (!empty($changes['removed'])) {
                outputColor(null, "    - $key = " . implode(';', $changes['removed']));
            }
        }
    }
    echo "\n";
}

// Function to ask the user a question
function askQuestion($Question) {
    outputColor('', "$Question (y/N): ", true, true);
    $handle = fopen('php://stdin', 'r');
    $line = fgets($handle);
    $answer = trim(strtolower($line));
    fclose($handle);
    return ($answer == 'y' || $answer == 'yes');
}

// Get architecture (moved above the loop)
$archCommand = 'flatpak --default-arch';
exec($archCommand, $archOutput);
$arch = trim(implode('', $archOutput));

// Step 1: Get the list of installed Flatpaks
$installedCommand = 'flatpak list -a --columns=application,branch,origin,options';
exec($installedCommand, $installedOutput);

// Remove header line if present
if (!empty($installedOutput) && strpos($installedOutput[0], 'Application ID') !== false) {
    array_shift($installedOutput);
}

// Build an associative array of installed applications
$installedApps = [];

foreach ($installedOutput as $line) {
    // The output is tab-delimited
    $columns = explode("\t", trim($line));
    if (count($columns) >= 4) {
        $appId = $columns[0];
        $branch = $columns[1];
        $origin = $columns[2];
        $options = $columns[3];

        $installedApps[$appId] = [
            'appId' => $appId,
            'branch' => $branch,
            'origin' => $origin,
            'options' => $options,
        ];
    }
}

// Get the list of updates
exec('echo n | flatpak update', $updatesOutput);
foreach ($updatesOutput as $Line) {
    if (preg_match('/^ *\d+\.\s+([\w\.-]+)/', $Line, $matches)) {
        $updatesAvailable[]=$matches[1];
    }
}

// Let the user confirm the updates
echo str_repeat('-', 80)."\n";
echo implode("\n", $updatesOutput)."\n";
echo str_repeat('-', 80)."\n";
if (empty($updatesAvailable)) {
    outputColor('', 'No updates available for installed Flatpaks.');
    exit(0);
}
if (!askQuestion('Found '.count($updatesAvailable).' updates. Continue?')) {
    echo "Updates canceled.\n";
    exit(0);
}

$permissionChanges = [];
foreach ($updatesAvailable as $appId) {
    if (!isset($installedApps[$appId])) {
        outputColor($appId, 'Installed app not found. Skipping.');
        continue;
    }

    $app = $installedApps[$appId];
    $branch = $app['branch'];
    $origin = $app['origin'];
    $options = $app['options'];

    // Determine if it's an app or runtime
    $isRuntime = strpos($options, 'runtime') !== false;

    // Paths to the metadata files
    if ($isRuntime) {
        $metadataPath = "/var/lib/flatpak/runtime/$appId/$arch/$branch/active/metadata";
    } else {
        $metadataPath = "/var/lib/flatpak/app/$appId/$arch/$branch/active/metadata";
    }

    // Check if the metadata file exists
    if (!file_exists($metadataPath)) {
        outputColor($appId, 'Metadata file not found. Skipping.');
        continue;
    }

    // Read current permissions from the metadata file
    $metadataContent = file_get_contents($metadataPath);
    $currentPermissions = parsePermissions($metadataContent);

    // Get new metadata from remote
    $ref = $appId . '/' . $arch . '/' . $branch;
    $remoteInfoCommand = 'flatpak remote-info --show-metadata ' . escapeshellarg($origin) . ' ' . escapeshellarg($ref);

    // Clear $remoteOutput before exec()
    $remoteOutput = [];
    exec($remoteInfoCommand, $remoteOutput);

    if (empty($remoteOutput)) {
        outputColor($appId, 'Failed to retrieve remote metadata. Skipping.');
        continue;
    }

    // Parse new permissions from the remote metadata
    $newPermissions = parsePermissions($remoteOutput);

    // Compare permissions
    $differences = comparePermissions($currentPermissions, $newPermissions);

    if (!empty($differences)) {
        $permissionChanges[$appId] = $differences;
        displayDifferences($appId, $differences);
    } else {
        outputColor($appId, 'No permission changes found.', false);
    }
}

// If there are no permission changes, inform the user
if (empty($permissionChanges)) {
    echo "No permission changes detected in the available updates.\n";
}

// Ask user if they want to proceed
if (askQuestion('Do you want to proceed with the updates?')) {
    // Proceed with updates
    echo "Updating Flatpaks...\n";
    passthru('flatpak update -y');
} else {
    echo "Updates canceled.\n";
}
?>
```