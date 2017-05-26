# engine
Job execution framework

Working as a responsible for job management in a financial institution I created a perl script that could execute scripts in a certain order determined by their names. It was soon used as a general addition to control-m jobs and to carry out installations. It was so successful that I decided to create a go version for me (and you). The new version is a complete reimplementation of the concepts found in the script, not just a translation. It would have been wasted time to just translate the script without taking into account what we learned by the application of it by its users.

USAGE

engine path

DESCRIPTION

Engine will look for files in the path give for files with naming scheme [A-Z][0-9][0-9]{3}_.* and executes them in a certain order.

EXECUTION ORDER

[A-Z]
The first part is the stream. All groups with a given stream are executed sequentially in alphabethical order.

First Number
The parallel flag. All files with different numbers will be executed in parallel within the stream. E.g. A0001 and A1001 will be executed in parallel.

Rest of numbers
The rest of the numbers are executed sequentially, A0002 after A0001 after A0000

To find files, subdirectories will be considered


