package main

import (
	"encoding/json"
	"evo1/core"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var gameLoopMutex sync.Mutex
var isGameLoopRunning bool = true
var GameSpeed int = 1

func RunAPI() {

	initGame()

	// Запускаем горутину для инкремента счётчика каждую секунду
	go gameLoopGoro()

	// Регистрируем обработчики
	http.HandleFunc("/api/getImage",
		getViewPortImage)
	http.HandleFunc("/api/getAllGameInfo",
		getAllGameInfo)

	http.HandleFunc("/api/toggleGameLoopRunning",
		toggleGameLoopRunning)
	http.HandleFunc("/api/changeLogEnergy",
		changeLogEnergy)
	http.HandleFunc("/api/changeMaxAge",
		changeMaxAge)
	http.HandleFunc("/api/changeGameSpeed",
		changeGameSpeed)
	http.HandleFunc("/api/resetWorld",
		resetWorld)
	// Запускаем сервер
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func gameLoopGoro() {
	for {
		if isGameLoopRunning {
			gameLoopMutex.Lock()
			gameLogicLoop()
			gameLoopMutex.Unlock()
		}
		time.Sleep(time.Microsecond * time.Duration(GameSpeed))
	}
}

func getViewPortImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLoopMutex.Lock()
	image := viewport.GetImage()
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := struct {
		Text string `json:"text"`
	}{
		Text: image,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func toggleGameLoopRunning(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLoopMutex.Lock()
	isGameLoopRunning = !isGameLoopRunning
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "success}`)
}

type changeLogEnergyRequestBody struct {
	LogEnergy int `json:"LogEnergy"`
}

func changeLogEnergy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody changeLogEnergyRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameLoopMutex.Lock()
	logEnergy = requestBody.LogEnergy
	fmt.Printf("logEnergy было изменено на %d\n", logEnergy)
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "success}`)
}

type getAllGameInfoData struct {
	World     *core.World `json:"World"`
	LogEnergy int         `json:"LogEnergy"`
	MaxAge    int         `json:"MaxAge"`
	TreeCount int         `json:"TreeCount"`
}

func getAllGameInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLoopMutex.Lock()
	response := getAllGameInfoData{
		World:     core.MainWorld,
		LogEnergy: logEnergy,
		MaxAge:    maxAge,
		TreeCount: len(core.Trees),
	}
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

type changeMaxAgeRequestBody struct {
	MaxAge int `json:"MaxAge"`
}

func changeMaxAge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody changeMaxAgeRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameLoopMutex.Lock()
	maxAge = requestBody.MaxAge
	fmt.Printf("MaxAge было изменено на %d\n", maxAge)
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "success}`)
}

type changeGameSpeedRequestBody struct {
	GameSpeed int `json:"GameSpeed"`
}

func changeGameSpeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody changeGameSpeedRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameLoopMutex.Lock()
	GameSpeed = requestBody.GameSpeed
	fmt.Printf("GameSpeed было изменено на %d\n", GameSpeed)
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "success}`)
}

func resetWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameLoopMutex.Lock()
	initGame()
	fmt.Printf("Мир был ресетнут \n")
	gameLoopMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "success}`)
}
