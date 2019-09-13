package v7action_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/actionerror"

	. "code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/actor/v7action/v7actionfakes"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Actions", func() {
	var (
		actor                     *Actor
		fakeUAAClient             *v7actionfakes.FakeUAAClient
		fakeCloudControllerClient *v7actionfakes.FakeCloudControllerClient
	)

	BeforeEach(func() {
		fakeUAAClient = new(v7actionfakes.FakeUAAClient)
		fakeCloudControllerClient = new(v7actionfakes.FakeCloudControllerClient)
		fakeConfig := new(v7actionfakes.FakeConfig)
		actor = NewActor(fakeCloudControllerClient, fakeConfig, nil, fakeUAAClient, nil)
	})

	Describe("CreateUser", func() {
		var (
			actualUser     User
			actualWarnings Warnings
			actualErr      error
		)

		JustBeforeEach(func() {
			actualUser, actualWarnings, actualErr = actor.CreateUser("some-new-user", "some-password", "some-origin")
		})

		When("no API errors occur", func() {
			var createdUser ccv3.User

			BeforeEach(func() {
				createdUser = ccv3.User{
					GUID: "new-user-cc-guid",
				}
				fakeUAAClient.CreateUserReturns(
					uaa.User{
						ID: "new-user-uaa-id",
					},
					nil,
				)
				fakeCloudControllerClient.CreateUserReturns(
					createdUser,
					ccv3.Warnings{
						"warning-1",
						"warning-2",
					},
					nil,
				)
			})

			It("creates a new user and returns all warnings", func() {
				Expect(actualErr).NotTo(HaveOccurred())

				Expect(actualUser).To(Equal(User(createdUser)))
				Expect(actualWarnings).To(ConsistOf("warning-1", "warning-2"))

				Expect(fakeUAAClient.CreateUserCallCount()).To(Equal(1))
				username, password, origin := fakeUAAClient.CreateUserArgsForCall(0)
				Expect(username).To(Equal("some-new-user"))
				Expect(password).To(Equal("some-password"))
				Expect(origin).To(Equal("some-origin"))

				Expect(fakeCloudControllerClient.CreateUserCallCount()).To(Equal(1))
				uaaUserID := fakeCloudControllerClient.CreateUserArgsForCall(0)
				Expect(uaaUserID).To(Equal("new-user-uaa-id"))
			})
		})

		When("a create user request to the UAA returns an error", func() {
			var returnedErr error

			BeforeEach(func() {
				returnedErr = errors.New("some UAA error")
				fakeUAAClient.CreateUserReturns(
					uaa.User{
						ID: "new-user-uaa-id",
					},
					returnedErr,
				)
			})

			It("returns the same error", func() {
				Expect(actualErr).To(MatchError(returnedErr))
			})
		})

		When("the CC API returns an error", func() {
			var returnedErr error

			BeforeEach(func() {
				returnedErr = errors.New("CC error")
				fakeUAAClient.CreateUserReturns(
					uaa.User{
						ID: "new-user-uaa-id",
					},
					nil,
				)
				fakeCloudControllerClient.CreateUserReturns(
					ccv3.User{},
					ccv3.Warnings{
						"warning-1",
						"warning-2",
					},
					returnedErr,
				)
			})

			It("returns the same error and all warnings", func() {
				Expect(actualErr).To(MatchError(returnedErr))
				Expect(actualWarnings).To(ConsistOf("warning-1", "warning-2"))
			})
		})
	})

	Describe("GetUser", func() {
		var (
			actualUser User
			actualErr  error
		)

		JustBeforeEach(func() {
			actualUser, actualErr = actor.GetUser("some-user", "some-origin")
		})

		When("when the API returns a success response", func() {
			When("the API returns one user", func() {
				BeforeEach(func() {

					fakeUAAClient.GetUsersReturns(
						[]uaa.User{
							{ID: "user-id"},
						},
						nil,
					)
				})

				It("returns the single user", func() {
					Expect(actualErr).NotTo(HaveOccurred())
					Expect(actualUser).To(Equal(User{GUID: "user-id"}))

					Expect(fakeUAAClient.GetUsersCallCount()).To(Equal(1))
					username, origin := fakeUAAClient.GetUsersArgsForCall(0)
					Expect(username).To(Equal("some-user"))
					Expect(origin).To(Equal("some-origin"))
				})
			})

			When("the API returns no user", func() {
				BeforeEach(func() {
					fakeUAAClient.GetUsersReturns(
						[]uaa.User{},
						nil,
					)
				})

				It("returns an error indicating user was not found in UAA", func() {
					Expect(actualUser).To(Equal(User{}))
					Expect(actualErr).To(Equal(actionerror.UAAUserNotFoundError{}))
					Expect(fakeUAAClient.GetUsersCallCount()).To(Equal(1))
				})
			})
		})

		When("the API returns an error", func() {
			var apiError error

			BeforeEach(func() {
				apiError = errors.New("uaa-api-get-users-error")
				fakeUAAClient.GetUsersReturns(
					nil,
					apiError,
				)
			})

			It("returns error coming from the API", func() {
				Expect(actualUser).To(Equal(User{}))
				Expect(actualErr).To(MatchError("uaa-api-get-users-error"))

				Expect(fakeUAAClient.GetUsersCallCount()).To(Equal(1))
			})
		})
	})
})