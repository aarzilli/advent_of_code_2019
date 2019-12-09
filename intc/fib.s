	.usestack
	
	// print the first ten fibonacci numbers
	// demonstrates pseudo instructions frame, push, call, ret and pseudo-addressing mode fb-relative
	
	// call label, n
	// 	call label with n arguments
	// frame nargs, nlocals
	// 	adds space for nlocals on the stack
	// push arg
	// 	pushes arg on the stack
	// pop arg
	// 	pops a value from the stack into arg
	// ret
	// 	returns to the calling function
	// [result]
	// 	contains the result value of the last called function
	// [fb-n]
	// 	accesses n-th argument (first argument has n==1)
	// [fb+n]
	// 	accesses n-th local variable (first local variable has n==1)
	// [fb+0]
	// 	return address for the current function call
	
loop:	push [i]
	call fib, 1
	out [result]
	add [i], 1, [i]
	eq [i], 10, [x]
	jz [x], loop
	end

fib:	frame 1, 2
	eq [fb-1], 0, [fb+1]
	jnz [fb+1], ret1
	eq [fb-1], 1, [fb+1]
	jnz [fb+1], ret1
	add [fb-1], -1, [fb+1]
	push [fb+1]
	call fib, 1
	add [result], 0, [fb+1]
	add [fb-1], -2, [fb+2]
	push [fb+2]
	call fib, 1
	add [result], [fb+1], [result]
	ret
ret1: add 1, 0, [result]
	ret
	
	.ord 500
	.var i, x