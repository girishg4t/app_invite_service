package service_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/girishg4t/app_invite_service/pkg/repo/repofakes"
	"github.com/girishg4t/app_invite_service/pkg/service"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Token tests", func() {
	var (
		adminToken model.Token = model.Token{
			TokenString: "some random string",
			Role:        "ADMIN",
			Username:    "admin",
		}
		appTokenSvc      *service.AppTokenService
		fakeAppTokenRepo *repofakes.FakeIAppToken
		req              *http.Request
		rr               *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		fakeAppTokenRepo = new(repofakes.FakeIAppToken)
		appTokenSvc = service.NewAppTokenService(fakeAppTokenRepo)
		rr = httptest.NewRecorder()
	})

	Describe("Check alert count", func() {
		BeforeEach(func() {
			os.Setenv("EXPIRE_IN_DAYS", "1")
			var err error
			req, err = http.NewRequest("POST", "http://localhost:8081/login", strings.NewReader(`{
					"username": "admin",
					"password": "admin",
				}`))
			req = req.WithContext(StubToken(adminToken))
			Expect(err).NotTo(HaveOccurred())
		})
		It("generate the app token", func() {
			req, err := http.NewRequest("GET", "http://localhost:8081/v1/api/genToken", nil)
			Expect(err).NotTo(HaveOccurred())
			req = req.WithContext(StubToken(adminToken))
			appTokenSvc.GenToken(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).NotTo(Equal("{}\n"))

			Expect(fakeAppTokenRepo.SaveAppTokenCallCount()).To(Equal(1))
		})
		It("validate the app token as success", func() {
			fakeAppTokenRepo.GetAppTokenReturns(model.AppToken{
				ExpDate:  time.Now().AddDate(0, 0, 1),
				IsActive: true,
			}, nil)
			req, err := http.NewRequest("GET", "http://localhost/validatetoken", nil)
			Expect(err).NotTo(HaveOccurred())

			req = mux.SetURLVars(req, map[string]string{"appToken": "ZKww2WFHQfL"})
			appTokenSvc.ValidateToken(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).NotTo(Equal("{}\n"))

			Expect(fakeAppTokenRepo.GetAppTokenCallCount()).To(Equal(1))
		})
		It("validate the app token as fail", func() {
			fakeAppTokenRepo.GetAppTokenReturns(model.AppToken{
				ExpDate:  time.Now(),
				IsActive: true,
			}, nil)
			req, err := http.NewRequest("GET", "http://localhost:8081/validatetoken", nil)
			req = mux.SetURLVars(req, map[string]string{"appToken": "ZKww2WFHQfL"})
			Expect(err).NotTo(HaveOccurred())
			appTokenSvc.ValidateToken(rr, req)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(fakeAppTokenRepo.GetAppTokenCallCount()).To(Equal(1))
		})
		It("invalidated the token", func() {
			fakeAppTokenRepo.GetAppTokenReturns(model.AppToken{
				ExpDate:  time.Now().AddDate(0, 0, 1),
				IsActive: true,
			}, nil)
			req, err := http.NewRequest("POST", "http://localhost:8081/v1/api/invalidateToken", strings.NewReader(`{
					"appToken": "admin"
				}`))
			Expect(err).NotTo(HaveOccurred())
			req = req.WithContext(StubToken(adminToken))
			appTokenSvc.InvalidateToken(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(fakeAppTokenRepo.UpdateAppTokenCallCount()).To(Equal(1))
		})
		It("get all the tokens", func() {
			var allTokens []model.AppToken
			allTokens = append(allTokens, model.AppToken{ExpDate: time.Now().AddDate(0, 0, 1),
				IsActive: true})
			fakeAppTokenRepo.GetAllAppTokenReturns(allTokens, nil)
			req, err := http.NewRequest("GET", "http://localhost:8081/v1/api/getAllToken", nil)
			Expect(err).NotTo(HaveOccurred())
			req = req.WithContext(StubToken(adminToken))
			appTokenSvc.GetAllAppToken(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())
			var resp []model.AppToken
			err = json.Unmarshal([]byte(data), &resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp)).To(Equal(1))
		})
	})
})
