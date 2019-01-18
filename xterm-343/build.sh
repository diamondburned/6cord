#!/bin/sh
./configure --enable-256-color --enable-sixel-graphics --enable-sixel
make
strip ./xterm
echo Done.
