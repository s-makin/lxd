package auth

import (
	"context"
	"net/http"

	"github.com/canonical/lxd/shared/api"
	"github.com/canonical/lxd/shared/entity"
)

const (
	// AuthenticationMethodCluster is set in the request context as request.CtxProtocol when the request is authenticated
	// via mTLS and the peer certificate is present in the trust store as type certificate.TypeServer.
	AuthenticationMethodCluster string = "cluster"

	// AuthenticationMethodUnix is set in the request context as request.CtxProtocol when the request is made over the
	// unix socket.
	AuthenticationMethodUnix string = "unix"

	// AuthenticationMethodPKI is set in the request context as request.CtxProtocol when a `server.ca` file exists in
	// LXD_DIR, the peer certificate of the request was signed by the CA file, and core.trust_ca_certificates is true.
	//
	// Note: If core.trust_ca_certificates is false, the peer certificate is additionally verified via mTLS and the
	// value of request.CtxProtocol is set to api.AuthenticationMethodTLS.
	//
	// Note: Regardless of whether `core.trust_ca_certificates` is enabled, we still check if the client certificate
	// fingerprint is in the identity cache. If they are found, standard TLS restrictions will apply.
	AuthenticationMethodPKI string = "pki"
)

// PermissionChecker is a type alias for a function that returns whether a user has required permissions on an object.
// It is returned by Authorizer.GetPermissionChecker.
type PermissionChecker func(entityURL *api.URL) bool

// Authorizer is the primary external API for this package.
type Authorizer interface {
	Driver() string

	CheckPermission(ctx context.Context, r *http.Request, entityURL *api.URL, entitlement Entitlement) error
	GetPermissionChecker(ctx context.Context, r *http.Request, entitlement Entitlement, entityType entity.Type) (PermissionChecker, error)
}

// IsDeniedError returns true if the error is not found or forbidden. This is because the CheckPermission method on
// Authorizer will return a not found error if the requestor does not have access to view the resource. If a requestor
// has view access, but not edit access a forbidden error is returned.
func IsDeniedError(err error) bool {
	return api.StatusErrorCheck(err, http.StatusNotFound, http.StatusForbidden)
}
