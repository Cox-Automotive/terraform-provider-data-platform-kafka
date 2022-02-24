# Kafka Manager Terraform Provider

Kafka Manager Provider is used to interact with your Kafka Confluent cluster within kafka manager API.
It allows to import existing Kafka resources as Terraform data sources and create new ones as Terraform resources.



## Example
```hcl
provider "kafkamanager" {
  url = "https://kafka-manager.foo.bar"
  version = "0.0.8"
}
```


## Authentication
Kafka Manager Provider uses multiple ways of providing credentials.
The following methods are supported:

### Static credentials
Static credentials can be provided via `url` and `access_token` in-line in the Kafka Manager provider block

```hcl
provider "kafkamanager" {
  url = "https://edp-microservices.awsedscinp.com/kafka-manager"
  access_token = "fakeoauth2tokenheader.fakeoauth2tokenpayload.fakeoauth2tokensignature"
}
```

### Environment variables
You can provide your credentials via the `KAFKAMANAGER_URL`, `KAFKAMANAGER_ACCESS_TOKEN` environment variables.


## Using the provider

After configuring credintals in your prefered way you can start using this provider to create kakfa resources.

Please note that Kafka Manager Provider only support creating/modifying/deleting topics for now.

### To create a topic:

```hcl
data "kafkamanager_cluster" "dev" {
  cluster_name = "dev-1"
}

resource "kafkamanager_topic" "my_topic" {
  name = "my_topic"
  cluster_id = data.kafkamanager_cluster.dev.id
  partitions = 12
  replication_factor = 3
}
```


## Argument Reference

* `url` - (Required) The URL to Kafka Manager. Also read from ENV.KAFKAMANAGER_URL
* `access_token` - (Required) The access token from Okta. Also read from ENV.KAFKAMANAGER_ACCESS_TOKEN. For more details on how to set up Okta access please contact the email id below.

---
### Supported Versions

| Terraform 0.15.x         |
| ------------------------ |
| Provider version < 0.0.3 |

For questions, please reach out to the Cox Automotive Data Platform team at CoxAutoDataPlatform@coxautoinc.com
