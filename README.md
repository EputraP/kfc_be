
# 📌Golang JWT (Access & Refresh Token)

## 📖 Overview
a golang JWT Auth

## 🚀 Features
- ✅ Sign Up
- ✅ Login
- ✅ Refresh

## 🛠️ Tech Stack  
This project is built using the following technologies:  
- **Go (Golang)**  
- **Gin**
- **GORM**
- **PostgreSQL**
- **JWT**


## 📌 API Endpoints  
| Method | Endpoint         | Description            |
|--------|----------------|------------------------|
| **POST**    | `/auth/register`    | Create user account     |
| **POST**   | `/auth/login`     | Login user   |
| **GET**   | `/auth/refresh`   | Renew access token      |
| **GET**    | `/auth/logout` | Logout |


## 📦 Installation
1. Clone the repository:  
   ```sh
   git clone https://github.com/EputraP/kfc_be_Golang_JWT.git

   go mod tidy
   ```
2. Copy and configure the environment file:
- This project includes a .env.example file. You need to copy and rename it to .env:
   ```sh
   cp .env.example .env
   ```
- Then, open .env and fill in the required environment variables:
   ```sh
   DB_HOST="localhost"
   DB_PORT="5435"
   DB_USER="postgres"
   DB_PASS="postgres"
   DB_NAME="postgres"
   TIMEZONE="Asia/Jakarta"
   ```
3. Start dependencies using Docker Compose:  
   ```sh
    docker compose up --build
   ```

## 📜 License  
This project is licensed under the **[MIT](https://choosealicense.com/licenses/mit/)**, which allows commercial and personal use, modification, and distribution.  