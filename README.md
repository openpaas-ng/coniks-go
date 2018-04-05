# The Trusternity Project

## 1. Overview

Public key server is a simple yet effective way of key manage-ment in secure end-to-end communication. 
To ensure the trustworthinessof a public key server, a ***Key Transparency*** solution employs a tamper-evident data structureon the server and a gossiping protocol among clients in order to detectcompromised servers. However, due to lack of incentive and vulnerabilityto malicious clients, a gossiping protocol is hard to implement in practice or too costly to operate.

We present Trusternity, an auditing scheme for Key Transparency server relying on Ethereum blockchain that is easy to implement, inexpensive to operate and resilientto malicious clients. Trusternity is implemented as an extension to [CONIKS](https://coniks.cs.princeton.edu/), a key management system that provides transparency and privacy for end-user public keys.

## 2. Extension over CONIKS-GO

We introduce several extensions to CONIKS-go.

### init

The CONIKS server init now creates an **eth.toml** config file.

```bash
coniksserver init -c
```

A new **eth.toml** is as follows.

```json
geth_rpc_endpoint_url = "127.0.0.1:"
eth_account_address = ""
trusternity_contract_address = ""
```

You must edit this file to add Geth RPC endpoint port (default 8545), an account address of the server wallet and the address of Trusternity contract. In a private network, you can deploy your own Trusternity contract [here](https://github.com/coast-team/trusternity-contract/blob/master/src/trusternity_log.sol)

Here is a sample config file after modified

```json
geth_rpc_endpoint_url = "127.0.0.1:8545"
eth_account_address = "0xcaaC50c6D025A1F7F912f9d120f90796C6Fb30aA"
trusternity_contract_address = "0x012bb0dC4E7ce56440d3AaDC68b6cDB240dC6b57"
```

### Trusternity Server

We can now run a Trusternity server by providing CONIKS server with 2 extra flags:

```bash
coniksserver -e
coniksserver -e -t "path to eth.toml"
```

The first is to run with a default **eth.toml** inside the working directory. The second is to provide the path to the config file.

For detail usage instructions for the CONIKS server, see the documentation in their respective packages: [CONIKS-server](keyserver)

### Trusternity Client

Similar to Trusternity Server, we need to turn on ethereum mode in the client as

```bash
coniksclient -e -t "path to eth.toml"
```

Then we can perform audit for a specific epoch by using a REPL command

```bash
audit $epoch
```

where we can replace $epoch with the corresponding epoch number. The client then download and extract the published **STR** via Geth.

### 3. Test Net

In order to build a private Ethereum Test Network, you can find the instructions [here](https://github.com/ethereum/go-ethereum)