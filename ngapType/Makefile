# Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
#
# SPDX-License-Identifier: Apache-2.0

.PHONY: install uninstall clean

ROOTDIR = ../..
ASN1CDIR = $(ROOTDIR)/asn1c

ASN1FILE = 38413-i10.asn1

COMPILELANG = -gen-Go
FLAGS = -pdu=all -fcompound-names -findirect-choice -fno-include-deps $(COMPILELANG)

all:
		$(ASN1CDIR)/asn1c/asn1c $(FLAGS) $(ASN1FILE)
install:
		cd $(ASN1CDIR) && test -f configure || autoreconf -iv && ./configure && make
clean:
		rm -f *.go

uninstall:
		cd $(ASN1CDIR) && $(MAKE) clean
