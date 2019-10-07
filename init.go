package xmpp

func init() {
	addCoreElements()
	addErrorElements()
	addTLSElements()
	addSASLElements()
	addFeaturesElements()
	initIQ()
	initPresence()
	initMessage()
}
