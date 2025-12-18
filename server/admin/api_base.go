package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/xlzd/gotp"
)

// Login 登陆接口
func Login(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	adminUser := r.PostFormValue("admin_user")
	adminPass := r.PostFormValue("admin_pass")

	// 启用otp验证
	if base.Cfg.AdminOtp != "" {
		pwd := adminPass
		pl := len(pwd)
		if pl < 6 {
			RespError(w, RespUserOrPassErr)
			base.Error(adminUser, "管理员otp错误")
			return
		}
		// 判断otp信息
		adminPass = pwd[:pl-6]
		otp := pwd[pl-6:]

		totp := gotp.NewDefaultTOTP(base.Cfg.AdminOtp)
		unix := time.Now().Unix()
		verify := totp.Verify(otp, unix)

		if !verify {
			RespError(w, RespUserOrPassErr)
			base.Error(adminUser, "管理员otp错误")
			return
		}
	}

	// 认证错误
	if !(adminUser == base.Cfg.AdminUser &&
		utils.PasswordVerify(adminPass, base.Cfg.AdminPass)) {
		RespError(w, RespUserOrPassErr)
		base.Error(adminUser, "管理员用户名或密码错误")
		return
	}

	// token有效期
	expiresAt := time.Now().Unix() + 3600*3
	jwtData := map[string]interface{}{"admin_user": adminUser}
	tokenString, err := SetJwtData(jwtData, expiresAt)
	if err != nil {
		RespError(w, 1, err)
		return
	}

	data := make(map[string]interface{})
	data["token"] = tokenString
	data["admin_user"] = adminUser
	data["expires_at"] = expiresAt

	ck := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, ck)

	RespSucess(w, data)
}

// SetPwd 修改密码接口
func SetPwd(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		OldPass string `json:"old_pass"`
		NewPass string `json:"new_pass"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		RespError(w, RespInternalErr, "参数解析失败")
		return
	}

	if reqData.OldPass == "" || reqData.NewPass == "" {
		RespError(w, RespParamErr, "密码不能为空")
		return
	}

	if len(reqData.NewPass) < 6 {
		RespError(w, RespParamErr, "新密码长度不能少于6位")
		return
	}

	if !utils.PasswordVerify(reqData.OldPass, base.Cfg.AdminPass) {
		RespError(w, RespUserOrPassErr, "原密码错误")
		return
	}

	newPassHash, err := utils.PasswordHash(reqData.NewPass)
	if err != nil {
		RespError(w, RespInternalErr, "密码加密失败")
		return
	}

	base.Cfg.AdminPass = newPassHash

	base.Info("管理员密码已修改")
	RespSucess(w, nil)
}

func authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		route := mux.CurrentRoute(r)
		name := route.GetName()
		if utils.InArrStr([]string{"login", "index", "static"}, name) {
			next.ServeHTTP(w, r)
			return
		}

		jwtToken := r.Header.Get("Jwt")
		if jwtToken == "" {
			jwtToken = r.FormValue("jwt")
		}
		if jwtToken == "" {
			cc, err := r.Cookie("jwt")
			if err == nil {
				jwtToken = cc.Value
			}
		}
		data, err := GetJwtData(jwtToken)
		if err != nil || base.Cfg.AdminUser != fmt.Sprint(data["admin_user"]) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func recoverHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				base.Error(err, string(stack))
				RespError(w, 500, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}