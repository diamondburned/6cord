## $XTermId: Makefile.in,v 1.245 2019/01/14 01:40:56 tom Exp $
# -----------------------------------------------------------------------------
# this file is part of xterm
#
# Copyright 1997-2018,2019 by Thomas E. Dickey
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

SHELL		= /bin/sh

#### Start of system configuration section. ####

srcdir		= .


x		= 
o		= .o

CC		= gcc
CPP		= gcc -E
AWK		= gawk
LINK		= $(CC) $(CFLAGS)

CTAGS		= ctags
ETAGS		= 

LN_S		= ln -s
RM              = rm -f
LINT		= 
LINT_OPTS	= 

INSTALL		= /usr/bin/install -c
INSTALL_PROGRAM	= ${INSTALL}
INSTALL_SCRIPT	= ${INSTALL}
INSTALL_DATA	= ${INSTALL} -m 644
transform	= s,x,x,

EXTRA_CFLAGS	= 
EXTRA_CPPFLAGS	= 
EXTRA_LOADFLAGS	= 

CPPFLAGS	= -I. -I$(srcdir) -DHAVE_CONFIG_H -D_DEFAULT_SOURCE -D_XOPEN_SOURCE=700 -DNARROWPROTO=1 -DFUNCPROTO=15 -DOSMAJORVERSION=4 -DOSMINORVERSION=18 -I/usr/include/freetype2 -I/usr/include/libpng16 -I/usr/include/harfbuzz -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include -I/usr/include/uuid -I/usr/include/freetype2 -I/usr/include/libpng16 -I/usr/include/harfbuzz -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include -I/usr/include/freetype2 -I/usr/include/libpng16 -I/usr/include/harfbuzz -I/usr/include/glib-2.0 -I/usr/include/uuid -DDEFCLASS=\"XTerm\" $(EXTRA_CPPFLAGS)
CFLAGS		= -g -O2 $(EXTRA_CFLAGS)
LDFLAGS		=  
LIBS		=  -lXft -lfontconfig -lfreetype -lutil -lXaw7 -lXmu -lXt -lX11 -lXinerama -lXpm -lSM -lICE -lncurses

prefix		= /usr/local
exec_prefix	= ${prefix}
datarootdir	= ${prefix}/share
datadir		= ${datarootdir}

manext		= 1
bindir		= ${exec_prefix}/bin
libdir		= ${exec_prefix}/lib
mandir		= ${datarootdir}/man/man$(manext)
appsdir		= ${exec_prefix}/lib/X11/app-defaults
icondir		= no
pixmapdir	= ${datadir}/pixmaps

#### End of system configuration section. ####

ICON_NAME	= mini.xterm
ICON_SYMLINK	= NONE

DESTDIR		=
BINDIR		= $(DESTDIR)$(bindir)
LIBDIR		= $(DESTDIR)$(libdir)
MANDIR		= $(DESTDIR)$(mandir)
APPSDIR		= $(DESTDIR)$(appsdir)

#ICONDIR		= $(DESTDIR)$(icondir)
PIXMAPDIR	= $(DESTDIR)$(pixmapdir)

INSTALL_DIRS    = $(BINDIR) $(APPSDIR) $(ICONDIR) $(PIXMAPDIR) $(MANDIR)

CLASS		= XTerm
EXTRAHDR	=  256colres.h Tekparse.h charclass.h precompose.h wcwidth.h graphics_sixel.h graphics.h
EXTRASRC	=  TekPrsTbl.c Tekproc.c charclass.c precompose.c wcwidth.c html.c svg.c graphics_sixel.c graphics.c
EXTRAOBJ	=  TekPrsTbl.o Tekproc.o charclass.o precompose.o wcwidth.o html.o svg.o graphics_sixel.o graphics.o

