package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var prefix = "0000"

func getComplexity() error {
	complexityStr := os.Getenv("COMPLEXITY")
	if complexityStr == "" {
		return nil
	}
	complexity, err := strconv.Atoi(complexityStr)
	if err != nil {
		return err
	}
	prefix = strings.Repeat("0", complexity)
	return nil
}

func main() {
	if err := getComplexity(); err != nil {
		log.Fatal(err)
	}
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		fmt.Println("SERVER_ADDRESS not set")
		return
	}

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	challenge, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading challenge:", err)
		return
	}
	challenge = strings.TrimSpace(challenge)
	fmt.Println("Challenge received:", challenge)

	start := time.Now()
	solution := solveChallenge(challenge)
	end := time.Now()

	fmt.Println("Solution found:", solution, end.Sub(start))
	_, err = conn.Write([]byte(solution + "\n"))
	if err != nil {
		fmt.Println("Error sending solution:", err)
		return
	}

	quote, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading quote:", err)
		return
	}
	fmt.Println("Quote:", quote)
}

func solveChallenge(challenge string) string {
	var solution int64
	for {
		solutionString := fmt.Sprintf("%d", solution)
		hash := sha256.Sum256([]byte(challenge + solutionString))
		hashString := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashString, prefix) {
			return solutionString
		}
		solution++
	}
}
