package api

import (
    "encoding/json"
    "github.com/labstack/echo"
    "gopkg.in/olahol/melody.v1"
    "log"
)

// =======================================================================================
// SocketApi
type requestLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Key       string  `json:"key"`
}

type responseWebSocket struct {
	Success   bool              `json:"success"`
	Locations []requestLocation `json:"locations"`
}

func SocketApi(e *echo.Echo) {
	m := melody.New()

	// Create new group routes ws api
	ws := e.Group("/api/v1/ws")

	// Routes
	ws.GET("/location", func(c echo.Context) error {
		m.HandleRequest(c.Response(), c.Request())
		return nil
	})

	requestLocations := make([]requestLocation, 0)

	// Response message
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//p := strings.Split(string(msg), " ")
		request := requestLocation{}
		err := json.Unmarshal(msg, &request)
		if err != nil {
			log.Printf("No se pudo convertir el json recibido: %v", err)
			return
		}

		// get connection
		//db := config.GetConnection()
		//defer db.Close()

        //if err := db.First(&mobile, mobile.ID).Error; err != nil {
        //    return err
        //}

		// get
		exist := false
		for k, v := range requestLocations {
			if v.Name == request.Name {
				requestLocations[k] = requestLocation{
					Name:      requestLocations[k].Name,
					Key:       requestLocations[k].Key,
					Latitude:  request.Latitude,
					Longitude: request.Longitude,
				}
				exist = true
			}
		}

		if !exist {
			requestLocations = append(requestLocations, request)
		}

		rws := responseWebSocket{
			Success:   true,
			Locations: requestLocations,
		}

		r, err := json.Marshal(rws)
		if err != nil {
			log.Printf("no se pudo procesar el objeto de respuesta: %v", err)
		}

		// Send New data
		m.Broadcast(r)
	})
}
