package main

import (
  "fmt"
  "strings"
  "net/http"
  "net/url"

  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/clientcredentials"

  oidc "github.com/coreos/go-oidc"
  "github.com/gin-gonic/gin"
  "github.com/atarantini/ginrequestid"

  "golang-cp-be/config"
  "golang-cp-be/environment"
  //"golang-cp-be/gateway/hydra"
  "golang-cp-be/authorizations"
)

const app = "cpbe"

func init() {
  config.InitConfigurations()
}

func main() {

  provider, err := oidc.NewProvider(context.Background(), config.Hydra.Url + "/")
  if err != nil {
    environment.DebugLog(app, "main", "[provider:hydra] " + err.Error(), "")
    return
  }

  // Setup the hydra client cpbe is going to use (oauth2 client credentials flow)
  hydraConfig := &clientcredentials.Config{
    ClientID:     config.CpBe.ClientId,
    ClientSecret: config.CpBe.ClientSecret,
    TokenURL:     config.Hydra.TokenUrl,
    Scopes:       config.CpBe.RequiredScopes,
    EndpointParams: url.Values{"audience": {"hydra"}},
    AuthStyle: 2, // https://godoc.org/golang.org/x/oauth2#AuthStyle
  }

  // Setup app state variables. Can be used in handler functions by doing closures see exchangeAuthorizationCodeCallback
  env := &environment.State{
    Provider: provider,
    HydraConfig: hydraConfig,
  }

  // Setup routes to use, this defines log for debug log
  routes := map[string]environment.Route{
    "/authorizations": environment.Route{
       URL: "/authorizations",
       LogId: "cpbe://authorizations",
    },
    "/authorizations/authorize": environment.Route{
      URL: "/authorizations/authorize",
      LogId: "cpfe://authorizations/authorize",
    },
    "/authorizations/reject": environment.Route{
      URL: "/authorizations/reject",
      LogId: "cpfe://authorizations/reject",
    },
  }

  r := gin.Default()
  r.Use(ginrequestid.RequestId())

  // ## QTNA - Questions that need answering before granting access to a protected resource
  // 1. Is the user or client authenticated? Answered by the process of obtaining an access token.
  // 2. Is the access token expired?
  // 3. Is the access token granted the required scopes?
  // 4. Is the user or client giving the grants in the access token authorized to operate the scopes granted?
  // 5. Is the access token revoked?

  // All requests need to be authenticated.
  r.Use(authenticationRequired())

  r.GET(routes["/authorizations"].URL, authorizationRequired(routes["/authorizations"], "cpbe.authorizations.get"), authorizations.GetCollection(env, routes["/authorizations"]))
  r.POST(routes["/authorizations"].URL, authorizationRequired(routes["/authorizations"], "cpbe.authorizations.post"), authorizations.PostCollection(env, routes["/authorizations"]))
  r.PUT(routes["/authorizations"].URL, authorizationRequired(routes["/authorizations"], "cpbe.authorizations.update"), authorizations.PutCollection(env, routes["/authorizations"]))

  r.POST(routes["/authorizations/authorize"].URL, authorizationRequired(routes["/authorizations/authorize"], "cpbe.authorize"), authorizations.PostAuthorize(env, routes["/authorizations/authorize"]))
  r.POST(routes["/authorizations/reject"].URL, authorizationRequired(routes["/authorizations/reject"], "cpbe.reject"), authorizations.PostAuthorize(env, routes["/authorizations/reject"]))

  r.RunTLS(":" + config.Self.Port, "/srv/certs/cpbe-cert.pem", "/srv/certs/cpbe-key.pem")
}

func authenticationRequired() gin.HandlerFunc {
  fn := func(c *gin.Context) {
    requestId := c.MustGet(environment.RequestIdKey).(string)
    environment.DebugLog(app, "authenticationRequired", "Checking Authorization: Bearer <token> in request", requestId)

    var token *oauth2.Token
    auth := c.Request.Header.Get("Authorization")
    split := strings.SplitN(auth, " ", 2)
    if len(split) == 2 || strings.EqualFold(split[0], "bearer") {
      environment.DebugLog(app, "authenticationRequired", "Authorization: Bearer <token> found for request.", requestId)
      token = &oauth2.Token{
        AccessToken: split[1],
        TokenType: split[0],
      }

      if token.Valid() == true {
        environment.DebugLog(app, "authenticationRequired", "Valid access token", requestId)
        c.Set(environment.AccessTokenKey, token)
        c.Next() // Authentication successful, continue.
        return;
      }

      // Deny by default
      environment.DebugLog(app, "authenticationRequired", "Invalid Access token", requestId)
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token."})
      c.Abort()
      return
    }

    // Deny by default
    environment.DebugLog(app, "authenticationRequired", "Missing access token", requestId)
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization: Bearer <token> not found in request."})
    c.Abort()
  }
  return gin.HandlerFunc(fn)
}

func authorizationRequired(route environment.Route, requiredScopes ...string) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    requestId := c.MustGet(environment.RequestIdKey).(string)
    environment.DebugLog(app, "authorizationRequired", "Checking Authorization: Bearer <token> in request", requestId)

    accessToken, accessTokenExists := c.Get(environment.AccessTokenKey)
    if accessTokenExists == false {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found. Hint: Is bearer token missing?"})
      c.Abort()
      return
    }

    // Sanity check: Claims
    fmt.Println(accessToken)

    foundRequiredScopes := true
    if foundRequiredScopes {
      environment.DebugLog(app, "authorizationRequired", "Valid scopes. WE DID NOT CHECK IT - TODO!", requestId)
      c.Next() // Authentication successful, continue.
      return;
    }

    // Deny by default
    environment.DebugLog(app, "authorizationRequired", "Missing required scopes: ", requestId)
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing required scopes: "})
    c.Abort()
  }
  return gin.HandlerFunc(fn)
}
