package xmlschema

import (
	"encoding/xml"
)

type TokenObject struct {
	XMLName           xml.Name          `xml:"getToken"`
	Item              TokenObjectItem   `xml:",omitempty"`
	GetTokenSignature GetTokenSignature `xml:",omitempty"`
}

type TokenObjectItem struct {
	XMLName xml.Name `xml:"item"`
	Semilla string   `xml:"Semilla,omitempty"`
}

type GetTokenSignedInfo struct {
	XMLName              xml.Name            `xml:"SignedInfo"`
	CanonalizationMethod string              `xml:"CanonalizationMethod,attr"`
	SignatureMethod      string              `xml:"SignatureMethod,attr"`
	Reference            SignedInfoReference `xml:",omitempty"`
}

type SignedInfoReference struct {
	XMLName              xml.Name             `xml:"Reference"`
	URI                  string               `xml:"URI,attr"`
	SignedInfoTransforms SignedInfoTransforms `xml:",omitempty"`
}

type SignedInfoTransforms struct {
	XMLName                       xml.Name `xml:"Transforms"`
	SignedInfoTransformsTransform `xml:",omitempty"`
	DigestMethod                  string `xml:"DigestMethod,attr"`
	DigestValue                   string `xml:",omitempty"`
}

type SignedInfoTransformsTransform struct {
	XMLName   xml.Name `xml:"Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
}

func (i *GetTokenSignedInfo) FillSignedInfo(hash string) {
	i.CanonalizationMethod = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	i.SignatureMethod = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	i.Reference.URI = ""
	i.Reference.SignedInfoTransforms.SignedInfoTransformsTransform.Algorithm = "http://www.w3.org/2000/09/xmldsig#enveloped-signature"
	i.Reference.SignedInfoTransforms.DigestMethod = "http://www.w3.org/2000/09/xmldsig#sha1"
	i.Reference.SignedInfoTransforms.DigestValue = hash
}

type GetTokenSignature struct {
	XMLName         xml.Name           `xml:"Signature"`
	XMLNS           string             `xml:"xmlns,attr"`
	SignedInfo      GetTokenSignedInfo `xml:",omitempty"`
	SignatureValue  string             `xml:"SignatureValue,attr"`
	GetTokenKeyInfo GetTokenKeyInfo    `xml:",omitempty"`
}

func (s *GetTokenSignature) Fill(info GetTokenSignedInfo, signatureValue string, keyInfo GetTokenKeyInfo) {
	s.XMLNS = "http://www.w3.org/2000/09/xmldsig#"
	s.SignedInfo = info
	s.SignatureValue = signatureValue
	s.GetTokenKeyInfo = keyInfo
}

type GetTokenKeyInfo struct {
	XMLName  xml.Name         `xml:"KeyInfo"`
	KeyValue GetTokenKeyValue `xml:",omitempty"`
	X509Data GetTokenX509Data `xml:",omitempty"`
}

type GetTokenKeyValue struct {
	XMLName     xml.Name            `xml:"KeyValue"`
	RSAKeyValue GetTokenRSAKeyValue `xml:",omitempty"`
}

type GetTokenX509Data struct {
	XMLName         xml.Name `xml:"X509Data"`
	X509Certificate string   `xml:"X509Certificate,omitempty"`
}

type GetTokenRSAKeyValue struct {
	XMLName  xml.Name `xml:"RSAKeyValue"`
	Modulus  string   `xml:"Modulus,omitempty"`
	Exponent string   `xml:"Exponent,omitempty"`
}

type GetTokenRequest struct {
	XMLName         xml.Name            `xml:"SOAP-ENV:Envelope"`
	XMLNSSoapEnv    string              `xml:"xmlns:SOAP-ENV,attr"`
	XMLNSSoapEnc    string              `xml:"xmlns:SOAP-ENC,attr"`
	XMLNSXSI        string              `xml:"xmlns:xsi,attr"`
	XMLNSXSD        string              `xml:"xmlns:xsd,attr"`
	SOAPENVEncStyle string              `xml:"SOAP-ENV:encodingStyle,attr"`
	Body            BodyGetTokenRequest `xml:",omitempty"`
}

type BodyGetTokenRequest struct {
	XMLName xml.Name  `xml:"SOAP-ENV:Body"`
	Content mGetToken `xml:",omitempty"`
}

type mGetToken struct {
	XMLName xml.Name `xml:"m:getToken"`
	XMLNSM  string   `xml:"xmlns:m,attr"`
	PSZXml  string   `xml:"pszXml,omitempty"`
}

func (body *GetTokenRequest) Fill() {
	body.XMLNSSoapEnv = httpSoapEnveloper
	body.XMLNSSoapEnc = httpSoapEnc
	body.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"
	body.XMLNSXSD = "http://www.w3.org/2001/XMLSchema"
	body.SOAPENVEncStyle = httpSoapEnc
	body.Body = BodyGetTokenRequest{
		Content: mGetToken{
			XMLNSM: "https://palena.sii.cl/DTEWS/GetTokenFromSeed.jws",
		},
	}
}

type GetTokenResponse struct {
	XMLName      xml.Name                 `xml:"Envelope"`
	SOAPEnv      string                   `xml:"soapenv,attr"`
	XMLNSSoapEnv string                   `xml:"xmlns:soapenv,attr"`
	XMLNSXSI     string                   `xml:"xmlns:xsi,attr"`
	XMLNSXSD     string                   `xml:"xmlns:xsd,attr"`
	Body         BodyGetTokenSOAPResponse `xml:",omitempty"`
}

type BodyGetTokenSOAPResponse struct {
	XMLName       xml.Name            `xml:"Body"`
	SOAPEnv       string              `xml:"soapenv,attr"`
	TokenResponse NS1GetTokenResponse `xml:",omitempty"`
}

type NS1GetTokenResponse struct {
	XMLName              xml.Name `xml:"getTokenResponse"`
	NS1                  string   `xml:"ns1,attr"`
	EncodedTokenResponse string   `xml:"getTokenReturn,omitempty"`
}

type GetTokenResponseSiiResponse struct {
	XMLName     xml.Name            `xml:"RESPUESTA"`
	SiiRespBody GetTokenSiiRespBody `xml:",omitempty"`
	SiiRespHdr  GetTokenSiiRespHdr  `xml:",omitempty"`
}

type GetTokenSiiRespBody struct {
	XMLName xml.Name `xml:"RESP_BODY"`
	Token   string   `xml:"TOKEN,omitempty"`
}

type GetTokenSiiRespHdr struct {
	XMLName xml.Name `xml:"RESP_HDR"`
	Status  string   `xml:"ESTADO,omitempty"`
	Message string   `xml:"GLOSA,omitempty"`
}
