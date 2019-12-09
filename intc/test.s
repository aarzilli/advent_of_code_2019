	// sums integers between 1 and 100 in a needlessly complicated way
	
	.const N 100 // input value
	
	addbas v
loop1:
	lt [i], N, [x]
	jz [x], step2
	add [i], 1, [base+0]
	add [i], 1, [i]
	addbas 1
	jz 0, loop1

step2:
	addbas -N
	add 0, 0, [i]

loop2:
	lt [i], N, [x]
	jz [x], step3
	add [base+0], [t], [t]
	add [i], 1, [i]
	addbas 1
	jz 0, loop2

step3:
	out [t]
	end
	
	.ord 100
	.var i, x, t
	.arr v N
