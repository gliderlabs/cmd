package console

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gliderlabs/cmd/lib/stripe"
	"github.com/gliderlabs/cmd/lib/web"
	"github.com/gliderlabs/cmd/pkg/auth0"
	"github.com/gliderlabs/comlab/pkg/events"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/mitchellh/mapstructure"
	stripelib "github.com/stripe/stripe-go"
)

type User struct {
	Name     string
	Nickname string
	Email    string
	Picture  string
	Account  AppMetadata `mapstructure:"app_metadata"`
	ID       string      `mapstructure:"user_id"`
}

type Account struct {
	CustomerID     string `mapstructure:"customer_id"`
	SubscriptionID string `mapstructure:"subscription_id"`
	Plan           string
}

type AppMetadata struct {
	Account `mapstructure:",squash"`
	Groups  map[string]Account
	Invites Invites
}

type Invites struct {
	Pending   []string
	InvitedBy string `mapstructure:"invited_by"`
}

func SessionUser(r *http.Request) *User {
	uid := web.SessionValue(r, "user_id")
	if uid == "" {
		return nil
	}
	var err error
	user, err := LookupUser(uid)
	if err != nil {
		log.Info(r, err, log.Fields{"uid": uid})
		return nil
	}
	return &user
}

func ContextUser(ctx context.Context) *User {
	u := ctx.Value("user")
	if u == nil {
		return nil
	}
	user, ok := u.(*User)
	if !ok {
		return nil
	}
	return user
}

func LookupUser(uid string) (User, error) {
	data, err := auth0.DefaultClient().User(uid)
	if err != nil {
		return User{}, err
	}
	var user User
	err = mapstructure.Decode(data, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func LookupNickname(nickname string) (User, error) {
	users, err := auth0.DefaultClient().SearchUsers(
		fmt.Sprintf("nickname:%s", nickname))
	if err != nil {
		return User{}, err
	}
	if len(users) < 1 {
		return User{}, fmt.Errorf("nickname not found")
	}
	var user User
	err = mapstructure.Decode(users[0], &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func RegisterUser(user *User) error {
	if user.Account.CustomerID != "" {
		return fmt.Errorf("user already registered")
	}
	if user.Email == "" {
		return fmt.Errorf("email required for registration")
	}
	params := &stripelib.CustomerParams{Email: user.Email}
	params.AddMeta("uid", user.ID)
	params.AddMeta("service", "cmd.io")
	customer, err := stripe.Client().Customers.New(params)
	if err != nil {
		return err
	}
	err = auth0.DefaultClient().PatchUser(user.ID, auth0.User{
		"app_metadata": map[string]interface{}{
			"customer_id": customer.ID,
		},
	})
	if err != nil {
		return err
	}
	events.Emit(events.Signal(EventFirstLogin))
	return nil
}
