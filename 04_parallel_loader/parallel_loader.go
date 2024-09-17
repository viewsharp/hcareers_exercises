package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

var (
	workersCount     = flag.Uint("workers", uint(runtime.NumCPU()*2), "Number of parallel loader workers")
	errorProbability = flag.Float64("error-probability", 0.01, "Error probability")
)

const (
	correctHost   = "https://jsonplaceholder.typicode.com"
	incorrectHost = "https://rl2v4.wiremockapi.cloud"

	retryCount = 2
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	posts, err := getWithRetry[[]Post](ctx, "/posts")
	if err != nil {
		return fmt.Errorf("get posts: %w", err)
	}

	postChannel := make(chan Post)
	eg, ctx := errgroup.WithContext(ctx)

	commentsByUser := make(map[string]int)
	var commentsByUserMutex sync.Mutex

	for i := 0; i < int(*workersCount); i++ {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case post, ok := <-postChannel:
					if !ok {
						return nil
					}

					url := fmt.Sprintf("/posts/%d/comments", post.Id)
					comments, err := getWithRetry[[]Comment](ctx, url)
					if err != nil {
						return fmt.Errorf("get comments for post id=%d: %w", post.Id, err)
					}

					commentsByUserMutex.Lock()
					for _, comment := range comments {
						commentsByUser[comment.Email] += 1
					}
					commentsByUserMutex.Unlock()
				}
			}
		})
	}

	eg.Go(func() error {
		for _, post := range posts {
			select {
			case postChannel <- post:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		close(postChannel)
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return err
	}

	for email, count := range commentsByUser {
		log.Printf("%s: %d", email, count)
	}

	return nil
}

func getWithRetry[T any](ctx context.Context, path string) (responseData T, err error) {
	for i := 0; i < retryCount; i++ {
		responseData, err = get[T](ctx, path)
		if err == nil {
			return responseData, nil
		}
	}
	return responseData, err
}

func get[T any](ctx context.Context, path string) (responseData T, err error) {
	url := correctHost + path
	if rand.Float64() < *errorProbability {
		url = incorrectHost + path
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return responseData, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return responseData, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("get %s, response status: %s", url, resp.Status)
	if resp.StatusCode != http.StatusOK {
		return responseData, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return responseData, fmt.Errorf("decode response: %w", err)
	}

	return responseData, nil
}
