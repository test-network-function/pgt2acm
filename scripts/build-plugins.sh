#!/bin/bash

ROOT=$(pwd)

rm -rf build | true
mkdir build
cd build
# Download ACM Policy Generator plugin repo
git clone  --depth 1 https://github.com/open-cluster-management-io/policy-generator-plugin.git
# Download PGT repo
git clone  --depth 1 git@github.com:openshift-kni/cnf-features-deploy.git
# Build latest Policy Generator Template plugin executable
cd cnf-features-deploy/ztp/policygenerator
make build
cp policygenerator "$ROOT"/kustomize/ran.openshift.io/v1/policygentemplate/PolicyGenTemplate
cd -

# Build ACM Policy Generator plugin executable
cd policy-generator-plugin/
API_PLUGIN_PATH="." make build
cp PolicyGenerator "$ROOT"/kustomize/policy.open-cluster-management.io/v1/policygenerator/PolicyGenerator