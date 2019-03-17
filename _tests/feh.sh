#!/bin/bash
[ -z $1 ] && exit 1

case "$1" in
*.jpg*|*.png*|*.PNG*|*.jpeg*|*.gif*)
	H="800"
	W="600"

	[[ "$1" = *"/emojis/"* ]] && {
		H="50"
		W="50"
	}

	feh -H $H -W $W -b trans --auto-zoom --xinerama-index 0 -B black -. -x "$1"
	;;
*)
	xdg-open "$1"
	;;
esac

