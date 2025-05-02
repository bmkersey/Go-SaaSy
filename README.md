
  # Go-SaaSy

  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

  ## Description

  A production-ready starter kit for building modern SaaS platforms with Go, Postgres, Stripe, and Next.js.

  ## Table of Contents
  * [Features](#features)
  * [Stack](#stack)
  * [Installation](#installation)
  * [Usage](#usage)
  * [License](#license)
  * [Contribution](#contribution)
  * [Tests](#tests)
  * [Questions](#questions)
  
  ## ✨ Features

  - 🔐 Auth (JWT, bcrypt)
  - 🏢 Organizations with owner roles
  - 💳 Stripe subscription billing (Pro, Ultimate tiers)
  - 📬 Email (welcome, reset password, plan changes, invites)
  - 👑 Admin dashboard
  - ⚙️ CLI utilities (create users, reset passwords, promote to admin)
  - 🐳 Dockerized (Go + Postgres + Next.js)
  - 🚀 Frontend: Next.js with Tailwind (JS only)

## 📦 Stack

  - Go (Chi router, SQLC, Goose)
  - PostgreSQL
  - Stripe API
  - Next.js (Pages Router)
  - Tailwind CSS
  - Docker & Makefile

  ## 📂 Directory Structure

  ```
  cmd/              # CLI commands (e.g., admin tasks)
  frontend/         # Next.js frontend
  internal/         # Go app logic (auth, orgs, billing, email)
  sql/              # SQLC + Goose migrations
  ```

## 🛠️ Installation

  1. **Clone the repo**  
   ```bash
   git clone https://github.com/bmkersey/Go-SaaSy.git && cd Go-SaaSy
   ```

  2. **Copy and edit env vars**  
   ```bash
   cp .env.example .env
   ```

  3. **Run the stack**  
   ```bash
   make up
   ```

  4. **Create an initial user**  
   ```bash
   make create-user EMAIL=admin@saasy.app PASSWORD=supersecure
   ```

  ## 🚀 Usage

  - Access the app: [http://localhost:3000](http://localhost:3000)
  - Backend API: [http://localhost:8080](http://localhost:8080)
  - CLI usage:
  ```bash
  make create-user ...
  make reset-password ...
  make promote-user ...
  ```
  ## License

  This project uses the The MIT License.  
  Please visit [https://opensource.org/licenses/MIT](https://opensource.org/licenses/MIT) to learn more.
  

  ## Contribution

  Sole contributor: [bmkersey](https://github.com/bmkersey)  
  Want to fork and build your own SaaS on top? Go for it — MIT license.
  
  ## Tests 

  TBD — coming soon (PRs welcome)
  
  ## Questions
  Questions? Comments? Concerns? Feel free to reach out!  
  Email: bmkersey@gmail.com  
  GitHub: [bmkersey](https://github.com/bmkersey)  
  
  