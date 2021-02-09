GOBIN=${PWD}/bin

SCRATCH?=n
NO_CACHE?=
ifeq ($(SCRATCH),y)
NO_CACHE=--no-cache
endif

# Default to build
.PHONY: all
all: sriov-dp sriov-cni multus dpdk-devel dpdk-centos

help:
	@echo "Make Targets:"
	@echo " make sriov-dp         - Make the docker image that runs the SR-IOV Device"
	@echo "                         Plugin with vDPA changes integrated. Append SCRATCH=y"
	@echo "                         re-download upstream repo and to build image using '--no-cache'."
	@echo " make multus           - Make the multus container image that deploys multus binary on all nodes"
	@echo " make sriov-cni        - Make the SR-IOV CNI binary with the vDPA changes"
	@echo "                         integrated. Binary needs to copied to proper location"
	@echo "                         once complete (i.e. - /opt/cni/bin/.)."
	@echo " make dpdk-devel       - Make the development container. It has dpdk as well as userspace utilities."
	@echo "                         Useful for testing or development purposes"
	@echo "                         Append SCRATCH=y to build image using '--no-cache'."
	@echo ""
	@echo " make dpdk-app         - Make the centos8-based DPDK app powered by app-netutils."
	@echo "                         This sample container is able to run l2fwd, l3fwd and testpmd by autodetecting"
	@echo "                         the configured network devices and creating the apropriate DPDK parameters"
	@echo "                         Append SCRATCH=y to build image using '--no-cache'."
	@echo ""
	@echo " make                  - Build all the local sub-projects locally."
	@echo " make clean            - Cleanup all build artifacts."
	@echo " make all              - Build all images for a deployment. Same as:"
	@echo "                           make sriov-dp; make sriov-cni; make multus; make dpdk-devel; make dpdk-centos"
	@echo " deploy                - Deploy the minimum pods and daemonsets. Same as:"
	@echo "                           for f in sriovdp-vdpa-daemonset.yaml multus-daemonset.yaml sriovcni-vdpa-daemonset.yaml; do"
	@echo "                              kubectl apply -f deployment/\$$f;"
	@echo "                           done;"
	@echo " show                   - Show a summary of the SR-IOV configuration" 
	@echo ""


#
# Archive or WIP targets
#
#httpd-init:
#	@cd seastar-httpd/init-container && go build -o ${GOBIN}/httpd-init -v
#
#scylla-init:
#	@cd scylla-init-container && go build -o ${GOBIN}/scylla-init -v
#
#dpdk-app:
#	@echo ""
#	@echo "dpdk-app $(NO_CACHE) ..."
#	@docker build $(NO_CACHE) --rm -t dpdk-app-centos -f ./dpdk-app-centos/Dockerfile .
#
#httpd-image:
#	@echo ""
#	@echo "Making httpd-image $(NO_CACHE) ..."
#	@docker build $(NO_CACHE) --rm -t seastar-httpd -f ./seastar-httpd/httpd/Dockerfile .
#
#httpd-init-image:
#	@echo ""
#	@echo "Making httpd-init-image $(NO_CACHE) ..."
#	@docker build $(NO_CACHE) --rm -t httpd-init-container -f ./seastar-httpd/init-container/Dockerfile .
#
#scylla-image:
#	@echo ""
#	@echo "Making scylla-image $(NO_CACHE) ..."
#	@docker build $(NO_CACHE) --rm -t scylla-init-container -f ./scylla-init-container/Dockerfile .

# SRI-IOV CNI and DP configuration
export ORG_PATH="github.com/intel"
export REPO_PATH_CNI="${ORG_PATH}/sriov-cni"
export REPO_PATH_DP="${ORG_PATH}/sriov-network-device-plugin"
export GOBIN=${PWD}/bin

## SR-IOV CNI
export ALT_CNI_REPO=https://github.com/amorenoz/sriov-cni.git
export ALT_CNI_REF=rfe/vdpa

.PHONY: clean-sriov-cni
clean-sriov-cni:
	@if [ -d gopath/src/$(REPO_PATH_CNI) ]; then \
	    pushd gopath/src/$(REPO_PATH_CNI) > /dev/null; \
	    make clean >/dev/null; \
	    popd > /dev/null; \
	    rm -fr gopath/src/$(REPO_PATH_CNI); \
	fi \

.PHONY: sriov-cni
ifeq ($(SCRATCH),y)
sriov-cni: clean-sriov-cni
 endif
sriov-cni: export GOPATH=${PWD}/gopath
sriov-cni:
	@if [ ! -d gopath/src/$(REPO_PATH_CNI) ]; then \
		echo ""; \
		echo "Making sriov-cni ..."; \
		echo "Downloading $(REPO_PATH_CNI)"; \
		mkdir -p gopath/src/$(ORG_PATH); \
		pushd gopath/src/ > /dev/null; \
		go get $(REPO_PATH_CNI) 2>&1 > /tmp/sriov-dp.log || echo "Can ignore no GO files."; \
		popd > /dev/null; \
		if [ -n "$(ALT_CNI_REPO)" ];  then \
		    echo "Checking out alternative repository $(ALT_CNI_REPO):$(ALT_CNI_REF)"; \
		    pushd gopath/src/$(REPO_PATH_CNI) > /dev/null; \
		    git remote add alt $(ALT_CNI_REPO) && git fetch alt; \
		    if [ -n "$(ALT_CNI_REF)" ]; then  \
			git checkout alt/$(ALT_CNI_REF) -b alt; \
		    fi; \
		    popd > /dev/null; \
		fi; \
		pushd gopath/src/$(REPO_PATH_CNI) > /dev/null; \
		echo "Build CNI image"; \
		make image; \
		popd > /dev/null; \
	fi

