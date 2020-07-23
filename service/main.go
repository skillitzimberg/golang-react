package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	api := http.Server{
		Addr:         "localhost:8080",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	// This is a channel listening for errors. We will be sending errors from Server.ListenAndServe.
	serverErrors := make(chan error, 1)

	// This goroutine starts the server (the service) which is listening for requests via Server.Handler.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		http.Handle("/favicon.ico", http.NotFoundHandler())
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Error from ListenAndServe: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown.")

		// Sets a deadline for requests to complete.
		const timeout = 5 * time.Second
		ctx, canel := context.WithTimeout(context.Background(), timeout)
		defer canel()

		// Ask the listener to shut down and load shed.
		// From https://medium.com/@dieswaytoofast/rate-limiting-vs-load-shedding-e4c41e854718: 'Load Shedding . . . [prevents] the system from getting overloaded . . . The idea here is . . . to ignore some requests, rather than having the system [fail] and not be able to serve any requests — think “Let 911 calls through, and ignore the rest”.'

		// From https://cloud.google.com/blog/products/gcp/using-load-shedding-to-survive-a-success-disaster-cre-life-lessons: "Load shedding is a technique that allows your system to serve nominal capacity, regardless of how much traffic is being sent to it, in order to maintain availability. To do this, you'll need to throw away some requests and make clients retry."
		err := api.Shutdown(ctx)
		if err != nil {
			log.Fatalf("main : Gracful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : Could not gracefully shut down the server: %v", err)
		}
	}

}

// Echo describes the request that was made.
func Echo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if req.Method == http.MethodPost {
		sum := add(req)
		log.Println(sum)
		fmt.Fprintf(w, "%v", sum)
		return
	}

	trncHex := generateHex()
	log.Println(trncHex)
	fmt.Fprintf(w, "%s", trncHex)
}

func check(err error, from string) {
	if err != nil {
		log.Printf("Error from %s: %s", from, err)
	}
}

func add(req *http.Request) int {
	num1 := req.FormValue("num1")
	num2 := req.FormValue("num2")
	one, err := strconv.Atoi(num1)
	check(err, "strconv.Atoi")
	two, err := strconv.Atoi(num2)
	check(err, "strconv.Atoi")
	return one + two
}

func generateHex() string {
	n := rand.Intn(1000)
	log.Println("Start", n)
	defer log.Println("End", n)
	randByte := make([]byte, 6)
	rand.Read(randByte)
	hexStr := hex.EncodeToString(randByte)
	return hexStr[0:6]
}
