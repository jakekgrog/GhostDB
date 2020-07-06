![Build](https://github.com/GhostDB/GhostDB/workflows/GhostDB%20Node%20Test/badge.svg)

# GhostDB

GhostDB is a distributed, in-memory, general purpose key-value data store that delivers microsecond performance at any scale.

GhostDB is designed to speed up dynamic database or API driven websites by storing data in RAM in order to reduce the number of times an external data source such as a databse or API must be read. GhostDB provides a very large hash table that is distributed across multiple machines and stores large numbers of key-value pairs within the hash table.

## GhostDB Cache Node Installation

To install a GhostDB Cache Node on a server you want to act as a cache for your application, you must first obtain a copy of the GhostDB Cache Node binary. 

You can also compile from source. To do this run the following:

```
> go build
```

GhostDB uses port 7991 so be sure to allow communication on that port for any servers GhostDB Cache Node is running on.
Once obtained, you must create a configuration file for your cache in the same directory as the cache binary.

### Cache Cluster Configuration

The configuration file is a json file and must be called “ghost_conf.json”.

There are multiple configuration options available to you. To start you must specify the keyspace size of the cache. This is the maximum number of key/value pairs you want to allow your cache to store. In the configuration file you set this value as follows:

```
“keyspaceSize”: 65536
```
If this value is not specified, the default is the same as the above.

The next configuration option you have available to you is the default time-to-live for key/value pairs. This is the number of seconds a key/value pair is considered “not-stale”. If, when using the SDK and you do not set a time-to-live on a key/value pair when storing, the default will be set to “-1”. This means the item does not expire. Otherwise, the value you provide will override the value in the configuration file. You set this value as follows in the configuration:
    “defaultTTL”: -1
The next configuration option is the crawler interval configuration. This is an integer value that represents how often the cache crawlers should run. The cache crawlers are concurrent crawlers that remove expired items from the cache. By default this is set to 300 seconds (5 minutes). This is set as follows in the cache:

```
“crawlerInterval”: 300
```

The next configuration option is the snitch metric interval configuration. This is an integer value that represents how often the snitch metrics sub-system should log system metrics in seconds. This is set to 300 by default (5 minutes). This is set as follows in the cache:

```
“snitchMetricInterval”: 300
```

The next configuration option is the watchdog metric interval configuration. This is an integer value that represents how often the watchdog metrics system should log application metrics in seconds. This is set to 300 by default (5 minutes). This is set as follows in the configuration file:
```
“watchdogMetricInterval”: 300
```

The next configuration option enables snapshotting within the cache. This is enabled by default and is set in the configuration file as follows:

```
“snapshotEnabled”: true
```

The next configuration option determines how often snapshotting should occur in seconds. By default this is 3600 seconds (1 hour). This is set as follows in the configuration file:

```
“snapshotInterval”: 3600
```

The next configuration option determines if your snapshots should be encrypted. This is set to true by default and is set in the configuration file as follows:

```
“enableEncryption”: true
```

In order to enable encryption you must also supply a passphrase for encrypting and decrypting the snapshots. This is set to “SUPPLY_ME” by default and should be updated. It is set as follows in the configuration:

```
“passphrase”: “SUPPLY_ME”
```

The next option available to you determines if append-only-file persistence should be used. By default this is set to false. This is set in the configuration file as follows: 

```
“persistenceAOF”: false
```

The next configuration option available determines the maximum byte size of the append-only-file can grow to before it is rotated. By default this is set to 5000000 (5MB). This is set in the configuration as follows:

```
“aofMaxByteSize”: 5000000
```

If both snapshots and append-only-file are enabled in the cache, snapshots will take precedence over append-only-file.
All snapshots are also compressed using gzip. This is not a configurable option. 
These are all the configuration options available. Below is an example of a complete configuration file:

```
{
    "keyspaceSize": 65536,
    "snitchMetricInterval": 300,
    "watchdogMetricInterval": 300,
    "defaultTTL": -1,
    "crawlerInterval": 300,
    "snapshotInterval": 3600,
    "snapshotEnabled": true,
    "persistenceAOF": false,
    "aofMaxByteSize": 50000000,
    "entryTimestamp": true,
    "enableEncryption": true,
    "passphrase": "SUPPLY_ME"
}
```

## Authors
* [Jake Grogan](https://github.com/jakekgrog)
* [Connor Mulready](https://github.com/nohclu)
