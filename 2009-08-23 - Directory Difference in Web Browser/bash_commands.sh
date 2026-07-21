#Output all files and their timestamps to “Dir1.txt”
find -type f -printf '%P\t%T+\n' > Dir1.txt

#Output all files and their md5sums to “Dir1.txt”
find -type f -print0 | xargs -0 md5sum | perl -pe 's/^(.*?) (.*)$/$2\t$1/g' > Dir1.txt