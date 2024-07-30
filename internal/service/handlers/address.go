package handlers

import (
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/service/errors/apierrors"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func SortBySender(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	db := helpers.DB(r)
	address, err := helpers.GetAddress(r, "sender")
	res, err := db.NewTransaction().FilterBySender(address)

	if res == nil {
		apierrors.ErrorConstructor(w, *logger, err, "404 not found", "404", "Not Found", "Not found transaction from this address")
		return
	}

	if err != nil {
		apierrors.ErrorConstructor(w, *logger, err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	logger.Infof("res: %s", res)
	ape.Render(w, res)
	return
}

func SortByRecipient(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)
	db := helpers.DB(r)

	address, err := helpers.GetAddress(r, "recipient")
	res, err := db.NewTransaction().FilterByRecipient(address)

	if res == nil {
		apierrors.ErrorConstructor(w, *logger, err, "404 not found", "404", "Not Found", "Not found transaction to this address")
		return
	}

	if err != nil {
		apierrors.ErrorConstructor(w, *logger, err, "Server error", "500", "Server error 500", "Unpredictable behavior")
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
	start, err := db.NewTransaction().FilterByRecipient(address)
	end, err := db.NewTransaction().FilterBySender(address)

	if end == nil && start == nil {
		apierrors.ErrorConstructor(w, *logger, err, "404 not found", "404", "Not Found", "Not found transaction at this address")
		return
	}

	if err != nil {
		apierrors.ErrorConstructor(w, *logger, err, "Server error", "500", "Server error 500", "Unpredictable behavior")
		return
	}

	response := map[string][]data.TransactionData{
		"send": start,
	}

	logger.Infof("res: %s", response)
	ape.Render(w, response)
	return
}
