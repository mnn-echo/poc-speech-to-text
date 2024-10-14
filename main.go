package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const API_URL = "http://192.168.1.117:8000/transcribe"

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading connection: %v", err)
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		// Process the audio data (message)
		fmt.Printf("Received audio data: %d bytes\n", len(message))

		// Send a response back to the client
		err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Audio received : %d bytes\n ", len(message))))

		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

/*
	// call REST service pour Transcribe

	func uploadAudioBytes(audioBytes []byte) (string, error) {
		// Create a buffer to store the multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Create the file field in the multipart form
		part, err := writer.CreateFormFile("file", "audioFile")
		if err != nil {
			return "", fmt.Errorf("could not create form file: %v", err)
		}

		// Write the byte slice (audio data) to the form
		_, err = io.Copy(part, bytes.NewReader(audioBytes))
		if err != nil {
			return "", fmt.Errorf("could not copy byte slice to form: %v", err)
		}

		// Close the multipart writer to finalize the form
		err = writer.Close()
		if err != nil {
			return "", fmt.Errorf("could not close writer: %v", err)
		}

		// Create a new HTTP POST request
		req, err := http.NewRequest("POST", API_URL, body)
		if err != nil {
			return "", fmt.Errorf("could not create POST request: %v", err)
		}

		// Set the Content-Type header to the correct multipart boundary
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("error making HTTP request: %v", err)
		}
		defer resp.Body.Close()

		// Check if the response status is OK (200)
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to transcribe audio: %v", resp.Status)
		}

		// Read the response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response body: %v", err)
		}

		// Extract the transcription from the response
		transcription := extractTranscriptionFromResponse(responseBody)

		return transcription, nil
	}

	func extractTranscriptionFromResponse(responseBody []byte) string {
		var result map[string]interface{}
		err := json.Unmarshal(responseBody, &result)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			return ""
		}

		// Assuming the transcription is in a field named "transcription"
		if transcription, ok := result["transcription"].(string); ok {
			return transcription
		}

		return ""
	}
*/
func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
