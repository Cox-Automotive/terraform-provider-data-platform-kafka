# kafkamanager_clusters (Data Source)


## Schema


### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **clusters** (List of Object) (see [below for nested schema](#nestedatt--clusters))

<a id="nestedatt--clusters"></a>
### Nested Schema for `clusters`

Read-Only:

- **availability** (String)
- **bootstrap_servers** (String)
- **cloud** (String)
- **cluster_type** (String)
- **confluent_id** (String)
- **default_s3_bucket_name** (String)
- **environment_id** (String)
- **id** (String)
- **name** (String)
- **networking_type** (String)
- **region** (String)
- **rest_proxy_url** (String)
- **secrets_manager_cloud_secret_arn** (String)
- **secrets_manager_secret_arn** (String)
