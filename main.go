package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ------------ REQUEST / RESPONSE STRUCTS --------------

type CodeRequest struct {
	Code string `json:"code"`
}

type HelpRequest struct {
	Query string `json:"query"`
}

type RunResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type AutoFixResponse struct {
	FixedCode string `json:"fixedCode"`
}

type HelpResponse struct {
	Answer string `json:"answer"`
}

// -------------------- MAIN ----------------------------

func main() {
	app := fiber.New()

	// Ye line: ./public folder ko root URL (/) par serve karega
	// Matlab public/index.html -> http://localhost:3000
	app.Static("/", "./public")

	// APIs
	app.Post("/run", handleRun)
	app.Post("/autofix", handleAutoFix)
	app.Post("/help", handleHelp)

	log.Println("Server running on http://localhost:3000")

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

// --------------- /run HANDLER -------------------------

func handleRun(c *fiber.Ctx) error {
	var req CodeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(RunResponse{
			Error: "invalid request body",
		})
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(RunResponse{
			Error: "code is empty",
		})
	}

	var outputLines []string
	lines := strings.Split(code, "\n")

	for _, line := range lines {
		l := strings.TrimSpace(line)
		low := strings.ToLower(l)

		// Agar "error" word mila to fake error
		if strings.Contains(low, "error") {
			return c.Status(fiber.StatusOK).JSON(RunResponse{
				Error: "Simulated runtime error: found the word 'error' in code",
			})
		}

		// println("...") style
		if strings.Contains(low, "println(") {
			start := strings.Index(l, "println(")
			if start != -1 {
				content := l[start+len("println("):]
				content = strings.TrimRight(content, "); ")
				outputLines = append(outputLines, content)
				continue
			}
		}

		// "print " se start
		if strings.HasPrefix(low, "print ") {
			outputLines = append(outputLines, strings.TrimSpace(l[6:]))
			continue
		}
	}

	if len(outputLines) == 0 {
		outputLines = append(outputLines,
			"(Simulated) Code ran successfully. No explicit output detected.",
		)
	}

	return c.Status(fiber.StatusOK).JSON(RunResponse{
		Output: strings.Join(outputLines, "\n"),
	})
}

// --------------- /autofix HANDLER ---------------------

func handleAutoFix(c *fiber.Ctx) error {
	var req CodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	fixed := autoFixCode(req.Code)

	return c.Status(fiber.StatusOK).JSON(AutoFixResponse{
		FixedCode: fixed,
	})
}

// Simple formatter:
// - extra spaces hatao
// - semicolons add karo
// - indentation { } se set karo
// - brackets balance karo
func autoFixCode(code string) string {
	lines := strings.Split(code, "\n")
	spaceRe := regexp.MustCompile(`\s+`)

	indentLevel := 0
	var fixedLines []string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			fixedLines = append(fixedLines, "")
			continue
		}

		line = spaceRe.ReplaceAllString(line, " ")

		if strings.HasPrefix(line, "}") || strings.HasPrefix(line, "});") {
			if indentLevel > 0 {
				indentLevel--
			}
		}

		if needsSemicolon(line) {
			line = line + ";"
		}

		indent := strings.Repeat("    ", indentLevel)
		fixedLines = append(fixedLines, indent+line)

		openBraces := strings.Count(line, "{")
		closeBraces := strings.Count(line, "}")
		indentLevel += openBraces - closeBraces
		if indentLevel < 0 {
			indentLevel = 0
		}
	}

	result := strings.Join(fixedLines, "\n")
	result = balanceBrackets(result)
	return result
}

func needsSemicolon(line string) bool {
	if line == "" {
		return false
	}

	last := line[len(line)-1]
	if strings.ContainsRune(";{}:(),[]", rune(last)) {
		return false
	}

	keywords := []string{
		"var ", "let ", "const ", "int ", "string ", "float", "return ",
	}
	for _, k := range keywords {
		if strings.HasPrefix(line, k) {
			return true
		}
	}

	if strings.Contains(line, "=") || strings.Contains(line, "(") {
		return true
	}

	return false
}

func balanceBrackets(code string) string {
	openPar := strings.Count(code, "(")
	closePar := strings.Count(code, ")")
	openCurly := strings.Count(code, "{")
	closeCurly := strings.Count(code, "}")
	openSquare := strings.Count(code, "[")
	closeSquare := strings.Count(code, "]")

	var sb strings.Builder
	sb.WriteString(code)

	for i := 0; i < openPar-closePar; i++ {
		sb.WriteString(")")
	}
	for i := 0; i < openCurly-closeCurly; i++ {
		sb.WriteString("}")
	}
	for i := 0; i < openSquare-closeSquare; i++ {
		sb.WriteString("]")
	}
	return sb.String()
}

// --------------- /help HANDLER ------------------------

func handleHelp(c *fiber.Ctx) error {
	var req HelpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	query := strings.ToLower(strings.TrimSpace(req.Query))
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query is empty",
		})
	}

	answer := matchHelpAnswer(query)

	return c.Status(fiber.StatusOK).JSON(HelpResponse{
		Answer: answer,
	})
}

func matchHelpAnswer(query string) string {
	type topic struct {
		keywords []string
		answer   string
	}

	topics := []topic{
		{
			keywords: []string{"for loop", "loop", "iterate"},
			answer:   "Tip: A typical for-loop has init, condition, and increment. Example (C-style): for (int i = 0; i < n; i++) { /* code */ }",
		},
		{
			keywords: []string{"if", "condition"},
			answer:   "Tip: Use if/else to branch logic. Keep conditions simple and readable, and avoid deeply nested ifs when possible.",
		},
		{
			keywords: []string{"function", "method"},
			answer:   "Tip: Functions should do one small task well. Keep them short, name them clearly, and avoid too many parameters.",
		},
		{
			keywords: []string{"variable", "var"},
			answer:   "Tip: Use meaningful variable names; avoid single letters except in short loops. Initialize variables close to where they are used.",
		},
		{
			keywords: []string{"debug", "error", "bug"},
			answer:   "Tip: Add print/log statements, check assumptions, and reduce the problem to the smallest reproducible example.",
		},
		{
			keywords: []string{"indent", "format", "style"},
			answer:   "Tip: Consistent indentation and formatting make code easier to read and debug. Use a formatter where possible.",
		},
	}

	for _, t := range topics {
		for _, k := range t.keywords {
			if strings.Contains(query, k) {
				return t.answer
			}
		}
	}

	return "General tip: Break problems into smaller pieces, write simple functions, and test often. If you share more details, you can get more specific help."
}
