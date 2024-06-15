all: delete_latest_tag recreate_tag push_tag

# Target to delete the latest tag
delete_latest_tag:
	@echo "Deleting latest local tag: latest"
	@git tag -d latest

# Target to recreate the tag
recreate_tag:
	@echo "Recreating tag: latest"
	@git tag latest

# Target to force push the tag
push_tag:
	@echo "Force pushing tag: latest"
	@git push --force origin latest

.PHONY: all delete_latest_tag recreate_tag push_tag
