getpos $.name $x $y

set $radius 1.6

set $xn $x
set $yn $y

#"Scale x into plane"
#"width 200 -> (-2, .47)"
divF $xn 200 $xn
mulF $xn 2.47 $xn
subF $xn 2 $xn

#"Scale y into plane"
#"height 200 -> (-1.12, 1.12)"
divF $yn 200 $yn
mulF $yn 2.24 $yn
subF $yn 1.12 $yn

set $x0 $xn
set $y0 $yn

#"Iterate -> if length > $radius jmp end -> yield -> goto startLoop"


call length
yield

set $a 20
jmpl sayLabel $x $a

END

sayLabel:
dblogf 2 "lt %s" $x
END


#"sets $length to the distance from the origin"
length:
mulF $yn $yn $yn2
mulF $xn $xn $xn2
addF $yn2 $xn2 $s
sqrtF $s $length
ret