# gomakedeps

Speed up your Makefile driven go project with auto-generated dependency Makefiles.

IMPORTANT: The program must be compiled on the target system, 
because import paths are resolved statically by the go/build package.

## Usage 

Install `go get github.com/Drachenfels-GmbH/go-makedeps`

Makefile inclusion:

```Makefile
BINS := helloworld foobar  

.%.deps.mk:
	echo -n "$*: " > $@
	gomakedeps $*.go >> $@

# include dependency files (except for the clean target)
ifneq ($(MAKECMDGOALS),clean)
-include $(BINS:%=.%.deps.mk)
endif
```
