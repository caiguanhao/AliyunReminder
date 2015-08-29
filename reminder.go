package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type GetAvailbleInstances struct {
	Code string `json:"code"`
	Data []struct {
		AliUID        int    `json:"aliUID"`
		AutoRenewal   bool   `json:"autoRenewal"`
		Bid           string `json:"bid"`
		BillingTag    string `json:"billingTag"`
		BuyerID       int    `json:"buyerID"`
		ChargeType    string `json:"chargeType"`
		CommodityCode string `json:"commodityCode"`
		CreatedTime   int    `json:"createdTime"`
		EndTime       int    `json:"endTime"`
		ExtraEndTime  int    `json:"extraEndTime"`
		GmtCreated    int    `json:"gmtCreated"`
		InstanceID    int    `json:"instanceID"`
		InstanceName  string `json:"instanceName"`
		Meta          struct {
			AutoRenewButtonStatus string `json:"autoRenewButtonStatus"`
			AutoRenewDuration     string `json:"autoRenewDuration"`
			AutoRenewNextTime     string `json:"autoRenewNextTime"`
			AutoRenewStatus       string `json:"autoRenewStatus"`
			EndDate               string `json:"endDate"`
			EndTime               string `json:"endTime"`
			HostName              string `json:"hostName"`
			InstanceName          string `json:"instanceName"`
			InternetIP            string `json:"internetIp"`
			IntranetIP            string `json:"intranetIp"`
			RegionName            string `json:"regionName"`
			RemainTime            string `json:"remainTime"`
		} `json:"meta"`
		PropInfo struct {
			AliUID            int    `json:"aliUid"`
			Bid               string `json:"bid"`
			BusinessStatus    string `json:"businessStatus"`
			Cores             int    `json:"cores"`
			Description       string `json:"description"`
			DiskSize          int    `json:"diskSize"`
			EcsBusinessStatus string `json:"ecsBusinessStatus"`
			GmtCreated        int    `json:"gmtCreated"`
			GmtModified       int    `json:"gmtModified"`
			GmtStarted        int    `json:"gmtStarted"`
			GmtSync           int    `json:"gmtSync"`
			GroupNo           string `json:"groupNo"`
			Hostname          string `json:"hostname"`
			ID                int    `json:"id"`
			ImageID           int    `json:"imageId"`
			ImageName         string `json:"imageName"`
			ImageNo           string `json:"imageNo"`
			ImageType         string `json:"imageType"`
			InstanceID        string `json:"instanceId"`
			InternetIP        string `json:"internetIp"`
			InternetRx        int    `json:"internetRx"`
			InternetTx        int    `json:"internetTx"`
			IntranetIP        string `json:"intranetIp"`
			IntranetRx        int    `json:"intranetRx"`
			IntranetTx        int    `json:"intranetTx"`
			IoOptimized       bool   `json:"ioOptimized"`
			IsWin             bool   `json:"isWin"`
			Iz                struct {
				CnName string `json:"cnName"`
				EnName string `json:"enName"`
				No     string `json:"no"`
			} `json:"iz"`
			IzID              int    `json:"izId"`
			Memory            int    `json:"memory"`
			NetWorkType       string `json:"netWorkType"`
			NetworkValidation bool   `json:"networkValidation"`
			OsName            string `json:"osName"`
			OsType            string `json:"osType"`
			RealHostname      string `json:"realHostname"`
			RecoverPolicy     string `json:"recoverPolicy"`
			Recoverable       bool   `json:"recoverable"`
			Region            struct {
				ID            int    `json:"id"`
				IsActive      string `json:"isActive"`
				RegionEnName  string `json:"regionEnName"`
				RegionName    string `json:"regionName"`
				RegionNo      string `json:"regionNo"`
				RegionNoAlias string `json:"regionNoAlias"`
			} `json:"region"`
			RegionID             int    `json:"regionId"`
			SerialNumber         string `json:"serialNumber"`
			Status               string `json:"status"`
			SystemDeviceCategory string `json:"systemDeviceCategory"`
			SystemDiskCategory   string `json:"systemDiskCategory"`
			VswitchInstanceID    string `json:"vswitchInstanceId"`
			Zone                 struct {
				ID          int    `json:"id"`
				IsActive    string `json:"isActive"`
				Writable    bool   `json:"writable"`
				ZoneName    string `json:"zoneName"`
				ZoneNo      string `json:"zoneNo"`
				ZoneNoAlias string `json:"zoneNoAlias"`
			} `json:"zone"`
			ZoneID int `json:"zoneId"`
		} `json:"propInfo"`
		ProviderID     int    `json:"providerID"`
		Region         string `json:"region"`
		ReleaseTime    int    `json:"releaseTime"`
		RemainTimeMeta struct {
			Day    int `json:"day"`
			Hour   int `json:"hour"`
			Minute int `json:"minute"`
		} `json:"remainTimeMeta"`
		RenewalDuration int    `json:"renewalDuration"`
		ResCreateTime   int    `json:"resCreateTime"`
		ResourceID      int    `json:"resourceID"`
		ResourceStatus  string `json:"resourceStatus"`
		ResourceType    string `json:"resourceType"`
	} `json:"data"`
	PageInfo struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		Total       int `json:"total"`
	} `json:"pageInfo"`
	RequestID       string `json:"requestId"`
	SuccessResponse bool   `json:"successResponse"`
}

