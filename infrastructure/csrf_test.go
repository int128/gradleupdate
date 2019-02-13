package infrastructure_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/infrastructure"
)

func TestCSRFMiddlewareFactory_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	credentials := gatewaysTestDoubles.NewMockCredentials(ctrl)
	credentials.EXPECT().
		Get(gomock.Not(nil)).
		AnyTimes().
		Return(&config.Credentials{CSRFKey: []byte("0123456789abcdef0123456789abcdef")}, nil)
	factory := infrastructure.CSRFMiddlewareFactory{
		Logger:      gatewaysTestDoubles.NewLogger(t),
		Credentials: credentials,
	}

	router := mux.NewRouter()
	router.Use(factory.New())
	router.Methods("GET").Path("/form").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := fmt.Fprint(w, csrf.TemplateField(r)); err != nil {
				t.Errorf("error while writing body: %s", err)
			}
		})
	router.Methods("POST").Path("/form").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("NoToken", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/form", "", nil)
		if err != nil {
			t.Fatalf("error while sending a request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("StatusCode wants %v but %v", http.StatusForbidden, resp.StatusCode)
		}
	})

	t.Run("WithToken", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/form")
		if err != nil {
			t.Fatalf("error while sending a request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %v but %v", http.StatusOK, resp.StatusCode)
		}
		cookies := resp.Cookies()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error while reading body: %s", err)
		}

		form := make(url.Values)
		form.Set(parseInputForm(t, body))
		r, err := http.NewRequest("POST", server.URL+"/form", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("error while creating a request: %s", err)
		}
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		r.Header.Set("content-type", "application/x-www-form-urlencoded")
		dumpRequestOut(t, r)
		resp, err = http.DefaultClient.Do(r)
		if err != nil {
			t.Fatalf("error while sending a request: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %v but %v", http.StatusOK, resp.StatusCode)
		}
	})
}

var formKey = regexp.MustCompile(`name="(.+?)"`)
var formValue = regexp.MustCompile(`value="(.+?)"`)

func parseInputForm(t *testing.T, body []byte) (string, string) {
	k := formKey.FindSubmatch(body)
	if len(k) != 2 {
		t.Errorf("could not find form name: %s", string(body))
	}
	v := formValue.FindSubmatch(body)
	if len(v) != 2 {
		t.Errorf("could not find form value: %s", string(body))
	}
	return string(k[1]), string(v[1])
}

func dumpRequestOut(t *testing.T, r *http.Request) {
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		t.Errorf("error while dump the request: %s", err)
	}
	t.Logf("sending request: %s", string(dump))
}
