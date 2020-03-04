package mdns

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"

	//log "github.com/giant-tech/go-service/base/plugin/logger"

	log "github.com/cihub/seelog"
)

func encode(txt *mdnsTxt) ([]string, error) {
	b, err := json.Marshal(txt)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var buf bytes.Buffer
	defer buf.Reset()

	w := zlib.NewWriter(&buf)
	if _, err := w.Write(b); err != nil {
		log.Error(err)
		return nil, err
	}
	w.Close()

	encoded := hex.EncodeToString(buf.Bytes())

	// individual txt limit
	if len(encoded) <= 255 {
		return []string{encoded}, nil
	}

	// split encoded string
	var record []string

	for len(encoded) > 255 {
		record = append(record, encoded[:255])
		encoded = encoded[255:]
	}

	record = append(record, encoded)

	return record, nil
}

func decode(record []string) (*mdnsTxt, error) {
	encoded := strings.Join(record, "")

	hr, err := hex.DecodeString(encoded)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	br := bytes.NewReader(hr)
	zr, err := zlib.NewReader(br)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	rbuf, err := ioutil.ReadAll(zr)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var txt *mdnsTxt

	if err := json.Unmarshal(rbuf, &txt); err != nil {
		log.Error(err)
		return nil, err
	}

	return txt, nil
}