type FlowdockRequest struct {
	UserName string `json:"external_user_name"`
	Content  string `json:"content"`
}

func getStatus() (*string, error) {
	url := "https://renew.console.aliyun.com/renew/getAvailbleInstances.json"
	url += "?commodityCode=vm&currentPage=1&pageSize=50"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", "login_aliyunid_ticket="+string(LOGIN_ALIYUNID_TICKET))

	/*
		*.console.aliyun.com has a relatively worse certificate,
		so we need to specify cipher and min tls version, for more info, see:
		https://www.ssllabs.com/ssltest/analyze.html?d=renew.console.aliyun.com&hideResults=on
	*/
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
				MinVersion:   tls.VersionTLS10,
			},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	instances := GetAvailbleInstances{}
	err = json.Unmarshal(body, &instances)

	if err != nil {
		return nil, err
	}

	if instances.Code != "200" || !instances.SuccessResponse {
		return nil, errors.New(string(body))
	}

	const format = "    %-10s %s\n"

	ret := fmt.Sprintf(format, "Server", "Expires In")

	for _, instance := range instances.Data {
		ret += fmt.Sprintf(format, instance.Meta.HostName, instance.Meta.RemainTime)
	}

	ret += "https://renew.console.aliyun.com/#/ecs"

	return &ret, nil
}

func notifyFlowdock(content *string) error {
	url := "https://api.flowdock.com/v1/messages/chat/" + string(FLOWDOCK_TOKEN)

	body, err := json.Marshal(&FlowdockRequest{
		UserName: "AliyunReminder",
		Content:  *content,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, respErr := ioutil.ReadAll(resp.Body)

	if respErr != nil {
		return respErr
	}

	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}

	return nil
}

var status *string

func keepAlive(duration time.Duration) {
	for {
		log.Println("Getting status...")
		var err error
		status, err = getStatus()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Status:\n" + *status)
		}
		log.Println("Wait for 30 minutes...")
		time.Sleep(duration)
	}
}

func notify(when string) {
	expected, err := time.Parse("15:04:05", when)
	if err != nil {
		log.Fatalln("Error date:", when)
	}
	log.Println("Will send status to flowdock at", when)
	h, m, s := expected.Hour(), expected.Minute(), expected.Second()
	for {
		time.Sleep(1 * time.Second)
		now := time.Now()
		if now.Hour() != h || now.Minute() != m || now.Second() != s {
			continue
		}
		if status == nil {
			log.Println("Status is empty! Can't send it to flowdock")
			continue
		}
		log.Println("Sending status to flowdock...")
		err := notifyFlowdock(status)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Successfully sent to flowdock")
		}
	}
}

func init() {
	flag.Usage = func() {
		fmt.Println("AliyunReminder - Reminds you when your ECS expires.")
		fmt.Println()
		fmt.Println("Built with Flowdock token", string(FLOWDOCK_TOKEN[:6]), string(MADE))
		fmt.Println("Source: https://github.com/caiguanhao/AliyunReminder")
	}
	flag.Parse()
}

func main() {
	go notify("09:00:00")
	keepAlive(30 * time.Minute)
}
