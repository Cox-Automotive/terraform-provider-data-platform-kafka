terraform {
  required_providers {
    kafkamanager = {
      source = "coxautoinc.com/data-platform/kafkamanager"
      version = "0.0.8"
    }
  }
}

provider "kafkamanager" {}

data "kafkamanager_environments" "all" {}
output "environments_all" {
  value = data.kafkamanager_environments.all.environments
}

data "kafkamanager_environment" "cinp_dp" {
  id = "1"
  # name = "data-platform-cinp-dp-environment"
}
output "environment_cinp_dp" {
  value = data.kafkamanager_environment.cinp_dp
}

data "kafkamanager_schema_registries" "all" {}
output "schema_registries_all" {
  value = data.kafkamanager_schema_registries.all.schema_registries
}
data "kafkamanager_schema_registry" "cinp_dp" {
  id = "3"
  # confluent_id = "lsrc-k3xzv"
}
output "schema_registry_cinp_dp" {
  value = data.kafkamanager_schema_registry.cinp_dp
}

data "kafkamanager_clusters" "all" {}
output "clusters_all" {
  value = data.kafkamanager_clusters.all.clusters
}
data "kafkamanager_cluster" "cinp_dp" {
  id = "1"
  # name = "data-platform-cinp-dp-cluster"
}
output "cluster_cinp_dp" {
  value = data.kafkamanager_cluster.cinp_dp
}

data "kafkamanager_topics" "all" {}
output "topics_all" {
  value = data.kafkamanager_topics.all.topics
}
data "kafkamanager_topic" "cinp_dp_connector_test" {
  # id = "5"
  name = "connector_test"
  cluster_id = "1"
}
output "topic_cinp_dp_connector_test" {
  value = data.kafkamanager_topic.cinp_dp_connector_test
}

resource "kafkamanager_topic" "test_topic" {
  name = "terraform_test_for_demo2"
  cluster_id = data.kafkamanager_cluster.cinp_dp.id
  partitions = 12
  replication_factor = 3
}
output "topic_cinp_dp_test_topic" {
  value = resource.kafkamanager_topic.test_topic
}
