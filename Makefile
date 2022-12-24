CURRENT_DIR=$(shell pwd)

submodule-gen:
	rm -rf modules/template_variables
	rm -rf modules/template_protos
	mkdir -p modules/template_variables
	mkdir -p modules/template_protos
	rsync -r --exclude '.git' template_variables/ modules/template_variables
	rsync -r --exclude '.git' template_protos/ modules/template_protos

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}