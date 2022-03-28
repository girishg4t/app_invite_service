package repo_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/girishg4t/app_invite_service/pkg/db"
	"github.com/girishg4t/app_invite_service/pkg/logging"
	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/girishg4t/app_invite_service/pkg/repo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	r  *repo.UserProcessor
	at *repo.AppTokenProcessor
)

func TestRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var tempDir string
var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	tempDir, err = ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())
	return []byte{}
}, func(data []byte) {
	dbCfg := model.Database{
		DataSource: filepath.Join(tempDir, "pulseid_test.db"),
		Debug:      len(os.Getenv("DEBUG")) != 0,
		Schema:     "../../scripts",
		Type:       "sqlite3",
	}
	dbconn, err := db.ConnectToDb(dbCfg)
	Expect(err).ShouldNot(HaveOccurred())
	log := logging.GetLogger()
	r = repo.NewUserProcessor(dbconn, log)
	Expect(r).NotTo(Equal(nil))
	at = repo.NewAppTokenProcessor(dbconn)
	Expect(at).NotTo(Equal(nil))
})
