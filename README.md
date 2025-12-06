Code Assistant – Go Fiber Backend

A simple backend service built using Go Fiber.
APIs included:

1. POST /run

Simulates code execution.
Input: Raw code
Output: Execution result or simulated messages.

2. POST /autofix

Automatically fixes common formatting issues:

Adds missing semicolons

Fixes indentation

Removes extra spaces

Corrects bracket issues

3. POST /help

Returns predefined help messages based on keyword matching.

How to Run
go run .


Server will start at:

http://localhost:3000

Testing

A simple HTML UI is provided in /public/index.html
Visit:

http://localhost:3000


You can test:

Run Code

Auto Fix

Help API

Endpoints (for Postman)
Run Code
POST http://localhost:3000/run
Body → raw text → (your code)

Auto Fix
POST http://localhost:3000/autofix
Body → raw text → (your code)

Help
POST http://localhost:3000/help
Body → raw text → (your query)
