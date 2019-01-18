#!/bin/sh
# $XTermId: run-tic.sh,v 1.7 2019/01/14 01:52:19 tom Exp $
# -----------------------------------------------------------------------------
# this file is part of xterm
#
# Copyright 2006-2007,2019 by Thomas E. Dickey
# 
#                         All Rights Reserved
# 
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the
# "Software"), to deal in the Software without restriction, including
# without limitation the rights to use, copy, modify, merge, publish,
# distribute, sublicense, and/or sell copies of the Software, and to
# permit persons to whom the Software is furnished to do so, subject to
# the following conditions:
# 
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
# OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
# IN NO EVENT SHALL THE ABOVE LISTED COPYRIGHT HOLDER(S) BE LIABLE FOR ANY
# CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
# TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
# SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
# 
# Except as contained in this notice, the name(s) of the above copyright
# holders shall not be used in advertising or otherwise to promote the
# sale, use or other dealings in this Software without prior written
# authorization.
# -----------------------------------------------------------------------------
#
# Run tic, either using ncurses' extension feature or filtering out harmless
# messages for the extensions which are otherwise ignored by other versions of
# tic.

TMP=run-tic$$.log
VER=`tic -V 2>/dev/null`
OPT=

PROG=tic
unset TERM
unset TERMINFO_DIRS

case "x$VER" in
*ncurses*)
	# Prefer ncurses 6 over 5, if we can get it, since some older 5.9
	# packages do not handle the extensions as well.
	case "$VER" in
	*\ [6789].*)
		;;
	*)
		VER=`tic6 -V 2>/dev/null`
		test -n "$VER" && PROG=tic6
		;;
	esac
	echo "** using tic from $VER"
	OPT="-x"
	;;
esac

echo "** $PROG $OPT" "$@"
$PROG $OPT "$@" 2>$TMP
RET=$?

fgrep -v 'Unknown Capability' $TMP | \
fgrep -v 'Capability is not recognized:' | \
fgrep -v 'tic: Warning near line ' >&2
rm -f $TMP

exit $RET
