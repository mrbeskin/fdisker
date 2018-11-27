# fdisker

`fdisker` is a simple command to run commands specificed in a file with fdisk. It has its own (very simple) syntax, so please see below for details on how to write your fdisk "script".

Requirements:
 - `fdisk` must be installed

## using fdisker
### cli

The cli is used as follows: 

```
$ fdisker -f /path/to/file 
```

The `-f` flag with a path to the filename is required.

If you do not wish to actually save the commands after you execute them, use the --write-off flag: 

```
$ fdisker -f /path/to/file --write-off=false
```

The above command would start `fdisk` and run all commands from your file, but quit without saving, as per the `fdisk` interactive `q` command.

### go lib

fdisker may also be used as a Go library. 

There is one function:

```
func RunFdiskCommandFile(path, mountPath string, writeFlag bool) error
```

Its usage is simple: 

```
err := fdisker.RunFdiskCommandFile("/path/to/file", "/dev/sda2", true)
if err != nil {
	// handle error
}
```

Executing the above code will run the commands described in `/path/to/file` using `fdisk` on `/dev/sda2`. 

`fdisk` output will be written to stdout and stderr. 

In the case of an error, or if `writeFlag` is set to `false`, no the commands will not be run.

## fdisker file syntax

### warning: do not quit in your file

`fdisker` will quit fdisk automatically. It will default to writing your command to disk, so please test it. If you want to run the commands but not save the outcome, please use the write flag. 

### default command

Newline characters are not interpreted as defaults. Instead there is a new default token that must be provided: `DEF`. This will allow a prompt to continue with the default.

### all non-default, non-quit commands



### comments

Basic comments have been implemented to allow you to describe your `fdisk` file. A comment starts with a `#`. Any line starting with `#` will be ignored. 

Prefix spaces are allowed and will be ignored. 

Any line not starting with `#` will be interpreted as a command.

### example file 

The following would delete a partitoin and then create a new partition using the default partition number, and default start and end locations. 

``` 
# This is an example fdisker file

d   # A comment that describes that I am deleting a partition
DEF # I am deleting the default provided by fdisk
n
p
DEF
DEF
DEF
```


