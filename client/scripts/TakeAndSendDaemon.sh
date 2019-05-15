#!/bin/sh
while :
do
    echo "Press [CTRL+C] to stop.."
    ./save-to-disk
    for f in $(ls data/)
    do
        echo
        echo Sending $f...
        ./send-file data/$f $(date '+%s')-$f asdf https://localhost/post 8080
    done
    rm -rf data/*
    sleep 300
done