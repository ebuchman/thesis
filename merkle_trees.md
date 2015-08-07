On Cryptographic State

#Preface

I recently decided to pivot my masters degree from machine learning to blockchains and to chart my progress in a series of blog posts on the technical matters which will be the focus of that work.

Part I of that series is a review of merkle trees as they relate to blockchains which maintain a global state. 
In particular, we study ethereum's merkle radix tree, and tendermint's merkle AVL tree.

#Intro

One of the core components of a blockchain's design is the cryptographic hash tree in which it stores its state. The purpose of the hash tree is two fold: it must account for all changes to the state, and it must enable the production of compact proofs of statements about substates. By compact we typically mean logarithmic in the size of the state. And by substate we mean something like a user's account balance, or a particular transaction.

There are many examples of hash trees in the wild. Typically they arise in peer-to-peer systems, where communications are less trusted and there is greater demand on the production of light weight proofs of data integrity or the occurence of an event.

Bittorrent is a classic example. Used to be, a torrent file contained a hash of all the torrent's data - once you donwloaded all the data, you could take the hash, and make sure you got the correct result. The protocol was dramatically improved by the introduction of merkle trees. Now, a torrent file may contain a merkle tree root hash, and a downloader can verify each piece of the data as they receive it, rather than having to wait until they get all of it.

#Merkle Trees

A simple merkle tree is a way to obtain a unique figerprint for some set of data, where the fingerprint is built up recursively by hashing pieces of the data together in pairs. Say we have some file 32kb file. We split it into 32 pieces (p1, p2, ..., p32), each 1kb. We take the hash of each piece (H(p1), etc.), leaving us with 32 hashes. Then we take each pair of hashes, concatenate them, and take the hash of that, leaving us with 16 hashes. If we keep taking pairs, concatenating, and hashing, we're eventually left with a single, final hash: the Merkle root hash.

Now we want to prove a particular piece is in the tree, given that we know the root hash. Say it's piece 3. All we need to prove piece 3 (`p3`) is in the dataset is: `H(p4)`, `H( H(p1) || H(p2) )`, and `H( H( H(p5) || H(p6) ) || H( H(p7) || H(p8) ) )`. So instead of needing all eight pieces, we only need ours and three hashes. While maybe not impressive in a dataset of 8, a dataset of 1 million requires only about twenty hashes, while 1 billion requires only thirty. The size of the proof is logarithmic in the number of items in the tree. This is part of what makes bittorrent so reliable. 

Bitcoin maintains its own merkle trees, one for each block, containing all the transactions in that block. To prove a transaction occured, you need the transaction root hash for the block in question (really you need the whole block header), and a logarithmic number of hashes from the merkle tree of transactions in that block.

One of the problems with merkle trees as we've been discussing them is that they are quite static; once the tree is computed, there is no way to add or remove items and yeild a new root hash in logarithmic time. This makes them unsuitable for applications with a changing state.

Bitcoin has no need to modify its merkle trees, since it keeps no global state, and offers no notion of an "account" or its balance. Rather, to determine how many bitcoin you can spend, find every transaction that you've received or sent, and take the difference. For a light client, this creates significant headaches: not only do they have to track the whole history of the chain's block headers to establish the correct chain (like any proof of work blockchain), they have to remember all their transactions and the blocks they occured in. Worse, there's no clean way to prove that a transaction output wasn't already spent.

Ethereum sought to improve on this design by explicitly introducing a different kind of merkle tree to manage the dynamic state, namely a merkle radix tree. In a radix tree, an item's key is its path through the tree. So if I have a binary radix tree and key `01100`, then to get to my item from the root I go left, right, right, left, left. This is different then a classical binary tree, where the keys are sorted in the tree according to their numeric value. In this case, assuming binary numerical values, our key corresponds to `12`, and it's position in the tree depends on the number of other keys greater or less than 12, `and` their order of insertertion.

In what follows, we will review the design of these trees, discuss their relative advantages and disadvantages, provide an overview of their implementation in golang, and benchmark the respective implementations. Enjoy.


# Ethereum's Merkle Radix Tree

## Intro

Quick bit on terminology. A `trie`, also known as a `prefix tree`, 
is the basic tree data structure in which a node's key represents the path to take through the tree to get to the node from the root. 
The word `trie` comes from retrieval. 
One researcher described tries as a `Practical Algorithm To Retrieve Information Coded In Alphanumeric`, 
so they are sometimes called patricia trees. But there is no Patricia.

For example, consider storing the words `car`, `cat`, and `camel` in the tree below:


Notice the inefficiency in storing the suffix of camel with one node for each letter, where each node is the only-child of its parent.
We can save space by instead storing only-children with their parents, so we get:


