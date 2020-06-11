# RedHat IDM Provider 

Summary of what the provider is for, including use cases and links to
app/service documentation.

## Example Usage

```hcl
provider "idm" {

}
```

## For terraform.tf

```hcl
terraform {
  required_providers {
    idm = {
      source = "DTherHtun/idm"
      version = "0.0.2"
    }
  }
}
```
