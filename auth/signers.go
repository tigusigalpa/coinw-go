// Package auth implements the two isolated CoinW signing profiles.
//
// Spot private REST:   sign = UPPERCASE MD5 of sorted "k=v&" pairs + "secret_key=…"
// Futures private REST: sign = base64(HMAC-SHA256(secret, ts+METHOD+path+query|body))
//
// The profiles are never mixed — each endpoint selects exactly one.
package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
// Spot
// ---------------------------------------------------------------------------

// SpotCanonicalString builds the documented sign input:
// sorted "k=v&" entries (api_key included) + "secret_key={secret}".
// Params must NOT already contain api_key — it is added here.
func SpotCanonicalString(params map[string]interface{}, apiKey, secretKey string) string {
	p := make(map[string]interface{}, len(params)+1)
	for k, v := range params {
		p[k] = v
	}
	p["api_key"] = apiKey

	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(spotStringify(p[k]))
		b.WriteByte('&')
	}
	b.WriteString("secret_key=")
	b.WriteString(secretKey)
	return b.String()
}

// SpotSign returns the uppercase MD5 hex signature.
func SpotSign(params map[string]interface{}, apiKey, secretKey string) string {
	sum := md5.Sum([]byte(SpotCanonicalString(params, apiKey, secretKey)))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// SpotSignedURL builds the full documented request URL:
// {base}{apiUrl}&sign={sign}&{urlencoded params}. The apiUrl already embeds
// "?command=…", hence "&" before sign.
func SpotSignedURL(baseURL, apiURL string, params map[string]interface{}, apiKey, secretKey string) string {
	sign := SpotSign(params, apiKey, secretKey)

	p := make(map[string]interface{}, len(params)+1)
	for k, v := range params {
		p[k] = v
	}
	p["api_key"] = apiKey

	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(baseURL)
	b.WriteString(apiURL)
	b.WriteString("&sign=")
	b.WriteString(sign)
	for _, k := range keys {
		b.WriteByte('&')
		b.WriteString(url.QueryEscape(k))
		b.WriteByte('=')
		b.WriteString(url.QueryEscape(spotStringify(p[k])))
	}
	return b.String()
}

// StringifyParam converts a scalar param to its documented string form.
func StringifyParam(v interface{}) string {
	return spotStringify(v)
}

func spotStringify(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case bool:
		if t {
			return "true"
		}
		return "false"
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		return fmt.Sprint(v)
	}
}

// ---------------------------------------------------------------------------
// Futures
// ---------------------------------------------------------------------------

// FuturesHeaders carries the documented auth headers.
type FuturesHeaders struct {
	Sign      string
	APIKey    string
	Timestamp string
}

// futuresSign computes base64(HMAC-SHA256(secret, ts + METHOD + apiURL + suffix)).
func futuresSign(method, apiURL, suffix, secretKey string, ts int64) FuturesHeaders {
	payload := strconv.FormatInt(ts, 10) + strings.ToUpper(method) + apiURL + suffix
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	return FuturesHeaders{
		Sign:      base64.StdEncoding.EncodeToString(mac.Sum(nil)),
		Timestamp: strconv.FormatInt(ts, 10),
	}
}

// FuturesCanonicalQuery builds the unescaped "k=v&k=v" signing input,
// preserving key order, as in CoinW's official Python authentication snippet.
// HTTP query escaping is applied separately after signing, like requests.get.
func FuturesCanonicalQuery(params map[string]interface{}, keys []string) string {
	var b strings.Builder
	first := true
	for _, k := range keys {
		v, ok := params[k]
		if !ok || v == nil {
			continue
		}
		if !first {
			b.WriteByte('&')
		}
		first = false
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(spotStringify(v))
	}
	return b.String()
}

// SortedParamKeys returns params' keys in deterministic order (nulls kept;
// they are skipped when building the query).
func SortedParamKeys(params map[string]interface{}) []string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FuturesSignedGet signs a GET: returns the literal query to append and the
// auth headers (sign, api_key, timestamp).
func FuturesSignedGet(apiURL string, params map[string]interface{}, apiKey, secretKey string, ts int64) (query string, headers FuturesHeaders) {
	canonical := FuturesCanonicalQuery(params, SortedParamKeys(params))
	values := url.Values{}
	for k, v := range params {
		if v != nil {
			values.Set(k, spotStringify(v))
		}
	}
	query = values.Encode()
	suffix := ""
	if canonical != "" {
		suffix = "?" + canonical
	}
	headers = futuresSign("GET", apiURL, suffix, secretKey, ts)
	headers.APIKey = apiKey
	return query, headers
}

// FuturesCanonicalBody produces the JSON body whose exact bytes are signed
// and sent. Pass a non-nil empty map/struct for a documented "{}" body.
func FuturesCanonicalBody(body interface{}) (string, error) {
	if body == nil {
		body = map[string]interface{}{}
	}
	if params, ok := body.(map[string]interface{}); ok && params == nil {
		body = map[string]interface{}{}
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// FuturesSignedBody signs a non-GET request: returns the exact JSON body and
// auth headers (+ Content-Type: application/json).
func FuturesSignedBody(method, apiURL string, body interface{}, apiKey, secretKey string, ts int64) (payload string, headers FuturesHeaders, err error) {
	payload, err = FuturesCanonicalBody(body)
	if err != nil {
		return "", FuturesHeaders{}, err
	}
	headers = futuresSign(method, apiURL, payload, secretKey, ts)
	headers.APIKey = apiKey
	return payload, headers, nil
}
