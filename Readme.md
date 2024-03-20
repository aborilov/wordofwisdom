# Word of Wisdom TCP Server with Proof of Work

This project implements a simple TCP server that provides quotes from a "Word of Wisdom" collection. To protect against unauthorized or excessive requests, it incorporates a Proof of Work (PoW) mechanism. The server is containerized using Docker, and the setup includes a client application also protected with a PoW challenge solution.

### Running the Project

Build and start the server and client services:

```
docker-compose up --build
```

To scale the number of client instances, use the --scale option:

```
docker-compose up --scale client=<number-of-instances>
```

### Proof of Work Algorithm

The project utilizes a simple Hashcash-like PoW algorithm. The server issues a challenge to the client, requiring the client to compute a hash that meets specific criteria (e.g., a hash with a predefined number of leading zeros) before it can receive a quote. This mechanism helps protect the server against DDoS attacks by ensuring that each client expends computational effort to solve a challenge, thus limiting the rate at which requests can be made.

### Choice of PoW Algorithm
The Hashcash-like algorithm was chosen for its simplicity and effectiveness in deterring spam and abuse. It requires the client to perform a computation that is moderately hard (but feasible) to solve, yet easy for the server to verify. This asymmetry ensures that clients cannot make excessive requests without incurring a computational cost, protecting the server's resources.

The complexity of the challenge (number of leading zeros required in the hash) can be adjusted through an environment variable, allowing the protection level to be scaled based on the threat level or server capacity.
