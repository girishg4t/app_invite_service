package repo_test

import (
	"github.com/girishg4t/app_invite_service/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("db operations", func() {
	Context("Check user", func() {
		It("is admin user present", func() {
			u, err := r.GetUser(&model.User{
				Username: "admin",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(u.Role).Should(Equal("ADMIN"))
		})
	})
})