AUTO_SOURCE	= \
		builtin_icons.h \
		VTparse.cin \
		Tekparse.cin \
		VTparse.hin \
		Tekparse.hin

          SRCS1 = button.c cachedGCs.c charproc.c charsets.c cursor.c \
	  	  data.c doublechr.c fontutils.c input.c \
		  linedata.c main.c menu.c misc.c \
		  print.c ptydata.c scrollback.c \
		  screen.c scrollbar.c tabs.c util.c version.c xstrings.c \
		  xtermcap.c VTPrsTbl.c $(EXTRASRC)
          OBJS1 = button$o cachedGCs$o charproc$o charsets$o cursor$o \
	  	  data$o doublechr$o fontutils$o input$o \
		  linedata$o main$o menu$o misc$o \
		  print$o ptydata$o scrollback$o \
		  screen$o scrollbar$o tabs$o util$o version$o xstrings$o \
		  xtermcap$o VTPrsTbl$o $(EXTRAOBJ)
          SRCS2 = resize.c version.c xstrings.c
          OBJS2 = resize$o version$o xstrings$o
           SRCS = $(SRCS1) $(SRCS2)
           OBJS = $(OBJS1) $(OBJS2)
           HDRS = VTparse.h data.h error.h fontutils.h main.h menu.h \
                  ptyx.h version.h xstrings.h xterm.h xtermcap.h $(EXTRAHDR)
       PROGRAMS = xterm$x resize$x

all :	$(PROGRAMS)
################################################################################
.SUFFIXES : .i .def .cin .hin .$(manext) .ms .man .txt  .html  .ps .pdf

.c$o :
	
	$(CC) $(CPPFLAGS) $(CFLAGS) -c $(srcdir)/$*.c

.c.i :
	
	$(CPP) -C $(CPPFLAGS) $*.c >$@

.def.cin :
	@echo "making $@ from $<"
	@$(AWK) 'BEGIN{printf "/* vile:cmode */\n";}/^CASE_/{printf "{ %d, \"%s\" },\n", n++, $$1; }' < $< >$@

.def.hin :
	@echo "making $@ from $<"
	@$(AWK) 'BEGIN{printf "/* vile:cmode */\n";}/^CASE_/{printf "#define %s %d\n", $$1, n++}' < $< >$@

.man.$(manext) :
	$(SHELL) ./minstall "$(INSTALL_DATA)" $< $@ $(appsdir) $(CLASS) $* $* $(pixmapdir)

#.$(manext).txt :
#	$(SHELL) -c "tbl $*.$(manext) | nroff -man | col -bx" >$@
#
#.ms.txt :
#	$(SHELL) -c "tbl $*.$(manext) | nroff -ms | col -bx" >$@
#

.$(manext).html :
	./man2html.tmp $* $(manext) man >$@

.$(manext).ps :
	$(SHELL) -c "tbl $*.$(manext) | groff -man" >$@

.$(manext).txt :
	GROFF_NO_SGR=stupid $(SHELL) -c "tbl $*.$(manext) | groff -Tascii -man | col -bx" >$@

.ms.html :
	./man2html.tmp $* ms ms >$@

.ms.ps :
	$(SHELL) -c "tbl $< | groff -ms" >$@

.ms.txt :
	GROFF_NO_SGR=stupid $(SHELL) -c "tbl $< | groff -Tascii -ms | col -bx" >$@

.ps.pdf :
	ps2pdf $*.ps
################################################################################

VTPARSE_H = VTparse.h VTparse.hin
TEKPARSE_H = Tekparse.h Tekparse.hin

main$o : main.h
misc$o : version.h

$(OBJS1) : xterm.h ptyx.h xtermcfg.h
main$o resize$o screen$o : xterm_io.h

xterm$x : $(OBJS1)
	$(SHELL) $(srcdir)/plink.sh $(LINK) $(LDFLAGS) -o $@ $(OBJS1) $(LIBS) $(EXTRA_LOADFLAGS)

resize$x : $(OBJS2)
	$(SHELL) $(srcdir)/plink.sh $(LINK) $(LDFLAGS) -o $@ $(OBJS2) $(LIBS)

256colres.h :
	-$(RM) $@
	perl $(srcdir)/256colres.pl > $@

88colres.h :
	-$(RM) $@
	perl $(srcdir)/88colres.pl > $@

charproc$o : $(VTPARSE_H) main.h  256colres.h
graphics_regis$o : $(VTPARSE_H) 
graphics_sixel$o : $(VTPARSE_H) 
misc$o : $(VTPARSE_H) 
VTPrsTbl$o : $(VTPARSE_H) 

TekPrsTbl$o : $(TEKPARSE_H)
Tekproc$o : $(TEKPARSE_H)

