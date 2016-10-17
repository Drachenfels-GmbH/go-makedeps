# Gomakefile

Generates makefile dependencies for a go main program.

Makefile inclusion:

```
#$(BINS:%=%.deps.mk)

.%.deps.mk:
	echo -n "$*: " > $@
	gomakedeps $*.go >> $@

ifneq ($(MAKECMDGOALS),clean)
-include $(BINS:%=.%.deps.mk)
endif
```
