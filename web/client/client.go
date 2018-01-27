package client

import (
	"fmt"
	"net/http"
	"net/url"
	"psychic-rat/types"
	"strconv"
)

type webClient struct {
	client http.Client
	url    string
}

func New(c http.Client, url string) *webClient {
	return &webClient{c, url}
}

func (w *webClient) AddNewItem(item types.NewItem) (*types.NewItem, error) {
	currency := fmt.Sprintf("%d", item.CurrencyID)
	value := fmt.Sprintf("%d", item.Value)
	v := url.Values{"company": {item.Company}, "make": {item.Make}, "model": {item.Model}, "currencyID": {currency}, "value": {value}}

	resp, err := w.client.PostForm(w.url+"/newitem", v)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with %d: %s", resp.StatusCode, resp.Status)
	}
	return &item, nil
}

func (w *webClient) AddItem(item types.Item) (*types.Item, error) {
	v := url.Values{
		"id[]":          []string{""},
		"action[]":      []string{"add"},
		"isPledge[]":    []string{"true"},
		"item[]":        []string{""},
		"company[]":     []string{""},
		"userID[]":      []string{""},
		"usercompany[]": []string{item.Company.Name},
		"usermake[]":    []string{item.Make},
		"usermodel[]":   []string{item.Model},
		"currencyID[]":  []string{strconv.FormatInt(int64(item.CurrencyID), 10)},
		"value[]":       []string{strconv.FormatInt(int64(item.Value), 10)},
	}

	resp, err := w.client.PostForm(w.url+"/admin/newitems", v)
	if err != nil {
		return nil, fmt.Errorf("server responded with %d: %s", resp.StatusCode, resp.Status)
	}

	return &item, nil
}
