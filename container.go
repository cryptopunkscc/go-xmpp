package xmpp

type Container struct {
	Children []Template
}

// ChildCount returns the number of children with a given name
func (c *Container) ChildCount(name string) int {
	if c.Children == nil {
		return 0
	}

	count := 0

	for _, child := range c.Children {
		n, _ := Identify(child)
		if n == name {
			count++
		}
	}

	return count
}

// Child returns the first child with a given name
func (c *Container) Child(name string) Template {
	if c.Children == nil {
		return nil
	}

	for _, child := range c.Children {
		n, _ := Identify(child)
		if n == name {
			return child
		}
	}

	return nil
}

func (c *Container) AddChild(child Template) {
	if c.Children == nil {
		c.Children = make([]Template, 0)
	}

	c.Children = append(c.Children, child)
}
