package constants

import "time"

var Version string = "dev"

var InstanceStartedAt time.Time = time.Now()

const (
	DonationAddress   = "0xEAee1B084bb5a041b639153F6D5946277Cf50462"
	UserTextPriceDown = "🔴"
	UserTextPriceUp   = "🟢"
)
