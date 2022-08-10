---
title: Explore your inventory
slug: inventory
weight: 3
---
{{< toc >}}

## List clients connected to your RPort server

Use `rportcli client list` to browse your inventory. With the `--search` flag you can narrow down the search.
`--search` takes a key value pair and it can be used multiple times. For example:

```shell
rportcli client list --search os_kernel=linux --search name=a*
```

The above search fetches all clients where the key `os_kernel` matches exactly the term `linux` and where the key
`name` starts with `a`.

{{< hint type=important title=Search >}}

* When using `--search` multiple times, the search criteria are combined with a logical `AND`.
A combination with `OR` is not yet supported.
* All searches are case-insensitive.
{{< /hint >}}

The following search keys are supported:

* `id`
* `name`
* `os`
* `os_arch`
* `os_family`
* `os_kernel`
* `os_full_name`
* `os_version`
* `os_virtualization_system`
* `os_virtualization_role`
* `cpu_family`
* `cpu_model`
* `cpu_model_name`
* `cpu_vendor`
* `num_cpus`
* `timezone`
* `hostname`
* `ipv4`
* `ipv6`
* `tags`
* `version`
* `address` `client_auth_id`
* `connection_state`
* `allowed_user_groups`
* `groups`

By default, a maximum of 50 rows are displayed. Use `--limit` to display more. Or use `--offset` to move to
the next page. Examples:

```shell
# Get the first 10 results
go run ./... client list --search name=a* --limit 10
# Get the second 10 results 
go run ./... client list --search name=a* --limit 10 --offset 10
```

Using the json output combined with `jq` you can create flexible reports. For example:

```shell
$ rportcli client list --search os_kernel=windows -o json|jq '.[]|[.name,.connection_state]|@tsv' -r
ITXC    connected
JONATHAN-LITTLE connected
Brooke-Hicks    connected
Testing-Win-SRV-2012R2  connected
DONALD-DAY-WS20 connected
ERIC-GARCIA-WS2 disconnected
ABRAHAM disconnected
Homer   disconnected
BOB-MORGAN-WS20 disconnected
Alan-Fernandez-WS2019   disconnected
Win11-Bonn      connected
ELLIE-DOUGLAS-W connected
BEATRICE        connected
Grandmother-W2012R2     disconnected
```
