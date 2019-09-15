
## rgf

A wrapper around [ripgrep](https://github.com/BurntSushi/ripgrep) to check for various common patterns.

Heavily inspired by the great [gf](https://github.com/tomnomnom/gf) project.

## Why?
While developing my [Osmedeus](https://github.com/j3ssie/Osmedeus) I found myself needing to grep a lot patterns for file and folder so I wrote `rgf`.

## Usage

#### Config rgf

``` bash
go get github.com/j3ssie/rgf
```

then add signatures to `~/.rgf/` folder.

#### Example command

```bash
rgf

rgf -dir /folder/to/grep/
rgf -file whateverfile
rgf -dir /folder/to/grep/ url
```

## Adding pattern

```
./rgf.py -add url 'https?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+
```

or just create it manually in `~/.rgf/your_pattern.json`


## Contact

[@j3ssiejjj](https://twitter.com/j3ssiejjj)
