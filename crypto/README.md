# 加密
* sha256
* ripemd
* base64
* base58
## 介绍
### sha256
***概述***  
***go语言实现***  
<code>h:=sha256.New()</code>  
<code>h.Write([]byte("crypto"))</code>  
<code>fmt.Printf("%x",h.Sum(nil))</code>
***
### ripemd
***概述***  
***go语言实现***
***
### base64
***概述***  
***go语言实现***  
>encoder:=base64.NewEncoder(base64.StdEncoding,os.Stdout)  
_, err := encoder.Write([]byte("crypto"))  
if err!=nil{  
    fmt.Println(err.Error())  
}  
*** 
### base58
***概述***  
在base64的基础上去掉了4个容易混淆的字符：'0','O','I','l'和两个字符'+','/' 
***go语言实现***
