// Package primitive contains definitions of the primitive types used
// in ag.
package primitive

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/arigatomachine/cli/base64"
	"github.com/arigatomachine/cli/identity"
	"github.com/arigatomachine/cli/pathexp"
)

// v1Schema embeds in other structs to indicate their schema version is 1.
type v1Schema struct{}

// Version returns the schema version of structs that embed this type.
func (v *v1Schema) Version() int {
	return 1
}

// User is the body of a user object
type User struct {
	v1Schema
	Username string        `json:"username"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	State    string        `json:"state"`
	Password *UserPassword `json:"password"`
	Master   *UserMaster   `json:"master"`
}

// UserMaster is the body.master object for a user
type UserMaster struct {
	Value *base64.Value `json:"value"`
	Alg   string        `json:"alg"`
}

// UserPassword is the body.password object for a user
type UserPassword struct {
	Salt  string        `json:"salt"`
	Value *base64.Value `json:"value"`
	Alg   string        `json:"alg"`
}

// Type returns the enumerated byte representation of User.
func (u *User) Type() byte {
	return byte(0x01)
}

// Signature is an immutable object, but not technically a payload. Its fields
// must be ordered properly so that ID generation is correct.
//
// If PublicKeyID is nil, the signature is self-signed.
type Signature struct {
	Algorithm   string        `json:"alg"`
	PublicKeyID *identity.ID  `json:"public_key_id"`
	Value       *base64.Value `json:"value"`
}

// Immutable object payloads. Their fields must be lexicographically ordered by
// the json value, so we can correctly calculate the signature.

// PrivateKeyValue holds the encrypted value of the PrivateKey.
type PrivateKeyValue struct {
	Algorithm string        `json:"alg"`
	Value     *base64.Value `json:"value"`
}

// PrivateKey is the private portion of an asymetric key.
type PrivateKey struct {
	v1Schema
	Key         PrivateKeyValue `json:"key"`
	OrgID       *identity.ID    `json:"org_id"`
	OwnerID     *identity.ID    `json:"owner_id"`
	PNonce      *base64.Value   `json:"pnonce"`
	PublicKeyID *identity.ID    `json:"public_key_id"`
}

// Type returns the enumerated byte representation of PrivateKey.
func (pk *PrivateKey) Type() byte {
	return byte(0x07)
}

// PublicKeyValue is the actual value of a PublicKey.
type PublicKeyValue struct {
	Value *base64.Value `json:"value"`
}

// PublicKey is the public portion of an asymetric key.
type PublicKey struct {
	v1Schema
	Algorithm string         `json:"alg"`
	Created   time.Time      `json:"created_at"`
	Expires   time.Time      `json:"expires_at"`
	Key       PublicKeyValue `json:"key"`
	OrgID     *identity.ID   `json:"org_id"`
	OwnerID   *identity.ID   `json:"owner_id"`
	KeyType   string         `json:"type"`
}

// Type returns the enumerated byte representation of PublicKey.
func (pk *PublicKey) Type() byte {
	return byte(0x06)
}

// Types of claims that can be made against public keys.
const (
	SignatureClaimType  = "signature"
	RevocationClaimType = "revocation"
)

// Claim is a signature or revocation claim against a public key.
type Claim struct {
	v1Schema
	Created     time.Time    `json:"created_at"`
	OrgID       *identity.ID `json:"org_id"`
	OwnerID     *identity.ID `json:"owner_id"`
	Previous    *identity.ID `json:"previous"`
	PublicKeyID *identity.ID `json:"public_key_id"`
	KeyType     string       `json:"type"`
}

// Type returns the enumerated byte representation of Claim.
func (c *Claim) Type() byte {
	return byte(0x08)
}

// NewClaim returns a new Claim, with the created time set to now
func NewClaim(orgID, ownerID, previous, pubKeyID *identity.ID,
	keyType string) *Claim {
	return &Claim{
		OrgID:       orgID,
		OwnerID:     ownerID,
		Previous:    previous,
		PublicKeyID: pubKeyID,
		KeyType:     keyType,
		Created:     time.Now().UTC(),
	}
}

// Credential is a secret value shared between a group of services based
// on users identity, operating environment, project, and organization
type Credential struct {
	v1Schema
	Credential        *CredentialValue `json:"credential"`
	KeyringID         *identity.ID     `json:"keyring_id"`
	Name              string           `json:"name"`
	Nonce             *base64.Value    `json:"nonce"`
	OrgID             *identity.ID     `json:"org_id"`
	PathExp           *pathexp.PathExp `json:"pathexp"`
	Previous          *identity.ID     `json:"previous"`
	ProjectID         *identity.ID     `json:"project_id"`
	CredentialVersion int              `json:"version"`
}

// Type returns the enumerated byte representation of Credential
func (c *Credential) Type() byte {
	return byte(0xb)
}

// CredentialValue is the secretbox encrypted value of the containing
// Credential.
type CredentialValue struct {
	Algorithm string        `json:"alg"`
	Nonce     *base64.Value `json:"nonce"`
	Value     *base64.Value `json:"value"`
}

// Keyring is a mechanism for sharing a shared secret between many different
// users and machines at a position in the credential path.
//
// Credentials belong to Keyrings
type Keyring struct {
	v1Schema
	Created        time.Time        `json:"created_at"`
	OrgID          *identity.ID     `json:"org_id"`
	PathExp        *pathexp.PathExp `json:"pathexp"`
	Previous       *identity.ID     `json:"previous"`
	ProjectID      *identity.ID     `json:"project_id"`
	KeyringVersion int              `json:"version"`
}

// Type returns the enumerated byte representation of Keyring
func (k *Keyring) Type() byte {
	return byte(0x09)
}

// KeyringMember is a record of sharing a master secret key with a user or
// machine.
//
// KeyringMember belongs to a Keyring
type KeyringMember struct {
	v1Schema
	Created         time.Time         `json:"created_at"`
	EncryptingKeyID *identity.ID      `json:"encrypting_key_id"`
	Key             *KeyringMemberKey `json:"key"`
	KeyringID       *identity.ID      `json:"keyring_id"`
	OrgID           *identity.ID      `json:"org_id"`
	OwnerID         *identity.ID      `json:"owner_id"`
	ProjectID       *identity.ID      `json:"project_id"`
	PublicKeyID     *identity.ID      `json:"public_key_id"`
}

// Type returns the enumerated byte representation of KeyringMember
func (km *KeyringMember) Type() byte {
	return byte(0x0a)
}

// KeyringMemberKey is the keyring master encryption key, encrypted for the
// owner of a KeyringMember
type KeyringMemberKey struct {
	Algorithm string        `json:"alg"`
	Nonce     *base64.Value `json:"nonce"`
	Value     *base64.Value `json:"value"`
}

// Org is a grouping of users that collaborate with each other
type Org struct {
	v1Schema
	Name string `json:"name"`
}

// Type returns the enumerated byte representation of Org
func (o *Org) Type() byte {
	return byte(0xd)
}

// Org Invitations exist in four states: pending, associated,
// accepted, and approved.
const (
	OrgInvitePendingState    = "pending"
	OrgInviteAssociatedState = "associated"
	OrgInviteAcceptedState   = "accepted"
	OrgInviteApprovedState   = "approved"
)

// OrgInvite is an invitation for an individual to join an organization
type OrgInvite struct {
	v1Schema
	OrgID      *identity.ID `json:"org_id"`
	Email      string       `json:"email"`
	InviterID  *identity.ID `json:"inviter_id"`
	InviteeID  *identity.ID `json:"invitee_id"`
	ApproverID *identity.ID `json:"approver_id"`
	State      string       `json:"state"`
	Code       *struct {
		Alg   string        `json:"alg"`
		Salt  *base64.Value `json:"salt"`
		Value *base64.Value `json:"value"`
	} `json:"code"`
	PendingTeams []identity.ID `json:"pending_teams"`
	Created      *time.Time    `json:"created_at"`
	Accepted     *time.Time    `json:"accepted_at"`
	Approved     *time.Time    `json:"approved_at"`
}

// Type returns the numerated byte representation of OrgInvite
func (o *OrgInvite) Type() byte {
	return byte(0x13)
}

// Project is an entity that represents a group of services
type Project struct {
	v1Schema
	Name  string       `json:"name"`
	OrgID *identity.ID `json:"org_id"`
}

// Type returns the enumerated byte representation of Project
func (t *Project) Type() byte {
	return byte(0x04)
}

// Policy is an entity that represents a group of statements for acl
type Policy struct {
	v1Schema
	PolicyType string       `json:"type"`
	Previous   *identity.ID `json:"previous"`
	OrgID      *identity.ID `json:"org_id"`
	Policy     struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Statements  []PolicyStatement `json:"statements"`
	} `json:"policy"`
}

// Type returns the enumerated byte representation of Policy
func (t *Policy) Type() byte {
	return byte(0x11)
}

// PolicyStatement is an acl statement on a policy object
type PolicyStatement struct {
	Effect   PolicyEffect `json:"effect"`
	Action   PolicyAction `json:"action"`
	Resource string       `json:"resource"`
}

// PolicyEffect is the effect type of the statement (allow or deny)
type PolicyEffect bool

// These are the two policy effect types
const (
	PolicyEffectAllow = true
	PolicyEffectDeny  = false
)

// MarshalText implements the encoding.TextMarshaler interface, used for JSON
// marshaling.
func (pe *PolicyEffect) MarshalText() ([]byte, error) {
	if *pe {
		return []byte("allow"), nil
	}
	return []byte("deny"), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface, used for
// JSON unmarshaling.
func (pe *PolicyEffect) UnmarshalText(b []byte) error {
	*pe = string(b) == "allow"
	return nil
}

// String returns a string representation of the PolicyEffect (allow or deny)
func (pe *PolicyEffect) String() string {
	b, _ := pe.MarshalText()
	return string(b)
}

// PolicyAction represents the user actions that are covered by a statement.
type PolicyAction byte

// These are all the possible PolicyActions
const (
	PolicyActionCreate = 1 << iota
	PolicyActionRead
	PolicyActionUpdate
	PolicyActionDelete
	PolicyActionList
)

var policyActionStrings = []string{
	"create",
	"read",
	"update",
	"delete",
	"list",
}

// MarshalJSON implements the json.Marshaler interface. A PolicyAction is
// encoded in JSON either the string representations of its actions in a list,
// or a single string when there is only one action.
func (pa *PolicyAction) MarshalJSON() ([]byte, error) {
	out := []string{}

	for i, v := range policyActionStrings {
		if (1<<uint(i))&byte(*pa) > 0 {
			out = append(out, v)
		}
	}

	if len(out) == 1 {
		return json.Marshal(out[0])
	}

	return json.Marshal(out)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (pa *PolicyAction) UnmarshalJSON(b []byte) error {
	raw := make([]string, 1)

	var err error
	if len(b) > 0 && b[0] == '"' {
		err = json.Unmarshal(b, &raw[0])
	} else {
		err = json.Unmarshal(b, &raw)
	}

	if err != nil {
		return err
	}

	for _, a := range raw {
		for i, v := range policyActionStrings {
			if a == v {
				*pa |= 1 << uint(i)
				continue
			}
		}
	}

	return nil
}

func (pa *PolicyAction) String() string {
	out := []string{}

	for i, v := range policyActionStrings {
		if (1<<uint(i))&byte(*pa) > 0 {
			out = append(out, v)
		}
	}

	return strings.Join(out, ", ")
}

// PolicyAttachment is an entity that represents the link between policies and teams
type PolicyAttachment struct {
	v1Schema
	OwnerID  *identity.ID `json:"owner_id"`
	PolicyID *identity.ID `json:"policy_id"`
	OrgID    *identity.ID `json:"org_id"`
}

// Type returns the enumerated byte representation of PolicyAttchment
func (t *PolicyAttachment) Type() byte {
	return byte(0x12)
}

// Service is an entity that represents a group of processes
type Service struct {
	v1Schema
	Name      string       `json:"name"`
	OrgID     *identity.ID `json:"org_id"`
	ProjectID *identity.ID `json:"project_id"`
}

// Type returns the enumerated byte representation of Service
func (t *Service) Type() byte {
	return byte(0x03)
}

// Environment is an entity that represents a group of processes
type Environment struct {
	v1Schema
	Name      string       `json:"name"`
	OrgID     *identity.ID `json:"org_id"`
	ProjectID *identity.ID `json:"project_id"`
}

// Type returns the enumerated byte representation of Environment
func (t *Environment) Type() byte {
	return byte(0x05)
}

// There are two types of teams: system and user. System teams are
// managed by the Arigato registry.
const (
	SystemTeam = "system"
	UserTeam   = "user"
)

// Team is an entity that represents a group of users
type Team struct {
	v1Schema
	Name     string       `json:"name"`
	OrgID    *identity.ID `json:"org_id"`
	TeamType string       `json:"type"`
}

// Type returns the enumerated byte representation of Team
func (t *Team) Type() byte {
	return byte(0x0f)
}

// Membership is an entity that represents whether a user or
// machine is a part of a team in an organization.
type Membership struct {
	v1Schema
	OrgID   *identity.ID `json:"org_id"`
	OwnerID *identity.ID `json:"owner_id"`
	TeamID  *identity.ID `json:"team_id"`
}

// Type returns the enumerated byte representation of Membership
func (m *Membership) Type() byte {
	return byte(0x0e)
}
