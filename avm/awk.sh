#!/bin/sh

#
# Copyright (c) 2021 aybehrouz <behrouz_ayati@yahoo.com>. This file is
# part of the go-avm repository: the Go implementation of the Argennon
# Virtual Machine (AVM).
#
# This program is free software: you can redistribute it and/or modify it
# under the terms of the GNU General Public License as published by the
# Free Software Foundation, either version 3 of the License, or (at your
# option) any later version.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
# Public License for more details.
#
# You should have received a copy of the GNU General Public License along
# with this program. If not, see <https://www.gnu.org/licenses/>.
#

awk 'BEGIN {
     	FS = "[.:, \t]+"
     }

     /0x/ {
     	print $2 "\t" $5
     }' controller.go > ../opcodes.txt
