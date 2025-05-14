package main

import (
    "crypto/sha256"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "time"
)

func main() {
    fmt.Println("AutoHelm is running!")

    remoteURL := os.Getenv("COMPOSE_URL")
    if remoteURL == "" {
        fmt.Println("Error: COMPOSE_URL environment variable is not set.")
        os.Exit(1)
    }

    localFile := os.Getenv("LOCAL_COMPOSE_FILE")
    if localFile == "" {
        localFile = "docker-compose.yml"
    }

    for {
        waitUntilNext3AM()
        fmt.Printf("\n[%s] Checking for compose updates...\n", time.Now().Format(time.RFC3339))
        err := updateComposeFile(remoteURL, localFile)
        if err != nil {
            fmt.Printf("Update error: %v\n", err)
        }
    }

    fmt.Println("Check complete.")
}

func waitUntilNext3AM() {
    now := time.Now()
    next3am := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
    if now.After(next3am) {
        next3am = next3am.Add(24 * time.Hour)
    }
    waitDuration := time.Until(next3am)
    fmt.Printf("Sleeping until 3AM (%s from now)...\n", waitDuration)
    //time.Sleep(waitDuration)
    time.Sleep(15*time.Second)
}

func updateComposeFile(url string, localFile string) error {
    localFile = "/data/" + localFile
    // Download remote file
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to fetch remote file: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
    }

    remoteContent, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read remote content: %w", err)
    }

    remoteHash := sha256.Sum256(remoteContent)

    // Read local file (if exists)
    localContent, err := os.ReadFile(localFile)
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("failed to read local file: %w", err)
    }

    localHash := sha256.Sum256(localContent)

    if remoteHash != localHash {
        fmt.Println("New version of docker-compose.yml detected. Updating...")
        err = os.WriteFile(localFile, remoteContent, 0644)
        if err != nil {
            if strings.Contains(err.Error(), "permission denied") {
                return fmt.Errorf("write failed: permission denied. If you're using a mounted file, ensure it is writable by UID 65532 (e.g., `chown 65532 %s`)", localFile)
            }
            return fmt.Errorf("failed to write updated file: %w", err)
        }

        fmt.Println("Compose file updated. Restarting containers...")

        if err := restartDockerCompose(localFile); err != nil {
            return fmt.Errorf("failed to restart containers: %w", err)
        }

        fmt.Println("Containers restarted successfully.")
    } else {
        fmt.Println("No changes detected in docker-compose.yml.")
    }

    return nil
}

func restartDockerCompose(composeFile string) error {
    pullCmd := exec.Command("/usr/local/bin/docker", "compose", "-f", composeFile, "pull")
    pullCmd.Stdout = os.Stdout
    pullCmd.Stderr = os.Stderr

    upCmd := exec.Command("/usr/local/bin/docker", "compose", "-f", composeFile, "up", "-d")
    upCmd.Stdout = os.Stdout
    upCmd.Stderr = os.Stderr

    fmt.Println("Running: docker compose pull")
    if err := pullCmd.Run(); err != nil {
        return fmt.Errorf("pull failed: %w", err)
    }

    fmt.Println("Running: docker compose up -d")
    if err := upCmd.Run(); err != nil {
        return fmt.Errorf("up -d failed: %w", err)
    }

    return nil
}
