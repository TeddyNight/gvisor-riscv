#include "funcdata.h"
#include "textflag.h"

// See kernel.go.
TEXT ·Halt(SB),NOSPLIT,$0
	//WFI
	RET