misc$o : builtin_icons.h

trace$o : VTparse.cin Tekparse.cin

# do this to quiet gcc -Wcast-qual warnings 
builtin_icons.h :
	@echo "#if OPT_BUILTIN_XPMS" >$@
	@echo "#include <icons/mini.xterm.xpms>" >>$@
	@echo "#include <icons/filled-xterm.xpms>" >>$@
	@echo "#include <icons/xterm.xpms>" >>$@
	@echo "#include <icons/xterm-color.xpms>" >>$@
	@echo "#else" >>$@
	@sed -e 's/static char \* /static const char * /' $(srcdir)/icons/mini.xterm_48x48.xpm >>$@
	@echo "#endif" >>$@
	@echo "made $@"

################################################################################
actual_xterm  = `echo xterm|        sed '$(transform)'`
actual_resize = `echo resize|       sed '$(transform)'`
actual_uxterm = `echo uxterm|       sed '$(transform)'`
actual_k8term = `echo koi8rxterm|   sed '$(transform)'`

binary_xterm  = $(actual_xterm)$x
binary_resize = $(actual_resize)$x
binary_uxterm = $(actual_uxterm)
binary_k8term = $(actual_k8term)

install \
install-bin \
install-full :: xterm$x resize$x $(BINDIR)
	$(SHELL) $(srcdir)/sinstall.sh  "$(INSTALL_PROGRAM)" xterm$x  /usr/bin/xterm $(BINDIR)/$(binary_xterm)
#	$(INSTALL_PROGRAM) xterm$x $(BINDIR)/$(binary_xterm)
	$(INSTALL_PROGRAM) -m  755 resize$x $(BINDIR)/$(binary_resize)

EDIT_SCRIPT = sed -e s,=xterm,=\$$name, -e s,XTerm,$(CLASS),

