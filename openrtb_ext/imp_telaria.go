package openrtb_ext

import(
	"net/url"
	"reflect"
)

// ExtImpTelaria defines the contract for bidrequest.imp[i].ext.telaria
type ExtImpTelaria struct {
	AdCode 			string `json:"adCode"`
	SupplyCode 		string `json:"supplyCode",omitempty`
	MediaId     	string `json:"mediaId",omitempty`
	MediaUrl 		string `json:"mediaUrl",omitempty`
	MediaTitle		string `json:"mediaTitle",omitempty`
	ContentLength		string `json:"contentLength",omitempty`
	Floor			string `json:"floor",omitempty`
	Efloor			string `json:"efloor",omitempty`
	Custom			string `json:"custom",omitempty`
	Categories		string `json:"categories",omitempty`
	Keywords		string `json:"keywords",omitempty`
	BlockDomains	string `json:"blockDomains",omitempty`
	C2				string `json:"c2",omitempty`
	C3				string `json:"c3",omitempty`
	C4				string `json:"c4",omitempty`
	Skip			string `json:"skip",omitempty`
	SkipMin			string `json:"skipmin",omitempty`
	SkipAfter		string `json:"skipafter",omitempty`
	Delivery		string `json:"delivery",omitempty`
	Placement		string `json:"placement",omitempty`
	VideoMinBitrate		string `json:"videoMinBitrate",omitempty`
	VideoMaxBitrate		string `json:"videoMaxBitrate",omitempty`
	IncIdSync		bool `json:"incIdSync",omitempty`
}

func (tExt *ExtImpTelaria) GenerateParams(params *url.Values) (*url.Values) {
	TelariaParamsMap := map[string]string {
		"AdCode": "adCode",
	}

	reflectedValue := reflect.ValueOf(s)
	typeOfExtImpTelaria := reflectedValue.Type()
	// todo: use range	
	for i := 0; i< reflectedValue.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v %s\n", typeOfExtImpTelaria.Field(i).Name, reflectedValue.Field(i).Interface(), TelariaParamsMap["AdCode"])
	}

}

// refer: https://console.telaria.com/examples/hb/headerbidding.jsp

// params: {
//     supplyCode: 'ssp-demo-rm6rh',
//     adCode: 'ssp-!demo!-lufip',
//     mediaId: 'MyCoolVideo',
//     mediaUrl: '',
//     mediaTitle: '',
//     contentLength: '',
//     floor: '',
//     efloor: '',
//     custom: '',
//     categories: '',
//     keywords: '',
//     blockDomains: '',
//     c2: '',
//     c3: '',
//     c4: '',
//     skip: '',
//     skipmin: '',
//     skipafter: '',
//     delivery: '',
//     placement: '',
//     videoMinBitrate: '',
//     videoMaxBitrate: '',
//     incIdSync: true,
//     gdpr: '',
//     gdpr_consent: ''
// }