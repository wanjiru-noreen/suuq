# Suuq

Suuq is a small-business management platform that helps shop owners manage their products, stock, debts, creditors, and business records in one place. It also uses AI to analyze business data and generate useful reports for better decision-making.

## Features

* User registration and login
* Business management
* Product management
* Stock tracking and stock updates
* Debtor management
* Creditor management
* Business reports
* AI-powered business data analysis
* Dashboard with key business information

## AI

Suuq uses AI to analyze business data and generate reports and insights.

AI is used for:

* Business performance analysis
* Sales and stock analysis
* Identifying important trends
* Generating business reports

**AI recommendations are not part of the current scope.**

## Tech Stack

### Frontend

* HTML
* CSS
* JavaScript

### Backend

* Go

### Database

* SQLite

### AI

* Google Gemini

### Development

* Git & GitHub
* Docker

## Project Structure

```text
suuq/
├── backend/
│   ├── cmd/
│   ├── database/
│   ├── internal/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   └── tests/
│
├── frontend/
│   ├── css/
│   ├── js/
│   ├── assets/
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   ├── dashboard.html
│   ├── products.html
│   ├── stock.html
│   ├── debtors.html
│   ├── creditors.html
│   ├── reports.html
│   └── ai.html
│
├── docs/
├── docker-compose.yml
└── README.md
```

## Architecture

```text
HTML / CSS / JavaScript
          ↓
       HTTP / JSON
          ↓
       Go Backend
          ↓
    Business Logic
          ↓
        SQLite
          ↓
     Business Data
          ↓
       Gemini AI
          ↓
   Reports & Analysis
```

## Project Scope

Suuq focuses on providing small businesses with a simple digital system for managing everyday business operations and understanding their business data.

The project does **not** include:

* Product price management
* AI-generated business recommendations

## Getting Started

Clone the repository:

```bash
git clone <repository-url>
cd suuq
```

### Backend

Navigate to the backend:

```bash
cd backend
```

Install dependencies:

```bash
go mod download
```

Run the backend:

```bash
go run ./cmd/server
```

### Frontend

The frontend consists of standard HTML, CSS, and JavaScript files.

Open the frontend using a local development server and access the application through your browser.

## Development

The project is developed collaboratively using Git and feature branches.

Example:

```bash
git checkout -b feature/frontend
```

Changes should be tested before being merged into the main branch.

## License

This project is developed for educational and project purposes.
