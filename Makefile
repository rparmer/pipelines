.PHONY: all
all: manifests generate install

.PHONY: manifests
manifests: controller-gen
	$(CONTROLLER_GEN) crd paths="./api/..." output:crd:artifacts:config=crds

.PHONY: generate
generate: controller-gen
	$(CONTROLLER_GEN) object paths="./api/..."

.PHONY: run-crd
run-crd:
	go run . crd

.PHONY: run-cm
run-cm:
	go run . cm

PHONY: run-labels
run-labels:
	go run . labels

.PHONY: run
run:
	go run . 

.PHONY: install
install: generate manifests
	kubectl apply -f crds

.PHONY: uninstall
uninstall: generate manifests
	kubectl delete -f crds

.PHONY: demo
demo:
	kustomize build demo | kubectl apply -f -

.PHONY: uninstall-demo
uninstall-demo:
	kustomize build demo | kubectl delete -f -

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CONTROLLER_TOOLS_VERSION ?= v0.9.0

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
