package main

import (
  "context"
  "fmt"
  "github.com/go-git/go-git/v5"
  "github.com/shurcooL/graphql"
  "golang.org/x/oauth2"
  "os"
)

type Repository struct {
  PullRequest struct {
    Number    int
    Title     string
    State     string
    CreatedAt string
    Author    struct {
      Login string
    }
    Comments struct {
      Nodes []struct {
        Author struct {
          Login string
        }
        Body      string
        CreatedAt string
      }
    } `graphql:"comments(first: 50)"`
    Files struct {
      Nodes []struct {
        Path      string
        Additions int
        Deletions int
      }
    } `graphql:"files(first: 100)"`
  } `graphql:"pullRequest(number: 69)"`
}

func main() {
  cwdGoZeroPath()

  repo, err := git.PlainOpen(".")
  if err != nil {
    fmt.Println("Error opening repository:", err)
    os.Exit(1)
  }

  remote, err := repo.Remote("origin")
  if err != nil {
    fmt.Println("Error getting remote:", err)
    os.Exit(1)
  }

  urls := remote.Config().URLs
  if len(urls) > 0 {
    fmt.Println("GitHub URL:", urls[0])
  } else {
    fmt.Println("No remote URL found.")
  }

  ghToken := os.Getenv("GH_TOKEN")
  if ghToken == "" {
    fmt.Println("GH_TOKEN environment variable is not set.")
    //} else {
    //  fmt.Printf("GH_TOKEN: %s\n", ghToken)
  }

  src := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: ghToken},
  )
  httpClient := oauth2.NewClient(context.Background(), src)

  client := graphql.NewClient("https://api.github.com/graphql", httpClient)

  var query struct {
    Repository Repository `graphql:"repository(owner: \"Goddard-Technologies-LLC\", name: \"ReprocessorAlpha\")"`
  }

  err = client.Query(context.Background(), &query, nil)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  // Process the retrieved data and generate the report
  fmt.Println("Pull Request Number:", query.Repository.PullRequest.Number)
  fmt.Println("Title:", query.Repository.PullRequest.Title)
  fmt.Println("State:", query.Repository.PullRequest.State)
  fmt.Println("Created At:", query.Repository.PullRequest.CreatedAt)
  fmt.Println("Author:", query.Repository.PullRequest.Author.Login)

  fmt.Println("Comments:")
  for _, comment := range query.Repository.PullRequest.Comments.Nodes {
    fmt.Println("  Author:", comment.Author.Login)
    fmt.Println("  Body:", comment.Body)
    fmt.Println("  Created At:", comment.CreatedAt)
    fmt.Println()
  }

  fmt.Println("Files:")
  for _, file := range query.Repository.PullRequest.Files.Nodes {
    fmt.Println("  Path:", file.Path)
    fmt.Println("  Additions:", file.Additions)
    fmt.Println("  Deletions:", file.Deletions)
    fmt.Println()
  }
}

func cwdGoZeroPath() {
  gozeroPath := os.Getenv("GOZERO_PWD")
  if gozeroPath != "" {
    err := os.Chdir(gozeroPath)
    if err != nil {
      fmt.Printf("Failed to change directory: %v\n", err)
      os.Exit(1)
    }
    fmt.Printf("Changed working directory to: %s\n", gozeroPath)
  } else {
    fmt.Println("GOZERO_PWD environment variable is not set.")
  }
}
