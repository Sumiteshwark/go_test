# Nomad Pack: Common Utilities

This pack is a collection of common utilities for Nomad packs of Covlant.


## Table of Contents
<!-- TOC -->
- [Nomad Pack: Common Utilities](#nomad-pack-common-utilities)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
<!-- TOC -->


## Usage

> [!NOTE]
> The `utilities` Pack is not designed to be run stand-alone.

The `utility` Pack is a collection of common utilities. These utilities are presented as [Template Helpers](https://developer.hashicorp.com/nomad/tutorials/nomad-pack/nomad-pack-writing-packs#write-the-templates) and should be loaded as a _dependency_:

```hcl
# see https://developer.hashicorp.com/nomad/tutorials/nomad-pack/nomad-pack-writing-packs#dependency
dependency "utilities" {
  source = "../utilities"
}
```

The included templates can be used like any other template:

```hcl
  [[ template "util_job_meta" . ]]
```

For an example of how to use the `utilities` Pack, see the [hello_world](../hello_world) Pack.