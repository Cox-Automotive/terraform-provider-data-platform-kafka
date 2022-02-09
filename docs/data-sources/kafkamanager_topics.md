# kafkamanager_topics (Data Source)


## Schema

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **topics** (List of Object) (see [below for nested schema](#nestedatt--topics))

<a id="nestedatt--topics"></a>
### Nested Schema for `topics`

Read-Only:

- **cleanup_policy** (String)
- **cluster_id** (String)
- **id** (String)
- **max_message_bytes** (Number)
- **name** (String)
- **partitions** (Number)
- **replication_factor** (Number)
- **retention_bytes** (Number)
- **retention_ms** (Number)


