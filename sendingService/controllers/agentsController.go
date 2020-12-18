package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sendingService/models"
	"strconv"
	"sync"
	"time"
)

type AgentsAuthData struct {
	URL      string
	Login    string
	Password string
}

type AgentsData struct {
	URL            string
	AgentsAuthData AgentsAuthData
	JSONToSend     string `json:"data"`
}

type AuthToken struct {
	Token          string    `json:"token"`
	ExpirationDate time.Time `json:"expirationDate"`
}

type AuthTokenStorageMap struct {
	storage map[AgentsAuthData]*AuthToken
	mx      sync.RWMutex
}

func (s *AuthTokenStorageMap) Get(key AgentsAuthData) *AuthToken {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.storage[key]
}

func (s *AuthTokenStorageMap) AddOrUpdate(key AgentsAuthData, value *AuthToken) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.storage[key] = value
}

//AgentsController is a root structure needed for functions invocation through reflection. This functions handle job sending to a certain group of outer systems
type AgentsController struct{}

var tokenStorage = AuthTokenStorageMap{
	storage: make(map[AgentsAuthData]*AuthToken),
	mx:      sync.RWMutex{},
}

func returnCriticalError(err error) *models.SendingError {
	errorMessage := err.Error()
	return &models.SendingError{
		Message:   &errorMessage,
		ErrorType: -1,
	}
}

func returnCriticalErrorByString(errorMessage string) *models.SendingError {
	return &models.SendingError{
		Message:   &errorMessage,
		ErrorType: -1,
	}
}

func returLogicalError(err error) *models.SendingError {
	errorMessage := err.Error()
	return &models.SendingError{
		Message:   &errorMessage,
		ErrorType: -2,
	}
}

func returLogicalErrorByString(errorMessage string) *models.SendingError {
	return &models.SendingError{
		Message:   &errorMessage,
		ErrorType: -2,
	}
}

func parseData(data *json.RawMessage) (*AgentsData, error) {
	dataObject := &AgentsData{}

	err := json.Unmarshal(*data, dataObject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return dataObject, nil
}

func getAuthToken(authData AgentsAuthData) (*AuthToken, *models.SendingError) {
	storageToken := tokenStorage.Get(authData)

	if storageToken != nil && storageToken.ExpirationDate.After(time.Now()) {
		return storageToken, nil
	}

	jsonStr, _ := json.Marshal(struct {
		Login    string `login`
		Password string `password`
	}{
		authData.Login,
		authData.Password,
	})

	request, _ := http.NewRequest("POST", authData.URL, bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, returnCriticalError(err)
	}

	if response.StatusCode != 200 {
		return nil, returnCriticalError(errors.New("Agent job authorization failed. Link " + authData.URL + " Response error:" + response.Status))
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	tokenObject := &AuthToken{}
	err = decoder.Decode(tokenObject)

	if err != nil {
		return nil, returnCriticalError(err)
	}
	tokenStorage.AddOrUpdate(authData, tokenObject)
	return tokenObject, nil
}

//SendOrder sends order data to outer system
func (agents *AgentsController) SendOrder(data *json.RawMessage) *models.SendingError {
	agentData, err := parseData(data)

	if err != nil {
		return returnCriticalError(err)
	}

	authTokenObject, agentsError := getAuthToken(agentData.AgentsAuthData)

	if agentsError != nil {
		return agentsError
	}

	req, _ := http.NewRequest("POST", agentData.URL, bytes.NewBuffer([]byte(agentData.JSONToSend)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AuthToken", authTokenObject.Token)

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return returnCriticalError(err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		return nil
	case 500:
		return returLogicalErrorByString("Server " + agentData.URL + " returned 500")
	default:
		return returnCriticalErrorByString("Server " + agentData.URL + " returned" + strconv.Itoa(response.StatusCode))
	}
}

//SendCancel sends data for canceling order for previously sent order to outer system
func (agents *AgentsController) SendCancel(data *json.RawMessage) *models.SendingError {
	_, _ = parseData(data)
	//send something somewhere(unimplemented)
	return nil
}
