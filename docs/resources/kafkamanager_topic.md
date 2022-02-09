# Resource: kafkamanager_topic

Creates a kafka topic in your confluent account.

## Example Usage

### Kafka topic creation

```hcl
resource "kafkamanager_topic" "topic" {
  name = "cox_topic"
  cluster_id = 5
  partitions = 12
  replication_factor = 3
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the kafka topic.
* `cluster_id` - (Required) The ID of your kafka cluser
* `partitions` - (Required) The number of the topic partitions.
* `replication_factor` - (Required) The number of servers will replicate each message.
* `retention_bytes` - (Computed) The maximum size a partition can grow to before discarding old log segments to free up space (default -1)
* `retention_ms` - (Computed) The maximum time a partition can retain a log to before discarding old log segments to free up space (default 86400000)
* `cleanup_policy` - (Computed) The retention policy to use on old log segments "delete" or "compact" (default delete)
* `max_message_bytes` - (Computed) The largest record batch size allowed by Kafka (default 1048588)
