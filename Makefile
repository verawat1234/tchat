SHELL := /bin/bash

.PHONY: bootstrap
bootstrap: ## Install toolchains and dependencies (manual per stack for now)
	@echo "→ Ensure Go, Node.js, Android Studio, and Xcode CLT are installed."
