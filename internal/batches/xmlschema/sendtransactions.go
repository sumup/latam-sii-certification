package xmlschema

import (
	"encoding/xml"
)

type IngresarTransaccionesConCanalResponse struct {
	XMLName                  xml.Name                 `xml:"ingresarTransaccionesConCanalResponse"`
	NS2                      string                   `xml:"ns2,attr"`
	RespuestaTransaccionesTo RespuestaTransaccionesTo `xml:"RespuestaTransaccionesTo,omitempty"`
}

type RespuestaTransaccionesTo struct {
	XMLName     xml.Name    `xml:"RespuestaTransaccionesTo"`
	CodResp     int32       `xml:"codResp,omitempty"`
	DescResp    string      `xml:"descResp,omitempty"`
	TrackIDs    []*int64    `xml:"tracksID,omitempty"`
	LineasError interface{} `xml:"lineasError,omitempty"`
}

type Response struct {
	XMLName xml.Name         `xml:"Envelope"`
	SOAP    string           `xml:"soap,attr"`
	Body    SOAPBodyResponse `xml:"Body,omitempty"`
}

func (response *Response) PreFillZeroValues() {
	response.Body.Resp.RespuestaTransaccionesTo.CodResp = -1
}

type SOAPBodyResponse struct {
	XMLName xml.Name                              `xml:"Body"`
	Resp    IngresarTransaccionesConCanalResponse `xml:",omitempty"`
}

type SendTransactionsBody struct {
	XMLName      xml.Name        `xml:"soapenv:Envelope"`
	XMLNSSoapEnv string          `xml:"xmlns:soapenv,attr"`
	XMLNSWS      string          `xml:"xmlns:ws,attr"`
	Header       string          `xml:"soapenv:Header"`
	Body         SOAPBodyRequest `xml:",omitempty"`
}

func (body *SendTransactionsBody) Fill(transactions []TaxTransaccion) {
	body.XMLNSSoapEnv = httpSoapEnveloper
	body.XMLNSWS = "http://ws.comprobanteboleta.oia.sdi.sii.cl"
	body.Body = SOAPBodyRequest{
		Content: IngresarTransaccionesConCanal{
			Transacciones: transactions,
		},
	}
}

type SOAPBodyRequest struct {
	XMLName xml.Name                      `xml:"soapenv:Body"`
	Content IngresarTransaccionesConCanal `xml:",omitempty"`
}

type IngresarTransaccionesConCanal struct {
	XMLName       xml.Name         `xml:"ws:ingresarTransaccionesConCanal"`
	Transacciones []TaxTransaccion `xml:"transaccion,omitempty" json:"transaccion,omitempty"`
}

type TaxTransaccion struct {
	RutInformante         string `xml:"rutInformante,omitempty" json:"rutInformante,omitempty"`
	DvInformante          string `xml:"dvInformante,omitempty" json:"dvInformante,omitempty"`
	RutContribuyente      string `xml:"rutContribuyente,omitempty" json:"rutContribuyente,omitempty"`
	DvContribuyente       string `xml:"dvContribuyente,omitempty" json:"dvContribuyente,omitempty"`
	TipoDocumento         string `xml:"tipoDocumento,omitempty" json:"tipoDocumento,omitempty"`
	FechaVenta            string `xml:"fechaVenta,omitempty" json:"fechaVenta,omitempty"`
	TotalMontoNeto        string `xml:"totalMontoNeto,omitempty" json:"totalMontoNeto,omitempty"`
	TotalMontoExento      string `xml:"totalMontoExento,omitempty" json:"totalMontoExento,omitempty"`
	TotalMontoTotal       string `xml:"totalMontoTotal,omitempty" json:"totalMontoTotal,omitempty"`
	TotalMontoPropina     string `xml:"totalMontoPropina,omitempty" json:"totalMontoPropina,omitempty"`
	TotalMontoVuelto      string `xml:"totalMontoVuelto,omitempty" json:"totalMontoVuelto,omitempty"`
	TotalMontoDonacion    string `xml:"totalMontoDonacion,omitempty" json:"totalMontoDonacion,omitempty"`
	TotalMontoTransaccion string `xml:"totalMontoTransaccion,omitempty" json:"totalMontoTransaccion,omitempty"`
	TotalValesEmitidos    string `xml:"totalValesEmitidos,omitempty" json:"totalValesEmitidos,omitempty"`
	IdentificadorEnvio    string `xml:"identificadorEnvio,omitempty" json:"identificadorEnvio,omitempty"`
	CanalTransaccion      string `xml:"canalTransaccion,omitempty" json:"canalTransaccion,omitempty"`
}
