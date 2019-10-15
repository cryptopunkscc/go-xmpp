package disco

import "github.com/cryptopunkscc/go-xmpp/xep0030"

// Info stores disco#info content
type Info struct {
	JID        string
	Features   []string
	Identities []Identity
}

// Identity stores identity information found in disco#info
type Identity struct {
	Category string
	Type     string
	Name     string
	Lang     string
}

// AddFeature adds an item to the features list
func (i *Info) AddFeature(f string) {
	if i.Features == nil {
		i.Features = make([]string, 0)
	}
	i.Features = append(i.Features, f)
}

// AddIdentity adds an item to the identities list
func (i *Info) AddIdentity(id Identity) {
	if i.Identities == nil {
		i.Identities = make([]Identity, 0)
	}
	i.Identities = append(i.Identities, id)
}

func (info *Info) queryInfo() *xep0030.QueryInfo {
	qi := &xep0030.QueryInfo{
		Identities: make([]xep0030.Identity, 0),
		Features:   make([]xep0030.Feature, 0),
	}
	for _, i := range info.Identities {
		qi.Identities = append(qi.Identities, xep0030.Identity{
			Category: i.Category,
			Type:     i.Type,
			Name:     i.Name,
			Lang:     i.Lang,
		})
	}
	for _, f := range info.Features {
		qi.Features = append(qi.Features, xep0030.Feature{
			Var: f,
		})
	}
	return qi
}

func queryInfoToInfo(q *xep0030.QueryInfo) *Info {
	info := &Info{
		Features:   make([]string, 0),
		Identities: make([]Identity, 0),
	}
	for _, f := range q.Features {
		info.Features = append(info.Features, f.Var)
	}
	for _, i := range q.Identities {
		info.Identities = append(info.Identities, Identity{
			Category: i.Category,
			Name:     i.Name,
			Type:     i.Type,
			Lang:     i.Lang,
		})
	}
	return info
}
