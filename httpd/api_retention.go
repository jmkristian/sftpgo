package httpd

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/drakkan/sftpgo/v2/common"
	"github.com/drakkan/sftpgo/v2/dataprovider"
)

func getRetentionChecks(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	render.JSON(w, r, common.RetentionChecks.Get())
}

func startRetentionCheck(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	username := getURLParam(r, "username")
	user, err := dataprovider.UserExists(username)
	if err != nil {
		sendAPIResponse(w, r, err, "", getRespStatus(err))
		return
	}
	var check common.RetentionCheck
	err = render.DecodeJSON(r.Body, &check.Folders)
	if err != nil {
		sendAPIResponse(w, r, err, "", http.StatusBadRequest)
		return
	}
	if err := check.Validate(); err != nil {
		sendAPIResponse(w, r, err, "Invalid folders to check", http.StatusBadRequest)
		return
	}
	c := common.RetentionChecks.Add(check, &user)
	if c == nil {
		sendAPIResponse(w, r, err, fmt.Sprintf("Another check is already in progress for user %#v", username),
			http.StatusConflict)
		return
	}
	go c.Start()
	sendAPIResponse(w, r, err, "Check started", http.StatusAccepted)
}
