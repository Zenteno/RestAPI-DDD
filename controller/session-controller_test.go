package controller

import (
	"api-ddd/entity"
	"api-ddd/repository"
	"api-ddd/service"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test404Route(t *testing.T) {
	sessionServiceTest := service.NewSessionService(nil)
	sessionControllerTest := NewSessionController(sessionServiceTest)
	appTest := fiber.New()

	appTest.Use(sessionControllerTest.NotFound)
	resp, err := appTest.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Status code")
}

func TestAddLocation(t *testing.T) {

	repoTest := repository.NewMemRepository()
	sessionServiceTest := service.NewSessionService(repoTest)
	session := NewSessionController(sessionServiceTest)
	appTest := fiber.New()
	v1 := appTest.Group("/api/v1")

	v1.Get("/current_shopper_location/:shopper_uuid", session.CurrentLocation)
	v1.Get("/session_location_history/:session_uuid", session.HistorySession)
	v1.Post("/location", session.AddLocation)

	res, err := appTest.Test(httptest.NewRequest("POST", "/api/v1/location", nil))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 400, res.StatusCode, "Posting empty body")

	a := entity.ShopperHistory{}
	gofakeit.Struct(&a)
	d, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/api/v1/location", strings.NewReader(string(d)))
	req.Header.Set("Content-type", "application/json")
	res, err = appTest.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 400, res.StatusCode, "Posting empty Reported_at")

	gofakeit.Struct(&a)
	a.ReportedAt = gofakeit.Date()
	d, err = json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	reqN, err := appTest.Test(httptest.NewRequest("POST", "/api/v1/location", strings.NewReader(string(d))))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 400, reqN.StatusCode, "Posting without Content-type")

	gofakeit.Struct(&a)
	a.ReportedAt = gofakeit.Date()
	d, err = json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest("POST", "/api/v1/location", strings.NewReader(string(d)))
	req.Header.Set("Content-type", "application/json")
	res, err = appTest.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 201, res.StatusCode, "Posting with Content-type and all required parameters")
}

func TestHistorySession(t *testing.T) {

	repoTest := repository.NewMemRepository()
	sessionServiceTest := service.NewSessionService(repoTest)
	session := NewSessionController(sessionServiceTest)

	appTest := fiber.New()
	v1 := appTest.Group("/api/v1")

	v1.Get("/current_shopper_location/:shopper_uuid", session.CurrentLocation)
	v1.Get("/session_location_history/:session_uuid", session.HistorySession)
	v1.Post("/location", session.AddLocation)

	var sessionUuids []string
	// First I create 20 randoms sessions uuids, insert 10, query them, and finally query 10 not inserted
	for i := 0; i < 20; i++ {
		sessionUuids = append(sessionUuids, gofakeit.UUID())
	}
	for i := 0; i < 100; i++ {
		a := entity.ShopperHistory{}
		gofakeit.Struct(&a)
		a.SessionUuid = sessionUuids[rand.Intn(10)]
		a.ReportedAt = gofakeit.Date()
		d, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		res, err := PostRequest("/api/v1/location", string(d), appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 201, res.StatusCode, "Posting new element")
	}

	//testing inserted elements
	for i := 0; i < 50; i++ {
		res, err := GetRequest("/api/v1/session_location_history/"+sessionUuids[rand.Intn(10)], appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, res.StatusCode, "Testing existing element")
	}

	//testing non inserted elements
	for i := 0; i < 50; i++ {
		res, err := GetRequest("/api/v1/session_location_history/"+sessionUuids[rand.Intn(10)+10], appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 404, res.StatusCode, "Testing non existing element")
	}

}

func TestCurrentLocation(t *testing.T) {
	repoTest := repository.NewMemRepository()
	sessionServiceTest := service.NewSessionService(repoTest)
	session := NewSessionController(sessionServiceTest)

	appTest := fiber.New()
	v1 := appTest.Group("/api/v1")

	v1.Get("/current_shopper_location/:shopper_uuid", session.CurrentLocation)
	v1.Get("/session_location_history/:session_uuid", session.HistorySession)
	v1.Post("/location", session.AddLocation)

	var shopperUuids []string
	// First I create 10 randoms shopper uuids, insert 20 and query them, after query 10 not inserted
	// 10 of the inserted will be in a range between now and ten minutes ago
	// 10 of the inserted will be older than 10 minutes ago
	// 10 will not be inserted, it should return 404 on their query

	for i := 0; i < 30; i++ {
		shopperUuids = append(shopperUuids, gofakeit.UUID())
	}
	//inserting elements older than 10 minutes ago (12 hours range of random time)
	for i := 0; i < 30; i++ {
		a := entity.ShopperHistory{}
		gofakeit.Struct(&a)
		a.ShopperUuid = shopperUuids[rand.Intn(10)+10]
		//random date older than 10 minutes + 1 second go
		randDate := time.Now().Add(time.Second * time.Duration(-rand.Intn(60*60*12)-60*10+1))
		a.ReportedAt = randDate
		d, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		res, err := PostRequest("/api/v1/location", string(d), appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 201, res.StatusCode, "Element older than 10 minutes ago")
	}
	//inserting elements between now and 10 minutes ago
	for i := 0; i < 30; i++ {
		a := entity.ShopperHistory{}
		gofakeit.Struct(&a)
		a.ShopperUuid = shopperUuids[rand.Intn(10)]
		//random date between now and ten minutes ago
		randDate := time.Now().Add(time.Second * time.Duration(-rand.Intn(60*10)))
		a.ReportedAt = randDate
		d, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		res, err := PostRequest("/api/v1/location", string(d), appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 201, res.StatusCode, "Element between now and 10 minutes ago")
	}

	//testing inserted elements newer than 10 minutes ago
	for i := 0; i < 50; i++ {
		dato := shopperUuids[rand.Intn(10)]
		res, err := GetRequest("/api/v1/current_shopper_location/"+dato, appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, res.StatusCode, "Testing existing element newer than 10 minutes")
	}

	//testing inserted elements older than 10 minutes
	for i := 0; i < 50; i++ {
		res, err := GetRequest("/api/v1/current_shopper_location/"+shopperUuids[rand.Intn(10)+10], appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 404, res.StatusCode, "Testing existing element older than 10 minutes")
	}
	//testing non inserted elements older than 10 minutes
	for i := 0; i < 50; i++ {
		res, err := GetRequest("/api/v1/current_shopper_location/"+shopperUuids[rand.Intn(10)+20], appTest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 404, res.StatusCode, "Testing non existing elements")
	}
}
func GetRequest(url string, appTest *fiber.App) (*http.Response, error) {
	res, err := appTest.Test(httptest.NewRequest("GET", url, nil))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PostRequest(url string, body string, appTest *fiber.App) (*http.Response, error) {
	req := httptest.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-type", "application/json")
	res, err := appTest.Test(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
