cloudbuild_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
gcp_build_project := einride-build

.PHONY: all
all: \
	triggers-create-review \

.PHONY: triggers-create-review
triggers-create-review: $(cloudbuild_cwd)/triggers/review.yaml $(cloudbuild_cwd)/configs/review.yaml
	gcloud beta builds triggers create github --project $(gcp_build_project) --trigger-config=$<
