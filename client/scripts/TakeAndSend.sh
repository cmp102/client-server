#!/bin/sh
./save-to-disk
for f in $(ls data/)
do
    echo
    echo Sending $f...
    ./send-file data/$f $(date '+%s')-$f $1 $2 $3
done
rm -rf data/*