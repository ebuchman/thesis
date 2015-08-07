
Eth P2P Crypto (a.k.a rlpX) Summary:
------------------
Establishing a secure p2p connection requires the two parties to first authenticate by establishing
a shared key. They can then use that key to encrypt further communications.

------------------

First, some background:

- Peers on the network are identified by a persistent elliptic curve public key
- The initiator of a connection dials a peer for which the public key is known
- *All* messages are encrypted: 
	- the initiator's first message is encrypted using the receiver's persistent public key, and the receiver's first response is encrypted using the initiator's persistent public key.
	- all further communications are encrypted using a symmetric key derrived from ephemeral keys that each peer generates only for this session. the use of ephemeral keys affords us perfect forward secrecy: if the persistent private keys are compromised, all past conversations are still secure.

------------------

The initiator's first message is the "authMsg". It contains proof that the initiator has possession of his persistent private key, as well as a commitment to use a particular ephemeral key pair for this session. This is the authentication step, that allows us to proceed using encryption derrived from ephemeral keys rather than persistent ones.

To create an authMsg, the initiator first performs the following:
(i.1) - generate a random nonce-init
(i.2) - generate an ephemeral key pair.
(i.3) - compute the initial-shared-secret from our persistent-private-key and the peer's peersistent-public-key
(i.4) - compute XOR(nonce-init, initial-shared-secret)
(i.5) - sign the result of (i.4) with the ephemeral private key (i.2)

Now, if the receiver does indeed have their persistent-private-key, they should be able to compute initial-shared-secret from our persistent-public-key. Using initial-shared-secret, nonce-init, and ECRECOVER, they can recover our ephemeral-public-key from the signature in (i.5). Thus we ought to send them:

(i.A) - the signature (i.5)
(i.B) - sha3(ephemeral-public-key) 
(i.C) - our persistent-public-key
(i.D) - nonce-init

(i.B) is used as a checksum so the receiver can ensure the recovered ephemeral pubkey is the one intended by the initiator. Note the ephemeral-public-key is never disclosed, and can only be recovered from the signature by someone who has both the nonce-init and the initial-shared-secret, the latter being only the two parties with the appropriate persistent-key-pairs.

(NOTE: we also include a flag at the end for whether or not we are resusing an old initial-shared-secret)

Now, the authMsg is encrypted via ECIES with the receiver's persistent-public-key and sent to the peer.
ECIES proceeds by first generating a new key pair and creating a shared secret from
that private key and the peer's public key. The generated public key is appended so the peer 
can compute the ECIES shared secret and decrypt the authMsg. Note this is not the same
as the shared secret we signed in the auth message, which is used to delegate encryption from the persistent-keys to the ephemeral-keys.

-----------------

Upon receiving the initializer's message, the receiver uses it's persistent-private-key to decrypt the data 
and recover the authMsg. The authMsg servers as a "request for authentication" from the initializer. It proposes that it, a particular identity (persistent-public-key), should communicate with the receiver (another persistent-public-key) with a commitment to a particular ephemeral key. Notice that the authMsg contains proof that the initializer has possession of both the persistent-private-key (needed to produce the initial-shared-secret) and of the ephemeral-private-key (which is used to sign xor(initial-shared-secret, nonce-init)). The receiver must now respond with an authResponse, which provides proof of ownership of the given persistent-private-key and a newly generated ephemeral-private-key. 

To compose the authResponse, the receiver performs the following:
(r.1) - generate a random nonce-resp
(r.2) - generate an ephemeral key-pair
(r.3) - decrypt the authMsg using persistent-private-key
(r.4) - decode the authMsg to recover the signature, the initializer's ephemeral-public-key checksum and persistent-public-key, and the nonce-init
(r.5) - compute initial-shared-secret using persistent-private-key and the initializer's persistent-public-key
(r.6) - run ECRECOVER on the signature using xor(nonce-init, initial-shared-secret) as the message
(r.7) - take the sha3 of the recovered ephemeral-public-key and validate against the checksum from (r.4)

Only if the receiver has the correct persistent-private-key will he be able to discover
the initializer's ephemeral-public-key through ECRECOVER. 

The receiver then creates an authResponse as follows:
(1) - receiver's ephemeral-public-key (r.2)
(2) - nonce-resp (r.1)

The authResponse is encrypted via ECIES using the initializer's persistent-public-key
