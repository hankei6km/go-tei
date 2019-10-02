# go-tei [![GoDoc](https://godoc.org/github.com/hankei6km/go-tei?status.svg)](https://godoc.org/github.com/hankei6km/go-tei) [![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org) [![Build Status](https://travis-ci.org/hankei6km/go-tei.svg?branch=master)](https://travis-ci.org/hankei6km/go-tei)

tei switch the piped input to another one if no data from the piped input, and simply use to just check no data from the piped input.

```
$ echo " = >" | tei string "NO DATA" | figlet -f banner
            #    
             #   
   #####      #  
               # 
   #####      #  
             #   
            #    

$ echo "" | tei string "NO DATA" | figlet -f banner
#     # #######    ######     #    #######    #    
##    # #     #    #     #   # #      #      # #   
# #   # #     #    #     #  #   #     #     #   #  
#  #  # #     #    #     # #     #    #    #     # 
#   # # #     #    #     # #######    #    ####### 
#    ## #     #    #     # #     #    #    #     # 
#     # #######    ######  #     #    #    #     # 

```

## Installation

`tei` binary file can be downloaded from [release page](https://github.com/hankei6km/go-tei/releases).

## Usage

```
tei [flags] <exit_code>

tei [flags] run <command> [command_args]...

tei [flags] file <input_file>

tei [flags] string [-n] <string>...

Global Flags:
  -l, --ignore-newline   ignore leading a newline while sniffing the input (default true)
```

## Example

swtich the piped input.
```console
$ echo "clipboard data" | xsel -i

$ tei run xsel -b -o | tr '[:lower:]' '[:upper:]'
CLIPBOARD DATA

$ echo "" | tei run xsel -b -o | tr '[:lower:]' '[:upper:]'
CLIPBOARD DATA

$ echo "input data" | tei run xsel -b -o | tr '[:lower:]' '[:upper:]'
INPUT DATA
```

check no data from the piped input.
```console
$ tei 1                              # exit code = 0

$ echo "" | tei 1                    # exit code = 0

$ echo "test" | tei 1                # exit code = 1

$ echo "" | tei -l=false 1           # exit code = 1

$ echo -n "" | tei -l=false 1        # exit code = 0

$ echo -en "\ntest" | tei 1          # exit code = 1

$ echo "test" | tei -l=false -p 1    # exit code = 1
test
```

## Code

```go
import "github.com/hankei6km/go-tei"
```

```go
output := func(builder tei.Builder, r io.Reader) {
	tei := builder.Standby(func() io.Reader {
		return strings.NewReader("standby-data")
	}).Build()

	r = tei.Switch(r)
	io.Copy(os.Stdout, r)
	fmt.Println()
}
output(tei.NewBuilder(), os.Stdin)
output(tei.NewBuilder(), strings.NewReader(""))
output(tei.NewBuilder(), strings.NewReader("\n"))
output(tei.NewBuilder(), strings.NewReader("input-data"))

// Output:
//
// standby-data
// standby-data
// standby-data
// input-data
```

```go
output := func(builder tei.Builder, r io.Reader) {
	tei := builder.Standby(func() io.Reader {
		return nil
	}).Build()

	if tei.Switch(r) == nil {
		fmt.Println("no-data")
		return
	}
	fmt.Println("data")

}
output(tei.NewBuilder(), strings.NewReader("\n"))
output(tei.NewBuilder(), strings.NewReader("\ninput-data"))
output(tei.NewBuilder().IgnoreLeadingNewline(false), strings.NewReader("\n"))

// Output:
//
// no-data
// data
// data
```