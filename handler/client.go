package handler

import (
	"fmt"
	"net/http"
)

type Client struct{}

func (c *Client) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a client")
}

func (c *Client) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all clients")
}

func (c *Client) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get a client by ID")
}

func (c *Client) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update a client by ID")
}

func (c *Client) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a client by ID")
}
