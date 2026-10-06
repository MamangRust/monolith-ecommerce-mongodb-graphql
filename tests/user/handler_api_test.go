package user_test

import (
	"fmt"
	"net/http"
	"testing"

	tests "github.com/MamangRust/monolith-ecommerce-test"

	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	tests.BaseTestSuite
	gql       http.Handler
	userID    int
	userEmail string
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()
	s.SetupUserService()

	cacheStore := s.GetGatewayCacheStore()
	s.gql = tests.NewGraphQLHandler(s.Conns, cacheStore, s.Log)
}

func (s *UserHandlerTestSuite) TestUserApiLifecycle() {
	// 1. Create
	s.userEmail = "handler.user@example.com"
	body, err := tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		createUser(input: {
			firstname: "Handler"
			lastname: "User"
			email: "%s"
			password: "password123"
			confirm_password: "password123"
		}) {
			status message
			data { id email firstname }
		}
	}`, s.userEmail))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "createUser")
	s.Require().NotNil(data)
	s.userID = int(data["id"].(float64))
	s.Require().NotZero(s.userID)
	s.Equal(s.userEmail, data["email"])

	// 2. FindAll
	body, err = tests.DoGraphQL(s.gql, `query { findAllUsers(input: { page: 1, page_size: 10 }) { status message data { id email } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findAllUsers"))

	// 3. FindById
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`query { findByIdUser(input: { id: %d }) { status message data { id email } } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "findByIdUser")
	s.Require().NotNil(data)
	s.Equal(float64(s.userID), data["id"])

	// 4. FindByActive
	body, err = tests.DoGraphQL(s.gql, `query { findByActiveUsers(input: { page: 1, page_size: 10 }) { status message data { id email } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByActiveUsers"))

	// 5. Update
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation {
		updateUser(input: {
			id: %d
			firstname: "Updated"
			lastname: "User"
			email: "%s"
			password: "password123"
			confirm_password: "password123"
		}) {
			status message
			data { id firstname email }
		}
	}`, s.userID, s.userEmail))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data = tests.GraphQLOpEnvelopeData(body, "updateUser")
	s.Require().NotNil(data)
	s.Equal("Updated", data["firstname"])

	// 6. Trash
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedUser(input: { id: %d }) { status message data { id } } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 7. FindByTrashed
	body, err = tests.DoGraphQL(s.gql, `query { findByTrashedUsers(input: { page: 1, page_size: 10 }) { status message data { id } } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	s.Require().NotNil(tests.GraphQLOpData(body, "findByTrashedUsers"))

	// 8. Restore
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { restoreUser(input: { id: %d }) { status message data { id } } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 9. DeletePermanent (trash first, like the old REST flow)
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { trashedUser(input: { id: %d }) { status message } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	body, err = tests.DoGraphQL(s.gql, fmt.Sprintf(`mutation { deleteUserPermanent(input: { id: %d }) { status message } }`, s.userID))
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 10. RestoreAll
	body, err = tests.DoGraphQL(s.gql, `mutation { restoreAllUser { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))

	// 11. DeleteAll
	body, err = tests.DoGraphQL(s.gql, `mutation { deleteAllUserPermanent { status message } }`)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
