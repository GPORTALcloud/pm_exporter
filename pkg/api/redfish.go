package api

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/GPORTALcloud/pm_exporter/pkg/metric"
)

type RedfishApi struct {
	auth string
	host string
}

func NewRedfishAPI(host string) *RedfishApi {
	r := RedfishApi{host: host}
	return &r
}

// Set redfish login credentials
func (r *RedfishApi) SetUser(user string, pass string) *RedfishApi {
	r.auth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))
	return r
}

// HTTP request wrapper adds authentication headers
// and parses the DellSystemMetric result
// TODO: rename
func (r *RedfishApi) Request(path string) (*DellSystemMetric, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := fmt.Sprintf("%s%s", r.host, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error building request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", r.auth))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Printf("request failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	obj := DellSystemMetric{}
	json.Unmarshal(body, &obj)
	return &obj, nil
}

func (r *RedfishApi) GetHost() string {
	return r.host
}

// utility function, converts "OK" and "" to 1
// note: if there is no power supply health check will result in 1
// only if state is critical / != null and != OK it will result in 0
func checkOK(ok string) int {
	if ok == "OK" || ok == "" {
		return 1
	}
	return 0
}

// FetchInventoryMetrics is used for populating metrics to the collector
// metrics.UpdateMetric stores those within the metric module
// no need to return anything except the error here
// also, meta metric "pm_platform_management_up" is used for indicating if
// the platform management is reachable or not
func (r *RedfishApi) FetchInventoryMetrics() error {
	response, err := r.Request("/redfish/v1/Dell/Systems/System.Embedded.1/DellSystem/System.Embedded.1")
	if err != nil {
		metric.UpdateMetric(metric.PMPlatformManagementUp, r.host, 0)
		return err
	}
	metric.UpdateMetric(metric.PMPlatformManagementUp, r.host, 1)
	metric.UpdateMetric(metric.PMPowerSupplyHealth, r.host, checkOK(response.Psrollupstatus))
	metric.UpdateMetric(metric.PMBatteryHealth, r.host, checkOK(response.Batteryrollupstatus))
	metric.UpdateMetric(metric.PMCpuHealth, r.host, checkOK(response.Cpurollupstatus))
	metric.UpdateMetric(metric.PMFanHealth, r.host, checkOK(response.Fanrollupstatus))
	metric.UpdateMetric(metric.PMStorageHealth, r.host, checkOK(response.Storagerollupstatus))
	metric.UpdateMetric(metric.PMTemperatureHealth, r.host, checkOK(response.Temprollupstatus))
	metric.UpdateMetric(metric.PMIntrusionHealth, r.host, checkOK(response.Intrusionrollupstatus))
	metric.UpdateMetric(metric.PMLicenceHealth, r.host, checkOK(response.Licensingrollupstatus))
	metric.UpdateMetric(metric.PMMemoryHealth, r.host, checkOK(response.Sysmemprimarystatus))

	// summarized metric across all health variables
	overallHealth := checkOK(response.Psrollupstatus) *
		checkOK(response.Batteryrollupstatus) *
		checkOK(response.Cpurollupstatus) *
		checkOK(response.Fanrollupstatus) *
		checkOK(response.Storagerollupstatus) *
		checkOK(response.Temprollupstatus) *
		checkOK(response.Intrusionrollupstatus) *
		checkOK(response.Licensingrollupstatus)
	metric.UpdateMetric(metric.PMOverallHealth, r.host, overallHealth)
	return nil
}

// Struct for the system metric endpoint
type DellSystemMetric struct {
	OdataContext                       string      `json:"@odata.context"`
	OdataID                            string      `json:"@odata.id"`
	OdataType                          string      `json:"@odata.type"`
	Biosreleasedate                    string      `json:"BIOSReleaseDate"`
	Baseboardchassisslot               string      `json:"BaseBoardChassisSlot"`
	Batteryrollupstatus                string      `json:"BatteryRollupStatus"`
	Bladegeometry                      string      `json:"BladeGeometry"`
	Cmcip                              interface{} `json:"CMCIP"`
	Cpurollupstatus                    string      `json:"CPURollupStatus"`
	Chassismodel                       string      `json:"ChassisModel"`
	Chassisname                        string      `json:"ChassisName"`
	Chassisservicetag                  string      `json:"ChassisServiceTag"`
	Chassissystemheightunit            int         `json:"ChassisSystemHeightUnit"`
	Currentrollupstatus                string      `json:"CurrentRollupStatus"`
	Description                        string      `json:"Description"`
	Estimatedexhausttemperaturecelsius int         `json:"EstimatedExhaustTemperatureCelsius"`
	Estimatedsystemairflowcfm          int         `json:"EstimatedSystemAirflowCFM"`
	Expressservicecode                 string      `json:"ExpressServiceCode"`
	Fanrollupstatus                    string      `json:"FanRollupStatus"`
	Idsdmrollupstatus                  interface{} `json:"IDSDMRollupStatus"`
	ID                                 string      `json:"Id"`
	Intrusionrollupstatus              string      `json:"IntrusionRollupStatus"`
	Isoembranded                       string      `json:"IsOEMBranded"`
	Lastsysteminventorytime            time.Time   `json:"LastSystemInventoryTime"`
	Lastupdatetime                     time.Time   `json:"LastUpdateTime"`
	Licensingrollupstatus              string      `json:"LicensingRollupStatus"`
	Maxcpusockets                      int         `json:"MaxCPUSockets"`
	Maxdimmslots                       int         `json:"MaxDIMMSlots"`
	Maxpcieslots                       int         `json:"MaxPCIeSlots"`
	Memoryoperationmode                string      `json:"MemoryOperationMode"`
	Name                               string      `json:"Name"`
	Nodeid                             string      `json:"NodeID"`
	Psrollupstatus                     string      `json:"PSRollupStatus"`
	Populateddimmslots                 int         `json:"PopulatedDIMMSlots"`
	Populatedpcieslots                 int         `json:"PopulatedPCIeSlots"`
	Powercapenabledstate               string      `json:"PowerCapEnabledState"`
	Sdcardrollupstatus                 interface{} `json:"SDCardRollupStatus"`
	Selrollupstatus                    string      `json:"SELRollupStatus"`
	Serverallocationwatts              interface{} `json:"ServerAllocationWatts"`
	Storagerollupstatus                string      `json:"StorageRollupStatus"`
	Sysmemerrormethodology             string      `json:"SysMemErrorMethodology"`
	Sysmemfailoverstate                string      `json:"SysMemFailOverState"`
	Sysmemlocation                     string      `json:"SysMemLocation"`
	Sysmemprimarystatus                string      `json:"SysMemPrimaryStatus"`
	Systemgeneration                   string      `json:"SystemGeneration"`
	Systemid                           int         `json:"SystemID"`
	Systemrevision                     string      `json:"SystemRevision"`
	Temprollupstatus                   string      `json:"TempRollupStatus"`
	Tempstatisticsrollupstatus         string      `json:"TempStatisticsRollupStatus"`
	UUID                               string      `json:"UUID"`
	Voltrollupstatus                   string      `json:"VoltRollupStatus"`
	Smbiosguid                         string      `json:"smbiosGUID"`
}