InstallLink = \
	if test NONE != NONE \
	&& test \$$source != NONE \
	&& test \$$source != \$$target ; then \
		cd \$$TARGET && ( \
			$(RM) \$$source ; \
			$(LN_S) \$$target \$$source ; \
			echo \"... created symbolic link:\" ; \
			ls -l \$$target \$$source ) ; \
	fi

InstallBinLink = TARGET=$(BINDIR); $(InstallLink)
InstallManLink = TARGET=$(MANDIR); $(InstallLink)

InstallScript = \
	echo \"... installing $(BINDIR)/\$$target\"; \
	name=$(binary_xterm); \
	$(EDIT_SCRIPT) $(srcdir)/\$$source >\$$source.tmp; \
	$(INSTALL_SCRIPT) -m  755 \$$source.tmp $(BINDIR)/\$$target; \
	$(RM) \$$source.tmp

install \
install-bin \
install-scripts \
install-full :: $(BINDIR)
	@$(SHELL) -c "source=\"NONE\"; \
		target=\"$(binary_xterm)\"; \
		$(InstallBinLink)"
	@$(SHELL) -c "source=\"resize\"; \
		target=\"$(binary_resize)\"; \
		$(InstallBinLink)"
	@$(SHELL) -c "source=\"uxterm\"; \
		target=\"$(binary_uxterm)\"; \
		$(InstallScript); \
		$(InstallBinLink)"
	@$(SHELL) -c "source=\"koi8rxterm\"; \
		target=\"$(binary_k8term)\"; \
		$(InstallScript); \
		$(InstallBinLink)"

install \
install-man \
install-full :: $(MANDIR)
	@-$(SHELL) -c "for source in xterm resize uxterm koi8rxterm ; \
		do \
			target=\`echo \"\$$source\" | sed 's,x,x,'\`; \
			$(SHELL) ./minstall \"$(INSTALL_DATA)\" \
				$(srcdir)/\$$source.man  \
				$(MANDIR)/\$$target.$(manext) \
				$(appsdir) \
				$(CLASS) \
				\$$source \
				\$$target \
				$(pixmapdir); \
		done"
	@-$(SHELL) -c "if test NONE != NONE ; then \
		source=$(actual_xterm).$(manext); \
		target=NONE.$(manext); \
		cd $(MANDIR) && ( \
			$(RM) \$$target ; \
			$(LN_S) \$$source \$$target ; \
			echo '... created symbolic link:' ; \
			ls -l \$$source \$$target ; \
			) \
		fi"

APP_NAMES = XTerm UXTerm KOI8RXTerm

install \
install-app \
install-full :: $(APPSDIR)
	@-$(SHELL) -c 'for s in $(APP_NAMES); \
	do \
		echo "** $$s"; \
		d=`echo $$s | sed -e s/XTerm/$(CLASS)/`; \
		echo installing $(APPSDIR)/$$d; \
		sed -e s/XTerm/$(CLASS)/ $(srcdir)/$$s.ad >XTerm.tmp; \
		$(INSTALL_DATA) XTerm.tmp $(APPSDIR)/$$d; \
		echo installing $(APPSDIR)/$$d-color; \
		sed -e s/XTerm/$$d/ $(srcdir)/XTerm-col.ad >XTerm.tmp; \
		$(INSTALL_DATA) XTerm.tmp $(APPSDIR)/$$d-color; \
	done'
	@$(RM) XTerm.tmp
#	@echo "... installed app-defaults"

#ICON_LIST = icons/mini.xterm_48x48.png icons/mini.xterm_48x48.xpm
#ICON_THEME = no
#install \
#install-icon \
#install-full :: $(ICONDIR)
#	ACTUAL_XTERM=$(actual_xterm) $(SHELL) -c '\
#		h=$(ICONDIR)/$(ICON_THEME); \
#		for n in $(ICON_LIST); \
#		do \
#			x=$$ACTUAL_XTERM; \
#			l=`echo "$$n" | cut -f1 -d:`; \
#			r=`echo "$$n" | cut -f2 -d: |sed -e "s,xterm,$$ACTUAL_XTERM,"`; \
#			test -z "$$r" && continue; \
#			d=$$h/`echo "$$r" | sed -e "s,/[^/]*$$,,"`; \
#			test -d "$$d" || mkdir -p "$$d"; \
#			echo "installing icon $$h/$$r"; \
#			$(INSTALL_DATA) $$l $$h/$$r; \
#			s=`echo "$$r" | sed -e '"'s,^.*\.,.,'"'`; \
#			t=$(ICON_SYMLINK)$$s; \
#			b=`basename $$n $$s | sed -e "s,_[1-9][0-9]*x.*,,"`; \
#			if test "$(ICON_SYMLINK)" != NONE ; then \
#				if test "$$r" != "$$t" ; then \
#					if test "x$$b" = "x$(ICON_NAME)" ; then \
#						echo "linking $$r -> $$t"; \
#						( cd $$h; $(RM) $$t; $(LN_S) $$r $$t; ) \
#					fi \
#				fi \
#			fi \
#		done'
#	@echo "... installed icons"

install \
install-icon \
install-full :: $(PIXMAPDIR)
	ACTUAL_XTERM=$(actual_xterm) $(SHELL) -c '\
		h=$(PIXMAPDIR); \
		for n in $(srcdir)/icons/*xterm*_32x32.xpm $(srcdir)/icons/*xterm*_48x48.xpm; \
		do \
			l=`basename $$n`; \
			r=`echo "$$l" | sed -e "s,xterm,$$ACTUAL_XTERM,"`; \
			echo "installing pixmap $$h/$$r"; \
			$(INSTALL_DATA) $(srcdir)/icons/$$l $$h/$$r; \
			s=`echo "$$r" | sed -e '"'s,^.*\.,.,'"'`; \
			t=$(ICON_SYMLINK)$$s; \
			b=`basename $$n $$s | sed -e "s,_[1-9][0-9]*x.*,,"`; \
			if test "$(ICON_SYMLINK)" != NONE ; then \
				if test "x$$r" != "$$t" ; then \
					if test "x$$b" = "x$(ICON_NAME)" ; then \
						echo "linking $$r -> $$t"; \
						( cd $$h; $(RM) $$t; $(LN_S) $$r $$t; ) \
					fi \
				fi \
			fi \
		done'
	@echo "... installed icons"

install ::
	@echo 'Completed installation of executables and documentation.'
	@echo 'Use "make install-ti" to install terminfo description.'

TERMINFO_DIR = $(DESTDIR)/usr/lib/terminfo
SET_TERMINFO = TERMINFO=$(TERMINFO_DIR)

install-full \
install-ti :: $(TERMINFO_DIR)
	@$(SHELL) -c "$(SET_TERMINFO) $(srcdir)/run-tic.sh $(srcdir)/terminfo"
	@echo 'Completed installation of terminfo description.'

install-full \
install-tc ::
	@-$(SHELL) -c "if test -f /etc/termcap ; then echo 'You must install the termcap entry manually by editing /etc/termcap'; fi"

installdirs : $(INSTALL_DIRS)
################################################################################

UninstallLink = \
	if test NONE != NONE \
	&& test \$$source != NONE \
	&& test \$$source != \$$target \
	&& test -h \$$TARGET/\$$source ; then \
		echo \"... removing \$$TARGET/\$$source\"; \
		cd \$$TARGET && \
		$(RM) \$$source; \
	fi

UninstallBinLink = TARGET=$(BINDIR); $(UninstallLink)
UninstallManLink = TARGET=$(MANDIR); $(UninstallLink)

UninstallBinary = \
	echo \"... removing $(BINDIR)/\$$target\"; \
	$(RM) $(BINDIR)/\$$target

uninstall \
uninstall-bin \
uninstall-full ::
	@-$(SHELL) -c "source=\"NONE\"; \
		target=\"$(binary_xterm)\"; \
		$(UninstallBinLink); \
		$(UninstallBinary)"
	@-$(SHELL) -c "source=\"resize\"; \
		target=\"$(binary_resize)\"; \
		$(UninstallBinLink); \
		$(UninstallBinary)"

uninstall \
uninstall-bin \
uninstall-scripts \
uninstall-full ::
	@-$(SHELL) -c "source=\"uxterm\"; \
		target=\"$(binary_uxterm)\"; \
		$(UninstallBinLink); \
		$(UninstallBinary)"
	@-$(SHELL) -c "source=\"koi8rxterm\"; \
		target=\"$(binary_k8term)\"; \
		$(UninstallBinLink); \
		$(UninstallBinary)"

uninstall \
uninstall-man \
uninstall-full ::
	@-$(SHELL) -c "\
		source=NONE.$(manext); \
		target=$(actual_xterm).$(manext); \
		$(UninstallManLink)"
	@-$(SHELL) -c "for source in \
		$(actual_xterm).$(manext) \
		$(actual_resize).$(manext) \
		$(actual_uxterm).$(manext) \
		$(actual_k8term).$(manext); \
		do \
			echo \"... removing $(MANDIR)/\$$source\"; \
			$(RM) $(MANDIR)/\$$source; \
		done"

uninstall \
uninstall-app \
uninstall-full ::
	-$(SHELL) -c 'for s in $(APP_NAMES); \
	do \
		echo "** $$s"; \
		d=`echo $$s | sed -e s/XTerm/$(CLASS)/`; \
		echo uninstalling $(APPSDIR)/$$d; \
		$(RM) $(APPSDIR)/$$d; \
		echo uninstalling $(APPSDIR)/$$d-color; \
		$(RM) $(APPSDIR)/$$d-color; \
	done'

#uninstall \
#uninstall-icon \
#uninstall-full ::
#	-$(SHELL) -c 'ACTUAL_XTERM=$(actual_xterm) ; \
#		for n in $(ICON_LIST); \
#		do \
#		r=`echo "$$n" | sed -e s,\^.\*:,, -e s,xterm,$$ACTUAL_XTERM,`; \
#		test -z "$$r" && continue; \
#		h=$(ICONDIR)/$(ICON_THEME); \
#		test -f $$h/$$r || continue; \
#		echo removing $$h/$$r; \
#		$(RM) $$h/$$r; \
#		done'
#	@echo "... removed icons"

uninstall \
uninstall-icon \
uninstall-full ::
	-@$(SHELL) -c 'ACTUAL_XTERM=$(actual_xterm) ; \
		for n in $(srcdir)/icons/*xterm*_32x32.xpm $(srcdir)/icons/*xterm*_48x48.xpm; \
		do \
			l=`basename $$n`; \
			r=`echo "$$l" | sed -e "s,xterm,$$ACTUAL_XTERM,"`; \
			echo removing $(PIXMAPDIR)/$$r; \
			$(RM) $(PIXMAPDIR)/$$r; \
		done'
	@echo "... removed icons"
################################################################################
# Desktop-utils does not provide an uninstall, and is not uniformly available.
DESKTOP_FILES = $(srcdir)/xterm.desktop $(srcdir)/uxterm.desktop
DESKTOP_FLAGS = 
install-desktop \
install-full ::
	ACTUAL_XTERM=$(actual_xterm) \
	$(SHELL) -c 'for n in $(DESKTOP_FILES); \
		do $(SHELL) df-install $$ACTUAL_XTERM $(ICON_NAME) $(DESKTOP_FLAGS) $$n; \
		done'
################################################################################
check :
	@ echo "There are no batch tests available (use vttest)"
################################################################################
mostlyclean :
	-$(RM) *$o *.[is] XTerm[1-9]*.* Xterm.log.* XtermLog.* .pure core *~ *.bak *.BAK *.out *.tmp

clean : mostlyclean
	-$(RM) $(PROGRAMS) $(AUTO_SOURCE)

sources : $(AUTO_SOURCE)

distclean :: clean
	-$(RM) Makefile config.status config.cache config.log xtermcfg.h
	-$(RM) df-install minstall

distclean \
docs-clean ::
	-$(RM) *.ps *.pdf *.png
	-$(SHELL) -c 'for p in xterm resize uxterm koi8rxterm; \
	do \
		$(RM) $$p.html $$p.$(manext) $$p.txt; \
	done'
	-$(RM) ctlseqs.html ctlseqs.$(manext)

distclean ::
	-$(RM) man2html.tmp

realclean : distclean
	-$(RM) tags TAGS

maintainer-clean : realclean
	-$(RM) 256colres.h 88colres.h
################################################################################
terminfo.out : terminfo		; tic -a -I -1 terminfo >$@
termcap.out : termcap		; tic -a -C -U termcap >$@
################################################################################
docs-ctlseqs \
docs :: $(srcdir)/ctlseqs.txt  ctlseqs.html  ctlseqs.pdf ctlseqs.ps

ctlseqs.html : $(srcdir)/ctlseqs.ms
ctlseqs.pdf : ctlseqs.ps
ctlseqs.ps : $(srcdir)/ctlseqs.ms
ctlseqs.txt : $(srcdir)/ctlseqs.ms
################################################################################
docs-resize \
docs ::  resize.txt  resize.html  resize.pdf resize.ps
resize.html : resize.$(manext)
resize.pdf : resize.ps
resize.ps : resize.$(manext)
resize.txt : resize.$(manext)
################################################################################
docs-xterm \
docs ::  xterm.txt  xterm.html  xterm.pdf xterm.ps
xterm.html : xterm.$(manext)
xterm.pdf : xterm.ps
xterm.ps : xterm.$(manext)
xterm.txt : xterm.$(manext)
################################################################################
docs-uxterm \
docs ::  uxterm.txt  uxterm.html  uxterm.pdf uxterm.ps
uxterm.html : uxterm.$(manext)
uxterm.pdf : uxterm.ps
uxterm.ps : uxterm.$(manext)
uxterm.txt : uxterm.$(manext)
################################################################################
docs-koi8rxterm \
docs ::  koi8rxterm.txt  koi8rxterm.html  koi8rxterm.pdf koi8rxterm.ps
koi8rxterm.html : koi8rxterm.$(manext)
koi8rxterm.pdf : koi8rxterm.ps
koi8rxterm.ps : koi8rxterm.$(manext)
koi8rxterm.txt : koi8rxterm.$(manext)
################################################################################
lint :
	$(LINT) $(LINT_OPTS) $(CPPFLAGS) $(SRCS1)
	$(LINT) $(LINT_OPTS) $(CPPFLAGS) $(SRCS2)

tags :
	$(CTAGS) $(SRCS) $(HDRS)

TAGS :
	$(ETAGS) $(SRCS) $(HDRS)

$(TERMINFO_DIR) $(INSTALL_DIRS) :
	mkdir -p $@

ALWAYS :

depend : $(TABLES)
	makedepend -- $(CPPFLAGS) -- $(SRCS)

# DO NOT DELETE THIS LINE -- make depend depends on it.
