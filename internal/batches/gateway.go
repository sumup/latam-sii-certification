package batches

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/sumup/sii-certification/internal/adapters"
	"github.com/sumup/sii-certification/internal/appconfig"
	"github.com/sumup/sii-certification/internal/batches/xmlschema"
	"github.com/sumup/sii-certification/internal/entities"
	"github.com/sumup/sii-certification/internal/utils"
)

type Gateway struct {
	client adapters.IAdapter
	config appconfig.TaxAutChile
}

func NewTaxAuthorityGateway(client adapters.IAdapter, config appconfig.TaxAutChile) *Gateway {
	return &Gateway{
		client,
		config,
	}
}

const successResponseCode = 0
const successResponseCodeString = "00"

func (g *Gateway) SendTransactions(ctx context.Context, token string, req xmlschema.SendTransactionsBody) (xmlschema.Response, string, error) {
	response := xmlschema.Response{}
	response.PreFillZeroValues()

	reqString, err := xml.MarshalIndent(req, " ", "  ")
	if err != nil {
		return response, "", errCantMarshalRequest
	}

	headers := map[string]string{
		"Cookie": fmt.Sprintf("TOKEN=%s", token),
	}

	respString, err := g.client.Post(ctx, g.config.BaseURL+g.config.SendTransactionsEndpoint, headers, reqString)

	if err != nil {
		return response, string(respString), errClientCallFailed
	}

	err = xml.Unmarshal(respString, &response)

	if err != nil {
		return response, string(respString), errCantUnmarshalResponse
	}

	return response, string(respString), nil
}

func (g *Gateway) GetSeed(ctx context.Context) (string, error) {
	var req xmlschema.GetSeedBody

	req.Fill()

	reqString, err := xml.Marshal(req)
	if err != nil {
		return "", errCantMarshalGetSeedRequest
	}

	headers := map[string]string{
		"SOAPAction": "",
	}

	respString, err := g.client.Post(ctx, g.config.AuthURL+g.config.GetSeedEndpoint, headers, reqString)
	if err != nil {
		return "", errClientCallFailed
	}

	response := xmlschema.GetSeedResponse{}
	if err = xml.Unmarshal(respString, &response); err != nil {
		return "", errCantUnmarshalResponse
	}

	var seedResponse xmlschema.GetSeedResponseSiiResponse
	if err = xml.Unmarshal([]byte(response.Body.SeedResponse.EncodedSeed), &seedResponse); err != nil {
		return "", errCantUnmarshalEncodedSeed
	}

	if seedResponse.SiiRespHdr.Status != successResponseCodeString {
		return "", errGetSeedReturnsError
	}

	return seedResponse.SiiRespBody.Seed, nil
}

func (g *Gateway) GetToken(ctx context.Context, seed string) (string, error) {
	var tokenObject xmlschema.TokenObject

	tokenObject.Item.Semilla = seed

	canonTokenObject, err := utils.CanonicalizeXML(tokenObject)
	if err != nil {
		return "", errCantCanonicalizeTokenObject
	}

	digestValue := utils.Sha1HashBase64(canonTokenObject)

	var signedInfo xmlschema.GetTokenSignedInfo

	signedInfo.FillSignedInfo(digestValue)

	canonSignedInfo, err := utils.CanonicalizeXML(signedInfo)
	if err != nil {
		return "", errCantCanonicalizeSignedInfo
	}

	signatureValue := utils.Sha1HashBase64(canonSignedInfo)

	var keyInfo xmlschema.GetTokenKeyInfo

	keyInfo.KeyValue.RSAKeyValue.Modulus = g.config.SumUpCertificateModulus
	keyInfo.KeyValue.RSAKeyValue.Exponent = g.config.SumUpCertificateExponent
	keyInfo.X509Data.X509Certificate = g.config.SumUpCertificate

	var signature xmlschema.GetTokenSignature

	signature.Fill(signedInfo, signatureValue, keyInfo)

	tokenObject.GetTokenSignature = signature

	marshalledToken, err := xml.Marshal(tokenObject)
	if err != nil {
		return "", errCantMarshalTokenObject
	}

	getTokenRequest := xmlschema.GetTokenRequest{}
	getTokenRequest.Fill()

	getTokenRequest.Body.Content.PSZXml = string(marshalledToken)

	marshalledRequest, err := xml.Marshal(getTokenRequest)
	if err != nil {
		return "", errCantMarshalGetTokenRequest
	}

	headers := map[string]string{
		"SOAPAction": " ",
	}

	tokenRespString, err := g.client.Post(ctx, g.config.AuthURL+g.config.GetTokenEndpoint, headers, marshalledRequest)
	if err != nil {
		return "", errClientCallFailed
	}

	getTokenResponse := xmlschema.GetTokenResponse{}
	if err = xml.Unmarshal(tokenRespString, &getTokenResponse); err != nil {
		return "", errCantUnmarshalResponse
	}

	var tokenResponse xmlschema.GetTokenResponseSiiResponse
	if err = xml.Unmarshal([]byte(getTokenResponse.Body.TokenResponse.EncodedTokenResponse), &tokenResponse); err != nil {
		return "", errCantUnmarshalEncodedToken
	}

	if tokenResponse.SiiRespHdr.Status != successResponseCodeString {
		return "", errGetTokenReturnsError
	}

	return tokenResponse.SiiRespBody.Token, nil
}

