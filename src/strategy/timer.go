package strategy

import "github.com/robfig/cron"

func init() {
	c := cron.New()

	c.AddFunc("0 0 0-4 * * ?", UpdateKlineList)

}

func UpdateKlineList() {

}
