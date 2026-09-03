SHELL := /usr/bin/env bash

.PHONY: check ci cohesion inventory repository-check

check:
	golib check --all

ci: repository-check cohesion check

cohesion:
	golib cohesion check

inventory repository-check:
	golib repository check
