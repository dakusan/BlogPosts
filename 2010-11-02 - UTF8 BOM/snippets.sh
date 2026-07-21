//Find PHP files with BOM (Cygwin compatible)
find . -name '*.php' -print0 | xargs -0 -n1000 grep -l $'^\xef\xbb\xbf'

//Remove the BOMs (Cygwin compatible)
find . -name '*.php' -print0 | xargs -0 -n1000 grep -l $'^\xef\xbb\xbf' | xargs -i perl -i.bak -pe 'BEGIN{ @d=@ARGV } s/^\xef\xbb\xbf//; END{ unlink map "$_$^I", @d }' "{}"

//Simpler remove BOMs (not Cygwin/Windows compatible)
find . -name '*.php' -print0 | xargs -0 -n1000 grep -l $'^\xef\xbb\xbf' | xargs -i perl -i -pe 's/^\xef\xbb\xbf//' "{}"