func (g *Gateway) SendMany(ctx context.Context, token string, batchesMatrix [][]entities.Batch) [][]entities.Batch {
	for i := range batchesMatrix {
		transactions := taxTransactions(batchesMatrix[i], g.config.SumUpRutBase, g.config.SumUpRutVerifier)

		var req xmlschema.SendTransactionsBody

		req.Fill(transactions)

		res, rAsString, err := g.SendTransactions(ctx, token, req)
		fmt.Println(rAsString)
		if err != nil {
			modifyBatchesWithResponse(&batchesMatrix[i], entities.FailedStatus, err.Error())

			continue
		}

		if res.Body.Resp.RespuestaTransaccionesTo.CodResp != successResponseCode {
			modifyBatchesWithResponse(&batchesMatrix[i], entities.SentButRejected, rAsString)

			continue
		}

		for j, trackID := range res.Body.Resp.RespuestaTransaccionesTo.TrackIDs {
			batchesMatrix[i][j].Status = entities.SentAndAccepted
			batchesMatrix[i][j].LastResponse = rAsString
			batchesMatrix[i][j].ExternalTrackID = strconv.FormatInt(*trackID, 10)
		}
	}

	return batchesMatrix
}

func modifyBatchesWithResponse(batches *[]entities.Batch, status int, lastResponse string) {
	for i := range *batches {
		(*batches)[i].Status = status
		(*batches)[i].LastResponse = lastResponse
	}
}

func taxTransactions(batches []entities.Batch, rutSumUp, dvSumUp string) []xmlschema.TaxTransaccion {
	transactions := make([]xmlschema.TaxTransaccion, 0)

	for _, batch := range batches {
		vatIDDivided := strings.Split(batch.VatID, "-")

		//docType := documentType(batch.HasTaxes)
		//channel := channelType(batch.IsCNP)

		//trackID := trackIDOrDefault(batch.ExternalTrackID)

		taxTx := xmlschema.TaxTransaccion{
			RutInformante:         rutSumUp,
			DvInformante:          dvSumUp,
			RutContribuyente:      vatIDDivided[0],
			DvContribuyente:       vatIDDivided[1],
			TipoDocumento:         batch.DocumentType,
			FechaVenta:            batch.Day,
			TotalMontoNeto:        strconv.FormatUint(batch.Amount, 10),
			TotalMontoExento:      "0",
			TotalMontoTotal:       strconv.FormatUint(batch.Amount, 10),
			TotalMontoPropina:     "0",
			TotalMontoVuelto:      "0",
			TotalMontoDonacion:    "0",
			TotalMontoTransaccion: strconv.FormatUint(batch.Amount, 10),
			TotalValesEmitidos:    strconv.Itoa(batch.NTransactions),
			IdentificadorEnvio:    batch.ExternalTrackID,
			CanalTransaccion:      batch.Channel,
		}
		transactions = append(transactions, taxTx)
	}

	return transactions
}
