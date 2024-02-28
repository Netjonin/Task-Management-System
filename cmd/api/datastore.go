package main

import (
	"time"

	"TMS.netjonin.net/internal/data"
)

var store = map[int]data.Task{
	1: {
		ID: 1,
		Title: "Laundry",
		Description: "Laundry",
		CreatedAt: time.Now(),
		Status: "To-Do",
		ExpiredAt: time.Now().Add(time.Hour * time.Duration(4)),
		Expired: false,
		Version: 1,

	},
	2: {
		ID: 2,
		Title: "Feeding",
		Description: "Feeding",
		CreatedAt: time.Now().Add(time.Minute * 20),
		Status: "To-Do",
		ExpiredAt: time.Now().Add(time.Hour * time.Duration(6)),
		Expired: false,
		Version: 1,

	},
	3: {
		ID: 3,
		Title: "Football",
		Description: "Football",
		CreatedAt: time.Now().Add(time.Minute * 20),
		Status: "To-Do",
		ExpiredAt: time.Now().Add(time.Hour * time.Duration(6)),
		Expired: false,
		Version: 1,

	},
}