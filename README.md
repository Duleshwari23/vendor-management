# Vendor Management System

A modern, full-stack application designed to streamline vendor onboarding, asset tracking, and document management. Built for performance and a great user experience, this system provides a centralized dashboard for administrators and a personalized portal for vendors.

## Key Features

- **Modern, Responsive UI**: A completely redesigned frontend built with **React**, **TypeScript**, and **Material-UI v7** for a beautiful and intuitive experience on any device.
- **Role-Based Access Control**: Separate, secure dashboards for **Admins** and **Vendors** with distinct permissions and features.
- **Interactive Dashboards**:
  - **Admin Dashboard**: At-a-glance overview of total vendors, assets, and documents, with quick links to management pages.
  - **Vendor Dashboard**: Clear, card-based view of attendance history.
- **Interactive Modals**: Seamlessly add new vendors and assets without leaving the page, thanks to interactive modals.
- **Card-Based Layouts**: Data is presented in clean, responsive cards for improved readability and a modern look.
- **Real-Time Data Handling**: Frontend is fully connected to a **Go (Golang)** backend for creating and fetching data in real-time.
- **Secure Authentication**: JWT-based authentication ensures that user data is secure.

## Tech Stack

| Area      | Technology                                    |
|-----------|-----------------------------------------------|
| **Frontend**  | React, TypeScript, Material-UI (v7), Axios    |
| **Backend**   | Go (Golang), Gorilla Mux, JWT for Go          |
| **Tooling**   | Vite, Go Modules                              |

## Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites

- Go (version 1.18+)
- Node.js (version 16+)
- npm

### Backend Setup

1.  **Navigate to the backend directory:**
    ```bash
    cd backend
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the server:**
    ```bash
    go run main.go
    ```
    The backend server will start on `http://localhost:8080`.

### Frontend Setup

1.  **Navigate to the frontend directory:**
    ```bash
    cd frontend
    ```

2.  **Install dependencies:**
    ```bash
    npm install
    ```

3.  **Run the development server:**
    ```bash
    npm start
    ```
    The frontend application will be available at `http://localhost:3000`.

## Login Credentials

Use the following credentials to access the different roles in the application.

### Admin Account

-   **Email:** `admin@company.com`
-   **Password:** `admin`

### Vendor Account

-   Please use the **Signup** page to create a new vendor account.

## Project Structure
```
vendor-management/
├── backend/           # Go backend application
│   ├── main.go       # Entry point
│   ├── handlers/     # API handlers
│   ├── models/       # Data models
│   ├── middleware/   # Authentication middleware
│   ├── utils/        # Helper functions
│   └── tests/        # Integration tests
│
└── frontend/         # React frontend application
    ├── src/
    │   ├── components/    # React components
    │   ├── pages/        # Page components
    │   ├── services/     # API services
    │   ├── utils/        # Helper functions
    │   └── context/      # React context
    └── public/           # Static files
```

## Getting Started

### Backend Setup
1. Navigate to the backend directory
2. Run `go mod init vendor-management`
3. Run `go mod tidy`
4. Start the server: `go run main.go`

### Frontend Setup
1. Navigate to the frontend directory
2. Run `npm install`
3. Start the development server: `npm start`

## Testing
- Backend integration tests: `cd backend && go test ./tests -v`
- Frontend unit tests: `cd frontend && npm test`

## Future Enhancements
- Database integration
- Real attendance API integration
- Enhanced reporting features
- Document expiry notifications
- Asset maintenance tracking


going to credate a new project for hackathon, my comapny hires some vendor for some specific time for some project, but they don't have any portal , so i want to create a vendor management app for my comaony, write a readme to have alll the feature, so first login/lgout/signup feature will be there, like there will be 2 role, one is admin and 2nd is user/vendor , admin can be anyone in the componay like hr or admin team, they can register a new vendor into the portal , they can upload thier documents, like pdf files and joining details paper, admin can also add alll the details about asset assigned to the vendor employee and joining date and all the details, also there will be cronjob everyday at 12pm it will run for each employee to update thier attendance data  for the preovious day, this cronjob will basically hit an api(in actual case) but for now this cronjob will hit a function which will give some 5 days attendance data for each employee, take some random attendance data for now, like 3 days present , 2 days absent for some employee for some employtee all present , for some 1/2 days and 4.5 days present, so that  attendance api will return login.logout time and prsent day like 1 or 0.5 day for the day as well of the employee, so that will also be visisble to employee when they login, 
and when user login they can see all the data , they can also check thier attendance , and they can also see asset assigned on them like compouter monitors, 

so this is it, write the readme, whole backend in go, use local data for now in place of db , write whole ui things use react , and create backend and frontend folder separetly and write integration test as well in backend to test it e2e and all features # vendor-management
