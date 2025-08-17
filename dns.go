// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"fmt"
	"os"
	"strings"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Record struct {
	Name   string
	Remark string
}

type RecordGroup struct {
	GroupId   int64
	GroupName string
	Records   []Record
}

var groups = [...]RecordGroup{
	{GroupId: 119833, GroupName: "VPS"},
	{GroupId: 119835, GroupName: "K3S"},
	{GroupId: 119834, GroupName: "K8S"},
	{GroupId: 119832, GroupName: "CNAME"},
}

// createClient of alibaba cloud api
func createClient() (_result *alidns20150109.Client, _err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Alidns
	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

// DomainGroupWithRecords list records of domain
func DomainGroupWithRecords(domain string, w func(text string)) []RecordGroup {
	var gs []RecordGroup
	for _, g := range groups {
		grs, err := describeDomainGroupRecords(domain, g.GroupId)
		if err != nil {
			w(fmt.Sprintf("describeDomainGroupRecords error: %s <br>", err.Error()))
			continue
		}
		var rs []Record
		for _, r := range grs {
			if r.Remark == nil || strings.HasPrefix(*r.RR, "_") {
				continue
			}
			rs = append(rs, Record{
				Name:   "`" + *r.RR + "`",
				Remark: *r.Remark,
			})
		}
		gs = append(gs, RecordGroup{
			GroupId:   g.GroupId,
			GroupName: g.GroupName,
			Records:   rs,
		})
	}
	return gs
}

// describeDomainGroupRecords list domain records of group
func describeDomainGroupRecords(domain string, groupId int64) ([]*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {
	client, _err := createClient()
	if _err != nil {
		return nil, _err
	}

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
		GroupId:    tea.Int64(groupId),
	}
	runtime := &util.RuntimeOptions{}
	// 复制代码运行请自行打印 API 的返回值
	_rs, _err := client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, runtime)
	if _err != nil {
		return nil, _err
	}
	return _rs.Body.DomainRecords.Record, nil
}
