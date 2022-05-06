# Scripting Language

The ability for maps and characters to interact with the world is done by a scripting language. It's technically turing complete I think but should have everything a pokemon trainer needs to do

## Syntax

The basic syntax of the language is roughly as follows

```
operation parameter "paramater"
```

operations correspond to the available functions that are registered (See Instructions)
paramaters can either be variable or string literals. String literals can either be within quotes or without quotes if there are no spaces

## Datatypes

Everything is stored in memory as a string. Memory is a map of strings (variable names) to other strings (variable values)

Some functions (addI) cast the string value in memory to the memory type it needs. If the string can not be parsed as that type, an error will occur

## Variables

Variables are referenced by the $ Syntax. A variable with the name x would be referred to as '$x'. 'x' would correspond to the string value.

## Built-in instructions

### END

---

```
END
```

Terminates the script

### set

---

```
set variable value
```

Sets the value of variable to value.
variable must be the name of a variable.
value can be the name of a variable if prefixed by $ or a literal string value

### addI

---

```
addI a b c
```

a and b can be variables or literals
c must be a variable

adds a + b and stores it to c

Internally it casts the strings to ints then adds then converts it back to a string to store.


### subI

See addI

### mulI

See addI

### divI

See addI

### addF
```
addF a b c
```

a and b can be variables or literals
c must be a variable

adds a + b and stores it to c

Internally it casts the strings to float64s then adds then converts it back to a string to store.

### subF

See addF

### mulF

See addF

### divF

See addF


### jmpe

----

```
jmpe label a b
```

Jump to label if a == b. a and b can be variables or literals

### jmpne

----

```
jmpne label a b
```

Jump to label if a != b. a and b can be variables or literals



### goto

---

```
goto label
```

unconditionally goes to that label

### call

---

```
call label
```

pushes the current location and goes to label. Used to enter subroutine

### ret

---

```
ret
```

returns to the top value of the stack. Used to exit subroutine
