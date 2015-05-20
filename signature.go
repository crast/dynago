package dynago

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

/*
The functions in this file implement the AWS v4 request signing algorithm.

This is required for all aws requests to ensure:

	1. request bodies and headers are not tampered with in flight.
	2. It also prevents replay attacks
	3. It also handles authentication without sending or revealing the shared secret
*/

const algorithm = "AWS4-HMAC-SHA256"

type awsInfo struct {
	AccessKey string
	SecretKey string
	Region    string
	Service   string
}

func (info *awsInfo) signRequest(request *http.Request, bodyBytes []byte) {
	now := time.Now().UTC()
	isoDateSmash := now.Format("20060102T150405Z")
	request.Header.Add("x-amz-date", isoDateSmash)
	canonicalHash, signedHeaders := canonicalRequest(request, bodyBytes)
	credentialScope := now.Format("20060102") + "/" + info.Region + "/" + info.Service + "/aws4_request"
	stringToSign := algorithm + "\n" + isoDateSmash + "\n" + credentialScope + "\n" + canonicalHash
	signingKey := signingKey(now, info)
	signature := hex.EncodeToString(hmacShort(signingKey, []byte(stringToSign)))
	authHeader := algorithm + " Credential=" + info.AccessKey + "/" + credentialScope + ", SignedHeaders=" + signedHeaders + ", Signature=" + signature
	fmt.Println(authHeader)
	request.Header.Add("Authorization", authHeader)
}

func canonicalRequest(request *http.Request, bodyBytes []byte) (string, string) {
	var canonical bytes.Buffer
	canonical.WriteString(request.Method)
	canonical.WriteByte('\n')
	canonical.WriteString(request.URL.Path)
	canonical.WriteRune('\n')
	canonical.WriteString(request.URL.RawQuery)
	canonical.WriteRune('\n')
	signedHeaders := canonicalHeaders(&canonical, request.Header)
	sum := sha256.Sum256(bodyBytes)
	canonical.WriteString(hex.EncodeToString(sum[:]))
	cBytes := canonical.Bytes()
	sum = sha256.Sum256(cBytes)
	return hex.EncodeToString(sum[:]), signedHeaders
}

func canonicalHeaders(buf *bytes.Buffer, headers http.Header) string {
	headerVals := make([]string, 0, len(headers))
	headerNames := make([]string, 0, len(headers))
	for key, val := range headers {
		name := strings.ToLower(key)
		s := name + ":" + strings.TrimSpace(val[0])
		headerVals = append(headerVals, s)
		headerNames = append(headerNames, name)
	}
	sort.Strings(headerVals)
	for _, cHeader := range headerVals {
		buf.WriteString(cHeader)
		buf.WriteRune('\n')
	}
	buf.WriteRune('\n')
	sort.Strings(headerNames)
	signedHeaders := strings.Join(headerNames, ";")
	buf.WriteString(signedHeaders)
	buf.WriteRune('\n')
	return signedHeaders
}

func signingKey(now time.Time, info *awsInfo) []byte {
	kSecret := "AWS4" + info.SecretKey
	kDate := hmacShort([]byte(kSecret), []byte(now.Format("20060102")))
	kRegion := hmacShort(kDate, []byte(info.Region))
	kService := hmacShort(kRegion, []byte(info.Service))
	kSigning := hmacShort(kService, []byte("aws4_request"))
	return kSigning
}

func hmacShort(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
