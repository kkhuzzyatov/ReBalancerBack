package handlers

import (
	"encoding/json"
	"gomod/pkg/entities"
	"net/http"
)

func GetData(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    return
  }

  var user entities.User
  json.NewDecoder(r.Body).Decode(&user)

  // TODO: newUser, err := <Метод, получающий данные о пользователе>(user)
  /*stub:*/ var newUser entities.User
  newUser.Email = "email"
  newUser.Password = "password"
  newUser.CurAllocation = " "
  newUser.TargetAllocation = " "
  newUser.TaxRate = 13
  /*
  if err != nil {
    log.Fatalf(err.Error())
    return
  }
  */
  
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(newUser)
}