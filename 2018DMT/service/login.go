package service

import (
	"../control"
	"../models"
	"../global"
	"../dao"
	"fmt"
	"net/http"
	"time"
	"net/url"
)

func Login(w http.ResponseWriter, r *http.Request) {
	msg := "账号与密码不匹配"
	r.ParseForm()
	emails, ok := r.PostForm["Email"]
	if !ok || len(emails) == 0 {
		models.SendRetJson2(0, "缺少Email参数", "手动滑稽", w)
		return
	}
	email := emails[0]
	pwds, ok := r.PostForm["Password"]
	if !ok || len(pwds) == 0 {
		models.SendRetJson2(0, "缺少Password参数", "手动滑稽", w)
		return
	}
	pwd := pwds[0]
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"),
		r.RemoteAddr, "登录：", email, pwd)
	if email == "" || pwd == "" {
		models.SendRetJson2(0, "用户名或密码不能为空", "手动滑稽", w)
		return
	}
	lr, user, err := control.CheckLogin(email, pwd)
	if lr == 1 {
		if err != nil {
			models.SendRetJson2(lr, "获取用户信息失败", err.Error(), w)
			return
		}
		msg = "登录成功"
		cookie := http.Cookie{
			Name:   "user",
			MaxAge: global.MaxCookieTime,
			Value:  dao.GenUserCookie(user.IdToString()),
		}
		http.SetCookie(w, &cookie)
		models.SendRetJson2(lr, msg, user, w)
		go dao.AddLoginRecord(user.Id,lr,time.Now(),r.RemoteAddr)
		return
	}
	models.SendRetJson2(lr, msg, "手动滑稽", w)
	go dao.AddLoginRecord(-1,lr,time.Now(),r.RemoteAddr)
}

func IsLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		models.SendRetJson2(0, "未登录", "", w)
		return
	}
	uid, t, err := dao.GetUserIdFromCookie(cookie.Value)
	dt := int(time.Now().Sub(t).Seconds())
	if dt > global.MaxCookieTime {
		models.SendRetJson2(1, "登录过期", uid, w)
		return
	}
	user, err := dao.GetUserById(uid)
	if err != nil {
		models.SendRetJson2(1, global.NoSuchUser.Error(),
			uid, w)
		return
	}
	models.SendRetJson2(1,
		t.Format("2006-01-02/15:04:05"),
		user, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		models.SendRetJson2(0, "未登录", "", w)
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	models.SendRetJson2(1, "退出成功", "(●'◡'●)", w)
}

//登录数量
func LoginRecordCount(w http.ResponseWriter, r *http.Request)  {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		models.SendRetJson2(0, "失败", err.Error(), w)
		return
	}
	d,err:=GetGetInt("Day",queryForm)
	if err != nil {
		models.SendRetJson2(0, "失败", err.Error(), w)
		return
	}
	c,err:=dao.LoginRecordCount(d)
	if err != nil {
		models.SendRetJson2(0, "失败", err.Error(), w)
		return
	}
	models.SendRetJson2(1, "成功", c, w)
}
