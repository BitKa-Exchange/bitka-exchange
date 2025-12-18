<div align="center">
  <img src="https://github.com/Bitka-Exchange/bitka-mobile/blob/e91b643231cc31b08c61c78ad1edd1c7d2e31504/assets/logo.png" alt="Bitka Backend Logo" width="120" />
  <h1>Bitka Exchange (Backend Core)</h1>
  <p>
    <strong>Microservices cryptocurrency exchange</strong><br>
    Built with Golang Fiber, Kafka, and PostgreSQL.
  </p>

  [![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
  [![Fiber](https://img.shields.io/badge/Fiber-v2-black?style=flat-square)](https://gofiber.io/)
  [![Architecture](https://img.shields.io/badge/Architecture-Clean%20Arch-orange?style=flat-square)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
  [![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
</div>

---

> [!CAUTION]
> This software is an **educational research project**.
>
> **We are NOT affiliated, associated, authorized, endorsed by, or in any way officially connected with [Bitkub Online Co., Ltd.](https://www.bitkub.com/) or any of its subsidiaries.**
>
> This software is **not intended for production financial use**. No real assets are handled. (For the live demo we'll use testnet tokens only.)

---

## 📖 Overview

**Bitka Exchange** is a proof-of-concept centralized crypto exchange (CEX) designed to handle high-concurrency trading with strict double-entry accounting.

### Key Features
* **Atomic Accounting:** Double-entry ledger system preventing race conditions and negative balances.
* **Event-Driven:** Asynchronous communication via Apache Kafka for scalability.
* **High Performance:** In-memory matching engine (coming soon) and internal communication via gRPC.
* **Secure:** RS256 JWT Authentication with strict role-based access control (RBAC).
* **Clean Architecture:** Strict separation of concerns (Domain, Usecase, Repository, Delivery).

> [!NOTE]
> We also have other repositories that you might want to check them out at [Bitka Exchange GitHub Organization](https://github.com/BitKa-Exchange).

---

## 🏗️ System Architecture

The system is composed of decoupled microservices communicating via **gRPC** for synchronous and **Kafka** for asynchronous. In the end we plan to deploy hybrid using **K3s** but for now it's on docker in on-premises with 1 **PostgreSQL** instance as container (still logically separate database for each service).

<img src="./docs/architecture/diagrams/out/docs/architecture/diagrams/c4/container-diagram/container-diagram.svg" alt="c4 container diagram" width="100%"/>

### 🛠️ Tech Stack

| Component | Technology | Description |
| --- | --- | --- |
| **Language** | Golang 1.24+ | - |
| **Framework** | Fiber v2 | High-performance HTTP web framework |
| **Database** | PostgreSQL 18 | Primary persistence |
| **Logging** | Zerolog | Structured logging |
| **ORM** | GORM | Data access layer |
| **Authentication** | JWT (RS256) | User identity & access control |
| **Internal Communication** | gRPC & Protobuf | Synchronous inter-service communication |
| **Messaging** | Apache Kafka | Event streaming |
| **Gateway** | Traefik | Reverse proxy & load balancing |
| **Containerization** | Docker & Docker Compose | Containerization & local orchestration |
| **CI/CD** | GitHub Actions | Automated testing (WIP) & deployment |
| **API Docs** | OpenAPI 3.0.3 | API Contract & Specifications |
| **Docs UI** | Docusaurus | React framework markdown based docs |

---

## 📂 Project Structure Overview

For more details, see [Bitka online docs](https://docs.bitka.polishstack.com/docs/developer-guide/project-structure) or [Markdown file](./portal/docs/developer-guide/project-structure.md).

```plaintext
bitka/
├── docs/                 # OpenAPI & AsyncAPI (WIP) specifications
├── deploy/               # Deployment configurations
├── pkg/                  # Shared libraries (Logger, JWT, Middleware, Etc.)
├── portal/               # Docusaurus documentation site
├── services/
│   ├── auth/             # Authentication & User Management
│   ├── user/             # User data & Profile Management
│   ├── ledger/           # Wallet, Assets, & Double-Entry Accounting
│   ├── order/            # Order Management System
│   ├── matching/         # Order Matching Engine
│   └── market/           # Market Data Aggregation
├── compose.yml    # Local development stack
├── Makefile              # Build & Run commands
└── go.work               # Go Workspace configuration

```

---

## 🚀 Getting Started

### Prerequisites

* **Docker** & **Docker Compose**
* **Make**
* **Node.js** (optionally, for docs site)

### 1. Clone the Repository

```bash
git clone https://github.com/BitKa-Exchange/bitka-exchange.git
cd bitka-exchange
```

### 2. Set Up Environment Variables
Copy the template `.env.template` to `.env` and adjust any necessary variables.


### 3. Start Local Development Stack

```bash
make docker-dev
```

### 4. Access API Documentation

You can test the APIs using the [interactive docs](http://localhost:3000/docs/openapi) once the services are running.

```bash
cd portal
npm install
cd ..
make docs # at root folder
```

---

## 🧪 Testing (Planned)

We plan to implement unit and integration tests on necessary part for each service. But we will focus on integration/e2e testing to allocate time more on other aspect of the project.

Currently, we are interested in [Testcontainer](https://www.testcontainers.org/) for integration/e2e tests. And [K6](https://k6.io/) for load testing.

---

## 🤝 Contributing

We welcome contributions from students and anyone who want to learn together! Because this project purpose is mainly for creating environment for learning and research in each person's area of interest, and combine every work to build a complete system. 

Right now, we don't have a formal contribution guide yet, only some developer guidelines at [Docs](https://docs.bitka.polishstack.com/docs/developer-guide). If you want to contribute, feel free to open issues or pull requests.

> [!NOTE]
> I also encourage you to reach out to me directly first if you want the experience of team collaboration. - @14two-77