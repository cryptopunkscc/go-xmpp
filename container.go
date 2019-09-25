package xmpp

type Container struct {
	Children []interface{} `xml:"-"`
}

// ChildCount returns the number of children with a given name
func (c *Container) ChildCount(t interface{}) int {
	if c.Children == nil {
		return 0
	}

	name := Identify(t)
	count := 0

	for _, child := range c.Children {
		id := Identify(child)
		if name == id {
			count++
		}
	}

	return count
}

// Child returns the first child with a given name
func (c *Container) Child(t interface{}) interface{} {
	if c.Children == nil {
		return nil
	}

	name := Identify(t)

	for _, child := range c.Children {
		id := Identify(child)
		if name == id {
			return child
		}
	}

	return nil
}

func (c *Container) AddChild(child interface{}) {
	if c.Children == nil {
		c.Children = make([]interface{}, 0)
	}

	c.Children = append(c.Children, child)
}
