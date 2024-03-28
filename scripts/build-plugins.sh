#!/bin/bash

ROOT=$(pwd)
POLICY_GENERATOR_TAG=release-2.10

rm -rf build || true
mkdir build
cd build || exit 1
# Download ACM Policy Generator plugin repo
git clone --branch="$POLICY_GENERATOR_TAG" --depth 1 https://github.com/stolostron/policy-generator-plugin
# Download PGT repo
git clone --depth 1 https://github.com/openshift-kni/cnf-features-deploy.git
# Build latest Policy Generator Template plugin executable
cd cnf-features-deploy/ztp/policygenerator || exit 1
make build
cp policygenerator "$ROOT"/kustomize/ran.openshift.io/v1/policygentemplate/PolicyGenTemplate
cd - || exit 1

# Build ACM Policy Generator plugin executable
cd policy-generator-plugin/ || exit 1
API_PLUGIN_PATH="." make build
cp PolicyGenerator "$ROOT"/kustomize/policy.open-cluster-management.io/v1/policygenerator/PolicyGenerator
