package disco

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/xep0030"
)

type Item struct {
	JID  xmpp.JID
	Node string
	Name string
}

type Items struct {
	Node  string
	Items []Item
}

// Add adds an item to the item list
func (i *Items) Add(item Item) {
	if i.Items == nil {
		i.Items = make([]Item, 0)
	}
	i.Items = append(i.Items, item)
}

func (items *Items) queryItems() *xep0030.QueryItems {
	qi := &xep0030.QueryItems{
		Items: make([]xep0030.Item, 0),
		Node:  items.Node,
	}
	for _, i := range items.Items {
		qi.Items = append(qi.Items, xep0030.Item{
			JID:  i.JID,
			Name: i.Name,
			Node: i.Node,
		})
	}
	return qi
}

func queryItemsToItems(q *xep0030.QueryItems) *Items {
	r := &Items{
		Items: make([]Item, 0),
	}
	for _, i := range q.Items {
		r.Items = append(r.Items, Item{
			JID:  i.JID,
			Name: i.Name,
			Node: i.Node,
		})
	}
	return r
}
