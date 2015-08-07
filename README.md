# Thesis Outline

What follows is a proposed outline for a thesis, which could be completed by December 2015.

## Intro

### Public keys, distributed consensus, profit
A brief but technical introduction to public-key cryptosystems, distributed consensus (Paxos, PBFT, Raft), and game theory

### Bitcoin (Proof of Work)
Summary of bitcoin's solution to the distributed consensus problem, and brief technical overview of the architecture's design.

### Ethereum (Turing-completeness and a better design)
Summary of how ethereum generalized the bitcoin architecture with turing-complete scripting, and where they made improvements in crypto currency design.

### Tendermint (Security-deposit based Proof-of-Stake)
The tendermint whitepaper. It's current form (written by Jae Kwon) reflects a much older codebase. I propose to re-write it with him,
and include it here.

### Multi-chain and payment channels
A summary of how to make one blockchain aware of another, 
and how to build channels of communication that work, securely and without trust in a third party, "off-chain"





## Technical Reviews

### Merkle state trees
A technical comparison between the merkle state tree implementations in ethereum and tendermint, 
as pertains to their binary encoding algorithms, performance under various read/write environments, and security.

### Encrypted wire protocols
A technical comparison between the Diffie-Helman key-exchange and encryption protocols used in the p2p layer of ethereum and tendermint,
as pertains to their overall structure, performance, and security.

### Distributed Consensus
A technical comparison between the distributed consensus implementations in etcb (ie. Raft) and tendermint,
as pertains to their performance and fault tolerance under various read/write/network environments.

### Eth mining




## Implementations

### Tendemint Codebase
A review of the tendermint codebase

### Permissions 
Details of the permissioning layer added to the tendermint state to better enable the use of blockchains in organizations and government

### Tendermint-tendermint multichain via light-client proofs 
An implementation of two parallel tendermint blockchains which are aware of eachother, including the ability to send tokens from one chain to another and back,
and smart contracts on each chain that can use data stored on the other.

### Micropayment channels between tendermint and ethereum
A native implementation of micro-payment channels in tendermint with a complementary one written to run on ethereum, allowing atomic swaps between tendermint and ethereum tokens


## Proofs

### Tendermint consensus ...

