// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// handleMemclrFault returns (the value stored in R0, the value stored in R1).
// Control is transferred to it when memclr below receives SIGSEGV or SIGBUS,
// with the faulting address stored in R0 and the signal number stored in R1.
//
// It must have the same frame configuration as memclr so that it can undo any
// potential call frame set up by the assembler.
TEXT handleMemclrFault(SB), NOSPLIT, $0-28
	MOV	A0, addr+16(FP)
	MOV	A1, sig+24(FP)
	RET

// See the corresponding doc in safecopy_unsafe.go
//
// The code is derived from runtime.memclrNoHeapPointers.
//
// func memclr(ptr unsafe.Pointer, n uintptr) (fault unsafe.Pointer, sig int32)
TEXT ·memclr(SB), NOSPLIT, $0-28
	// Store 0 as the returned signal number. If we run to completion,
	// this is the value the caller will see; if a signal is received,
	// handleMemclrFault will store a different value in this address.
	MOV ZERO, sig+24(FP)
	MOV ptr+0(FP), A0
	MOV n+8(FP), A1

	// If less than 8 bytes, do single byte zeroing.
	MOV	$8, X9
	BLT	X11, X9, check4

	// Check alignment
	AND	$7, X10, X5
	BEQZ	X5, aligned

	// Zero one byte at a time until we reach 8 byte alignment.
	SUB	X5, X9, X5
	SUB	X5, X11, X11
align:
	SUB	$1, X5
	MOVB	ZERO, 0(X10)
	ADD	$1, X10
	BNEZ	X5, align

aligned:
	// X9 already contains $8
	BLT	X11, X9, check4
	MOV	$16, X9
	BLT	X11, X9, zero8
	MOV	$32, X9
	BLT	X11, X9, zero16
	MOV	$64, X9
	BLT	X11, X9, zero32
loop64:
	MOV	ZERO, 0(X10)
	MOV	ZERO, 8(X10)
	MOV	ZERO, 16(X10)
	MOV	ZERO, 24(X10)
	MOV	ZERO, 32(X10)
	MOV	ZERO, 40(X10)
	MOV	ZERO, 48(X10)
	MOV	ZERO, 56(X10)
	ADD	$64, X10
	SUB	$64, X11
	BGE	X11, X9, loop64
	BEQZ	X11, done

check32:
	MOV	$32, X9
	BLT	X11, X9, check16
zero32:
	MOV	ZERO, 0(X10)
	MOV	ZERO, 8(X10)
	MOV	ZERO, 16(X10)
	MOV	ZERO, 24(X10)
	ADD	$32, X10
	SUB	$32, X11
	BEQZ	X11, done

check16:
	MOV	$16, X9
	BLT	X11, X9, check8
zero16:
	MOV	ZERO, 0(X10)
	MOV	ZERO, 8(X10)
	ADD	$16, X10
	SUB	$16, X11
	BEQZ	X11, done

check8:
	MOV	$8, X9
	BLT	X11, X9, check4
zero8:
	MOV	ZERO, 0(X10)
	ADD	$8, X10
	SUB	$8, X11
	BEQZ	X11, done

check4:
	MOV	$4, X9
	BLT	X11, X9, loop1
zero4:
	MOVB	ZERO, 0(X10)
	MOVB	ZERO, 1(X10)
	MOVB	ZERO, 2(X10)
	MOVB	ZERO, 3(X10)
	ADD	$4, X10
	SUB	$4, X11

loop1:
	BEQZ	X11, done
	MOVB	ZERO, 0(X10)
	ADD	$1, X10
	SUB	$1, X11
	JMP	loop1

done:
	RET

// func addrOfMemclr() uintptr
TEXT ·addrOfMemclr(SB), $0-8
	MOV	$·memclr(SB), A0
	MOV	A0, ret+0(FP)
	RET
