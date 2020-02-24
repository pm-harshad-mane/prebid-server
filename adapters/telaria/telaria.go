package telaria

import(
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

const (
	REQUEST_METHOD = "GET"
	ERR_URI_CONF = "Incorrect Telaria request URI %s, please check the configuration."
	ERR_NO_IMPS = "Telaria:No impressions found in the bid request"
	ERR_NO_VALID_IMPS = "Telaria:No valid impression in the bid request"
	ERR_NO_VID_MEDIA_TYPE = "Telaria: only supports video media type. Ignoring imp id=%s"
	ERR_EXT_NOT_FOUND = "Telaria: ext.bidder not found"
	ERR_EXT_BIDDER_PARSE_FAIL = "Telaria: Failed to parse ext.bidder.publisher"
	ERR_UNEXPEXCTED_ERR_CODE = "Telaria: Unexpected response status code: %d. Run with request.debug = 1 for more info"
	ERR_BAD_SERVER = "Telaria: bad server response: %d."
)

type TelariaAdapter struct {
	http 	*adapters.HTTPAdapter
	URI 	string
	testing bool
	httpHeaders *http.Header // using common headers per request as headers will not change per request
}

func NewTelariaBidder(client *http.Client, endpointURL string) *TelariaAdapter {	
	if _, err := url.Parse(endpointURL), err != nil {
		panic(fmt.Sprintf(ERR_URI_CONF, endpointURL))
	}

	return &TelariaAdapter{
		http: &adapters.HTTPAdapter{Client: client},
		URI: endpoint,
		testing: false,
		httpHeaders: &http.Header{}
	}
}

func (adapter *TelariaAdapter) MakeRequests (request *openrtb.BidRequest, reqInfo *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	errs := make([]error, 0, len(request.Imp))

	if len(request.Imp) == 0 {
		err := &errortypes.BadInput{
			Message: ERR_NO_IMPS,
		}
		errs = append(errs, err)
		return nil, errs
	}
	
	// using common headers per request as headers will not change per request
	adapter.buildCommonHttpHeaders(request)
	
	var adapterRequests []*adapters.RequestData 
	const IMPRESSION_COUNT_IN_REQUEST = len(request.Imp)
	countOfInvalidImp := 0

	for _, imp := range request.Imp {
		if imp.Video == nil {
			err := &errortypes.BadInput{
				Message: fmt.Sprintf(ERR_NO_VID_MEDIA_TYPE, imp.ID),
			}
			errs = append(errs, err)
			countOfInvalidImp++
		} else {
			adapterReq, errors := adapter.buildSingleRequestData(request, imp)
			if adapterReq != nil {
				adapterRequests = append(adapterRequests, adapterReq)
			}
			errs = append(errs, errors...)
		}
	}

	if IMPRESSION_COUNT_IN_REQUEST == countOfInvalidImp {
		err := &errortypes.BadInput{
			Message: ERR_NO_VALID_IMPS,
		}
		errs = append(errs, err)
		return nil, errs
	}

	return adapterRequests, errs
}

func (adapter *TelariaAdapter) buildCommonHttpHeaders(request *openrtb.BidRequest) {	
	adapter.httpHeaders = &http.Header{} // always reset
	// todo: add the headers to const
	adapter.httpHeaders.Add("Accept", "*/*")
	adapter.httpHeaders.Add("Connection", "keep-alive")
	adapter.httpHeaders.Add("cache-control", "no-cache")
	adapter.httpHeaders.Add("Accept-Encoding", "gzip, deflate")
	if request.Device != nil {
		addHeaderIfNotEmpty(adapter.httpHeaders, "User-Agent", request.Device.UA)
		addHeaderIfNotEmpty(adapter.httpHeaders, "X-Forwarded-For", request.Device.IP)
		addHeaderIfNotEmpty(adapter.httpHeaders, "Accept-Language", request.Device.Language)
		if request.Device.DNT != nil {
			addHeaderIfNotEmpty(adapter.httpHeaders, "DNT", strconv.Itoa(int(*request.Device.DNT)))
		}
	}
}

func addHeaderIfNotEmpty(headers http.Header, headerName string, headerValue string) {
	if len(headerValue) > 0 {
		headers.Add(headerName, headerValue)
	}
}

func (adapter *TelariaAdapter) buildSingleRequestData(request *openrtb.BidRequest, imp *openrtb.Imp) (*adapters.RequestData, []error) {
	var errors []error
	var bidderExt adapters.ExtImpBidder

	err := json.Unmarshal(imp.Ext, &bidderExt)
	if err != nil {
		err = &errortypes.BadInput{
			Message: ERR_EXT_NOT_FOUND,
		}
		errors = append(errors, err)
		return nil, errors
	}

	var telariaExt openrtb_ext.ExtImpTelaria
	err = json.Unmarshal(bidderExt.Bidder, &telariaExt)
	if err != nil {
		err = &errortypes.BadInput{
			Message: ERR_EXT_BIDDER_PARSE_FAIL,
		}
		errors = append(errors, err)
		return nil, errors
	}

	err = checkParams(telariaExt)
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}

	return &adapters.RequestData{
		Method:  REQUEST_METHOD,
		Uri:     adapter.buildSingleRequestURI(request, imp, telariaExt),
		Headers: adapter.httpHeaders,
	}, errors
}

func checkParams(telariaExt openrtb_ext.ExtImpTelaria) error {
	return nil
}

func (adapter *TelariaAdapter) buildSingleRequestURI(request *openrtb.BidRequest, imp *openrtb.Imp, telariaExt *openrtb_ext.ExtImpTelaria) string {
	var uri *url.URL
	uri, _ = url.Parse(adapter.URI)
	parameters := url.Values{}
	parameters.Add("CC", "1")
	uri.RawQuery = parameters.Encode()
	return uri.String()
}

func (adapter *TelariaAdapter) MakeBids (internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if response.StatusCode == http.StatusBadRequest || response.StatusCode != http.StatusOK {
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf(ERR_UNEXPEXCTED_ERR_CODE, response.StatusCode),
		}}
	}

	var bidResp openrtb.BidResponse

	if err := json.Unmarshal(response.Body, &bidResp); err != nil {
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf(ERR_BAD_SERVER, err),
		}}
	}

	bidResponse := adapters.NewBidderResponseWithBidsCapacity(len(bidResp.SeatBid[0].Bid))
	for _, sb := range bidResp.SeatBid {
		for i := range sb.Bid {
			bidResponse.Bids = append(bidResponse.Bids, &adapters.TypedBid{
				Bid:     &sb.Bid[i],
				BidType: openrtb_ext.BidTypeVideo,
			})
		}
	}

	return bidResponse, nil
}