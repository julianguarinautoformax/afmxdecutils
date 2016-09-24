package afmxdec

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type AutoFMXImageMetaInfo struct {
	OriginatorHardwareAddress       net.HardwareAddr
	OriginatorTimeStampUTC          time.Time
	OriginatorHardwareAddressString string
	OriginatorTimeStampUTCString    string
}

func ParseAutoFMXImageMetaInfo(jsonString []byte) (*AutoFMXImageMetaInfo, error) {

	type parsero_ struct {
		HwAdd string `json:"hwAdd"`
		TsUtc string `json:"tsUtc"`
	}
	var o_ parsero_
	json.Unmarshal(jsonString, &o_)

	/* Convert abcdef010204 to ab:cd:ef:01:02:04 */
	hwAdd_ := string(o_.HwAdd[0])
	for i := 1; i < len(o_.HwAdd); i++ {

		if i%2 == 0 {
			hwAdd_ += ":"
		}
		hwAdd_ += string(o_.HwAdd[i])
	}

	afmxMetaInfo := AutoFMXImageMetaInfo{}

	/* UTC yy.MM.dd.HH.mm.ss.zzz*/
	afmxMetaInfo.OriginatorTimeStampUTCString = o_.TsUtc
	Y, _ := strconv.Atoi(o_.TsUtc[0:2])
	M, _ := strconv.Atoi(o_.TsUtc[2:4])
	D, _ := strconv.Atoi(o_.TsUtc[4:6])
	H, _ := strconv.Atoi(o_.TsUtc[6:8])
	m, _ := strconv.Atoi(o_.TsUtc[8:10])
	s, _ := strconv.Atoi(o_.TsUtc[10:12])
	n, _ := strconv.Atoi(o_.TsUtc[12:15])
	afmxMetaInfo.OriginatorTimeStampUTC = time.Date(Y+2000, time.Month(M), D, H, m, s, n*1000000, time.UTC)

	/* HwAddr */
	if hwAdd, err := net.ParseMAC(hwAdd_); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	} else {
		afmxMetaInfo.OriginatorHardwareAddress = hwAdd
		afmxMetaInfo.OriginatorHardwareAddressString = o_.HwAdd
	}
	fmt.Println("Strings:[", afmxMetaInfo.OriginatorHardwareAddressString, "][", afmxMetaInfo.OriginatorTimeStampUTCString, "]")
	fmt.Println("Native Types:[", afmxMetaInfo.OriginatorHardwareAddress, "][", afmxMetaInfo.OriginatorTimeStampUTC, "]")
	return &afmxMetaInfo, nil

}
func (metaInfo AutoFMXImageMetaInfo) ParseIdentificactionString() (string, error) {
	return metaInfo.OriginatorHardwareAddressString + "/" + metaInfo.OriginatorTimeStampUTCString, nil
} /*
func main() {

	json_byte := []byte(`{"hwAdd":"ABCDEF012345","tsUtc":"160921134511320"}`)
	a, _ := parseAutoFMXImageMetaInfo(json_byte)
	fmt.Println(a.OriginatorHardwareAddress, a.OriginatorHardwareAddressString, a.OriginatorTimeStampUTC, a.OriginatorTimeStampUTCString)
}*/
