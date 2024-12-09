package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"sync"
	"syscall/js"
	"time"
)

var (
	verifyCodeMap = make(map[string]string)
	mutex         sync.RWMutex
)

func main() {
	c := make(chan struct{}, 0)

	// 注册全局函数
	js.Global().Set("sendVerifyCode", js.FuncOf(sendVerifyCode))
	js.Global().Set("register", js.FuncOf(register))

	<-c
}

func sendVerifyCode(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	phoneInput := document.Call("getElementById", "phone").Get("value").String()

	// 验证手机号格式
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phoneInput)
	if !matched {
		js.Global().Call("alert", "请输入正确的手机号")
		return nil
	}

	// 生成验证码
	code := generateVerifyCode()

	// 存储验证码
	mutex.Lock()
	verifyCodeMap[phoneInput] = code
	mutex.Unlock()

	// 模拟发送短信
	go sendSMS(phoneInput, code)

	// 禁用发送按钮60秒
	btn := document.Call("getElementById", "sendCodeBtn")
	btn.Set("disabled", true)
	go countDown(btn)

	return nil
}

func register(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	phone := document.Call("getElementById", "phone").Get("value").String()
	code := document.Call("getElementById", "verifyCode").Get("value").String()

	mutex.RLock()
	savedCode, exists := verifyCodeMap[phone]
	mutex.RUnlock()

	if !exists {
		js.Global().Call("alert", "请先获取验证码")
		return nil
	}

	if code != savedCode {
		js.Global().Call("alert", "验证码错误")
		return nil
	}

	// 验证成功，清除验证码
	mutex.Lock()
	delete(verifyCodeMap, phone)
	mutex.Unlock()

	js.Global().Call("alert", "注册成功")
	return nil
}

func generateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

func sendSMS(phone, code string) {
	// 这里应该调用实际的短信发送API
	fmt.Printf("向手机号 %s 发送验证码: %s\n", phone, code)
}

func countDown(btn js.Value) {
	for i := 60; i > 0; i-- {
		time.Sleep(time.Second)
		js.Global().Get("document").Call("getElementById", "sendCodeBtn").Set("innerText", fmt.Sprintf("%d秒后重试", i))
	}
	btn.Set("disabled", false)
	js.Global().Get("document").Call("getElementById", "sendCodeBtn").Set("innerText", "发送验证码")
}
