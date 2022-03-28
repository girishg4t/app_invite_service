package repo_test

import (
	"time"

	"github.com/girishg4t/app_invite_service/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("db operations", func() {
	var (
		token    string
		username string
		date     time.Time
	)
	BeforeEach(func() {
		token = "23sg325sdsg5"
		username = "test"
	})
	Context("Check app token", func() {
		It("save app token", func() {
			date = time.Now()
			err := at.SaveAppToken(&model.AppToken{
				Username: username,
				Token:    token,
				ExpDate:  date.AddDate(0, 0, 7),
				IsActive: true,
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("get app token", func() {
			result, err := at.GetAppToken(&model.AppToken{
				Token: token,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Username).Should(Equal(username))
			Expect(result.IsActive).Should(Equal(true))
			Expect(result.ExpDate.Day()).To(Equal(date.AddDate(0, 0, 7).Day()))
		})
		It("update and check is active status", func() {
			result, err := at.UpdateAppToken(&model.AppToken{
				Token:    token,
				IsActive: false,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).ShouldNot(BeNil())

			apptoken, err := at.GetAppToken(result)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(apptoken).ShouldNot(BeNil())
			Expect(apptoken.IsActive).Should(Equal(false))
		})
		It("get all token", func() {
			result, err := at.GetAllAppToken()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).ShouldNot(BeNil())
			Expect(len(result)).Should(Equal(1))
		})
	})
})
