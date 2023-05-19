# Terraform Provider Metaplane
Welcome to the Terraform provider for Metaplane repository! This provider
allows you to manage your Metaplane resources using the popular infrastructure
as code tool, Terraform.

Metaplane is a robust monitoring platform designed to ensure the health and
integrity of your data. By offering real-time data insights and anomaly
detection, Metaplane is instrumental in data-driven decision making and
maintaining data reliability.

This Terraform provider makes it easier to manage your Metaplane resources
programmatically, allowing you to streamline your workflow and keep your data
environments in line with your infrastructure code.

## Project Structure

```
.
├── docs                                # auto-generated docs on the usage of custom provider
├── examples                            # provides example code of custom provider
├── GNUmakefile                         # script for testing
├── go.mod                              # project dependencies
├── internal
│   ├── api                             # interfacing with the Metaplane API
│   └── provider                        # integration with provider plugin framework
├── main.go                             # Entry point for local development
├── terraform-registry-manifest.json    # Metadata for Terraform registry
└── tools                               # supplementary tools, such as  doc generation
```

## Developing the Provider

Install `go` and `terraform` in your local machine. Create a file called
`.terraformrc` in the home directory, and add the following:
```
provider_installation {
  dev_overrides {
      "klaviyo/metaplane" = "~/go/bin"
  }
  direct {}
}

```
This will allow local provider builds by setting a dev_overrides block in the
configuration file. This block overrides all other configured installation
methods.

[Here](https://learn.hashicorp.com/collections/terraform/providers-plugin-framework)
is a must read for developing terraform provider using terraform plugin
framework.

## Common Commands
All the commands are executed in the root directory.

`go install` Build the provider and put the provider binary in the
`$GOPATH/bin` directory.

`go generate` Generate or update documentation, run.

`go get github.com/author/dependency` Add a new dependency.

`go mod tidy` Clean up module's dependencies, ensuring they accurately reflect
what your code actually uses

`make testacc` Run the full suite of Acceptance tests.

## Release
TBD
