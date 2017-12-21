package user

import (
	"fmt"
	"github.com/go-macaron/captcha"
	"github.com/hani17/chtq/models"
	"github.com/hani17/chtq/models/errors"
	"github.com/hani17/chtq/pkg/config"
	"github.com/hani17/chtq/pkg/context"
	"github.com/hani17/chtq/pkg/form"

	"net/http"
)

const (
	SIGNUP = "user/auth/signup"
	LOGIN  = "user/auth/login"
)

func Signup(c *context.Context) {
	c.Data["Title"] = "Sign Up"

	isSucceed, err := AutoLogin(c)
	if err != nil {
		c.Handle(500, "AutoLogin", err)
	}

	if isSucceed {
		c.Redirect("/")
	}

	c.HTML(http.StatusOK, SIGNUP)
}

func SignUpPost(c *context.Context, cpt *captcha.Captcha, f form.Register) {
	c.Data["Title"] = "Sign Up"
	if !cpt.VerifyReq(c.Req) {
		c.Data["Err_Captcha"] = true
		c.Data["ErrorMsg"] = "Captcha didn't match."
		SetSignUpFormData(c, &f)
		c.HTML(http.StatusOK, SIGNUP)
		return
	}

	u := &models.User{
		UserName: f.UserName,
		Email:    f.Email,
		Password: f.Password,
	}

	if err := models.CreateUser(u); err != nil {
		switch {
		case models.IsErrUserAlreadyExist(err):
			c.Data["Err_UserName"] = true
			c.Data["ErrorMsg"] = "Username is taken."
			SetSignUpFormData(c, &f)
			c.HTML(http.StatusOK, SIGNUP)
		case models.IsErrEmailAlreadyUsed(err):
			c.Data["Err_Email"] = true
			c.Data["ErrorMsg"] = "Email is Already used."
			SetSignUpFormData(c, &f)
			c.HTML(http.StatusOK, SIGNUP)
		default:
			fmt.Println("Error Error")
			c.HTML(http.StatusInternalServerError, SIGNUP)
		}
		return
	}
	c.Redirect("/user/login")
}

//SetSignUpFormData sets signup form data if signing failed
func SetSignUpFormData(c *context.Context, f *form.Register) {
	c.Data["user_name"] = f.UserName
	c.Data["email"] = f.Email
	c.Data["password"] = f.Password
}

func SetSignInFormData(c *context.Context, f *form.SignIn) {
	c.Data["user_name"] = f.UserName
}

func LoginPost(c *context.Context, f form.SignIn) {
	c.Data["Title"] = "Sign In"

	u, err := models.UserSignIn(f.UserName, f.Password)
	if err != nil {
		if errors.IsUserNotExist(err) {
			c.Data["Err_Username_Password_Incorrect"] = true
			c.Data["ErrorMsg"] = "Username or Password is not correct"
			c.HTML(http.StatusOK, LOGIN)
		} else {
			c.ServerError("UserSignIn", err)
		}
		return
	}

	afterLogin(c, u, f.Remember)

}

func Login(c *context.Context) {
	c.Data["Title"] = "Login"

	isSucceed, err := AutoLogin(c)
	if err != nil {
		c.Handle(500, "AutoLogin", err)
		return
	}

	if isSucceed {
		c.Redirect("/")
	}
	c.HTML(200, LOGIN)
}

func AutoLogin(c *context.Context) (bool, error) {
	username := c.GetCookie(config.CookieUsername)
	if len(username) == 0 {
		return false, nil
	}

	isSucceed := false
	defer func() {
		if !isSucceed {
			//log.Trace("auto-login cookie cleared: %s", uname)
			c.SetCookie(config.CookieUsername, "", -1)
			c.SetCookie(config.CookieRemember, "", -1)
			//c.SetCookie(setting.LoginStatusCookieName, "", -1)
		}
	}()

	u, err := models.GetUserByName(username)
	if err != nil {
		if !errors.IsUserNotExist(err) {
			return false, fmt.Errorf("GetUserByName: %v", err)
		}
		return false, nil
	}

	if val, ok := c.GetSuperSecureCookie(u.Password, config.CookieRemember); !ok || val != u.UserName {
		return false, nil
	}

	isSucceed = true
	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.UserName)
	return true, nil
}

func afterLogin(c *context.Context, u *models.User, remember bool) {
	if remember {
		days := 86400 * 7
		c.SetCookie(config.CookieUsername, u.UserName, days, "/", "", config.CookieSecure, true)
		c.SetSuperSecureCookie(u.Password, config.CookieRemember, u.UserName, days, "/", "", config.CookieSecure, true)
	}
	c.Session.Set("uid", u.ID)
	c.Session.Set("uname", u.UserName)

	c.Redirect("/")
}

func Signout(c *context.Context) {
	c.Session.Delete("uid")
	c.Session.Delete("uname")
	c.SetCookie(config.CookieUsername, "", -1, "/")
	c.SetCookie(config.CookieRemember, "", -1, "/")
	//c.SetCookie(config.CSRFCookieName, "", -1, setting.AppSubURL)
	c.Redirect("/")
}
