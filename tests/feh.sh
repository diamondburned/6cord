#!/bin/bash
[ -z $1 ] && exit 1
case "$1" in
*.jpg*|*.png*|*.jpeg*|*.gif*)
	feh -H 800 -W 600 -b trans --auto-zoom --xinerama-index 0 -B black -. -x "$1"
	;;
*)
	xdg-open "$1"
	;;
esac

