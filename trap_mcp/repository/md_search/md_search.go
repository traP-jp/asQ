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

var (
	store   util.CacheStore[map[uint32][]string]
	innerFn func(context.Context) (map[uint32][]string, error) = util.GetWithCache(
		updateCache,
		&store,
		time.Hour,
	)
)

func updateCache(ctx context.Context, cache *map[uint32][]string) error {
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
	*cache = make(map[uint32][]string)
	for _, v := range parsedJson["notes"].([]interface{}) {
		fmt.Println("Processing entry:", v)
		id := v.(map[string]interface{})["id"].(string)
		title := v.(map[string]interface{})["text"].(string)
		if id == "" || title == "" {
			continue
		}
		for i := 0; i < len(title)-7; i++ {
			substr := title[i : i+8]
			h := myhash(substr)
			if _, exists := (*cache)[h]; !exists {
				(*cache)[h] = []string{}
			}
			(*cache)[h] = append((*cache)[h], id)
		}
	}
	fmt.Println("Cache updated with", len(*cache), "hashes")
	return nil
}

func Search(ctx context.Context, keyword string) ([]string, error) {
	fmt.Println("Searching for keyword:", keyword)
	dict, err := innerFn(ctx)
	fmt.Println("Dictionary loaded with", len(dict), "entries")
	if err != nil {
		return nil, err
	}
	var ret []string
	for i := 0; i < len(keyword)-7; i++ {
		substr := keyword[i : i+8]
		h := myhash(substr)
		if ids, exists := dict[h]; exists {
			fmt.Println("Found", len(ids), "matches for substring:", substr)
			ret = append(ret, ids...)
		} else {
			fmt.Println("No matches found for substring:", substr)
		}
	}
	var retUnique []string
	seen := make(map[string]bool)
	for _, id := range ret {
		if !seen[id] {
			seen[id] = true
			retUnique = append(retUnique, id)
		}
	}
	return retUnique, nil
}
