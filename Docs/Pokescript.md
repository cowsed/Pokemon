# Functions specific to the pokemon recreation

### say
---
```
say message
```
Prints the message to the dialogue window. Will hold the script until enter is pressed

### sayf
---
```
say n format ...
```
Prints the message to the dialogue window. Will hold the script until enter is pressed. 
- n - the number of tokens used. 1+number of variadic arguments (format string + args)
- format - the format string of the message. Use %s in every situation
- ... - stands for a variable number of arguments. if the number of arguments does not = the number of format specifiers in the format string, bad things will happen. Will not block.

### dblog 
---
```
dblog message
```
Prints the message to the debug console

### dblogf 
---
```
say n format ...
```
does the same thing as sayf but instead of going to the dialogue box goes to the log. Will not block.