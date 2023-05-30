package test

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/coreservices/user"
	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserService", func() {

	var (
		f *fuzz.Fuzzer
	)

	BeforeEach(func() {
		f = fuzz.New()
	})

	Describe("IsemailpasswordValidate", func() {
		It("should return no error if email and password are valid", func() {

			email := "socialtest@tftus.com"
			password := "Tftus@1234"

			user := &user.UserInfo{
				Email:    email,
				Password: password,
			}

			err := user.IsemailpasswordValidate()
			Expect(err).To(BeNil())
		})
		It("should return an error if email and password are not valid", func() {
			var email, password string
			f.Fuzz(&email)
			f.Fuzz(&password)

			user := &user.UserInfo{
				Email:    email,
				Password: password,
			}

			err := user.IsemailpasswordValidate()

			Expect(err).ToNot(BeNil())
		})
	})
	Describe("SignupService", func() {
		It("should create user and return no error if all required details are not missing ", func() {

			name := "tft"
			email := "socialtest@tftus.com"
			password := "Test@12345"

			user := &user.UserInfo{
				Name:     name,
				Email:    email,
				Password: password,
			}

			err := user.SignupService()
			Expect(err).To(BeNil())
		})
		It("should return error if user already exists ", func() {

			name := "tft"
			email := "socialtest@tftus.com"
			password := "Test@12345"

			user := &user.UserInfo{
				Name:     name,
				Email:    email,
				Password: password,
			}

			err := user.SignupService()
			Expect(err).ToNot(BeNil())
		})
		It("should return error if any required details are missing", func() {

			name := ""
			email := "socialtest@tftus.com"
			password := "Test@12345"

			user := &user.UserInfo{
				Name:     name,
				Email:    email,
				Password: password,
			}

			err := user.SignupService()
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("Login", func() {
		It("should return error if login credentials mis match", func() {

			email := "socialtestdummy@tftus.com"
			password := "Testdummy@12345"

			userinfo := &user.UserInfo{
				Email:    email,
				Password: password,
			}

			_, err := userinfo.Login()
			Expect(err).ToNot(BeNil())
		})
		It("should return userdetails and no error if valid credentials given", func() {

			email := "socialtest@tftus.com"
			password := "Test@12345"

			testUser := &types.User{}

			userinfo := &user.UserInfo{
				Email:    email,
				Password: password,
			}

			userdetails, err := userinfo.Login()
			testUser.Key = userdetails.Key
			testUser.Token, _ = user.GenerateJwt(int(testUser.Key))

			Expect(userdetails).To(Equal(testUser))
			Expect(err).To(BeNil())
		})
	})
	Describe("IsUserFoundByEmail", func() {
		It("should return no error if email exists in our records", func() {
			email := "socialtest@tftus.com"
			testresp := true
			user := &user.UserInfo{
				Email: email,
			}

			res, err := user.IsUserFoundByEmail()

			Expect(res).To(Equal(testresp))
			Expect(err).To(BeNil())

		})
		It("should return an error if email not exists in our record", func() {

			email := "socialtestdummy@tftus.com"
			testresp := false

			user := &user.UserInfo{
				Email: email,
			}

			res, err := user.IsUserFoundByEmail()
			Expect(res).To(Equal(testresp))
			Expect(err).ToNot(BeNil())

		})
	})
	Describe("IsUserSignupFoundByEmail", func() {
		It("should return true if user email exists", func() {

			email := "socialtest@tftus.com"
			testresp := true

			user := &user.UserInfo{
				Email: email,
			}

			res, _ := user.IsUserSignupFoundByEmail()
			Expect(res).To(Equal(testresp))
		})
		It("should return false if user email does not exist", func() {

			email := "socialtestdummy@tftus.com"
			testresp := false

			user := &user.UserInfo{
				Email: email,
			}

			res, _ := user.IsUserSignupFoundByEmail()

			Expect(res).To(Equal(testresp))
		})
	})
	Describe("UpdateUserProfile", func() {
		var (
			id       = 1
			name     = "tft"
			email    = "test@tftus.com"
			password = "Test@12345"
		)
		testuser := &types.DbUser{
			ID:       int64(id),
			Name:     name,
			Email:    email,
			Password: password,
		}
		testuser, err := postgres.CreateUser(name, email, password)
		Expect(err).To(BeNil())
		It("should update user profile and return no error if all parameters are provided", func() {
			err := user.UpdateUserProfile(testuser.ID, name, email, password)
			Expect(err).To(BeNil())
		})
		It("should update user profile with existing values and return no error if any parameter is empty", func() {
			// Set the name parameter as empty to test if it retains the existing value
			err := user.UpdateUserProfile(testuser.ID, "", email, password)
			Expect(err).To(BeNil())
		})
		It("should return an error if user does not exist", func() {
			err := user.UpdateUserProfile(999, name, email, password)
			Expect(err).ToNot(BeNil())
		})
		AfterSuite(func() {
			_, err := postgres.DeleteUser(testuser.ID)
			Expect(err).To(BeNil())
		})
	})
	// This handler is still not yet implemented .
	//This test case is included to only delete the testuser that we have created
	Describe("Delete", func() {
		It("should delete the user and return no error", func() {

			email := "socialtest@tftus.com"

			user, _ := postgres.GetUserByEmail(email)
			_, err := postgres.DeleteUser(user.ID)
			Expect(err).To(BeNil())
		})
	})

})