export ALT_DP_REPO=https://github.com/amorenoz/sriov-network-device-plugin.git
export ALT_DP_REF=vdpaInfoProvider

.PHONY: clean-sriov-dp
clean-sriov-dp:
	@if [ -d gopath/src/$(REPO_PATH_DP) ]; then \
	    pushd gopath/src/$(REPO_PATH_DP) > /dev/null; \
	    make clean >/dev/null; \
	    popd > /dev/null; \
	    rm -fr gopath/src/$(REPO_PATH_DP); \
	fi \

.PHONY: sriov-dp
ifeq ($(SCRATCH),y)
sriov-dp: clean-sriov-dp
 endif
sriov-dp: export GOPATH=${PWD}/gopath
sriov-dp:
	@if [ ! -d gopath/src/$(REPO_PATH_DP) ]; then \
		echo ""; \
		echo "Making sriov-dp ..."; \
		echo "Downloading $(REPO_PATH_DP)"; \
		mkdir -p gopath/src/$(ORG_PATH); \
		pushd gopath/src/ > /dev/null; \
		go get $(REPO_PATH_DP) 2>&1 > /tmp/sriov-dp.log || echo "Can ignore no GO files."; \
		popd > /dev/null; \
		if [ -n "$(ALT_DP_REPO)" ];  then \
		    pushd gopath/src/$(REPO_PATH_DP) > /dev/null; \
		    git remote add alt $(ALT_DP_REPO) && git fetch alt; \
		    if [ -n "$(ALT_DP_REF)" ]; then  \
			git checkout alt/$(ALT_DP_REF) -b alt; \
		    fi; \
		    popd > /dev/null; \
		fi; \
		pushd gopath/src/$(REPO_PATH_DP) > /dev/null; \
		echo "Build binary"; \
		make; \
		echo "Build docker image \"sriov-device-plugin\""; \
		make image; \
		popd > /dev/null; \
	fi

.PHONY: clean-multus
REPO_PATH_MULTUS=github.com/intel/multus-cni
ALT_MULTUS_REPO=https://github.com/amorenoz/multus-cni.git
ALT_MULTUS_REF=bug/runtime

clean-multus:
	@docker image rm -f nfvpe/multus || true
	@if [ -d gopath/src/$(REPO_PATH_MULTUS) ]; then \
	    rm -fr gopath/src/$(REPO_PATH_MULTUS); \
	fi \

.PHONY: multus
ifeq ($(SCRATCH),y)
multus: clean-multus
endif
multus: export GOPATH=${PWD}/gopath
multus:
	@if [ ! -d gopath/src/$(REPO_PATH_MULTUS) ]; then \
		echo ""; \
		echo "Making multus ..."; \
		echo "Downloading $(REPO_PATH_MULTUS)"; \
		mkdir -p gopath/src/$(ORG_PATH); \
		pushd gopath/src/ > /dev/null; \
		go get $(REPO_PATH_MULTUS) 2>&1 > /tmp/multus.log || echo "Can ignore no GO files."; \
		popd > /dev/null; \
		if [ -n "$(ALT_MULTUS_REPO)" ];  then \
		    echo "Checking out alternative repository $(ALT_MULTUS_REPO):$(ALT_MULTUS_REF)"; \
		    pushd gopath/src/$(REPO_PATH_MULTUS) > /dev/null; \
		    git remote add alt $(ALT_MULTUS_REPO) && git fetch alt; \
		    if [ -n "$(ALT_MULTUS_REF)" ]; then  \
			git checkout alt/$(ALT_MULTUS_REF) -b alt; \
		    fi; \
		    popd > /dev/null; \
		fi; \
		pushd gopath/src/$(REPO_PATH_MULTUS) > /dev/null; \
		echo "Build MULTUS image"; \
		docker build -t nfvpe/multus -f deployments/Dockerfile .; \
		popd > /dev/null; \
	fi

.PHONY: clean
clean: clean-sriov-dp clean-sriov-cni clean-multus
	@export GOPATH=${PWD}/gopath && go clean --modcache
	@rm -rf bin/
	@rm -rf gopath/


.PHONY: dpdk-devel
dpdk-devel:
	@echo ""
	@echo "dpdk-app-devel $(NO_CACHE) ..."
	@cd dpdk-app-devel; docker build $(NO_CACHE) --rm -t dpdk-app-devel .

.PHONY: dpdk-centos
dpdk-centos:
	@echo ""
	@echo "dpdk-app-centos $(NO_CACHE) ..."
	@cd dpdk-app-centos; docker build $(NO_CACHE) --rm -t dpdk-app-centos .

.PHONY: deploy
deploy:
	@for f in sriovdp-vdpa-daemonset.yaml multus-daemonset.yaml sriovcni-vdpa-daemonset.yaml; do\
		kubectl delete -f deployment/$${f} 2>/dev/null || true;\
		kubectl apply -f deployment/$${f};\
	done;\
.PHONY: show
show:
	@echo "SR-IOV Device Plugin: Discovered Devices:"
	@for node in $$(kubectl get nodes | grep Ready | awk '{print $$1}' ); do \
		echo "Node $$node:" ; kubectl get node $$node -o json | jq '.status.allocatable'; \
	done
	@echo ""
	@echo "Network Attachment Definitions"
	@ for na in $$(kubectl get network-attachment-definition | grep -v NAME | awk '{print $$1}'); do \
		kubectl get network-attachment-definition $$na -o jsonpath='Name = {.metadata.name} => Resource: {.metadata.annotations.k8s\.v1\.cni\.cncf\.io/resourceName}'; echo "" ;\
	done

	
