package handler

import model "GeoLoginSer/models"

// JsonError is a generic error in JSON format
//
// swagger:response jsonError
type jsonError struct {
	// in: body
	Message string `json:"message"`
}

// PersonResponse contains a single person information
//
// swagger:response personResponse
type personResponse struct {
	// in: body
	Payload *model.User `json:"User"`
}

// PeopleResponse constains all people from database information
//
// swagger:response peopleResponse
type peopleResponse struct {
	// in: body
	Payload *[]model.User `json:"user"`
}
