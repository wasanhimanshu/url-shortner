package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

type UrlHandler struct {
	db *gorm.DB
}

func NewUrlHandler(db *gorm.DB) *UrlHandler {
	return &UrlHandler{db: db}
}

func (u *UrlHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/short", u.handlerUrl).Methods("POST")
}

func (u *UrlHandler) handlerUrl(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("missing url"))
		return
	}
	var payload *CreateShortUrlPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	shortId, err := shortid.Generate()
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println(shortId)
	var url = Url{}
	rows := u.db.Where("redirect_url = ?", payload.UrlAddr).First(&url)
	var result = make(map[string]string)
	if rows.RowsAffected > 0 {
		url.Clicks = url.Clicks + 1
		rows := u.db.Save(&url)
		if rows.Error != nil {
			WriteError(w, http.StatusInternalServerError, rows.Error)
			return
		}
		result = map[string]string{"shortId": url.Shorturl, "redirectUrl": url.RedirectUrl, "clicks": strconv.FormatUint(uint64(url.Clicks), 20)}
	} else {
		url = Url{RedirectUrl: payload.UrlAddr, Shorturl: shortId, Clicks: 1}
		rows := u.db.Create(&url)
		if rows.Error != nil {
			WriteError(w, http.StatusInternalServerError, rows.Error)
			return
		}
		result = map[string]string{"shortId": url.Shorturl, "redirectUrl": url.RedirectUrl, "clicks": strconv.FormatUint(uint64(url.Clicks), 20)}
	}
	WriteJSON(w, http.StatusCreated, result)
}
