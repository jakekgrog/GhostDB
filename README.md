![GhostDB logo](https://imgur.com/ZEGFVo6.png)

![Build](https://github.com/GhostDB/GhostDB/workflows/GhostDB%20Node%20Test/badge.svg)
[![Twitter](https://img.shields.io/badge/ghostdb-v2.0.0-green)](http://www.ghostdbcache.com)

[![Discord](https://img.shields.io/badge/chat-Join%20us!-green?style=for-the-badge&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2)](https://discord.gg/ZkT5Sdf)

## Update - 03/03/2021

### Where we've been

GhostDB stemmed from a University project. Due to the nature of these projects (time constraints etc.), we feel some corners were cut. For example, we opted for the memcached model of distribution to save on time as it was easier to implement. However, this wasn't the original vision of GhostDB. Myself and Connor also started new jobs and these took up a good chunk of our time. This combined with just finishing a really busy final year in Univeristy, we decided to mothball the project for a while. We're finally returning to it and hopefully transforming it into what we had originally planned. 

### A new roadmap

We are revising our roadmap below and plan to release an updated version soon but before we do here is a brief rundown on what we want

- Transition away from the memcached model and move to a consistent, partition tolerant system (with limited fault tolerance too) by implementing the raft concensus protocol. (This is almost complete)
- Release a CLI to allow users to easily manage their clusters
- Re-build our SDKs from the ground up to allow users to interact with GhostDB with more ease than is currently possible.
- Implement new data types to broaden GhostDBs use cases.
- Local caching to give an even greater performance boost to users.
- Release AWS Amazon Machine Images (AMIs) and Google Compute Engine Images to allow users to easily create GhostDB clusters in the cloud with only a few clicks.
- Updates to the website that include a download centre and documentation improvements.

### Contributing

Unfortunately, with work and life we simply don't have the time at the moment to manage pull requests from anyone else. However, we are still accepting issues and are encouraging them.

And of course, we also want to continue improving on our performance :)

## :books: Overview

GhostDB is a distributed, in-memory, general purpose key-value data store that delivers microsecond performance at any scale.

GhostDB is designed to speed up dynamic database or API driven websites by storing data in RAM in order to reduce the number of times an external data source such as a database or API must be read. GhostDB provides a very large hash table that is distributed across multiple machines and stores large numbers of key-value pairs within the hash table.

## :car: Roadmap

> GhostDB was a university project - it is not fully featured but we're getting there!

This is a high-level roadmap of what we want GhostDB to become by the end of 2020. If you have any feature requests please create one from the [template](https://github.com/jakekgrog/GhostDB/blob/master/docs/FEATURE_REQUEST.md) and label it as `feature request`!

- First hand support for list, set, stack and queue data structures
- Atomic command queues
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

## :pencil2: Authors

**Jake Grogan**

- Email: <jake.kgrogan@gmail.com>
- Github: [@jakekgrog](https://github.com/jakekgrog)

**Connor Mulready**

- Github: [@nohclu](https://github.com/nohclu)

## :star: Show your support

Give a :star: if this project helped you!
