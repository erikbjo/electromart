package categories

import (
	"Database_Project/internal/db"
	"Database_Project/internal/structs"
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

var detailImplementedMethods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

func HandleCategoryDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleGetDetailRequest(w, r)

	case http.MethodPut:
		handleUpdateDetailRequest(w, r)

	case http.MethodDelete:
		handleDeleteDetailRequest(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				detailImplementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleGetDetailRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
		return
	}

	// Get the category with the given ID
	category, err := db.GetCategoryByID(id)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "Error getting categories from database") {
		return
	}

	// Return the category
	if marshalledCategory, err := json.Marshal(category); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error during encoding response") {
		return
	} else {
		if _, err := w.Write(marshalledCategory); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error writing response") {
			return
		}
	}
}

func handleUpdateDetailRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
		return
	}

	// Decode the request body into a category
	var updatedCategory structs.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); utils.HandleError(w, r, http.StatusBadRequest, err, "Error decoding request body") {
		return
	}

	if updatedCategory.ID != id {
		utils.HandleError(w, r, http.StatusBadRequest, fmt.Errorf("ID in request body does not match ID in URL"), "ID in request body does not match ID in URL")
		return
	}

	// Update the category with the given ID
	if err := db.UpdateCategory(updatedCategory); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error updating category in database") {
		return
	}

	// Return no content
	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteDetailRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
		return
	}

	// Get the category with the given ID
	if err := db.DeleteCategoryByID(id); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error deleting category from database") {
		return
	}

	// Return no content
	w.WriteHeader(http.StatusNoContent)
}
