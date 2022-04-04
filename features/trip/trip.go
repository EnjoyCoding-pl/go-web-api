package trip

import (
	"net/http"
	"time"

	"go-web-api/features/trip/domain"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TripService struct {
	Db *gorm.DB
}

func (t *TripService) MuxRegister(r *mux.Router) {
	r.HandleFunc("/trips", t.addHandler).Methods("GET")
}

func (t *TripService) addHandler(w http.ResponseWriter, r *http.Request) {

	t.Db.Create(&domain.Trip{
		Country: "Poland",
		Begin:   time.Now(),
		End:     time.Now().Add(time.Duration(10)),
	})

}
