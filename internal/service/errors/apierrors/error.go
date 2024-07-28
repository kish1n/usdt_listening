package apierrors

import (
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

func ErrorConstructor(w http.ResponseWriter, logger logan.Entry, err error, debugInfo string,
	errStatus string, errTitle string, errDetail string) {
	logger.WithError(err).Debug(debugInfo)
	ape.Render(w, &jsonapi.ErrorObject{
		Status: errStatus,
		Title:  errTitle,
		Detail: errDetail,
	})
	return
}
