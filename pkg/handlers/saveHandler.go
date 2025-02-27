package handlers

import (
	"encoding/json"
	"gomod/pkg/entities"
	"net/http"
)

func Save(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    return
  }

  var user entities.User
  json.NewDecoder(r.Body).Decode(&user)

  // TODO: userEmail, err := <Метод, сохраняющий данные пользователя в базу данных>(user)
  /*stub:*/ userEmail := "email"
  /*
  if err != nil {
    log.Fatalf(err.Error())
    return
  }
  */
  
  w.Write([]byte(userEmail))
}
