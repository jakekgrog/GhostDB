![GhostDB logo](https://imgur.com/ZEGFVo6.png)

![Build](https://github.com/GhostDB/GhostDB/workflows/GhostDB%20Node%20Test/badge.svg)
[![Twitter](https://img.shields.io/badge/ghostdb-v2.0.0-green)](http://www.ghostdbcache.com)

[![Discord](https://img.shields.io/badge/chat-Join%20us!-green?style=for-the-badge&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2)](https://discord.gg/ZkT5Sdf)



## :books: Overview

GhostDB is a distributed, in-memory, general purpose key-value data store that delivers microsecond performance at any scale.

GhostDB is designed to speed up dynamic database or API driven websites by storing data in RAM in order to reduce the number of times an external data source such as a database or API must be read. GhostDB provides a very large hash table that is distributed across multiple machines and stores large numbers of key-value pairs within the hash table.

## :car: Roadmap

> GhostDB was a university project - it is not fully featured but we're getting there!

This is a high-level roadmap of what we want GhostDB to become by the end of 2020. If you have any feature requests please create one from the [template](https://github.com/jakekgrog/GhostDB/blob/master/docs/FEATURE_REQUEST.md) and label it as `feature request`!

- First hand support for list, set, stack and queue data structures (among others)
- Transactions
- Batch read/write
- Subscribable streams
- Monitoring & administration dashboard
- Enhanced security features
- Transition to TCP sockets as transport protocol
- CLI
- Support for a wide range of programming languages

## :wrench: Installation

To install GhostDB please consult the [installation guide](https://github.com/jakekgrog/GhostDB/blob/master/docs/INSTALL.md) for a quick walkthrough on setting up the system.

## :hammer: Cluster Configuration

To configure a GhostDB cluster please follow the instructions in the [configuration guide](https://github.com/jakekgrog/GhostDB/blob/master/docs/INSTALL.md)

## :wave: Interacting with your Cluster

### Using our SDKs

The way you interact with your GhostDB cluster is through our SDKs. We currently have two SDKs available. Follow the SDK specific guides on the respective SDK repositories for installation instructions:

- [GhostDB SDK for Python](https://github.com/jakekgrog/GhostDB-SDK-Python)
- [GhostDB SDK for Javascript](https://github.com/jakekgrog/ghostdb-sdk-js)

We also have a number of SDKs in active development and will be available soon:

- [GhostDB SDK for Java](https://github.com/jakekgrog/GhostDB-SDK-Java)
- [GhostDB SDK for Golang](https://github.com/jakekgrog/ghostdb-sdk-golang)

We are also looking for developers who would like to contribute to GhostDB by developing SDKs for languages they would like to see an SDK for!

### Using our CLI

We are currently developing a CLI for GhostDB which will provide you with the ability to interact with your cluster using the same commands you would with our SDKs, the ability to perform administrative actions on your cluster, and get more in-depth information on your clusters state. 

## :green_book: FAQ

**What is GhostDB?**

GhostDB is a distributed, in-memory, general purpose key-value data store that delivers microsecond performance at any scale.

GhostDB is designed to speed up dynamic database or API driven websites by storing data in RAM in order to reduce the number of times an external data source such as a database or API must be read. GhostDB provides a very large hash table that is distributed across multiple machines and stores large numbers of key-value pairs within the hash table.

**How does GhostDB compare to other systems like Redis?**

In it's current state, GhostDB isn't all that different but GhostDB is still in very early development and we have a tonne of features to add that you can checkout in our Roadmap.

If you have any features you'd love to see in GhostDB then open a feature request using the [template](https://github.com/jakekgrog/GhostDB/blob/master/docs/FEATURE_REQUEST.md)!

**What are some of GhostDB's use cases?**

- In-memory data lookup
- Relational and Non-relational database speedup
- Managing spikes in web/mobile apps
- Session-store
- Token caching
- Gaming - Player profiles & leaderboards
- Web page caching
- Global ID or counter generation
- Fast access to any suitable data

## :pencil2: Authors

**Jake Grogan**

- Email: <jake.kgrogan@gmail.com>
- Github: [@jakekgrog](https://github.com/jakekgrog)

**Connor Mulready**

- Github: [@nohclu](https://github.com/nohclu)

## :star: Show your support

Give a :star: if this project helped you!
