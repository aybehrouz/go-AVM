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

/*
Package avm implements the core functionalities of the Argennon Virtual
Machine.

Performance Considerations:

We avoid using interfaces for frequently used methods. However, using
ordinary functions for operations does not affect performance, and we use
them freely. Instead of using general purpose functions we try to use
separate functions for different operations as much as possible. This will
reduce the number of if-then-else checks, but could harm the
maintainability of the code. To mitigate this problem we use code
generation techniques to generate repetitive codes and keep the code
maintainable.
*/
package avm
