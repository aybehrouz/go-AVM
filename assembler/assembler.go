// Copyright (c) 2021 aybehrouz <behrouz_ayati@yahoo.com>. This file is
// part of the go-avm repository: the Go implementation of the Argennon
// Virtual Machine (AVM).
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program. If not, see <https://www.gnu.org/licenses/>.

package assembler

import (
	"bufio"
	"fmt"
	"go-AVM/avm/binary"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var opcodes = make(map[string]byte, 256)

const OpcodesFile = "../opcodes.txt"

func init() {
	f, err := os.Open(OpcodesFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		if e := f.Close(); e != nil {
			log.Fatal(e)
		}
	}(f)

	var (
		opcode      byte
		instruction string
	)
	for {
		_, err := fmt.Fscanln(f, &opcode, &instruction)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		opcodes[instruction] = opcode
	}
}

func AssembleFile() {
}

func AssembleString(program string) []byte {
	return assemble(strings.NewReader(program))
}

func assemble(r io.Reader) (bytecode []byte) {
	wordScanner := bufio.NewScanner(r)
	wordScanner.Split(bufio.ScanWords)
	for wordScanner.Scan() {
		token := wordScanner.Text()
		v, err := strconv.ParseInt(token, 0, 64)
		if err == nil {
			b := make([]byte, 8)
			binary.PutInt64(b, 0, v)
			bytecode = append(bytecode, b...)
		} else if matched, _ := regexp.MatchString("[1-8][d][0-9|-]", token); matched {
			parts := strings.Split(token, "d")
			v, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			bytes, _ := strconv.Atoi(parts[0])

			b := make([]byte, 8)
			binary.PutInt64(b, 0, v)
			bytecode = append(bytecode, b[:bytes]...)
		} else {
			opcode, ok := opcodes[token]
			if !ok {
				log.Fatal("unknown instruction: " + token)
			}
			bytecode = append(bytecode, opcode)
		}
	}

	if err := wordScanner.Err(); err != nil {
		log.Fatal(err)
	}
	return bytecode
}
