all: delete_latest_tag recreate_tag push_tag

delete_latest_tag:
	@echo "Deleting latest local tag: latest"
	@git tag -d latest

recreate_tag:
	@echo "Recreating tag: latest"
	@git tag latest

push_tag:
	@echo "Force pushing tag: latest"
	@git push --force origin latest

.PHONY: all delete_latest_tag recreate_tag push_tag
