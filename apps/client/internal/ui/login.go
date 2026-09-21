package ui

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/api"
	"github.com/Yogesh-Kumar-Mallik-dev/campus-os/apps/client/internal/auth"
)

// LoginView encapsulates the state and UI rendering for institutional user authentication.
type LoginView struct {
	client              *api.Client
	window              fyne.Window
	sessionStore        auth.SessionStore
	content             *fyne.Container
	onSuccess           func(session *auth.AuthSession)
	onNavigateToOnboard func()

	StatusText string
	IsError    bool
}

// BLOCK_UI_LOGIN_NEW_001
// Purpose: Constructs a new LoginView instance using native shadcn primitives.
func NewLoginView(
	client *api.Client,
	window fyne.Window,
	sessionStore auth.SessionStore,
	onSuccess func(session *auth.AuthSession),
	onNavigateToOnboard func(),
) *LoginView {
	v := &LoginView{
		client:              client,
		window:              window,
		sessionStore:        sessionStore,
		content:             container.NewVBox(),
		onSuccess:           onSuccess,
		onNavigateToOnboard: onNavigateToOnboard,
	}
	v.render()
	return v
}

// CanvasObject returns the renderable canvas object.
func (v *LoginView) CanvasObject() fyne.CanvasObject {
	return v.content
}

func (v *LoginView) render() {
	v.content.Objects = nil

	// Status / Error Banner
	if v.StatusText != "" {
		variant := AlertSuccess
		svgIcon := ResourceFromSVG("check.svg", SVGCheckVerified)
		title := "Authentication Success"
		if v.IsError {
			variant = AlertDestructive
			svgIcon = ResourceFromSVG("alert.svg", SVGClockGrace)
			title = "Authentication Error"
		}
		alert := NewShadcnAlert(title, v.StatusText, variant, svgIcon)
		v.content.Add(alert)
	}

	identField, identEntry := NewFormField(FormField{
		Label:       "Institutional Identifier",
		Placeholder: "Username, official email, or PRN (e.g. chairperson.2024)",
		HelperText:  "Assigned institutional account username or university PRN",
	})

	passField, passEntry := NewFormField(FormField{
		Label:       "Access Password",
		Placeholder: "Enter your account password",
		HelperText:  "Secured via argon2id cryptographic hash verification",
		IsPassword:  true,
	})

	loginBtn := NewShadcnButton("Sign In to Institution", ButtonDefault, ButtonSizeDefault, WhiteResourceFromSVG("login.svg", LucideLogIn), func() {
		ident := identEntry.Text
		pass := passEntry.Text
		if ident == "" || pass == "" {
			v.IsError = true
			v.StatusText = "Please enter both identifier and password"
			if v.window != nil {
				ShowToast(v.window, "Missing Credentials", "Identifier and password are required.", AlertDestructive, 3*time.Second)
			}
			v.render()
			return
		}

		go func() {
			resp, err := v.client.Login(context.Background(), ident, pass, "")
			if err == nil {
				role := "USER"
				if len(resp.RoleCodes) > 0 {
					role = resp.RoleCodes[0]
				}
				sess := &auth.AuthSession{
					AccessToken:  resp.AccessToken,
					RefreshToken: resp.RefreshToken,
					UserID:       resp.UserID,
					Username:     resp.Username,
					FullName:     resp.FullName,
					RoleCode:     role,
					RoleCodes:    resp.RoleCodes,
				}
				if v.sessionStore != nil {
					_ = v.sessionStore.Save(sess)
				}
				v.IsError = false
				v.StatusText = "Authentication successful"
				if v.onSuccess != nil {
					v.onSuccess(sess)
				}
			} else if api.IsUnreachable(err) {
				// Offline simulated login for local testing
				sess := &auth.AuthSession{
					AccessToken:  "mock_jwt_access_offline",
					RefreshToken: "mock_jwt_refresh_offline",
					UserID:       "u-mock-offline",
					Username:     ident,
					FullName:     "Offline User",
					RoleCode:     "SUPER_ADMIN",
					RoleCodes:    []string{"SUPER_ADMIN"},
				}
				if v.sessionStore != nil {
					_ = v.sessionStore.Save(sess)
				}
				v.IsError = false
				v.StatusText = "Authenticated (Offline Demo Mode)"
				if v.onSuccess != nil {
					v.onSuccess(sess)
				}
			} else {
				v.IsError = true
				v.StatusText = err.Error()
				if v.window != nil {
					ShowToast(v.window, "Sign In Failed", err.Error(), AlertDestructive, 3*time.Second)
				}
				v.render()
			}
		}()
	})

	onboardBtn := NewShadcnButton("Activate with Sealed QR Docket", ButtonOutline, ButtonSizeDefault, ResourceFromSVG("qr.svg", LucideQrCode), func() {
		if v.onNavigateToOnboard != nil {
			v.onNavigateToOnboard()
		}
	})

	loginCard := NewShadcnCard(CardParts{
		Badge:       NewBadge("INSTITUTIONAL ACCESS", BadgeDefault, BadgeShapePill),
		Title:       "Sign In to Campus OS",
		Description: "Enter your institutional credentials or activate a new sealed admission / executive docket",
		Content: container.NewVBox(
			identField,
			passField,
		),
		Footer: container.NewVBox(
			loginBtn,
			NewShadcnSeparator(true),
			onboardBtn,
		),
	})

	v.content.Add(loginCard)
}
