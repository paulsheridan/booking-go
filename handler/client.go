package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/paulsheridan/booking-go/database/client"
	"github.com/paulsheridan/booking-go/model"
)

type Client struct {
	Repo *client.RedisDatabase
}

func (c *Client) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name             string `json:"name"`
		Pronouns         string `json:"pronouns"`
		Over18           bool   `json:"over_18"`
		PreferredContact string `json:"preferred_contact"`
		PhoneNumber      int64  `json:"phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	client := model.Client{
		ClientID:         uuid.New(),
		Name:             body.Name,
		Pronouns:         body.Pronouns,
		Over18:           body.Over18,
		PreferredContact: body.PreferredContact,
		PhoneNumber:      body.PhoneNumber,
		CreatedAt:        &now,
	}

	err := c.Repo.Insert(r.Context(), client)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(client)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (c *Client) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := c.Repo.FindAll(r.Context(), client.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Client `json:"items"`
		Next  uint64         `json:"next,omitempty"`
	}
	response.Items = res.Clients
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (c *Client) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	clientID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	o, err := c.Repo.FindByID(r.Context(), clientID)
	if errors.Is(err, client.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Client) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	clientID, err := uuid.ParseBytes(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theClient, err := c.Repo.FindByID(r.Context(), clientID)
	if errors.Is(err, client.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.Repo.Update(r.Context(), theClient)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(theClient); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Client) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a client by ID")
}
