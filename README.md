# GobbleGo Backend

## Overview
GobbleGo revolutionizes restaurant dining by integrating a collaborative ordering system that enhances group dining experiences. Each table is uniquely identified with a QR code, ensuring easy access to a shared digital ordering platform. Guests can view the menu, place individual orders, and organize checkout and billing in a streamlined manner.

## Product Flow

### 1. QR Code Scanning
- Each table has a unique QR code.
- Guests scan the QR code to join a virtual room linked to their table.

### 2. Assigning Leadership
- The first person to scan the QR code is designated as the Table Leader.
- The Table Leader has exclusive authority to finalize the order and proceed to checkout.

### 3. Collaborative Menu Access
- Guests can:
  - View the complete menu, including tags (e.g., Appetizers, Drinks, Main Course, Desserts).
  - Place individual items in their personal cart.

### 4. Checkout Management
- When all members are ready, the Table Leader reviews the combined cart and clicks Checkout.
- Priority of serving is determined based on item tags.
- Multiple checkouts are allowed during a session.

### 5. Billing and Payment
- Once dining concludes, all checkout orders are summed up for a final bill.

### Tech Stack
- **GoFiber** - High-performance web framework in Golang
- **PostgreSQL** - Relational database for storing menu, orders, and user data
- **GORM** - ORM for database management in Golang

## Installation and Setup

### Prerequisites
- Golang installed on your machine
- PostgreSQL database setup

### Steps to Run Locally
1. Fork the repository from GitHub.
2. Clone your forked repository:
   ```sh
   git clone https://github.com/yourusername/gobblego-backend.git
   cd GobbleGo
   ```
3. Install dependencies:
   ```sh
   go mod tidy
   ```
4. Set up environment variables for database configuration.
- Copy `.env.example` to `.env`:
   ```sh
   cp .env.example .env
   ```
- Update `.env` with your PostgreSQL credentials:
   ```sh
   DB_HOST=<your_host>
   DB_PORT=<your_port>
   DB_USER=<your_username>
   DB_PASSWORD=<your_password>
   DB_NAME=<your_database>
   ```
5. Run the application:
   ```sh
   go run main.go
   ```
## Connecting to Backend
Ensure your frontend is running by following the setup steps in the **[GobbleGo Frontend Repository](https://github.com/saanvi-iyer/gobblego/frontend)**.
---
<p align="center">Made with ❤️ by Saanvi Iyer</p>
