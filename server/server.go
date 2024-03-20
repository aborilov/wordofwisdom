package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var prefix = "0000"

var quotes = []string{
	"Life is about making an impact, not making an income. –Kevin Kruse",
	"Whatever the mind of man can conceive and believe, it can achieve. –Napoleon Hill",
	"Strive not to be a success, but rather to be of value. –Albert Einstein",
}

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
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Issue challenge
	challenge := generateChallenge()
	_, err := conn.Write([]byte(challenge + "\n"))
	if err != nil {
		fmt.Println("Error sending challenge:", err)
		return
	}

	// Set a deadline for reading the solution from the client
	timeoutDuration := 30 * time.Second
	if err := conn.SetReadDeadline(time.Now().Add(timeoutDuration)); err != nil {
		log.Fatal(err)
	}

	// Wait for solution
	solution, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("Connection timed out waiting for solution")
		} else {
			fmt.Println("Error reading solution:", err)
		}
		return
	}

	// Process solution
	if verifySolution(challenge, strings.TrimSpace(solution)) {
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		quote := quotes[rng.Intn(len(quotes))]
		if _, err := conn.Write([]byte(quote + "\n")); err != nil {
			return
		}
	} else {
		if _, err := conn.Write([]byte("Incorrect solution.\n")); err != nil {
			return
		}
	}
}

func generateChallenge() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return fmt.Sprintf("%x", rng.Int31())
}

func verifySolution(challenge, solution string) bool {
	hash := sha256.Sum256([]byte(challenge + solution))
	hashString := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hashString, prefix)
}
