# Superwallet
Index blockchain transaction, currently support only Ethereum chain. More chains to come.

## Overview
- **Indexer**: `indexer.go` define the interface for parser's implementation. 
  - **EthIndexer**: `ethindexer.go` is the implementation of the parser interface. It parse the raw transaction fetch by `EthClient` from Geth node, and map to internal domain `Transaction`. 
- **Transaction**: An abstraction of all different transactions type (`RawTransaction`) from different blockchains. All `RawTransaction` types are mapped to `Transaction`.
- **RawTransaction**: `rawtransaction.go` is blockchain specific. Different `RawTransaction` type live in each indexer rpc package.  
- **Storage**: `storage.go` define interface for database operations. Help us to easily switch to any database if we want to, by just implementing the storage interface. 
- **InMemoryStorage**: `inmemorystorage.go` implements the storage interface, interact with the simple `InMemoryDatabase`.
- **InMemoryDatabase**: `inmemorydatabase.go` simple key-value store in memory.  