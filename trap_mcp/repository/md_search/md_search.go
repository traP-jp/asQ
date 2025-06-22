package mdsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/traP-jp/h25s_05/trap_mcp/util"
)

func myhash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type NoteInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

var (
	store   util.CacheStore[map[uint32][]NoteInfo]
	innerFn func(context.Context) (map[uint32][]NoteInfo, error) = util.GetWithCache(
		updateCache,
		&store,
		time.Hour,
	)
)

func updateCache(ctx context.Context, cache *map[uint32][]NoteInfo) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://md.trap.jp/notes", nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:  "traP_token",
		Value: os.Getenv("TRAP_TOKEN"),
	})
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	fmt.Println(len(os.Getenv("TRAP_TOKEN")), "characters in TRAP_TOKEN")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("Response status:", resp.Status)
	/*
	 */
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var parsedJson map[string]interface{}
	json.Unmarshal(bodyBytes, &parsedJson)
	*cache = make(map[uint32][]NoteInfo)
	for _, v := range parsedJson["notes"].([]interface{}) {
		fmt.Println("Processing entry:", v)
		id := v.(map[string]interface{})["id"].(string)
		title := []rune(v.(map[string]interface{})["text"].(string))
		if id == "" || len(title) == 0 {
			continue
		}
		for i := 0; i < len(title)-2; i++ {
			substr := title[i : i+3]
			h := myhash(string(substr))
			if _, exists := (*cache)[h]; !exists {
				(*cache)[h] = []NoteInfo{}
			}
			(*cache)[h] = append((*cache)[h], NoteInfo{
				Id:    id,
				Title: string(title),
			})
		}
	}
	fmt.Println("Cache updated with", len(*cache), "hashes")
	return nil
}

func Search(ctx context.Context, keyword string) ([]NoteInfo, error) {
	fmt.Println("Searching for keyword:", keyword)
	dict, err := innerFn(ctx)
	fmt.Println("Dictionary loaded with", len(dict), "entries")
	if err != nil {
		return nil, err
	}
	var ret []NoteInfo
	key_runes := []rune(keyword)
	for i := 0; i < len(key_runes)-2; i++ {
		substr := string(key_runes[i : i+3])
		h := myhash(substr)
		if ids, exists := dict[h]; exists {
			fmt.Println("Found", len(ids), "matches for substring:", substr)
			ret = append(ret, ids...)
		} else {
			fmt.Println("No matches found for substring:", substr)
		}
	}
	var retUnique []NoteInfo
	seen := make(map[string]bool)
	for _, note := range ret {
		if !seen[note.Id] {
			seen[note.Id] = true
			retUnique = append(retUnique, note)
		}
	}
	return retUnique, nil
}
