# Code Assistant – Go Fiber Backend

A simple backend service built using **Go Fiber**.  
It includes APIs to:

- Run code (simulated)
- Auto-fix formatting issues
- Provide help suggestions
- A small HTML UI to test everything

---

## How to Run

```
go run .
```

Server starts at:

```
http://localhost:3000
```

Visit this URL to open the HTML testing page.

---

## API Endpoints

### 1. POST /run  
Simulates code execution and returns output.

**Request Body (text):**
```
println("Hello Sanskar")
```

**Response:**
```json
{
  "output": "\"Hello Sanskar\""
}
```

---

### 2. POST /autofix  
Fixes common formatting issues:

- Missing semicolons  
- Braces formatting  
- Indentation  
- Extra spaces  

**Input:**
```
int x = 5
if(x > 3) {
console.log("hi")
}
```

**Output:**
```
int x = 5;
if (x > 3) {
    console.log("hi");
}
```

---

### 3. POST /help  
Returns predefined suggestions.

**Input:**
```
how to write loop?
```

**Output:**
```
Tip: Use a for-loop like: for(int i = 0; i < n; i++) { }
```

---

## Testing with Postman

### Run Code
```
POST http://localhost:3000/run
Body → raw → (your code)
```

### Auto Fix
```
POST http://localhost:3000/autofix
Body → raw → (your code)
```

### Help
```
POST http://localhost:3000/help
Body → raw → (your query)
```

---

## Project Structure

```
code-assistant-fiber/
│── main.go
│── go.mod
│── go.sum
│── public/
│     └── index.html
```

---

## Contact

Sanskar Khare  
Email: sanskarkhare3128@gmail.com  
GitHub: https://github.com/Sanskarkhare31

