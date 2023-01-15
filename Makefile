.PHONY: deploy-demo help

help: ## Show this help
	@echo 'usage: make [target]'
	@echo
	@echo 'available targets:'
	@egrep '^.+\:.*##\ .+' ${MAKEFILE_LIST} | sed -r 's/^(.+): .*##\ (.+)/  \1#\2/' | column -t -c 2 -s '#'

deploy-demo: ## Deploys demo instance using GCP Cloud Build
	gcloud config set project authz-374814
	gcloud builds submit --config ./cloudbuild.demo.yaml
