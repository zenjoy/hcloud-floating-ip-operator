NAME=hcloud-floating-ip-operator
NAMESPACE=zenjoy
OS ?= linux
ifeq ($(strip $(shell git status --porcelain 2>/dev/null)),)
  GIT_TREE_STATE=clean
else
  GIT_TREE_STATE=dirty
endif
COMMIT ?= $(shell git rev-parse HEAD)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
PKG ?= github.com/$(NAMESPACE)/$(NAME)

CODE_GENERATOR_IMAGE := slok/kube-code-generator:v1.9.1
DIRECTORY := $(PWD)

## Bump the version in the version file. Set BUMP to [ patch | major | minor ]
BUMP := patch
VERSION ?= $(shell cat VERSION)

all: test

publish: compile build push clean

.PHONY: bump-version
bump-version: 
	@go get -u github.com/jessfraz/junk/sembump # update sembump tool
	$(eval NEW_VERSION = $(shell sembump --kind $(BUMP) $(VERSION)))
	@echo "Bumping VERSION from $(VERSION) to $(NEW_VERSION)"
	@echo $(NEW_VERSION) > VERSION

generate:
	docker run --rm -it \
	-v $(DIRECTORY):/go/src/$(PKG) \
	-e PROJECT_PACKAGE=$(PKG) \
	-e CLIENT_GENERATOR_OUT=$(PKG)/client/k8s \
	-e APIS_ROOT=$(PKG)/apis \
	-e GROUPS_VERSION="hcloud:v1alpha1" \
	-e GENERATION_TARGETS="deepcopy,client" \
	$(CODE_GENERATOR_IMAGE)

.PHONY: build
build:
	@echo "==> Building the docker image"
	@docker build -t $(NAMESPACE)/$(NAME):$(VERSION) -t $(NAMESPACE)/$(NAME) .

.PHONY: push
push:
ifeq ($(shell [[ $(BRANCH) != "master" && $(VERSION) != "dev" ]] && echo true ),true)
	@echo "ERROR: Publishing image with a SEMVER version '$(VERSION)' is only allowed from master"
else
	@echo "==> Publishing $(NAMESPACE)/$(NAME):$(VERSION)"
	@docker push $(NAMESPACE)/$(NAME)
	@docker push $(NAMESPACE)/$(NAME):$(VERSION)
	@echo "==> Your image is now available at $(NAMESPACE)/$(NAME):$(VERSION)"
endif
