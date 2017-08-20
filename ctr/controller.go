package ctr

type Controller interface {
	Pledge() PledgeController
	Item() ItemController
}

type controller struct {
	pc *pledgeController
	ic *itemController
}

var ctr = controller{ &pledgeController{}, &itemController{}}

func GetController() Controller {
	return &ctr
}

func (c *controller) Pledge() PledgeController {
	return c.pc
}

func (c *controller) Item() ItemController {
	return c.ic
}

