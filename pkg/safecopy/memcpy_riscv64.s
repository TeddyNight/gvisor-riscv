// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// handleMemcpyFault returns (the value stored in R0, the value stored in R1).
// Control is transferred to it when memcpy below receives SIGSEGV or SIGBUS,
// with the faulting address stored in R0 and the signal number stored in R1.
//
// It must have the same frame configuration as memcpy so that it can undo any
// potential call frame set up by the assembler.
TEXT handleMemcpyFault(SB), NOSPLIT, $0-36
	MOV A0, addr+24(FP)
	MOV A1, sig+32(FP)
	RET

// memcpy copies data from src to dst. If a SIGSEGV or SIGBUS signal is received
// during the copy, it returns the address that caused the fault and the number
// of the signal that was received. Otherwise, it returns an unspecified address
// and a signal number of 0.
//
// Data is copied in order, such that if a fault happens at address p, it is
// safe to assume that all data before p-maxRegisterSize has already been
// successfully copied.
//
// The code is derived from the Go source runtime.memmove.
//
// func memcpy(dst, src unsafe.Pointer, n uintptr) (fault unsafe.Pointer, sig int32)
TEXT ·memcpy(SB), NOSPLIT, $-8-36
	// Store 0 as the returned signal number. If we run to completion,
	// this is the value the caller will see; if a signal is received,
	// handleMemcpyFault will store a different value in this address.
	MOV ZERO, sig+32(FP)

	MOV  dst+0(FP), X10
	MOV  src+8(FP), X11
	MOV  n+16(FP), X12

	// X10 = to
	// X11 = from
	// X12 = n
	BEQ	X10, X11, done
	BEQZ	X12, done

	// If the destination is ahead of the source, start at the end of the
	// buffer and go backward.
	BGTU	X10, X11, backward

	// If less than 8 bytes, do single byte copies.
	MOV	$8, X9
	BLT	X12, X9, f_loop4_check

	// Check alignment - if alignment differs we have to do one byte at a time.
	AND	$7, X10, X5
	AND	$7, X11, X6
	BNE	X5, X6, f_loop8_unaligned_check
	BEQZ	X5, f_loop_check

	// Move one byte at a time until we reach 8 byte alignment.
	SUB	X5, X9, X5
	SUB	X5, X12, X12
f_align:
	ADDI    $1, ZERO, T1
	SUB	T1, X5
	MOVB	0(X11), X14
	MOVB	X14, 0(X10)
	ADD	$1, X10
	ADD	$1, X11
	BNEZ	X5, f_align

f_loop_check:
	MOV	$16, X9
	BLT	X12, X9, f_loop8_check
	MOV	$32, X9
	BLT	X12, X9, f_loop16_check
	MOV	$64, X9
	BLT	X12, X9, f_loop32_check
f_loop64:
	MOV	0(X11), X14
	MOV	8(X11), X15
	MOV	16(X11), X16
	MOV	24(X11), X17
	MOV	32(X11), X18
	MOV	40(X11), X19
	MOV	48(X11), X20
	MOV	56(X11), X21
	MOV	X14, 0(X10)
	MOV	X15, 8(X10)
	MOV	X16, 16(X10)
	MOV	X17, 24(X10)
	MOV	X18, 32(X10)
	MOV	X19, 40(X10)
	MOV	X20, 48(X10)
	MOV	X21, 56(X10)
	ADD	$64, X10
	ADD	$64, X11
	ADDI    $64, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, f_loop64
	BEQZ	X12, done

f_loop32_check:
	MOV	$32, X9
	BLT	X12, X9, f_loop16_check
f_loop32:
	MOV	0(X11), X14
	MOV	8(X11), X15
	MOV	16(X11), X16
	MOV	24(X11), X17
	MOV	X14, 0(X10)
	MOV	X15, 8(X10)
	MOV	X16, 16(X10)
	MOV	X17, 24(X10)
	ADD	$32, X10
	ADD	$32, X11
	ADDI    $32, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, f_loop32
	BEQZ	X12, done

f_loop16_check:
	MOV	$16, X9
	BLT	X12, X9, f_loop8_check
f_loop16:
	MOV	0(X11), X14
	MOV	8(X11), X15
	MOV	X14, 0(X10)
	MOV	X15, 8(X10)
	ADD	$16, X10
	ADD	$16, X11
	ADDI    $16, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, f_loop16
	BEQZ	X12, done

f_loop8_check:
	MOV	$8, X9
	BLT	X12, X9, f_loop4_check
f_loop8:
	MOV	0(X11), X14
	MOV	X14, 0(X10)
	ADD	$8, X10
	ADD	$8, X11
	ADDI    $8, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, f_loop8
	BEQZ	X12, done
	JMP	f_loop4_check

f_loop8_unaligned_check:
	MOV	$8, X9
	BLT	X12, X9, f_loop4_check
f_loop8_unaligned:
	MOVB	0(X11), X14
	MOVB	1(X11), X15
	MOVB	2(X11), X16
	MOVB	3(X11), X17
	MOVB	4(X11), X18
	MOVB	5(X11), X19
	MOVB	6(X11), X20
	MOVB	7(X11), X21
	MOVB	X14, 0(X10)
	MOVB	X15, 1(X10)
	MOVB	X16, 2(X10)
	MOVB	X17, 3(X10)
	MOVB	X18, 4(X10)
	MOVB	X19, 5(X10)
	MOVB	X20, 6(X10)
	MOVB	X21, 7(X10)
	ADD	$8, X10
	ADD	$8, X11
	ADDI    $8, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, f_loop8_unaligned

