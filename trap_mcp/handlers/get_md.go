package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetMdTool() mcp.Tool {
	tool := mcp.NewTool("GetMarkdown",
		mcp.WithDescription("Get markdown documents by id"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the markdown document to retrieve available in md.trap.jp/<id> or returns of SearchMarkdown tool"),
		),
	)
	return tool
}

type AAAResponse struct {
	Revisions []Revision `json:"revision"`
}

type Revision struct {
	UnixTime int64 `json:"time"`
	Length   int   `json:"length"`
}

func GetMdHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	revUrl := fmt.Sprintf("https://md.trap.jp/%s/revision", id)

	req, err := http.NewRequestWithContext(ctx, "GET", revUrl, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.AddCookie(&http.Cookie{
		Name:  "traP_token",
		Value: os.Getenv("TRAP_TOKEN"),
	})
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected status: %s", resp.Status)), nil
	}
	defer resp.Body.Close()
	var respJson AAAResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	err = json.Unmarshal(bodyBytes, &respJson)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	revs := respJson.Revisions
	newest := revs[0]
	fmt.Println("Newest revision:", newest)
	contUrl := fmt.Sprintf("https://md.trap.jp/%s/revision/%d", id, newest.UnixTime)
	fmt.Println("Content URL:", contUrl)
	req, err = http.NewRequestWithContext(ctx, "GET", contUrl, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req.AddCookie(&http.Cookie{
		Name:  "traP_token",
		Value: os.Getenv("TRAP_TOKEN"),
	})
	jar, _ = cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}
	resp, err = client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected status: %s", resp.Status)), nil
	}
	defer resp.Body.Close()
	var conts map[string]interface{}
	bodyBytes, err = io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &conts)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fmt.Println("Content response:", string(bodyBytes))
	content := conts["content"]
	if content == nil {
		return mcp.NewToolResultError("content not found in response"), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", content)), nil
}
