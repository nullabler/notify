.PHONY: help
help: help_message

help_message:
	@awk 'BEGIN {FS = ":.*##"; \
    	printf "$(<b>)$(<cYellow>) Usage:$(</>)$(</>) make $(<cGreen>)<target>$(</>) $(<cBlue>)\"<arguments>\"$(</>)️$(<br>)"\
		}/^[a-zA-Z_-]+:.*?##/\
		{ printf "  $(<cGreen>)%-30s$(</>) %s$(<br>)", $$1, $$2 } /^##@/\
		{ printf "$(<br>)$(<b>) $(<cYellow>)%s$(</>)$(<br>)", substr($$0, 5) }'$(MAKEFILE_LIST)