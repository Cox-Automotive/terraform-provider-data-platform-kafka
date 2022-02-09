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
Static credentials can be provided via an `key`, `user_id`, `supplier` and `okta_groups` in-line in the Kafka Manager provider block

```hcl
provider "kafkamanager" {
  url = "https://kafka-manager.foo.bar"
  key = "fakedataplatformkeyfakedataplatformkeyfakedataplatformkeyfakedat"
  okta_groups = "OKTA_GROUP"
  supplier = "DEV"
  user_id = "FAKEUSER"
}
```

### Environment variables
You can provide your credentials via the `KAFKAMANAGER_URL`, `KAFKAMANAGER_KEY`, `KAFKAMANAGER_SUPPLIER`, `KAFKAMANAGER_USER_ID` and `KAFKAMANAGER_OKTA_GROUPS` environment variables.


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
* `key` - (Required) The access key from loadbalancer. Also read from ENV.KAFKAMANAGER_KEY
* `okta_groups` - (Required) The OKTA group name of your team. Also read from ENV.KAFKAMANAGER_OKTA_GROUPS
* `supplier` - (Required) Your data platform supplier code. Also read from ENV.KAFKAMANAGER_SUPPLIER
* `user_id` - (Required) Your data platform user ID. Also read from ENV.KAFKAMANAGER_USER_ID
