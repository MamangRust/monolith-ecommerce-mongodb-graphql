package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	tests "github.com/MamangRust/monolith-ecommerce-test"

	"github.com/stretchr/testify/suite"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/roles"
)

type AuthHandlerApiTestSuite struct {
	tests.BaseTestSuite
	gql         http.Handler
	tokenMgr    auth.TokenManager
	email       string
	password    string
	accessToken string
	userID      int
}

func (s *AuthHandlerApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()

	s.SetupAuthService()

	cacheStore := s.GetGatewayCacheStore()
	tokenManager, err := auth.NewManager("mysecret")
	s.Require().NoError(err)
	s.tokenMgr = tokenManager

	s.gql = tests.NewAuthenticatedGraphQLHandler(s.Conns, cacheStore, s.Log, tokenManager)

	// Seed a role so registration flows that assign default roles can resolve it.
	_, _ = pbrole.NewRoleCommandServiceClient(s.Conns["role"]).CreateRole(
		s.Ctx,
		&pbrole.CreateRoleRequest{Name: "ROLE_ADMIN"},
	)

	s.email = "auth.handler.api.test@example.com"
	s.password = "password123"
}

func (s *AuthHandlerApiTestSuite) register(email, password string) int {
	body, err := tests.DoGraphQLAuth(s.gql, fmt.Sprintf(`mutation {
		registerUser(input: {
			firstname: "Auth"
			lastname: "API"
			email: "%s"
			password: "%s"
			confirm_password: "%s"
		}) {
			status message
			data { id email }
		}
	}`, email, password, password), "")
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "registerUser")
	s.Require().NotNil(data)
	s.Require().Equal(email, data["email"])
	id, ok := data["id"].(float64)
	s.Require().True(ok)
	s.Require().NotZero(int(id))
	return int(id)
}

func (s *AuthHandlerApiTestSuite) Test1_Register() {
	s.userID = s.register(s.email, s.password)
}

func (s *AuthHandlerApiTestSuite) Test2_Login() {
	body, err := tests.DoGraphQLAuth(s.gql, fmt.Sprintf(`mutation {
		loginUser(input: { email: "%s", password: "%s" }) {
			status message
			data { access_token refresh_token }
		}
	}`, s.email, s.password), "")
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "loginUser")
	s.Require().NotNil(data)
	accessToken, ok := data["access_token"].(string)
	s.Require().True(ok)
	s.Require().NotEmpty(accessToken)
	s.accessToken = accessToken
}

func (s *AuthHandlerApiTestSuite) Test3_GetMe() {
	s.Require().NotEmpty(s.accessToken)
	s.Require().NotZero(s.userID)

	body, err := tests.DoGraphQLAuth(s.gql, `query { getMe(input: { access_token: "" }) { status message data { id email } } }`, s.accessToken)
	s.Require().NoError(err)
	s.Require().Empty(tests.GraphQLErrorMessages(body), tests.GraphQLErrorMessages(body))
	data := tests.GraphQLOpEnvelopeData(body, "getMe")
	s.Require().NotNil(data)
	s.Equal(float64(s.userID), data["id"])
	s.Equal(s.email, data["email"])
}

func (s *AuthHandlerApiTestSuite) Test4_LoginLockout() {
	email := "locked.api@example.com"
	password := "wrongpassword"

	s.register(email, "correctpassword")

	// Fail login 5 times; each attempt must surface as an error.
	for i := 0; i < 5; i++ {
		body, err := tests.DoGraphQLAuth(s.gql, fmt.Sprintf(`mutation {
			loginUser(input: { email: "%s", password: "%s" }) {
				status message
				data { access_token }
			}
		}`, email, password), "")
		s.Require().NoError(err)
		graphqlErrors := tests.GraphQLErrorMessages(body)
		s.Require().NotEmpty(graphqlErrors, "failed login %d must return an error, got: %v", i+1, body)
	}

	// 6th attempt must be rejected because the account is locked.
	body, err := tests.DoGraphQLAuth(s.gql, fmt.Sprintf(`mutation {
		loginUser(input: { email: "%s", password: "%s" }) {
			status message
			data { access_token }
		}
	}`, email, password), "")
	s.Require().NoError(err)
	graphqlErrors := tests.GraphQLErrorMessages(body)
	s.Require().NotEmpty(graphqlErrors, "locked account must return an error, got: %v", body)
	s.Require().Contains(graphqlErrors[0], "locked", "expected an account-locked error, got: %s", graphqlErrors[0])
}

func TestAuthHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerApiTestSuite))
}
