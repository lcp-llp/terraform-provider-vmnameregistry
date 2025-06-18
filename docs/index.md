# VM Name Registry Provider

The VM Name Registry provider can be used to query api for automated name generation and logging for virtual machines

## Example Usage



    terraform {
    required_providers {
        vmnameregistry = {
        version        = "~> 0.0.1"
        source         = "lcp-llp/vmnameregistry"
        }
      }
    }


    provider "vmnameregistry" {
        url = "www.example.com/apo/vmnameregistry"
    }

## Schema

### Required
-   `url` (string) the endpoint of your vmnameregistry