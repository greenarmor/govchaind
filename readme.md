# OpenGovChain

**A tokenless, public good blockchain for government data transparency and accountability.**

---

## 🎯 Mission

OpenGovChain is a decentralized, permissionless, and tokenless blockchain network designed to store and manage government datasets with complete transparency. Our mission is to create an open, immutable, and accessible platform where government data is a public good, freely accessible to all citizens.

## 📚 Documentation

This README provides a general overview. For detailed guides, please refer to our full documentation:

- **[Getting Started](./documentation/GETTING_STARTED.md)**: A comprehensive guide to setting up your development environment and running a node.
- **[Running with Docker](./documentation/DOCKER.md)**: General instructions for using Docker.
- **[Standalone Node Guide](./documentation/RUNNING_WITHOUT_TAILSCALE.md)**: How to run a node without Tailscale, ideal for local development or VPS setups.
- **[Tailscale Guide](./documentation/DOCKER_COMPOSE_TAILSCALE.md)**: How to join the network securely from behind a NAT using Tailscale.
- **[Network Configuration](./documentation/NETWORK_CONFIG.md)**: Details on network settings, peers, and ports.
- **[Technical Implementation](./documentation/TECHNICAL_IMPLEMENTATION.md)**: A deeper dive into the blockchain's architecture and custom modules.

## 🚀 Getting Started (for Developers)

These instructions are for setting up a local development environment.

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18+)
- [Ignite CLI](https://docs.ignite.com/guide/install)

### 1. Verify Your Environment

Before you begin, run the environment check script to ensure you have all the necessary dependencies and compatible versions installed.

```bash
./scripts/check-build-env.sh
```

### 2. Quick Setup

1.  **Set up the environment:**
    ```bash
    ./scripts/setup-env.sh
    ```

2.  **Build the binary:**
    ```bash
    ignite chain build
    ```

3.  **Start the local blockchain:**
    ```bash
    ignite chain serve
    ```

## 🌐 Joining as a Volunteer Validator

Help secure the network by running a validator node. The entrypoint script now automatically detects your public IP on VPS environments or can be used with Tailscale for home connections.

### 1. Using Docker (Recommended)

Running a node with Docker is the easiest and most maintainable method.

- For nodes on a VPS or with a public IP, see the **[Standalone Node Guide](./documentation/RUNNING_WITHOUT_TAILSCALE.md)**.
- For nodes behind a firewall or on a home network, see the **[Tailscale Guide](./documentation/DOCKER_COMPOSE_TAILSCALE.md)**.

### 2. From Source

For advanced users who prefer to build from source, follow the **[Getting Started](./documentation/GETTING_STARTED.md)** guide and then use the `join-as-volunteer.sh` script:

```bash
./scripts/join-as-volunteer.sh <your-node-name> <genesis-url>
```

## 🌟 Key Features

- **Tokenless Architecture**: No economic barriers to participation. The network is a public good operated by volunteers.
- **Immutable Government Data**: Datasets are stored permanently on the blockchain with IPFS integration for efficient file storage.
- **Rich Metadata & Queries**: Datasets include comprehensive metadata, with capabilities to search by agency, category, and file type.
- **Decentralized & Secure**: Built on the Cosmos SDK and secured by the Tendermint BFT consensus engine, operated by a distributed network of volunteer validators.

## 🤝 Contributing

We welcome contributions from everyone! Whether you are a developer, a data provider, or a citizen advocate, you can help.

1.  **Run a Node**: The most direct way to support the network.
2.  **Contribute Code**: Help us build new features and fix bugs.
3.  **Improve Documentation**: Clear documentation is crucial for adoption.
4.  **Spread the Word**: Share our mission with others.

## 📜 License

This project is open source and available under the [MIT License](./LICENSE).

---

**OpenGovChain by BetterGov.ph**: Empowering transparency through decentralized government data.
