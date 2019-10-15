package caps

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cryptopunkscc/go-xmpp/ext/disco"
	"hash"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/xep0115"
)

type Caps struct {
	session xmpp.Session
	Disco   *disco.Disco
}

type CapsCache interface {
}

type UserCaps struct {
	JID  string
	Ver  string
	Hash string
	Node string
}

func (r *Caps) Online(s xmpp.Session) {
	r.session = &xmpp.Callbacks{Session: s}
}

func (r *Caps) Offline(error) {
	r.session = nil
}

func (r *Caps) HandleStanza(s xmpp.Stanza) {
	xmpp.HandleStanza(r, s)
}

func (caps *Caps) onDiscoInfo(info *disco.Info) {
	h := DiscoHash(info, sha1.New())
	fmt.Println("DiscoInfo calculated ver:", h)
}

func (caps *Caps) HandlePresence(s *xmpp.Presence) {
	if c, ok := s.Child(&xep0115.Capability{}).(*xep0115.Capability); ok {
		b64 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(c.Ver))
		bytes, _ := ioutil.ReadAll(b64)
		sha := hex.EncodeToString(bytes)

		log.Println(s.From, "sent us caps!", c.Ver, sha)

		// &UserCaps{
		// 	JID:  s.From,
		// 	Ver:  c.Ver,
		// 	Hash: c.Hash,
		// 	Node: c.Node,
		// }
	}
}

func DiscoHash(d *disco.Info, hash hash.Hash) string {
	builder := strings.Builder{}

	ids := make([]string, 0)
	for _, id := range d.Identities {
		s := fmt.Sprintf("%s/%s/%s/%s", id.Category, id.Type, id.Lang, id.Name)
		ids = append(ids, s)
	}
	sort.Strings(ids)
	for _, id := range ids {
		builder.WriteString(id)
		builder.WriteByte('<')
	}

	feats := append([]string{}, d.Features...)
	sort.Strings(feats)
	for _, f := range feats {
		builder.WriteString(f)
		builder.WriteByte('<')
	}

	rawbytes := builder.String()
	fmt.Println("Concatenated:", rawbytes)
	hash.Write([]byte(rawbytes))
	sum := hash.Sum(nil)
	fmt.Println("SHA1:", hex.EncodeToString(sum))

	buff := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, buff)
	enc.Write(sum)
	enc.Close()
	return buff.String()
}
