package ctr

type Controller interface {
	Pledge() PledgeController
	Item() ItemController
	Company() CompanyController
}

type controller struct {
	pc *pledgeController
	ic *itemController
	cc *companyController
}

var ctr = controller{ &pledgeController{}, &itemController{}, &companyController{}}

func GetController() Controller {
	return &ctr
}

func (c *controller) Pledge() PledgeController {
	return c.pc
}

func (c *controller) Item() ItemController {
	return c.ic
}

func (c *controller) Company() CompanyController {
	return c.cc
}
