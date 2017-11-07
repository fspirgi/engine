# engine
## Job execution framework

Working as a responsible for job management in a financial institution I created a perl script that could execute scripts in a certain order determined by their names. It was soon used as a general addition to control-m jobs and to carry out installations. It was so successful that I decided to create a go version for me (and you). The new version is a complete reimplementation of the concepts found in the script, not just a translation. It would have been wasted time to just translate the script without taking into account what we learned by the application of it by its users.
The original script is still heavily used for automating a portifolio accounting application.

## USAGE

engine --path path [--stream regexp]

## DESCRIPTION

Engine will look for files in the path give for files with naming scheme [A-Z][0-9][0-9]{3}_.* and executes them in a certain order.

## EXECUTION ORDER

### [A-Z]
The first part is the stream. All groups with a given stream are executed sequentially in alphabethical order.

### First Number
The parallel flag. All files with different numbers will be executed in parallel within the stream. E.g. A0001 and A1001 will be executed in parallel.

### Rest of numbers
The rest of the numbers are executed sequentially, A0002 after A0001 after A0000

To find files, subdirectories will be considered

## CONFIGURATION FILES

For each part of the filename a corresponding configuration file is searched for in the whole tree above the directory. For example a executed file named /a/b/A0001_gugus configuration files Arc, A0rc, A00rc, A0001rc will be searched in all directories called etc in the path, e.g. /a/etc, /a/b/etc.

Configuration entries are in the form key = value, you can use keys as variables and extract their values with $. E.g.

a = b c

d = c $a f # d will be "c b c f"

Parts after # will be treated as comments and empty lines ignored.

Every value is pushed into the environment and available to the called programs or scripts.

