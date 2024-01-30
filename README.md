# PGT2ACM translation tool 

The goal of the pgt2acm tool is to ease the conversion between [Policy Generator Templates](https://github.com/openshift-kni/cnf-features-deploy/tree/master/ztp/policygenerator) and [ACM Generator Templates](https://github.com/open-cluster-management-io/policy-generator-plugin)

## Pre-requisites
An input directory containing PGTs: 
- The path to the PGT manifests must be local, e.g. not in subdirectory, for instance:
`PtpConfigSlave`.yaml and not mysubdirectory/`PtpConfigSlave`.yaml

An output directory should:
- contain a directory containing all the source manifest used by the PGT and named `source-crs`. No subdirectory should be present in this directory.

## Usage
```console
Usage of ./pgt2acm:
  -c string
        the optional comma delimited list of reference source CRs templates
  -g    optionally generates ACM policies for PGT and ACM Gen templates
  -i string
        the PGT input file
  -k string
        the optional list of manifest kinds for which to pre-render patches
  -n string
        the optional ns.yaml file path (default "ns.yaml")
  -o string
        the ACMGen output Directory
  -s string
        the optional schema for all non base CRDs
```

The -g option also requires `PolicyGenerator` and `PolicyGenTemplate` executables in the proper subdirectory as specified by Kustomize plugin API. The kustomize subdirectory of this project provides both executable in the correct relative path. To use them, use the following kustomize environment variable when running pgt2acm:
```
KUSTOMIZE_PLUGIN_HOME=$(pwd)/kustomize 
```
so for instance, a example to generate policies for PGT and ACM Gen would be: 
```
KUSTOMIZE_PLUGIN_HOME=$(pwd)/../../pgt2acm-new/kustomize   ../../pgt2acm-new/pgt2acm -i policygentemplates -o acmgentemplates -s ../../pgt2acm-new/test/newptpconfig-schema.json -k PtpConfig

```

## conversion steps

### Convert simple fields
The `apiVersion`, `kind`, `metadata->name`, `namespace`, and `binding rules` are mapped as follow between the PGT and ACM Gen templates

![image](./docs/images/simple-fields.svg)

### Convert PGT policies to ACM Gen policies

![image](./docs/images/policies.svg)

In PGT, the policies are listed under the sourceFiles field. To each file name in the list is associated a single policy. The same policy can be mapped to different file names

In ACM Gen Templates, the policies are listed under the policies field. The files implementing the policies are listed under the policies->manifests section. See the converted ACM Gen Template policy below. So in short in ACM policies we have a list of policies with a list of manifests. In PGT we have a list of source files that are tagged with a policy name,

The manifest name and path is changed as follows:

`PtpConfigSlave`.yaml --> source-crs/`PtpConfigSlave`-MCP-worker.yaml  

The policy wave is retrieved from the first manifest defined in the policy. All manifest must have the same wave number 

```
     policyAnnotations:
       ran.openshift.io/ztp-deploy-wave: "10"
```
`PolicyDefaults.RemediationAction` is set to inform

### Policy patches

PGT and ACM Gen templates can apply patches to reference manifest using Kustomize merge method. PGT can arbitrarily merge CRDs containing list of objects but not ACM Gen. See the following issue in the [policy-generator-plugin](https://github.com/open-cluster-management-io/policy-generator-plugin) project: https://github.com/open-cluster-management-io/policy-generator-plugin/issues/142
As a workaround, pgt2acm supports passing a open API schema that allows kustomizing CRDs that contain list of objects such as in PtpConfig. In addition to the schema, pgt2acmalso needs a list of CRD kinds that should be pre-kustomized. Internally, pgt2acm will use the schema to Kustomize the reference manifest with the patches contained in the policy. The resulting "Kustomized" manifest is used to create a new patch that will render properly when using [policy-generator-plugin](https://github.com/open-cluster-management-io/policy-generator-plugin) without an open API schema.

* Take this example of PGT for a PTP policy: [policies-pgt.yaml](./docs/examples/policies-pgt.yaml)  
* The same policy converted with prerendering of patches to ACM Gen is:  
[policies-acmGenPreRender.yaml](./docs/examples/policies-acmGenPreRender.yaml)  
* If not using pre-rendering we en up with the following incorrect ACM Gen Template since [policy-generator-plugin](https://github.com/open-cluster-management-io/policy-generator-plugin) cannot Kustomize the PtpConfig CRD:  [policies-acmGen.yaml](./docs/examples/policies-acmGen.yaml)

In the picture below, the lines in yellow were merged from the original patch, the lines in grey are coming untouched from the source manifest.
![image](./docs/images/pre-rendered-patches.svg)

### MCP field

The PGT contains a MCP field indicating whether the it applies to `worker` or `master`" nodes. PGT manifest contain a `$mcp` that must be replaced by a `master` or `worker` string when the PGT manifest is converted to ACM Gen. pgt2acm generates new versions of manifests depending of whether they are used for `master` of `worker` nodes.  
For instance, a PGT declaring a mcp field equal to worker will generate extra manifest files for the manifest containing the `$mcp` string. 
If the `PtpConfigSlave`.yaml manifest contains a `$mcp` keyword, a new manifest is generated for the ACM Gen named source-crs/`PtpConfigSlave`-MCP-worker.yaml. The new path is properly referenced in the ACM Gen template.    

### Miscellaneous $names in manifests

Some manifests contain `$name`, `$namespace`, etc... These keywords are used by PGT but are overwritten by the patch merging mechanism defined by [policy-generator-plugin](https://github.com/open-cluster-management-io/policy-generator-plugin), so they can be ignored.