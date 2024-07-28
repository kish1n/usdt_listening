package handlers

import (
	"github.com/google/jsonapi"
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func SortBySender(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	db := helpers.DB(r)
	address, err := helpers.GetAddress(r, "from_address")
	logger.Infof("from_address %s:", address)
	res, err := db.Link().SortByParameter(address, "from_address")

	if res == nil {
		logger.WithError(err).Debug("Not Found 404")
		ape.Render(w, &jsonapi.ErrorObject{
			Status: "404",
			Title:  "Not Found 404",
			Detail: "Nonexistent link",
		})
		return
	}

	if err != nil {
		logger.WithError(err).Debug("Server error")
		ape.Render(w, &jsonapi.ErrorObject{
			Status: "500",
			Title:  "Server error 500",
			Detail: "Unpredictable behavior",
		})
		return
	}

	logger.Infof("res: %s", res)
	ape.Render(w, res)
	return
}

func SortByRecipient(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	db := helpers.DB(r)

	address, err := helpers.GetAddress(r, "to_address")
	res, err := db.Link().SortByParameter(address, "to_address")

	if res == nil {
		logger.WithError(err).Debug("Not Found 404")
		ape.Render(w, &jsonapi.ErrorObject{
			Status: "404",
			Title:  "Not Found 404",
			Detail: "Nonexistent link",
		})
		return
	}

	if err != nil {
		logger.WithError(err).Debug("Server error")
		ape.Render(w, &jsonapi.ErrorObject{
			Status: "500",
			Title:  "Server error 500",
			Detail: "Unpredictable behavior",
		})
		return
	}

	logger.Infof("res: %s", res)
	ape.Render(w, res)
	return
}

func SortByAddress(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	db := helpers.DB(r)

	address, err := helpers.GetAddress(r, "address")
	start, err := db.Link().SortByParameter(address, "to_address")
	end, err := db.Link().SortByParameter(address, "from_address")

	logger.Infof("res: %s", start)
	logger.Infof("res: %s", end)

	if end == nil && start == nil {
		ape.Render(w, &jsonapi.ErrorObject{
			Status: "404",
			Title:  "Not Found 404",
			Detail: "Nonexistent link",
		})
		return
	}

	if err != nil {
		ape.Render(w, &jsonapi.ErrorObject{
			Status: "500",
			Title:  "Server error 500",
			Detail: "Unpredictable behavior",
		})
		return
	}

	response := map[string][]data.TransactionData{
		"send": start,
	}

	ape.Render(w, response)
}
