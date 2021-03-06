package service

import (
	"../control/EmailVerify"
	"../dao"
	"../models"
	"fmt"
	"net/http"
	"time"
	"net/url"
)

func SendCode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	emails, ok := r.PostForm["Email"]
	if !ok || len(emails) == 0 {
		models.SendRetJson2(0, "缺少Email参数", "手动滑稽", w)
		return
	}
	email := emails[0]
	fmt.Println(time.Now(), r.RemoteAddr, email, "请求验证码")
	if dao.ExistLogin(email) {
		models.SendRetJson2(0, "该邮箱已被注册", "", w)
		return
	}
	go EmailVerify.SendCode(email)
	models.SendRetJson2(1, "验证码已发送", "", w)
}

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	emails, ok := r.PostForm["Email"]
	if !ok || len(emails) == 0 {
		models.SendRetJson2(0, "缺少Email参数", "手动滑稽", w)
		return
	}
	email := emails[0]
	if email == "" {
		models.SendRetJson2(0, "缺少Email参数", "手动滑稽", w)
		return
	}
	pwds, ok := r.PostForm["Password"]
	if !ok || len(pwds) == 0 {
		models.SendRetJson2(0, "缺少Password参数", "手动滑稽", w)
		return
	}
	pwd := pwds[0]
	if pwd == "" {
		models.SendRetJson2(0, "缺少Password参数", "手动滑稽", w)
		return
	}
	codes, ok := r.PostForm["VerifyCode"]
	if !ok || len(codes) == 0 {
		models.SendRetJson2(0, "缺少VerifyCode参数", "手动滑稽", w)
		return
	}
	code := codes[0]
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"),
		r.RemoteAddr, "注册：", email, pwd, code)
	if code == "" {
		models.SendRetJson2(0, "注册失败", "请填写验证码", w)
		return
	}
	status, msg := EmailVerify.CheckCode(email, code)
	if !status {
		models.SendRetJson2(0, msg, "注册失败", w)
		return
	}
	if dao.ExistLogin(email) {
		models.SendRetJson2(0, "该邮箱已被注册", "", w)
		return
	}
	id, err := dao.AddUser(&models.User{Email: email}, pwd)
	if err != nil {
		models.SendRetJson2(0, "注册失败", err.Error(), w)
		return
	}
	go dao.AddRegisterRecord(id,time.Now(),r.RemoteAddr)
	models.SendRetJson2(1, "注册成功", id, w)
}

//注册数量
func RegisterRecordCount(w http.ResponseWriter, r *http.Request)  {
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
	c,err:=dao.RegisterRecordCount(d)
	if err != nil {
		models.SendRetJson2(0, "失败", err.Error(), w)
		return
	}
	models.SendRetJson2(1, "成功", c, w)
}