This kind of trie is known as a `radix tree`. It is a slightly space-optimized trie. 
Furthermore, we can establish the radix, which is the size of the alphabet used by the keys.
In this case, the English alphabet has 26 characters, so our radix is 26. 
A binary radix tree has radix of 2, which minimizes spareseness in the tree (no node has more than one empty child) at the expense of depth.

Common advantages of radix trees are their prefix-matching search capacity (making them very relevant to spell-checking software),
and their lack of collissions when to compared to hash tables. Disadvantages include their lack of balance (one branch of a radix tree can be much deeper than others).

## Implementation

Ethereum uses a hexary radix tree, meaning the alphabet is the sixteen hex characters `0123456789abcdef`.
Of course it's a merkle tree, so nodes in the tree must be referenced by their cryptographic hash, 
leaving the hash of the root node as a cryptographic fingerprint for the entire tree.

The nodes in ethereum's merkle radix tree come in a number of types:

- ValueNode: the standard leaf node, which is just [key, value]
- ShortNode: also a [key, value] pair, but where the value is another node (with total size <= XXX)
- HashNode: also a [key, value] pair, but where the value is the hash of another node (of size > XXX) that must be fetched from the database
- FullNode: a branch node, containing a list of 17 nodes, where each of the first 16 correspond to the possible options for the next letter in the key, and the 17th corresponds to a node whose key actually terminates at the given branch node.

I have given a detailed example on how this works in a [previous blog post](https://easythereentropy.wordpress.com/2014/06/04/understanding-the-ethereum-trie/}

TODO: hex prefix encoding (note https://github.com/ethereum/go-ethereum/pull/1594 and https://github.com/ethereum/go-ethereum/pull/1603)

One caveat to using a radix tree is that they are not balanced by default. 
A determiend attacker could conspire to create a series of accounts (or storage slots) with addresses sharing a common long prefix,
such that their depth in the tree is much greater than anything else (as the chances of matching prefixes of increasing length drop off exponentially).
This would force all validators to expend extra storage space and processing time to fetch or update those accounts, 
and would cost the attacker virtually nothing - 
a potential denial of service attack against the network.

To deal with this threat, all keys are hashed by the SHA3 algorithm before hitting the tree. Since the output of the SHA3 is (supposedly) uniformly distributed, the resulting radix tree will be approximately balanced. To execute the attack, an attacker would now have to expend considerable work to find the partial hash collission - work which would probably be more profitable if it were instead directed at honestly mining the chain.
The downside of this solution is the difficulty it presents to inspecting the tree, since the keys have become obfuscated, its not possible to determine what keys are there and/or what their values are without keeping track externally.


# Tendermint's Merkle AVL Tree

## Intro

An AVL tree, named after its inventors (Georgy Adelson-Velsky and Evgenii Landis), is a self-balancing binary tree.
Thus, its a binary tree with an additional `rotation` operation such that whenever the difference in depth between neighbouring branches 
is greater than 1, the tree is `rotated` in such a way as to eliminate that difference. 
AVL trees thus have a built in defence to the attack described in the previous section, at the expense of longer average write times.
However, the rotation operation is also logarithmic in the number of items in the tree, mitigating the overall effect on write times.

In a typical AVL tree, values may be stored anywhere in the tree (ie. in the leaf nodes or in the inner nodes). 
Tendermint uses a modified form of AVL tree, a so called AVL+ tree, which only stores values in leaf nodes.
The reason for this is to accomodate the production of compact proofs: 
if values are also stored on inner nodes, proofs become much longer, as they have to include the value of parent nodes in addition to the hash of neighbouring branches.

The implementation of the AVL+ tree in the tendermint code base permits one to take immutable snapshots of the tree, such that additional insertions/deletions can be made without affecting previous states. For this reason, the full name of the datastructure is the Immutable AVL+ tree, or IAVL+ tree.

## Implementation



# Benchmarks

In what follows we graph a set of benchmarks comparing the two trees, and provide links to source code from which the data can be reproduced.

For each benchmark, we consider variations on the key size, value size, and the number of items we prefill the tree with.

## Time to a million

Wherein we benchmark the time it takes to fill a tree with one million entries, hash the whole thing, and persist to disk.

## Insert+Remove

Wherein we benchmark the time it takes to insert and remove an item in a tree prefilled with some amount of data.

## Retrieve

Wherein we benchmark the time it takes to retrieve a single entry from a tree prefilled with some amount of data

## Proofs

Wherein we benchmark the time it takes to produce a merkle proof for an item in the tree, the size of the proof, and the time it takes to verify it
(Note the code for merkle proofs in the ethereum trie has not yet been written)











