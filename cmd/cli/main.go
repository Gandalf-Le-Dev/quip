package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
)

type CLI struct {
	File     string        `arg:"" optional:"" help:"File to share"`
	Language string        `short:"l" help:"Language for syntax highlighting"`
	TTL      time.Duration `short:"t" default:"24h" help:"Time to live"`
	Edit     bool          `short:"e" help:"Open editor for text"`
	Server   string        `default:"http://localhost:8080" help:"Server URL"`
}

func (c *CLI) Run() error {
	// Check if we have piped input
	stat, _ := os.Stdin.Stat()
	isPiped := (stat.Mode() & os.ModeCharDevice) == 0

	if c.Edit {
		return c.createPasteWithEditor()
	}

	if isPiped && c.File == "" {
		// Handle piped input as paste
		return c.createPasteFromStdin()
	}

	if c.File != "" {
		// Check if file exists
		if _, err := os.Stat(c.File); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", c.File)
		}
		// Upload file
		return c.uploadFile()
	}

	return fmt.Errorf("no input provided")
}

func (c *CLI) uploadFile() error {
	file, err := os.Open(c.File)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add file
	fw, err := w.CreateFormFile("file", filepath.Base(c.File))
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, file); err != nil {
		return err
	}

	// Add TTL
	if err := w.WriteField("ttl", c.TTL.String()); err != nil {
		return err
	}

	w.Close()

	// Make request
	resp, err := http.Post(c.Server+"/api/file", w.FormDataContentType(), &b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
		Download string `json:"download"`
		View     string `json:"view"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Print results
	fmt.Printf("📤 Uploaded: %s\n", result.Filename)
	fmt.Printf("🔗 Download: curl -J -O %s/api/file/%s\n", c.Server, result.ID)
	fmt.Printf("👀 View: %s%s\n", c.Server, result.View)

	return nil
}

func (c *CLI) createPasteFromStdin() error {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	return c.createPaste(string(content))
}

func (c *CLI) createPasteFromText(text string) error {
	return c.createPaste(text)
}

func (c *CLI) createPaste(content string) error {
	payload := map[string]string{
		"content":  content,
		"language": c.Language,
		"ttl":      c.TTL.String(),
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(c.Server+"/api/paste", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		ID       string `json:"id"`
		Language string `json:"language"`
		Raw      string `json:"raw"`
		View     string `json:"view"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Print results
	fmt.Printf("📋 Created paste\n")
	fmt.Printf("🔗 Raw: curl %s%s\n", c.Server, result.Raw)
	fmt.Printf("👀 View: %s%s\n", c.Server, result.View)

	return nil
}

func (c *CLI) createPasteWithEditor() error {
	// Would open $EDITOR here
	fmt.Println("Editor mode not implemented in this example")
	return nil
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("quip"),
		kong.Description("Simple file sharing and pastebin"),
		kong.UsageOnError(),
	)

	if err := cli.Run(); err != nil {
		ctx.FatalIfErrorf(err)
	}
}
