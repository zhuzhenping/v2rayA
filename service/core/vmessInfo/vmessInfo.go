package vmessInfo

import (
	"encoding/base64"
	"fmt"
	"github.com/json-iterator/go"
	"net/url"
	"reflect"
	"v2rayA/common"
)

type VmessInfo struct {
	Ps            string `json:"ps"`
	Add           string `json:"add"`
	Port          string `json:"port"`
	ID            string `json:"id"`
	Aid           string `json:"aid"`
	Net           string `json:"net"`
	Type          string `json:"type"`
	Host          string `json:"host"`
	Path          string `json:"path"`
	TLS           string `json:"tls"`
	V             string `json:"v"`
	AllowInsecure bool   `json:"allowInsecure"`
	Protocol      string `json:"protocol"`
}

func (v *VmessInfo) ExportToURL() string {
	switch v.Protocol {
	case "", "vmess":
		//去除info中的protocol，减少URL体积
		it := reflect.TypeOf(*v)
		iv := reflect.ValueOf(*v)
		m := make(map[string]interface{})
		for i := 0; i < it.NumField(); i++ {
			f := it.Field(i)
			chKey := f.Tag.Get("json")
			if chKey == "protocol" { //不转换protocol
				continue
			}
			m[chKey] = iv.FieldByName(f.Name).Interface()
		}
		b, _ := jsoniter.Marshal(m)
		return "vmess://" + base64.URLEncoding.EncodeToString(b)
	case "ss":
		/* ss://BASE64(method:password)@server:port#name */
		nameField := ""
		if v.Ps != "" {
			nameField = "#" + v.Ps
		}
		return fmt.Sprintf(
			"ss://%v@%v:%v%v",
			base64.URLEncoding.EncodeToString([]byte(v.Net+":"+v.ID)),
			v.Add,
			v.Port,
			nameField,
		)
	case "ssr":
		/* ssr://server:port:proto:method:obfs:URLBASE64(password)/?remarks=URLBASE64(remarks)&protoparam=URLBASE64(protoparam)&obfsparam=URLBASE64(obfsparam)) */
		return fmt.Sprintf("ssr://%v", base64.URLEncoding.EncodeToString([]byte(
			fmt.Sprintf(
				"%v:%v:%v:%v:%v:%v/?remarks=%v&protoparam=%v&obfsparam=%v",
				v.Add,
				v.Port,
				v.Type,
				v.Net,
				v.TLS,
				base64.URLEncoding.EncodeToString([]byte(v.ID)),
				base64.URLEncoding.EncodeToString([]byte(v.Ps)),
				base64.URLEncoding.EncodeToString([]byte(v.Host)),
				base64.URLEncoding.EncodeToString([]byte(v.Path)),
			),
		)))
	case "pingtunnel":
		// pingtunnel://server:URLBASE64(passwd)#URLBASE64(remarks)
		return fmt.Sprintf("pingtunnel://%v", base64.URLEncoding.EncodeToString([]byte(
			fmt.Sprintf("%v:%v#%v",
				v.Add,
				base64.URLEncoding.EncodeToString([]byte(v.ID)),
				common.UrlEncoded(v.Ps),
			),
		)))
	case "trojan":
		// trojan://passwd@server:port#URLESCAPE(remarks)
		u, _ := url.Parse(fmt.Sprintf(
			"trojan://%v@%v:%v",
			v.ID,
			v.Add,
			v.Port,
		))
		u.Fragment = v.Ps
		q := u.Query()
		if v.Host != "" {
			q.Set("peer", v.Host)
		}
		if v.AllowInsecure {
			q.Set("allowInsecure", "1")
		}
		u.RawQuery = q.Encode()
		return u.String()
	}
	return ""
}
