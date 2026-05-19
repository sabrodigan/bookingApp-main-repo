# Booking Application User Guide

Welcome to the Bookings and Reservations Application! This guide will help you understand how to start the application, navigate its features, and make reservations.

## 1. Overview
This is a modern web application designed for a Bed & Breakfast or hotel. It allows users to view different rooms, check availability for specific dates, and book reservations. 

## 2. Getting Started

### Prerequisites
- Go version 1.15 or higher must be installed on your system.

### Starting the Application
You can start the web server in two ways:

**Option 1: Using the start script (Linux/macOS)**
1. Open your terminal in the root directory of the project.
2. Run the provided startup script:
   ```bash
   ./run.sh
   ```

**Option 2: Using the Go CLI**
1. Open your terminal in the root directory.
2. Run the following command:
   ```bash
   go run ./cmd/web
   ```

Once started, the application will be available in your web browser at:
**http://localhost:8080**

---

## 3. Navigating the Application

The application consists of several main sections, accessible via the navigation menu:

### Home
The landing page of the application, welcoming guests to the Bed & Breakfast.

### About
Learn more about the establishment and what makes it special.

### Rooms
You can view details and photos for specific rooms:
- **General's Quarters**: A detailed view of our premier suite.
- **Major's Suite**: A detailed view of our secondary suite.

### Search Availability
This feature allows you to select a starting and ending date to check if rooms are available for booking during that period. 

### Contact
A page providing contact information for the establishment.

---

## 4. Making a Reservation

1. Navigate to the **Search Availability** page.
2. Select your desired arrival and departure dates.
3. If dates are available, you can proceed to the **Make a Reservation** page.
4. Fill out the required personal information:
   - First Name (minimum 3 characters)
   - Last Name
   - Email Address (must be a valid email format)
   - Phone Number
5. Submit the form. If everything is correct, you will be redirected to the **Reservation Summary** page, confirming your booking details.

---

## 5. Technical Details (For Developers)
- **Framework & Routing:** Built in Go using the `chi` router.
- **Session Management:** Uses `alexedwards/scs` for secure, server-side session state.
- **Security:** Implements Cross-Site Request Forgery (CSRF) protection using `nosurf`.
- **Templates:** Server-side rendered HTML using standard Go `html/template` caching.
