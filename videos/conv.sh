#!/bin/sh

if [$2 -ne 3 ]; then
    echo "Example usage: optvideo atom.mp4 1920 output_folder"
    exit 1
fi

FILENAME=${1%%.*}

ffmpeg -i $1 -c:v libvpx-vp9 -crf 40 -vf scale=$2:-2 -an $3/${FILENAME}.webm
ffmpeg -i $1 -c:v libx264 -crf 24 -vf scale=$2:-2 -movflags faststart -an $3/${FILENAME}_h264.mp4
ffmpeg -i $1 -c:v libx264 -crf 24 -vf scale=$2:-2 -tag:v hvc1 -movflags faststart -an $3/${FILENAME}_h265.mp4
