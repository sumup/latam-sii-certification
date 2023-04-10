package xmlschema

import (
	"encoding/xml"
)

const (
	httpSoapEnveloper = "http://schemas.xmlsoap.org/soap/envelope/"
	httpSoapEnc       = "http://schemas.xmlsoap.org/soap/encoding/"
)

type GetSeedBody struct {
	XMLName         xml.Name               `xml:"SOAP-ENV:Envelope"`
	XMLNSSoapEnv    string                 `xml:"xmlns:SOAP-ENV,attr"`
	XMLNSSoapEnc    string                 `xml:"xmlns:SOAP-ENC,attr"`
	XMLNSXSI        string                 `xml:"xmlns:xsi,attr"`
	XMLNSXSD        string                 `xml:"xmlns:xsd,attr"`
	SOAPENVEncStyle string                 `xml:"SOAP-ENV:encodingStyle,attr"`
	Body            GetSeedSOAPBodyRequest `xml:",omitempty"`
}

func (body *GetSeedBody) Fill() {
	body.XMLNSSoapEnv = httpSoapEnveloper
	body.XMLNSSoapEnc = httpSoapEnc
	body.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"
	body.XMLNSXSD = "http://www.w3.org/2001/XMLSchema"
	body.SOAPENVEncStyle = httpSoapEnc
	body.Body = GetSeedSOAPBodyRequest{
		Content: mGetSeed{
			XMLNSM: "https://palena.sii.cl/DTEWS/CrSeed.jws",
		},
	}
}

type GetSeedSOAPBodyRequest struct {
	XMLName xml.Name `xml:"SOAP-ENV:Body"`
	Content mGetSeed `xml:",omitempty"`
}

type mGetSeed struct {
	XMLName xml.Name `xml:"m:getSeed"`
	XMLNSM  string   `xml:"xmlns:m,attr"`
}

type GetSeedResponse struct {
	XMLName      xml.Name                `xml:"Envelope"`
	SOAP         string                  `xml:"soapenv,attr"`
	XMLNSSoapEnv string                  `xml:"xmlns:soapenv,attr"`
	XMLNSXSI     string                  `xml:"xmlns:xsi,attr"`
	XMLNSXSD     string                  `xml:"xmlns:xsd,attr"`
	Body         BodyGetSeedSOAPResponse `xml:",omitempty"`
}

type BodyGetSeedSOAPResponse struct {
	XMLName      xml.Name           `xml:"Body"`
	SOAP         string             `xml:"soapenv,attr"`
	SeedResponse NS1GetSeedResponse `xml:",omitempty"`
}

type NS1GetSeedResponse struct {
	XMLName     xml.Name `xml:"getSeedResponse"`
	NS1         string   `xml:"ns1,attr"`
	EncodedSeed string   `xml:"getSeedReturn,omitempty"`
}

type GetSeedResponseEncodedSeed struct {
	XMLName     xml.Name                   `xml:"getSeedReturn"`
	SiiResponse GetSeedResponseSiiResponse `xml:",omitempty"`
}

type GetSeedResponseSiiResponse struct {
	XMLName     xml.Name                   `xml:"RESPUESTA"`
	SiiRespBody GetSeedResponseSiiRespBody `xml:",omitempty"`
	SiiRespHdr  GetSeedResponseSiiRespHdr  `xml:",omitempty"`
}

type GetSeedResponseSiiRespBody struct {
	XMLName xml.Name `xml:"RESP_BODY"`
	Seed    string   `xml:"SEMILLA,omitempty"`
}

type GetSeedResponseSiiRespHdr struct {
	XMLName xml.Name `xml:"RESP_HDR"`
	Status  string   `xml:"ESTADO,omitempty"`
}
