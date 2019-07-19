package authorizations

import (
  "net/http"
  "strings"
  "fmt"

  "github.com/gin-gonic/gin"
  "golang-cp-be/environment"
  "golang-cp-be/gateway/cpbe"
  //"golang-cp-be/gateway/hydra"
)

type ConsentRequest struct {
  Subject string `json:"sub" binding:"required"`
  App string `json:"app" binding:"required"`
  ClientId string `json:"client_id,omitempty"`
  GrantedScopes []string `json:"granted_scopes,omitempty"`
  RevokedScopes []string `json:"revoked_scopes,omitempty"`
  RequestedScopes []string `json:"requested_scopes,omitempty"`
}

type ConsentResponse struct {

}

func GetCollection(env *environment.State, route environment.Route) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    requestId := c.MustGet(environment.RequestIdKey).(string)
    environment.DebugLog(route.LogId, "GetCollection", "", requestId)

    id, _ := c.GetQuery("id")
    if id == "" {
      c.JSON(http.StatusNotFound, gin.H{
        "error": "Not found. Hint: Are you missing id in request?",
      })
      c.Abort()
      return;
    }

    app, _ := c.GetQuery("app")
    if app == "" {
      c.JSON(http.StatusNotFound, gin.H{
        "error": "Not found. Hint: Are you missing app in request?",
      })
      c.Abort()
      return;
    }

    clientId, _ := c.GetQuery("client_id")
    if clientId == "" {
      c.JSON(http.StatusNotFound, gin.H{
        "error": "Not found. Hint: Are you missing client_id in request?",
      })
      c.Abort()
      return;
    }

    var permissions []cpbe.Permission
    requestedScopes, _ := c.GetQuery("scope")
    if requestedScopes != "" {
      scopes := strings.Split(requestedScopes, ",")
      for _, scope := range scopes {
        permissions = append(permissions, cpbe.Permission{ Name:scope,})
      }
    }

    identity := cpbe.Identity{
      Subject: id,
    }
    application := cpbe.App{
      Name: app,
    }
    applicationIdentity := cpbe.Identity{
      Subject: clientId,
    }
    permissionList, err := cpbe.FetchConsentsForIdentityToApplication(env.Driver, identity, application, applicationIdentity, permissions)
    if err == nil {

      var grantedPermissions []string
      for _, permission := range permissionList {
        grantedPermissions = append(grantedPermissions, permission.Name)
      }

      c.JSON(http.StatusOK, grantedPermissions)
      return
    }

    // Deny by default
    c.JSON(http.StatusNotFound, gin.H{
      "error": "Not found",
    })
    c.Abort()
  }
  return gin.HandlerFunc(fn)
}

func PostCollection(env *environment.State, route environment.Route) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    requestId := c.MustGet(environment.RequestIdKey).(string)
    environment.DebugLog(route.LogId, "PostCollection", "", requestId)

    var input ConsentRequest
    err := c.BindJSON(&input)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      c.Abort()
      return
    }

    if len(input.RequestedScopes) <= 0 {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Missing granted_scopes"})
      c.Abort()
      return
    }

    var grantPermissions []cpbe.Permission
    for _, scope := range input.GrantedScopes {
      grantPermissions = append(grantPermissions, cpbe.Permission{ Name:scope,})
    }

    var revokePermissions []cpbe.Permission
    for _, scope := range input.RevokedScopes {
      revokePermissions = append(revokePermissions, cpbe.Permission{ Name:scope,})
    }

    identity := cpbe.Identity{
      Subject: input.Subject,
    }
    application := cpbe.App{
      Name: input.App,
    }
    applicationIdentity := cpbe.Identity{
      Subject: input.ClientId,
    }
    permissionList, err := cpbe.CreateConsentsForIdentityToApplication(env.Driver, identity, application, applicationIdentity, grantPermissions, revokePermissions)
    if err != nil {
      fmt.Println(err)
    }
    if err == nil {

      var grantedPermissions []string
      for _, permission := range permissionList {
        grantedPermissions = append(grantedPermissions, permission.Name)
      }

      c.JSON(http.StatusOK, grantedPermissions)
      return
    }

    // Deny by default
    c.JSON(http.StatusNotFound, gin.H{
      "error": "Not found",
    })
    c.Abort()
  }
  return gin.HandlerFunc(fn)
}

func PutCollection(env *environment.State, route environment.Route) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    requestId := c.MustGet(environment.RequestIdKey).(string)
    environment.DebugLog(route.LogId, "PutCollection", "", requestId)

    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  }
  return gin.HandlerFunc(fn)
}