f_loop4_check:
	MOV	$4, X9
	BLT	X12, X9, f_loop1
f_loop4:
	MOVB	0(X11), X14
	MOVB	1(X11), X15
	MOVB	2(X11), X16
	MOVB	3(X11), X17
	MOVB	X14, 0(X10)
	MOVB	X15, 1(X10)
	MOVB	X16, 2(X10)
	MOVB	X17, 3(X10)
	ADD	$4, X10
	ADD	$4, X11
	ADDI    $4, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, f_loop4

f_loop1:
	BEQZ	X12, done
	MOVB	0(X11), X14
	MOVB	X14, 0(X10)
	ADD	$1, X10
	ADD	$1, X11
	ADDI    $1, ZERO, T1
	SUB	T1, X12
	JMP	f_loop1

backward:
	ADD	X10, X12, X10
	ADD	X11, X12, X11

	// If less than 8 bytes, do single byte copies.
	MOV	$8, X9
	BLT	X12, X9, b_loop4_check

	// Check alignment - if alignment differs we have to do one byte at a time.
	AND	$7, X10, X5
	AND	$7, X11, X6
	BNE	X5, X6, b_loop8_unaligned_check
	BEQZ	X5, b_loop_check

	// Move one byte at a time until we reach 8 byte alignment.
	SUB	X5, X12, X12
b_align:
	ADDI    $1, ZERO, T1
	SUB	T1, X5
	SUB	T1, X10
	SUB	T1, X11
	MOVB	0(X11), X14
	MOVB	X14, 0(X10)
	BNEZ	X5, b_align

b_loop_check:
	MOV	$16, X9
	BLT	X12, X9, b_loop8_check
	MOV	$32, X9
	BLT	X12, X9, b_loop16_check
	MOV	$64, X9
	BLT	X12, X9, b_loop32_check
b_loop64:
	ADDI    $64, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOV	0(X11), X14
	MOV	8(X11), X15
	MOV	16(X11), X16
	MOV	24(X11), X17
	MOV	32(X11), X18
	MOV	40(X11), X19
	MOV	48(X11), X20
	MOV	56(X11), X21
	MOV	X14, 0(X10)
	MOV	X15, 8(X10)
	MOV	X16, 16(X10)
	MOV	X17, 24(X10)
	MOV	X18, 32(X10)
	MOV	X19, 40(X10)
	MOV	X20, 48(X10)
	MOV	X21, 56(X10)
	ADDI    $64, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, b_loop64
	BEQZ	X12, done

b_loop32_check:
	MOV	$32, X9
	BLT	X12, X9, b_loop16_check
b_loop32:
	ADDI    $32, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOV	0(X11), X14
	MOV	8(X11), X15
	MOV	16(X11), X16
	MOV	24(X11), X17
	MOV	X14, 0(X10)
	MOV	X15, 8(X10)
	MOV	X16, 16(X10)
	MOV	X17, 24(X10)
	ADDI    $32, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, b_loop32
	BEQZ	X12, done

b_loop16_check:
	MOV	$16, X9
	BLT	X12, X9, b_loop8_check
b_loop16:
	ADDI    $16, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOV	0(X11), X14
	MOV	8(X11), X15
	MOV	X14, 0(X10)
	MOV	X15, 8(X10)
	ADDI    $16, ZERO, T1
	SUB	T1, X12
	BGE	X12, X9, b_loop16
	BEQZ	X12, done

b_loop8_check:
	MOV	$8, X9
	BLT	X12, X9, b_loop4_check
b_loop8:
	ADDI    $8, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOV	0(X11), X14
	MOV	X14, 0(X10)
	SUB	T1, X12
	BGE	X12, X9, b_loop8
	BEQZ	X12, done
	JMP	b_loop4_check

b_loop8_unaligned_check:
	MOV	$8, X9
	BLT	X12, X9, b_loop4_check
b_loop8_unaligned:
	ADDI    $8, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOVB	0(X11), X14
	MOVB	1(X11), X15
	MOVB	2(X11), X16
	MOVB	3(X11), X17
	MOVB	4(X11), X18
	MOVB	5(X11), X19
	MOVB	6(X11), X20
	MOVB	7(X11), X21
	MOVB	X14, 0(X10)
	MOVB	X15, 1(X10)
	MOVB	X16, 2(X10)
	MOVB	X17, 3(X10)
	MOVB	X18, 4(X10)
	MOVB	X19, 5(X10)
	MOVB	X20, 6(X10)
	MOVB	X21, 7(X10)
	SUB	T1, X12
	BGE	X12, X9, b_loop8_unaligned

b_loop4_check:
	MOV	$4, X9
	BLT	X12, X9, b_loop1
b_loop4:
	ADDI    $4, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOVB	0(X11), X14
	MOVB	1(X11), X15
	MOVB	2(X11), X16
	MOVB	3(X11), X17
	MOVB	X14, 0(X10)
	MOVB	X15, 1(X10)
	MOVB	X16, 2(X10)
	MOVB	X17, 3(X10)
	SUB	T1, X12
	BGE	X12, X9, b_loop4

b_loop1:
	BEQZ	X12, done
	ADDI    $1, ZERO, T1
	SUB	T1, X10
	SUB	T1, X11
	MOVB	0(X11), X14
	MOVB	X14, 0(X10)
	SUB	T1, X12
	JMP	b_loop1

done:
	RET

// func addrOfMemcpy() uintptr
TEXT ·addrOfMemcpy(SB), $0-8
	MOV	$·memcpy(SB), A0
	MOV	A0, ret+0(FP)
	RET
