TARGET ?= 6cord
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
GOFLAGS = -buildmode=pie -ldflags="-s -w -extldflags=$(LDFLAGS)" -gcflags=all=-trimpath=$(PWD) -asmflags=all=-trimpath=$(PWD)
CGO_ENABLED = 0

all: build
.PHONY: all build install clean

build:
	go build $(GOFLAGS) -o $(TARGET) .

install:
	install -Dm755 $(TARGET) -t $(DESTDIR)$(BINDIR)/

clean:
	$(RM) $(TARGET)

