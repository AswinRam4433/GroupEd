# MOOC Website README

## Overview

This repository contains the code for a MOOC (Massive Open Online Course) website that facilitates students studying the same courses to connect with each other and manage their schedules. The website includes features such as user registration, sign-in, OTP (One-Time Password) verification, and a schedule display.

## Features

1. **User Registration:**
   - Collects user information such as name, email, Discord ID, phone number, college, gender, and password.
   - Validates the password to ensure it meets certain criteria.

2. **RSA Key Pair Generation:**
   - Generates an RSA key pair for encrypting and decrypting sensitive data.

3. **Twilio Integration:**
   - Utilizes Twilio API for OTP verification via SMS.

4. **MongoDB Integration:**
   - Stores user registration data and exam schedule information in a MongoDB database.

5. **Password Security:**
   - Hashes and stores passwords securely using SHA-256.

6. **Schedule Display:**
   - Retrieves exam schedule information from MongoDB and displays it on the user's next page.

7. **Sign-In:**
   - Allows users to sign in using their Discord ID and password.

8. **AES Encryption (Optional):**
   - Provides optional AES encryption for sensitive data.

## Dependencies

Ensure you have the following dependencies installed:

- Go programming language
- MongoDB
- Twilio Account (for OTP verification)
- Environment variables: TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, VERIFY_SERVICE_SID, and MongoDB connection details.

## Setup

1. Clone the repository:


2. Create a `.env` file with the following environment variables:

   ```env
   TWILIO_ACCOUNT_SID=your_twilio_account_sid
   TWILIO_AUTH_TOKEN=your_twilio_auth_token
   VERIFY_SERVICE_SID=your_verify_service_sid
   ```

3. Ensure MongoDB is running, and update MongoDB connection details in the code if necessary.

4. Run the application:

   ```bash
   go run sign_up.go
   ```

5. Access the website at [http://localhost:8080](http://localhost:8080).

## Usage

1. Navigate to the registration page and fill in the required information.
2. After successful registration, you will be redirected to the OTP verification page. Enter the OTP received via SMS.
3. Upon successful OTP verification, you will be redirected to your next page, displaying your schedule.
4. Use the sign-in page to log in with your Discord ID and password.

## Contributing

Feel free to contribute to the project by opening issues, suggesting improvements, or submitting pull requests.
