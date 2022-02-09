# Kafka Manager Terraform Provider
This Terraform provider allows to use Infrastructure as Code (IaC) approach with [Kafka Manager](https://ghe.coxautoinc.com/ETS-EDSS/data-platform-kafka-manager).

It allows to import existing Kafka resources as Terraform [data sources](https://www.terraform.io/docs/language/data-sources/index.html) and create new ones as Terraform [resources](https://www.terraform.io/docs/language/resources/index.html).

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.15.x
- [Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

## Implemented resources and datasources
### Datasources
- `kafkamanager_environment`
- `kafkamanager_environments`
- `kafkamanager_schema_registry`
- `kafkamanager_schema_registries`
- `kafkamanager_cluster`
- `kafkamanager_clusters`
- `kafkamanager_topic`
- `kafkamanager_topics`

### Resources
- `kafkamanager_topic`

## Building the provider
Clone repository
```sh
git clone https://ghe.coxautoinc.com/ETS-EDSS/data-platform-terraform-providers
```

Enter the provider directory and build the provider
```sh
cd ./data-platform-terraform-providers/kafka-manager
make build
```

To install the provider locally, run
```sh
make install
```

## Testing the provider
> **NOTE**: To execute testing operations, you need to be in the provider directory.

### Unit tests
[Unit tests](https://www.terraform.io/plugin/sdkv2/testing/unit-testing) check various parts of the provider in isolation and **do not** provision any infrastructure.

To execute acceptance tests, run
```sh
make test
```

### Acceptance tests
[Acceptance tests](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests) provision **real infrastructure** and check Terraform resources and Kafka Manager resources for consistency.
They require that respective environment variables are set so the Terraform provider and the test module can connect to Kafka Manager.

To execute acceptance tests, run
```sh
export KAFKAMANAGER_URL=https://edp-microservices.awsedscinp.com/kafka-manager
export KAFKAMANAGER_KEY=fakedataplatformkeyfakedataplatformkeyfakedataplatformkeyfakedat
export KAFKAMANAGER_OKTA_GROUPS=EDP_DEV_OPS_DELETE
export KAFKAMANAGER_SUPPLIER=DP
export KAFKAMANAGER_USER_ID=FAKEUSER

make testacc
```

### Sweepers
[Sweepers](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/sweepers) are used to remove dangling infrastructure that may not be deleted due to errors in the provider or tests.
They rely on all resources created through tests to have the `tf_acc_test_` prefix.

To execute sweepers, run
```sh
make sweep
```

## Using the provider
1. Add the `kafkamanager` section to the `required_providers`:
```hcl
terraform {
  required_providers {
    kafkamanager = {
      source = "coxautoinc.com/data-platform/kafkamanager"
      version = "0.0.8"
    }
  }
}
```

2. Configure the provider
    - Through environment variables:
        - Local use:
            ```sh
            export KAFKAMANAGER_URL=https://edp-microservices.awsedscinp.com/kafka-manager
            export KAFKAMANAGER_KEY=fakedataplatformkeyfakedataplatformkeyfakedataplatformkeyfakedat
            export KAFKAMANAGER_OKTA_GROUPS=EDP_DEV_OPS_DELETE
            export KAFKAMANAGER_SUPPLIER=DP
            export KAFKAMANAGER_USER_ID=FAKEUSER
            ```
        - Through the Data Platform Gateway:
            ```sh
            export KAFKAMANAGER_URL=https://edp-microservices.awsedscinp.com/kafka-manager
            export KAFKAMANAGER_ACCESS_TOKEN=fakeoauth2tokenheader.fakeoauth2tokenpayload.fakeoauth2tokensignature
            ```
    - Or directly in the Terraform code:
        - Local use:
            ```hcl
            provider "kafkamanager" {
              url = "https://edp-microservices.awsedscinp.com/kafka-manager"
              key = "fakedataplatformkeyfakedataplatformkeyfakedataplatformkeyfakedat"
              okta_groups = "EDP_DEV_OPS_DELETE"
              supplier = "DP"
              user_id = "FAKEUSER"
            }
            ```
        - Through the Data Platform Gateway:
            ```hcl
            provider "kafkamanager" {
              access_token = "fakeoauth2tokenheader.fakeoauth2tokenpayload.fakeoauth2tokensignature"
            }
            ```
3. Use a resource:
    ```hcl
    resource "kafkamanager_topic" "my_super_topic" {
      name = "my_super_topic"
      cluster_id = data.kafkamanager_cluster.cinp_dp.id
      partitions = 12
      replication_factor = 3
    }
    ```
4. Or data source:
    ```hcl
    data "kafkamanager_topic" "some_topic" {
      name = "some_topic"
      cluster_id = "1"
    }
    data "kafkamanager_topic" "another_topic" {
      id = "1"
      cluster_id = "1"
    }
    ```

See the `examples` directory for more examples.
